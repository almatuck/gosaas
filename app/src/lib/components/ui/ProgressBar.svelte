<!--
  Progress Bar Component
  Linear progress indicator with labels
-->

<script lang="ts">
	interface Props {
		value: number;
		max?: number;
		label?: string;
		showPercent?: boolean;
		size?: 'sm' | 'md' | 'lg';
		variant?: 'default' | 'success' | 'warning' | 'error';
		striped?: boolean;
		animated?: boolean;
	}

	let {
		value,
		max = 100,
		label,
		showPercent = true,
		size = 'md',
		variant = 'default',
		striped = false,
		animated = false
	}: Props = $props();

	let percentage = $derived(Math.min(100, Math.max(0, (value / max) * 100)));

	const sizeClasses = {
		sm: 'h-2',
		md: 'h-3',
		lg: 'h-4'
	};

	const variantClasses = {
		default: 'bg-gradient-to-r from-indigo-600 to-purple-600',
		success: 'bg-gradient-to-r from-green-600 to-emerald-600',
		warning: 'bg-gradient-to-r from-yellow-600 to-orange-600',
		error: 'bg-gradient-to-r from-red-600 to-rose-600'
	};
</script>

<div class="w-full">
	{#if label || showPercent}
		<div class="mb-2 flex items-center justify-between">
			{#if label}
				<span class="text-sm font-medium text-slate-400">{label}</span>
			{/if}
			{#if showPercent}
				<span class="text-sm font-semibold text-white">{percentage.toFixed(0)}%</span>
			{/if}
		</div>
	{/if}

	<div class="w-full overflow-hidden rounded-full bg-slate-700">
		<div
			class="{sizeClasses[size]} {variantClasses[variant]} transition-all duration-500 ease-out dynamic-width
				{striped ? 'bg-stripes' : ''}
				{animated ? 'animate-stripes' : ''}"
			style:--width="{percentage}%"
			role="progressbar"
			aria-valuenow={value}
			aria-valuemin="0"
			aria-valuemax={max}
		></div>
	</div>
</div>

<style>
	@reference "$src/app.css";

	@layer components.progress-bar {
		.bg-stripes {
			background-image: linear-gradient(
				45deg,
				color-mix(in srgb, var(--color-base-50) 15%, transparent) 25%,
				transparent 25%,
				transparent 50%,
				color-mix(in srgb, var(--color-base-50) 15%, transparent) 50%,
				color-mix(in srgb, var(--color-base-50) 15%, transparent) 75%,
				transparent 75%,
				transparent
			);
			background-size: 1rem 1rem;
		}

		@keyframes stripes {
			from {
				background-position: 1rem 0;
			}
			to {
				background-position: 0 0;
			}
		}

		.animate-stripes {
			animation: stripes 1s linear infinite;
		}
	}
</style>
