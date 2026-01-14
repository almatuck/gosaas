<!--
  Collapsible Component
  Simple collapsible container (headless - you provide trigger and content)
-->

<script lang="ts">
	interface Props {
		open?: boolean;
		onOpenChange?: (open: boolean) => void;
		trigger?: any;
		children: any;
	}

	let {
		open = $bindable(false),
		onOpenChange,
		trigger,
		children
	}: Props = $props();

	function toggleOpen() {
		open = !open;
		onOpenChange?.(open);
	}
</script>

<div class="collapsible">
	{#if trigger}
		<div class="collapsible-trigger" onclick={toggleOpen} role="button" tabindex="0" onkeydown={(e) => e.key === 'Enter' && toggleOpen()}>
			{@render trigger()}
		</div>
	{/if}

	{#if open}
		<div class="collapsible-content">
			{@render children()}
		</div>
	{/if}
</div>

<style>
@reference "$src/app.css";

@layer components.collapsible {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
