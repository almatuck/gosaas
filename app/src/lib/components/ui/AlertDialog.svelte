<!--
  AlertDialog Component
  Alert/confirmation dialog with actions
-->

<script lang="ts">
	import { Button } from '$lib/components/ui';

	interface Props {
		open?: boolean;
		title?: string;
		description?: string;
		actionLabel?: string;
		cancelLabel?: string;
		actionType?: 'primary' | 'danger';
		children?: any;
		onAction?: () => void;
		onCancel?: () => void;
		onOpenChange?: (open: boolean) => void;
	}

	let {
		open = $bindable(false),
		title = '',
		description = '',
		actionLabel = 'Continue',
		cancelLabel = 'Cancel',
		actionType = 'primary',
		children,
		onAction,
		onCancel,
		onOpenChange
	}: Props = $props();

	function handleAction() {
		onAction?.();
		close();
	}

	function handleCancel() {
		onCancel?.();
		close();
	}

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
		class="alert-dialog-backdrop"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="button"
		tabindex="0"
		aria-label="Close dialog"
	>
		<div
			class="alert-dialog"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.key === 'Enter' && e.stopPropagation()}
			role="alertdialog"
			aria-modal="true"
			tabindex="0"
		>
			<!-- Header -->
			<div class="alert-dialog-header">
				<div class="flex items-start justify-between">
					{#if title}
						<h3 class="alert-dialog-title">{title}</h3>
					{/if}
					<button
						onclick={close}
						class="alert-dialog-close"
						aria-label="Close"
					>
						<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M6 18L18 6M6 6l12 12"
							/>
						</svg>
					</button>
				</div>
			</div>

			<!-- Body -->
			<div class="alert-dialog-body">
				{#if description}
					<p class="alert-dialog-description">{description}</p>
				{/if}
				{#if children}
					{@render children()}
				{/if}
			</div>

			<!-- Footer -->
			<div class="alert-dialog-footer">
				<div class="flex justify-end gap-3">
					<Button type="secondary" onclick={handleCancel}>
						{cancelLabel}
					</Button>
					<Button type={actionType} onclick={handleAction}>
						{actionLabel}
					</Button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.alert-dialog {
		.alert-dialog-backdrop {
			@apply fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm;
			background: color-mix(in srgb, black 50%, transparent);
		}

		.alert-dialog {
			@apply w-full max-w-md rounded-xl backdrop-blur-sm shadow-2xl flex flex-col;
			background: color-mix(in srgb, var(--color-base-800) 95%, transparent);
			border: 1px solid color-mix(in srgb, var(--color-base-700) 50%, transparent);
		}

		.alert-dialog-header {
			@apply flex-shrink-0 px-6 py-4;
			border-bottom: 1px solid color-mix(in srgb, var(--color-base-700) 50%, transparent);
		}

		.alert-dialog-title {
			@apply text-xl font-bold text-white;
		}

		.alert-dialog-close {
			@apply ml-auto rounded-lg p-1.5 transition-colors;
			color: var(--color-base-400);
		}

		.alert-dialog-close:hover {
			@apply bg-base-700 text-white;
		}

		.alert-dialog-body {
			@apply flex-1 p-6;
		}

		.alert-dialog-description {
			@apply text-sm;
			color: var(--color-base-400);
		}

		.alert-dialog-footer {
			@apply flex-shrink-0 px-6 py-4;
			background: color-mix(in srgb, var(--color-base-800) 40%, transparent);
			border-top: 1px solid color-mix(in srgb, var(--color-base-700) 50%, transparent);
		}
	}
</style>
