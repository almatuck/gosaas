<script lang="ts">
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { AppNav } from '$lib/components/navigation';
	import { auth, isAuthenticated, subscription } from '$lib/stores';
	import { getWebSocketClient } from '$lib/websocket/client';

	let { children }: { children: Snippet } = $props();
	let mounted = $state(false);

	const authenticated = $derived($isAuthenticated);

	onMount(async () => {
		// Redirect if not authenticated
		if (!authenticated) {
			goto('/auth/login');
			return;
		}

		// Initialize WebSocket
		const ws = getWebSocketClient();
		ws.connect();

		// Fetch user profile and subscription data
		await Promise.all([
			auth.fetchCurrentUser(),
			subscription.loadSubscription(),
			subscription.loadUsage()
		]);

		mounted = true;
	});
</script>

{#if mounted && authenticated}
	<div class="layout-app">
		<AppNav />
		<main id="main-content" class="flex-1 p-6">
			<div class="max-w-[1400px] mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
{:else}
	<div class="layout-app flex items-center justify-center">
		<div class="spinner"></div>
	</div>
{/if}
