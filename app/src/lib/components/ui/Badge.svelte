<!--
  Badge Component
  For status indicators, counts, labels
-->

<script lang="ts">
	interface Props {
		variant?: 'default' | 'primary' | 'success' | 'warning' | 'error' | 'info';
		size?: 'sm' | 'md' | 'lg';
		rounded?: boolean;
		children: any;
	}

	let {
		variant = 'default',
		size = 'md',
		rounded = false,
		children
	}: Props = $props();

</script>

<span
	class="badge badge-{variant} badge-{size} {rounded ? 'badge-rounded' : 'badge-square'}"
>
	{@render children()}
</span>

<style>
	@reference "$src/app.css";
	@layer components.badge {
		.badge {
			@apply inline-flex items-center gap-1 font-semibold;
		}

		.badge-rounded { @apply rounded-full; }
		.badge-square { @apply rounded-lg; }

		.badge-sm { @apply px-2 py-0.5 text-xs; }
		.badge-md { @apply px-2.5 py-1 text-sm; }
		.badge-lg { @apply px-3 py-1.5 text-base; }

		.badge-default {
			background: color-mix(in srgb, var(--color-base-700) 50%, transparent);
			color: var(--color-base-300);
		}

		.badge-primary {
			@apply text-white;
			background: linear-gradient(to right, var(--color-primary), var(--color-secondary));
		}

		.badge-success {
			background: color-mix(in srgb, var(--color-success) 20%, transparent);
			color: var(--color-success);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-success) 30%, transparent);
		}

		.badge-warning {
			background: color-mix(in srgb, var(--color-warning) 20%, transparent);
			color: var(--color-warning);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
		}

		.badge-error {
			background: color-mix(in srgb, var(--color-error) 20%, transparent);
			color: var(--color-error);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-error) 30%, transparent);
		}

		.badge-info {
			background: color-mix(in srgb, var(--color-info) 20%, transparent);
			color: var(--color-info);
			@apply ring-1;
			ring-color: color-mix(in srgb, var(--color-info) 30%, transparent);
		}
	}
</style>
