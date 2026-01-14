<!--
  Toast Notification Component
  Temporary notification messages
-->

<script lang="ts">
	interface Props {
		message: string;
		type?: 'success' | 'error' | 'warning' | 'info';
		duration?: number;
		show?: boolean;
		onclose?: () => void;
	}

	let {
		message,
		type = 'info',
		duration = 5000,
		show = $bindable(false),
		onclose
	}: Props = $props();

	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	$effect(() => {
		if (show && duration > 0) {
			timeoutId = setTimeout(() => {
				handleClose();
			}, duration);
		}

		return () => {
			if (timeoutId) clearTimeout(timeoutId);
		};
	});

	function handleClose() {
		show = false;
		if (onclose) {
			onclose();
		}
	}

	const typeConfig = {
		success: {
			icon: 'check-circle'
		},
		error: {
			icon: 'exclamation-circle'
		},
		warning: {
			icon: 'exclamation-triangle'
		},
		info: {
			icon: 'info-circle'
		}
	};

	let config = $derived(typeConfig[type]);
</script>

{#if show}
	<div
		class="toast toast-{type}"
		role="alert"
	>
		<i class="toast-icon fas fa-{config.icon}"></i>
		<p class="toast-message">{message}</p>
		<button
			type="button"
			onclick={handleClose}
			class="toast-close"
			aria-label="Close notification"
		>
			<i class="fas fa-times"></i>
		</button>
	</div>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.toast {
		.toast {
			@apply fixed bottom-4 right-4 z-50 flex items-center gap-3 rounded-xl px-5 py-4 shadow-xl ring-1;
			backdrop-filter: blur(4px);
			animation: slide-up 300ms ease-out;
		}

		.toast-success {
			background: color-mix(in srgb, var(--color-success) 20%, transparent);
			border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
		}

		.toast-error {
			background: color-mix(in srgb, var(--color-error) 20%, transparent);
			border-color: color-mix(in srgb, var(--color-error) 30%, transparent);
		}

		.toast-warning {
			background: color-mix(in srgb, var(--color-warning) 20%, transparent);
			border-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
		}

		.toast-info {
			background: color-mix(in srgb, var(--color-info) 20%, transparent);
			border-color: color-mix(in srgb, var(--color-info) 30%, transparent);
		}

		.toast-icon {
			@apply text-xl;
		}

		.toast-success .toast-icon {
			color: var(--color-success);
		}

		.toast-error .toast-icon {
			color: var(--color-error);
		}

		.toast-warning .toast-icon {
			color: var(--color-warning);
		}

		.toast-info .toast-icon {
			color: var(--color-info);
		}

		.toast-message {
			@apply text-base font-medium;
			color: var(--color-text);
		}

		.toast-close {
			@apply ml-4 transition-colors;
			color: var(--color-base-400);

			&:hover {
				color: var(--color-text);
			}
		}
	}
</style>
