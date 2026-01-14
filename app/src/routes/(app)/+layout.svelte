<script lang="ts">
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { AppNav } from '$lib/components/navigation';
	import { auth, isAuthenticated, subscription } from '$lib/stores';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();
	let mounted = $state(false);

	const authenticated = $derived($isAuthenticated);

	onMount(async () => {
		mounted = true;

		// Redirect if not authenticated
		if (!authenticated) {
			goto('/auth/login');
			return;
		}

		// Fetch user profile and subscription data
		await Promise.all([
			auth.fetchCurrentUser(),
			subscription.loadSubscription(),
			subscription.loadUsage()
		]);
	});
</script>

{#if mounted && authenticated}
	<div class="layout-app">
		<AppNav />
		<main id="main-content" class="flex-1 p-6">
			<div class="max-w-7xl mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
{:else}
	<div class="layout-app flex items-center justify-center">
		<div class="spinner"></div>
	</div>
{/if}
