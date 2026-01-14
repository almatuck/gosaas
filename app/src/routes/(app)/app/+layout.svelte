<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { AppNav } from '$lib/components/navigation';
	import AuthGuard from '$lib/components/AuthGuard.svelte';
	import { getWebSocketClient } from '$lib/websocket/client';

	let { children } = $props();

	// Initialize WebSocket connection when app layout mounts
	onMount(() => {
		console.log('[AppLayout] Initializing WebSocket connection...');
		const ws = getWebSocketClient();
		ws.connect();

		// Subscribe to status changes
		const unsubscribe = ws.onStatus((status) => {
			console.log('[AppLayout] WebSocket status:', status);
		});

		return () => {
			unsubscribe();
		};
	});

	onDestroy(() => {
		// Optionally disconnect when layout unmounts (e.g., logout)
		// getWebSocketClient().disconnect();
	});
</script>

<AuthGuard>
	<div class="layout-app">
		<AppNav />
		<main id="main-content" class="flex-1 p-6">
			<div class="max-w-[1400px] mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
</AuthGuard>
