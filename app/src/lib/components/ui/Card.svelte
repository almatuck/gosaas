<!--
  Reusable Card Component
  Ensures 100% consistency for all card-based layouts
-->

<script lang="ts">
	let {
		title,
		subtitle,
		children,
		class: extraClass = '',
		hover = false,
		onclick
	}: {
		title?: string;
		subtitle?: string;
		children: any;
		class?: string;
		hover?: boolean;
		onclick?: () => void;
	} = $props();

	const isClickable = $derived(!!onclick);
	const className = $derived(`card ${hover ? 'card-hover' : ''} ${isClickable ? 'cursor-pointer' : ''} ${extraClass}`.trim());
</script>

{#snippet cardContent()}
	{#if title}
		<div class="card-header">
			<h2 class="card-title">{title}</h2>
			{#if subtitle}
				<p class="card-subtitle">{subtitle}</p>
			{/if}
		</div>
	{/if}
	<div class="card-body">
		{@render children()}
	</div>
{/snippet}

{#if isClickable}
	<button type="button" class={className} onclick={onclick}>
		{@render cardContent()}
	</button>
{:else}
	<div class={className}>
		{@render cardContent()}
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.card {
		/* Card styles are defined in app.css */
	}
</style>
