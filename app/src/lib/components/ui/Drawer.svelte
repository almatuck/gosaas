<!--
  Drawer/Slide-over Component
  Sliding panel from edge of screen
-->

<script lang="ts">
	interface Props {
		open?: boolean;
		title?: string;
		position?: 'left' | 'right';
		size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		onclose?: () => void;
		children: any;
	}

	let {
		open = $bindable(false),
		title,
		position = 'right',
		size = 'md',
		onclose,
		children
	}: Props = $props();

	function handleClose() {
		open = false;
		if (onclose) {
			onclose();
		}
	}

	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl',
		full: 'max-w-full'
	};
</script>

{#if open}
	<div class="drawer-container" role="dialog" aria-modal="true">
		<!-- Backdrop -->
		<button
			onclick={handleClose}
			class="drawer-backdrop"
			aria-label="Close drawer"
		></button>

		<!-- Panel -->
		<div class="drawer-panel-wrapper drawer-{position}">
			<div class="drawer-panel drawer-{size}">
				<div class="drawer-content">
					<!-- Header -->
					{#if title}
						<div class="drawer-header">
							<h2 class="drawer-title">{title}</h2>
							<button
								type="button"
								onclick={handleClose}
								class="drawer-close"
								aria-label="Close drawer"
							>
								<i class="fas fa-times text-xl"></i>
							</button>
						</div>
					{/if}

					<!-- Content -->
					<div class="drawer-body">
						{@render children()}
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.drawer {
		.drawer-container {
			@apply fixed inset-0 z-50 overflow-hidden;
		}

		.drawer-backdrop {
			@apply fixed inset-0 transition-opacity;
			background: color-mix(in srgb, black 50%, transparent);
			backdrop-filter: blur(4px);
		}

		.drawer-panel-wrapper {
			@apply fixed inset-y-0 flex;
		}

		.drawer-right {
			@apply right-0;
		}

		.drawer-left {
			@apply left-0;
		}

		.drawer-panel {
			@apply pointer-events-auto w-screen transform transition-transform duration-300 ease-in-out;
		}

		.drawer-sm {
			@apply max-w-sm;
		}

		.drawer-md {
			@apply max-w-md;
		}

		.drawer-lg {
			@apply max-w-lg;
		}

		.drawer-xl {
			@apply max-w-xl;
		}

		.drawer-full {
			@apply max-w-full;
		}

		.drawer-content {
			@apply flex h-full flex-col overflow-y-auto shadow-xl;
			background: var(--color-base-900);
		}

		.drawer-header {
			@apply flex items-center justify-between px-6 py-4;
			border-bottom: 1px solid color-mix(in srgb, var(--color-base-700) 50%, transparent);
		}

		.drawer-title {
			@apply text-xl font-semibold;
			color: var(--color-text);
		}

		.drawer-close {
			@apply rounded-lg p-2 transition-colors;
			color: var(--color-base-400);

			&:hover {
				background: var(--color-base-800);
				color: var(--color-text);
			}
		}

		.drawer-body {
			@apply flex-1 overflow-y-auto p-6;
		}
	}
</style>
