<script lang="ts">
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Users, CreditCard, TrendingUp, Calendar, LogOut } from 'lucide-svelte';
	import { login, getAdminStats, adminListUsers } from '$lib/api/gosaas';
	import type { AdminStatsResponse, AdminUser } from '$lib/api/gosaasComponents';

	let username = $state('');
	let password = $state('');
	let isAuthenticated = $state(false);
	let authError = $state('');
	let loading = $state(false);

	let stats = $state<AdminStatsResponse | null>(null);
	let recentUsers = $state<AdminUser[]>([]);

	// Check if admin token exists
	function getAdminToken(): string | null {
		if (typeof window === 'undefined') return null;
		return localStorage.getItem('gosaas_admin_token');
	}

	function storeAdminToken(token: string) {
		localStorage.setItem('gosaas_token', token);
		localStorage.setItem('gosaas_admin_token', token);
	}

	function clearAdminToken() {
		localStorage.removeItem('gosaas_token');
		localStorage.removeItem('gosaas_admin_token');
		isAuthenticated = false;
		stats = null;
		recentUsers = [];
	}

	async function handleLogin() {
		if (!username || !password) {
			authError = 'Please enter username and password';
			return;
		}

		loading = true;
		authError = '';

		try {
			// Use the standard login API
			const response = await login({ email: username, password });

			// Check if this is an admin login (backend returns /backoffice as checkoutUrl)
			if (response.checkoutUrl !== '/backoffice') {
				authError = 'Admin access required';
				return;
			}

			// Store the token for subsequent API calls
			storeAdminToken(response.token);
			isAuthenticated = true;

			// Load admin data
			await loadAdminData();
		} catch (err) {
			authError = err instanceof Error ? err.message : 'Login failed';
		} finally {
			loading = false;
		}
	}

	async function loadAdminData() {
		loading = true;
		try {
			// Use auto-generated API functions
			const statsData = await getAdminStats();
			stats = statsData;

			const usersData = await adminListUsers({ page: 1, pageSize: 10 });
			recentUsers = usersData.users;
		} catch (err) {
			console.error('Failed to load admin data:', err);
			// If we get 401/403, clear auth state
			if (err instanceof Error && err.message.includes('401')) {
				clearAdminToken();
			}
		} finally {
			loading = false;
		}
	}

	async function refreshData() {
		await loadAdminData();
	}

	onMount(async () => {
		// Check for existing admin token
		const token = getAdminToken();
		if (token) {
			isAuthenticated = true;
			await loadAdminData();
		}
	});

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function formatCurrency(amount: number, currency: string): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: currency.toUpperCase()
		}).format(amount / 100);
	}
</script>

<svelte:head>
	<title>Backoffice - SaaS Admin</title>
	<meta name="robots" content="noindex, nofollow" />
</svelte:head>

{#if !isAuthenticated}
	<div class="min-h-screen flex items-center justify-center p-4">
		<Card class="w-full max-w-md">
			<div class="text-center mb-6">
				<h1 class="font-display text-2xl font-bold text-white mb-2">Backoffice Login</h1>
				<p class="text-sm text-base-400">Enter your admin credentials to access the dashboard</p>
			</div>

			<form onsubmit={(e) => { e.preventDefault(); handleLogin(); }} class="space-y-4">
				{#if authError}
					<div class="p-3 rounded-lg bg-error/10 border border-error/20 text-error text-sm">
						{authError}
					</div>
				{/if}

				<div>
					<label for="username" class="block text-sm font-medium text-base-300 mb-1.5">Username</label>
					<input
						type="text"
						id="username"
						bind:value={username}
						class="input w-full"
						placeholder="admin"
						autocomplete="username"
					/>
				</div>

				<div>
					<label for="password" class="block text-sm font-medium text-base-300 mb-1.5">Password</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						class="input w-full"
						placeholder="••••••••"
						autocomplete="current-password"
					/>
				</div>

				<Button type="primary" class="w-full" disabled={loading}>
					{#if loading}
						Logging in...
					{:else}
						Login
					{/if}
				</Button>
			</form>
		</Card>
	</div>
{:else}
	<div class="p-6">
		<div class="max-w-7xl mx-auto">
			<div class="flex items-center justify-between mb-8">
				<div>
					<h1 class="font-display text-2xl font-bold text-white mb-1">Backoffice Dashboard</h1>
					<p class="text-sm text-base-400">Your SaaS metrics at a glance</p>
				</div>
				<div class="flex items-center gap-3">
					<Button type="secondary" onclick={refreshData} disabled={loading}>
						{loading ? 'Refreshing...' : 'Refresh'}
					</Button>
					<Button type="ghost" onclick={clearAdminToken}>
						<LogOut class="w-4 h-4 mr-2" />
						Logout
					</Button>
				</div>
			</div>

			{#if stats}
				<div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
					<Card>
						<div class="flex items-center justify-between mb-4">
							<div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center">
								<Users class="w-5 h-5 text-primary" />
							</div>
						</div>
						<p class="text-sm font-medium text-base-400 mb-1">Total Users</p>
						<p class="font-display text-2xl font-bold text-white">{stats.totalUsers}</p>
						<div class="mt-3 pt-3 border-t border-base-700 text-xs text-base-500">
							+{stats.newUsersToday} today
						</div>
					</Card>

					<Card>
						<div class="flex items-center justify-between mb-4">
							<div class="w-10 h-10 rounded-xl bg-secondary/10 flex items-center justify-center">
								<CreditCard class="w-5 h-5 text-secondary" />
							</div>
						</div>
						<p class="text-sm font-medium text-base-400 mb-1">Active Subscriptions</p>
						<p class="font-display text-2xl font-bold text-white">{stats.activeSubscriptions}</p>
						<div class="mt-3 pt-3 border-t border-base-700 text-xs text-base-500">
							{stats.trialSubscriptions} on trial
						</div>
					</Card>

					<Card>
						<div class="flex items-center justify-between mb-4">
							<div class="w-10 h-10 rounded-xl bg-tertiary/10 flex items-center justify-center">
								<TrendingUp class="w-5 h-5 text-tertiary" />
							</div>
						</div>
						<p class="text-sm font-medium text-base-400 mb-1">This Week</p>
						<p class="font-display text-2xl font-bold text-white">+{stats.newUsersThisWeek}</p>
						<div class="mt-3 pt-3 border-t border-base-700 text-xs text-base-500">
							New users this week
						</div>
					</Card>

					<Card>
						<div class="flex items-center justify-between mb-4">
							<div class="w-10 h-10 rounded-xl bg-success/10 flex items-center justify-center">
								<Calendar class="w-5 h-5 text-success" />
							</div>
						</div>
						<p class="text-sm font-medium text-base-400 mb-1">This Month</p>
						<p class="font-display text-2xl font-bold text-white">+{stats.newUsersThisMonth}</p>
						<div class="mt-3 pt-3 border-t border-base-700 text-xs text-base-500">
							New users this month
						</div>
					</Card>
				</div>
			{/if}

			{#if recentUsers.length > 0}
				<Card>
					<h2 class="font-display text-lg font-bold text-white mb-4">Recent Users</h2>
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead>
								<tr class="border-b border-base-700">
									<th class="text-left py-3 px-2 text-sm font-medium text-base-400">Email</th>
									<th class="text-left py-3 px-2 text-sm font-medium text-base-400">Name</th>
									<th class="text-left py-3 px-2 text-sm font-medium text-base-400">Plan</th>
									<th class="text-left py-3 px-2 text-sm font-medium text-base-400">Status</th>
									<th class="text-left py-3 px-2 text-sm font-medium text-base-400">Joined</th>
								</tr>
							</thead>
							<tbody>
								{#each recentUsers as user}
									<tr class="border-b border-base-700/50 hover:bg-base-800/50">
										<td class="py-3 px-2 text-sm text-white">{user.email}</td>
										<td class="py-3 px-2 text-sm text-base-300">{user.name || '-'}</td>
										<td class="py-3 px-2 text-sm">
											<span class="px-2 py-0.5 rounded-full text-xs font-medium capitalize
												{user.plan === 'free' ? 'bg-base-700 text-base-300' :
												 user.plan === 'pro' ? 'bg-primary/20 text-primary' :
												 'bg-secondary/20 text-secondary'}">
												{user.plan}
											</span>
										</td>
										<td class="py-3 px-2 text-sm">
											<span class="px-2 py-0.5 rounded-full text-xs font-medium capitalize
												{user.status === 'active' ? 'bg-success/20 text-success' :
												 user.status === 'trialing' ? 'bg-warning/20 text-warning' :
												 'bg-base-700 text-base-400'}">
												{user.status}
											</span>
										</td>
										<td class="py-3 px-2 text-sm text-base-400">{formatDate(user.createdAt)}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</Card>
			{/if}
		</div>
	</div>
{/if}
