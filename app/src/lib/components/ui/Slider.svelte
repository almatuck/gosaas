<!--
  Slider/Range Component
  For selecting numeric values with a slider
-->

<script lang="ts">
	interface Props {
		value?: number;
		min?: number;
		max?: number;
		step?: number;
		disabled?: boolean;
		label?: string;
		showValue?: boolean;
		formatValue?: (value: number) => string;
		onchange?: (value: number) => void;
		id?: string;
	}

	let {
		value = $bindable(50),
		min = 0,
		max = 100,
		step = 1,
		disabled = false,
		label,
		showValue = true,
		formatValue = (v) => v.toString(),
		onchange,
		id
	}: Props = $props();

	const sliderId = $derived(id || `slider-${Math.random().toString(36).slice(2, 9)}`);

	function handleChange(event: Event) {
		const target = event.target as HTMLInputElement;
		value = Number(target.value);
		if (onchange) {
			onchange(value);
		}
	}

	let percentage = $derived(((value - min) / (max - min)) * 100);
</script>

<div class="w-full">
	{#if label}
		<div class="mb-2 flex items-center justify-between">
			<label for={sliderId} class="slider-label">{label}</label>
			{#if showValue}
				<span class="slider-value">{formatValue(value)}</span>
			{/if}
		</div>
	{/if}
	<div class="relative">
		<input
			type="range"
			id={sliderId}
			bind:value
			{min}
			{max}
			{step}
			{disabled}
			onchange={handleChange}
			oninput={handleChange}
			class="slider-input"
			style:--slider-percentage="{percentage}%"
		/>
	</div>
</div>

<style>
	@reference "$src/app.css";
	@layer components.slider {
		.slider-label {
			@apply text-base font-medium;
			color: var(--color-text);
		}

		.slider-value {
			@apply text-sm font-semibold;
			color: var(--color-primary-light);
		}

		.slider-input {
			@apply w-full h-2 rounded-lg appearance-none cursor-pointer;
			@apply disabled:opacity-50 disabled:cursor-not-allowed;
			background: linear-gradient(
				to right,
				var(--color-primary) 0%,
				#a855f7 var(--slider-percentage),
				var(--color-base-700) var(--slider-percentage),
				var(--color-base-700) 100%
			);
		}

		.slider-input::-webkit-slider-thumb {
			@apply appearance-none w-5 h-5 rounded-full cursor-pointer transition-all;
			@apply bg-gradient-to-r from-primary to-[#a855f7];
			box-shadow: 0 0 12px color-mix(in srgb, var(--color-primary) 30%, transparent);
		}

		.slider-input:hover::-webkit-slider-thumb {
			box-shadow: 0 0 20px color-mix(in srgb, var(--color-primary) 50%, transparent);
		}

		.slider-input::-moz-range-thumb {
			@apply appearance-none w-5 h-5 rounded-full border-0 cursor-pointer transition-all;
			@apply bg-gradient-to-r from-primary to-[#a855f7];
			box-shadow: 0 0 12px color-mix(in srgb, var(--color-primary) 30%, transparent);
		}

		.slider-input:hover::-moz-range-thumb {
			box-shadow: 0 0 20px color-mix(in srgb, var(--color-primary) 50%, transparent);
		}
	}
</style>
