<!--
  CollapsibleCard Component
  Card with collapsible content and smooth height transitions
-->

<script lang="ts">
	interface Props {
		children: any;
		header: any;
		isExpanded?: boolean;
		onToggle?: (expanded: boolean) => void;
		class?: string;
		onCopy?: () => Promise<void>;
		copyTooltip?: string;
		autoExpandOnSearch?: boolean;
		statusRibbon?: any;
	}

	let {
		children,
		header,
		isExpanded = $bindable(true),
		onToggle,
		class: className = '',
		onCopy,
		copyTooltip = "Copy",
		autoExpandOnSearch = false,
		statusRibbon
	}: Props = $props();

	let copied = $state(false);
	let containerRef = $state<HTMLDivElement>();
	let contentHeight = $state(0);
	let shouldRender = $state(true);

	// Handle height transitions
	$effect(() => {
		if (isExpanded) {
			shouldRender = true;
			// Measure height after render
			setTimeout(() => {
				if (containerRef) {
					const contentElement = containerRef.firstElementChild as HTMLElement;
					if (contentElement) {
						contentHeight = contentElement.scrollHeight;
					}
				}
			}, 10);
		} else {
			contentHeight = 0;
			// Delay unmounting to allow close animation
			const timer = setTimeout(() => shouldRender = false, 300);
			return () => clearTimeout(timer);
		}
	});

	// Re-measure on content changes with ResizeObserver
	$effect(() => {
		if (isExpanded && shouldRender && containerRef) {
			const resizeObserver = new ResizeObserver((entries) => {
				for (const entry of entries) {
					const newHeight = entry.contentRect.height;
					if (newHeight > 0) {
						contentHeight = newHeight;
					}
				}
			});

			const contentElement = containerRef.firstElementChild as HTMLElement;
			if (contentElement) {
				resizeObserver.observe(contentElement);
			}

			return () => resizeObserver.disconnect();
		}
	});

	function handleToggle() {
		const newExpanded = !isExpanded;
		isExpanded = newExpanded;
		onToggle?.(newExpanded);
	}

	async function handleCopy() {
		if (!onCopy) return;

		try {
			await onCopy();
			copied = true;
			setTimeout(() => copied = false, 2000);
		} catch (error) {
			console.error('Copy failed:', error);
		}
	}
</script>

<div class="rounded-xl shadow-lg overflow-hidden bg-white dark:bg-gray-950 ring-1 ring-gray-200 dark:ring-gray-800 {className}">
	<!-- Status Ribbon -->
	{#if statusRibbon}
		{@render statusRibbon()}
	{/if}

	<!-- Header - Always Visible -->
	<div class="w-full py-1.5 px-2 bg-gray-50/80 dark:bg-gray-900/60 border-b border-gray-200/50 dark:border-gray-800/50 flex items-center justify-between">
		<!-- Header Content -->
		<div class="flex items-center gap-2 flex-1">
			{@render header()}
		</div>

		<!-- Controls -->
		<div class="flex items-center gap-2">
			<!-- Copy Button -->
			{#if onCopy}
				<button
					type="button"
					onclick={handleCopy}
					title={copyTooltip}
					class="gap-1 px-1 text-sm h-8 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
				>
					{#if copied}
						<svg class="h-4 w-4 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
						</svg>
					{:else}
						<svg class="h-4 w-4 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
						</svg>
					{/if}
					Copy
				</button>
			{/if}

			<!-- Collapse/Expand Button -->
			<button
				type="button"
				onclick={handleToggle}
				title={isExpanded ? "Collapse" : "Expand"}
				aria-label={isExpanded ? "Collapse" : "Expand"}
				class="w-6 h-8 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors flex items-center justify-center"
			>
				<svg
					class="h-4 w-4 transition-transform duration-200 {isExpanded ? 'rotate-180' : 'rotate-0'}"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
				</svg>
			</button>
		</div>
	</div>

	<!-- Collapsible Content -->
	<div
		bind:this={containerRef}
		class="collapsible-container dynamic-height"
		style:--height="{contentHeight}px"
	>
		{#if shouldRender}
			<div class="transition-opacity duration-150 {isExpanded ? 'opacity-100' : 'opacity-0'}">
				{@render children()}
			</div>
		{/if}
	</div>
</div>

<style>
@reference "$src/app.css";

@layer components.collapsible-card {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
