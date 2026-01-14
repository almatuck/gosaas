<!--
  Divider Component
  Section dividers with optional labels
-->

<script lang="ts">
	interface Props {
		label?: string;
		orientation?: 'horizontal' | 'vertical';
		variant?: 'solid' | 'dashed' | 'dotted';
	}

	let {
		label,
		orientation = 'horizontal',
		variant = 'solid'
	}: Props = $props();

	const variantClasses = {
		solid: 'border-solid',
		dashed: 'border-dashed',
		dotted: 'border-dotted'
	};
</script>

{#if orientation === 'horizontal'}
	{#if label}
		<div class="relative my-6">
			<div class="absolute inset-0 flex items-center" aria-hidden="true">
				<div class="w-full border-t border-slate-700 {variantClasses[variant]}"></div>
			</div>
			<div class="relative flex justify-center">
				<span class="bg-slate-900 px-3 text-sm font-medium text-slate-400">{label}</span>
			</div>
		</div>
	{:else}
		<hr class="my-6 border-t border-slate-700 {variantClasses[variant]}" />
	{/if}
{:else}
	<div class="inline-block h-full w-px border-l border-slate-700 {variantClasses[variant]}"></div>
{/if}

<style>
@reference "$src/app.css";

@layer components.divider {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
