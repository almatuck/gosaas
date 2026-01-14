<!--
  Reusable Empty State Component
-->

<script lang="ts">
	let {
		icon = 'inbox',
		title,
		message
	}: {
		icon?: string;
		title: string;
		message?: string;
	} = $props();
</script>

<div class="p-12 text-center">
	<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-slate-700/30">
		<i class="fas fa-{icon} text-2xl text-slate-500"></i>
	</div>
	<p class="text-sm font-medium text-white">{title}</p>
	{#if message}
		<p class="mt-1 text-xs text-slate-400">{message}</p>
	{/if}
</div>

<style>
@reference "$src/app.css";

@layer components.empty-state {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
