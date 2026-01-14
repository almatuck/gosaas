<!--
  Sheet Component
  Slide-in panel from the side (similar to Drawer but with different API)
-->

<script lang="ts">
	interface Props {
		open?: boolean;
		side?: 'top' | 'bottom' | 'left' | 'right';
		title?: string;
		description?: string;
		children?: any;
		onOpenChange?: (open: boolean) => void;
	}

	let {
		open = $bindable(false),
		side = 'right',
		title = '',
		description = '',
		children,
		onOpenChange
	}: Props = $props();

	const sideClasses = {
		top: 'inset-x-0 top-0 border-b',
		bottom: 'inset-x-0 bottom-0 border-t',
		left: 'inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm',
		right: 'inset-y-0 right-0 h-full w-3/4 border-l sm:max-w-sm'
	};

	const slideClasses = {
		top: 'data-[state=closed]:slide-out-to-top data-[state=open]:slide-in-from-top',
		bottom: 'data-[state=closed]:slide-out-to-bottom data-[state=open]:slide-in-from-bottom',
		left: 'data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left',
		right: 'data-[state=closed]:slide-out-to-right data-[state=open]:slide-in-from-right'
	};

	function close() {
		open = false;
		onOpenChange?.(false);
	}

	function handleBackdropClick() {
		close();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			close();
		}
	}
</script>

{#if open}
	<div
		class="sheet-backdrop"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="button"
		tabindex="0"
		aria-label="Close sheet"
	>
		<div
			class="sheet-panel sheet-{side}"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Enter' && e.stopPropagation()}
			role="dialog"
			aria-modal="true"
			tabindex="0"
		>
			<!-- Close button -->
			<button
				type="button"
				onclick={close}
				class="sheet-close"
			>
				<svg class="sheet-close-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
				</svg>
				<span class="sr-only">Close</span>
			</button>

			<!-- Header -->
			{#if title || description}
				<div class="sheet-header">
					{#if title}
						<h2 class="sheet-title">{title}</h2>
					{/if}
					{#if description}
						<p class="sheet-description">{description}</p>
					{/if}
				</div>
			{/if}

			<!-- Content -->
			{#if children}
				{@render children()}
			{/if}
		</div>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.sheet {
		.sheet-backdrop {
			@apply fixed inset-0 z-50;
			background: color-mix(in srgb, black 80%, transparent);
			animation: fade-in 200ms ease-out;
		}

		.sheet-panel {
			@apply fixed z-50 gap-4 p-6 shadow-lg transition ease-in-out;
			background: var(--color-bg);
			animation-duration: 500ms;
		}

		.sheet-top {
			@apply inset-x-0 top-0;
			border-bottom: 1px solid var(--color-border);
		}

		.sheet-bottom {
			@apply inset-x-0 bottom-0;
			border-top: 1px solid var(--color-border);
		}

		.sheet-left {
			@apply inset-y-0 left-0 h-full w-3/4 sm:max-w-sm;
			border-right: 1px solid var(--color-border);
		}

		.sheet-right {
			@apply inset-y-0 right-0 h-full w-3/4 sm:max-w-sm;
			border-left: 1px solid var(--color-border);
		}

		.sheet-close {
			@apply absolute right-4 top-4 rounded-sm opacity-70 transition-opacity;
			@apply hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-offset-2;
			@apply disabled:pointer-events-none;
			color: var(--color-text);
			ring-color: var(--color-primary);
		}

		.sheet-close-icon {
			@apply h-4 w-4;
		}

		.sheet-header {
			@apply flex flex-col space-y-2 text-center sm:text-left;
		}

		.sheet-title {
			@apply text-lg font-semibold;
			color: var(--color-text);
		}

		.sheet-description {
			@apply text-sm;
			color: var(--color-text-secondary);
		}
	}
</style>
