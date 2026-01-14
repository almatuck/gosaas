<!--
  Separator Component
  Horizontal or vertical divider line
-->

<script lang="ts">
	interface Props {
		orientation?: 'horizontal' | 'vertical';
		decorative?: boolean;
		class?: string;
	}

	let {
		orientation = 'horizontal',
		decorative = true,
		class: className = ''
	}: Props = $props();

	const orientationClasses = {
		horizontal: 'h-[1px] w-full',
		vertical: 'h-full w-[1px]'
	};
</script>

<div
	class="shrink-0 bg-slate-200 dark:bg-slate-700 {orientationClasses[orientation]} {className}"
	role={decorative ? 'presentation' : 'separator'}
	aria-orientation={orientation}
></div>

<style>
@reference "$src/app.css";

@layer components.separator {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
