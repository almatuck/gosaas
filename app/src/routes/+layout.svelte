<script lang="ts">
	import type { Snippet } from 'svelte';
	import favicon from '$lib/assets/favicon.svg';
	import '$src/app.css';

	let { children }: { children: Snippet } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

{@render children()}
