<!--
  Tooltip Component
  Hover tooltip with positioning
-->

<script lang="ts">
	interface Props {
		content: string;
		position?: 'top' | 'bottom' | 'left' | 'right';
		children: any;
	}

	let {
		content,
		position = 'top',
		children
	}: Props = $props();

	let showTooltip = $state(false);

	const positionClasses = {
		top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
		bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
		left: 'right-full top-1/2 -translate-y-1/2 mr-2',
		right: 'left-full top-1/2 -translate-y-1/2 ml-2'
	};

	const arrowClasses = {
		top: 'top-full left-1/2 -translate-x-1/2 border-l-transparent border-r-transparent border-b-transparent',
		bottom: 'bottom-full left-1/2 -translate-x-1/2 border-l-transparent border-r-transparent border-t-transparent',
		left: 'left-full top-1/2 -translate-y-1/2 border-t-transparent border-b-transparent border-r-transparent',
		right: 'right-full top-1/2 -translate-y-1/2 border-t-transparent border-b-transparent border-l-transparent'
	};
</script>

<div
	class="tooltip-container"
	role="group"
	onmouseenter={() => showTooltip = true}
	onmouseleave={() => showTooltip = false}
>
	{@render children()}

	{#if showTooltip}
		<div
			class="tooltip-wrapper tooltip-{position}"
			role="tooltip"
		>
			<div class="tooltip-content">
				{content}
			</div>
			<div class="tooltip-arrow tooltip-arrow-{position}"></div>
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.tooltip {
		.tooltip-container {
			@apply relative inline-block;
		}

		.tooltip-wrapper {
			@apply absolute z-50;
			animation: fade-in 200ms ease-out;
		}

		.tooltip-top {
			@apply bottom-full left-1/2 -translate-x-1/2 mb-2;
		}

		.tooltip-bottom {
			@apply top-full left-1/2 -translate-x-1/2 mt-2;
		}

		.tooltip-left {
			@apply right-full top-1/2 -translate-y-1/2 mr-2;
		}

		.tooltip-right {
			@apply left-full top-1/2 -translate-y-1/2 ml-2;
		}

		.tooltip-content {
			@apply rounded-lg px-3 py-2 text-sm font-medium shadow-xl whitespace-nowrap;
			background: var(--color-base-900);
			color: var(--color-text);
			border: 1px solid color-mix(in srgb, var(--color-text) 10%, transparent);
		}

		.tooltip-arrow {
			@apply absolute border-4;
			border-color: var(--color-base-900);
		}

		.tooltip-arrow-top {
			@apply top-full left-1/2 -translate-x-1/2;
			border-left-color: transparent;
			border-right-color: transparent;
			border-bottom-color: transparent;
		}

		.tooltip-arrow-bottom {
			@apply bottom-full left-1/2 -translate-x-1/2;
			border-left-color: transparent;
			border-right-color: transparent;
			border-top-color: transparent;
		}

		.tooltip-arrow-left {
			@apply left-full top-1/2 -translate-y-1/2;
			border-top-color: transparent;
			border-bottom-color: transparent;
			border-right-color: transparent;
		}

		.tooltip-arrow-right {
			@apply right-full top-1/2 -translate-y-1/2;
			border-top-color: transparent;
			border-bottom-color: transparent;
			border-left-color: transparent;
		}
	}
</style>
