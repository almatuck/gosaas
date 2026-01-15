<script lang="ts">
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { AppNav } from '$lib/components/navigation';
	import { auth, isAuthenticated, subscription } from '$lib/stores';
	import { getWebSocketClient } from '$lib/websocket/client';

	let { children }: { children: Snippet } = $props();

	const authenticated = $derived($isAuthenticated);

	onMount(async () => {
		if (!authenticated) {
			goto('/auth/login');
			return;
		}

		getWebSocketClient().connect();

		await Promise.all([
			auth.fetchCurrentUser(),
			subscription.loadSubscription(),
			subscription.loadUsage()
		]);
	});
</script>

<div class="layout-app">
	<AppNav />
	<main id="main-content" class="flex-1 p-6">
		<div class="max-w-[1400px] mx-auto">
			{@render children()}
		</div>
	</main>
</div>
