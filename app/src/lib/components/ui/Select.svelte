<!--
  Reusable Select Component
-->

<script lang="ts">
	let {
		value = $bindable(),
		size = 'md',
		disabled = false,
		children,
		class: extraClass = ''
	}: {
		value?: string;
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		children: any;
		class?: string;
	} = $props();

	let safeValue = $derived(value ?? '');

	const sizeClasses = {
		sm: 'select-sm',
		md: 'select-md',
		lg: 'select-lg'
	};

	const className = $derived(`select ${sizeClasses[size]} ${extraClass}`.trim());
</script>

<select class={className} value={safeValue} onchange={(e) => value = e.currentTarget.value} {disabled}>
	{@render children()}
</select>

<style>
	@reference "$src/app.css";
	@layer components.select {
		.select {
			@apply w-full px-4 py-3 rounded-xl transition-all duration-200;
			@apply disabled:opacity-50 disabled:cursor-not-allowed;
			background: var(--color-base-700);
			color: var(--color-text);
			border: 1px solid color-mix(in srgb, white 10%, transparent);
			font-size: 0.9375rem;
		}

		.select:focus {
			@apply outline-none;
			border-color: var(--color-primary);
			box-shadow: 0 0 0 3px color-mix(in srgb, var(--color-primary) 20%, transparent);
		}

		.select-sm {
			@apply h-8 px-3 py-1 text-sm;
		}

		.select-md {
			/* Default size */
		}

		.select-lg {
			@apply h-12 px-5 py-4 text-base;
		}
	}
</style>
