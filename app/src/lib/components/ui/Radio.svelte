<!--
  Radio Component
  Styled radio button with label
-->

<script lang="ts">
	interface Props {
		value: string;
		checked?: boolean;
		name: string;
		disabled?: boolean;
		label?: string;
		description?: string;
		onchange?: (value: string) => void;
		id?: string;
	}

	let {
		value,
		checked = false,
		name,
		disabled = false,
		label,
		description,
		onchange,
		id
	}: Props = $props();

	const radioId = $derived(id || `radio-${Math.random().toString(36).slice(2, 9)}`);

	function handleChange() {
		if (onchange) {
			onchange(value);
		}
	}
</script>

<div class="flex items-start gap-3">
	<div class="flex h-6 items-center">
		<input
			type="radio"
			id={radioId}
			{name}
			{value}
			{checked}
			{disabled}
			onchange={handleChange}
			class="radio-input"
		/>
	</div>
	{#if label}
		<div class="flex-1">
			<label for={radioId} class="radio-label {disabled ? 'opacity-50' : 'cursor-pointer'}">
				{label}
			</label>
			{#if description}
				<p class="radio-description {disabled ? 'opacity-50' : ''}">{description}</p>
			{/if}
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.radio {
		.radio-input {
			@apply h-5 w-5 border-2 transition-all;
			@apply focus:ring-2 focus:ring-offset-2;
			@apply disabled:cursor-not-allowed disabled:opacity-50;
			border-color: var(--color-base-600);
			background: color-mix(in srgb, var(--color-base-800) 60%, transparent);
			color: var(--color-primary);
		}

		.radio-input:focus {
			@apply ring-primary ring-offset-base-900;
		}

		.radio-input:checked {
			@apply bg-gradient-to-r from-primary to-[#a855f7] border-transparent;
		}

		.radio-label {
			@apply text-base font-medium;
			color: var(--color-text);
		}

		.radio-description {
			@apply text-sm;
			color: var(--color-text-muted);
		}
	}
</style>
