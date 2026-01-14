<!--
  Checkbox Component (Switch Style)
  Toggle switch with label
-->

<script lang="ts">
	interface Props {
		checked?: boolean;
		disabled?: boolean;
		label?: string;
		description?: string;
		indeterminate?: boolean;
		onchange?: (checked: boolean) => void;
		id?: string;
	}

	let {
		checked = $bindable(),
		disabled = false,
		label,
		description,
		indeterminate = false,
		onchange,
		id
	}: Props = $props();

	const switchId = $derived(id || `switch-${Math.random().toString(36).slice(2, 9)}`);
	let safeChecked = $derived(checked ?? false);

	function handleToggle() {
		if (disabled) return;
		checked = !checked;
		if (onchange) {
			onchange(checked ?? false);
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			handleToggle();
		}
	}
</script>

<div class="flex items-start gap-3">
	<button
		type="button"
		role="switch"
		aria-checked={safeChecked}
		aria-labelledby={label ? `${switchId}-label` : undefined}
		id={switchId}
		{disabled}
		onclick={handleToggle}
		onkeydown={handleKeydown}
		class="checkbox-switch {safeChecked ? 'checked' : ''}"
	>
		<!-- Track glow effect when checked -->
		{#if safeChecked}
			<span class="checkbox-glow"></span>
		{/if}

		<!-- Thumb -->
		<span class="checkbox-thumb {safeChecked ? 'checked' : ''}">
			<!-- Check icon when on -->
			<span class="checkbox-icon {safeChecked ? 'checked' : ''}">
				<svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
					<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
				</svg>
			</span>
		</span>
	</button>

	{#if label}
		<div class="flex-1 min-w-0">
			<label
				id="{switchId}-label"
				for={switchId}
				class="checkbox-label {disabled ? 'opacity-50' : 'cursor-pointer'}"
			>
				{label}
			</label>
			{#if description}
				<p class="checkbox-description {disabled ? 'opacity-50' : ''}">{description}</p>
			{/if}
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.checkbox {
		.checkbox-switch {
			@apply relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full transition-all duration-300 ease-out;
			@apply focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2;
			@apply disabled:cursor-not-allowed disabled:opacity-50;
		}

		.checkbox-switch:focus-visible {
			@apply ring-primary ring-offset-base-900;
		}

		.checkbox-switch.checked {
			@apply bg-gradient-to-r from-primary to-[#a855f7];
			box-shadow: 0 0 20px color-mix(in srgb, var(--color-primary) 25%, transparent);
		}

		.checkbox-switch:not(.checked) {
			@apply bg-base-700;
		}

		.checkbox-switch:not(.checked):hover:not(:disabled) {
			@apply bg-base-600;
		}

		.checkbox-glow {
			@apply absolute inset-0 rounded-full opacity-0 blur-sm transition-opacity duration-300;
			@apply bg-gradient-to-r from-primary-light to-[#c084fc];
		}

		.checkbox-switch.checked:hover .checkbox-glow {
			@apply opacity-50;
		}

		.checkbox-thumb {
			@apply pointer-events-none relative inline-block h-5 w-5 transform rounded-full shadow-lg ring-0 transition-all duration-300 ease-out;
		}

		.checkbox-thumb.checked {
			@apply translate-x-5 bg-white;
		}

		.checkbox-thumb:not(.checked) {
			@apply translate-x-0.5 bg-base-300;
		}

		.checkbox-icon {
			@apply absolute inset-0 flex items-center justify-center transition-opacity duration-200;
			color: var(--color-primary);
		}

		.checkbox-icon.checked {
			@apply opacity-100;
		}

		.checkbox-icon:not(.checked) {
			@apply opacity-0;
		}

		.checkbox-label {
			@apply text-sm font-medium;
			color: var(--color-text);
		}

		.checkbox-description {
			@apply text-xs mt-0.5;
			color: var(--color-text-muted);
		}
	}
</style>
