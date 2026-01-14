<!--
  Toggle/Switch Component
  For boolean on/off settings
-->

<script lang="ts">
	interface Props {
		checked?: boolean;
		disabled?: boolean;
		label?: string;
		description?: string;
		onchange?: (checked: boolean) => void;
		size?: 'sm' | 'md' | 'lg';
	}

	let {
		checked = $bindable(false),
		disabled = false,
		label,
		description,
		onchange,
		size = 'md'
	}: Props = $props();

	function handleToggle() {
		if (disabled) return;
		checked = !checked;
		if (onchange) {
			onchange(checked);
		}
	}

	const sizeClasses = {
		sm: 'toggle-sm',
		md: 'toggle-md',
		lg: 'toggle-lg'
	};

	const dotSizeClasses = {
		sm: 'toggle-dot-sm',
		md: 'toggle-dot-md',
		lg: 'toggle-dot-lg'
	};
</script>

<div class="flex items-center gap-3">
	<button
		type="button"
		role="switch"
		aria-checked={checked}
		aria-label={label || 'Toggle switch'}
		{disabled}
		onclick={handleToggle}
		class="toggle {sizeClasses[size]} {checked ? 'checked' : ''} {disabled ? 'disabled' : ''}"
	>
		<span class="toggle-dot {dotSizeClasses[size]} {checked ? 'checked' : ''}"></span>
	</button>
	{#if label}
		<div class="flex-1">
			<button
				type="button"
				{disabled}
				onclick={handleToggle}
				class="text-left {disabled ? 'cursor-not-allowed' : 'cursor-pointer'}"
			>
				<p class="toggle-label">{label}</p>
				{#if description}
					<p class="toggle-description">{description}</p>
				{/if}
			</button>
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.toggle {
		.toggle {
			@apply relative inline-flex flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent;
			@apply transition-colors duration-200 ease-in-out;
			@apply focus:outline-none focus:ring-2 focus:ring-offset-2;
		}

		.toggle:focus {
			@apply ring-primary ring-offset-base-900;
		}

		.toggle-sm { @apply h-5 w-9; }
		.toggle-md { @apply h-6 w-11; }
		.toggle-lg { @apply h-7 w-14; }

		.toggle.checked {
			@apply bg-gradient-to-r from-primary to-[#a855f7];
		}

		.toggle:not(.checked) {
			@apply bg-base-700;
		}

		.toggle.disabled {
			@apply opacity-50 cursor-not-allowed;
		}

		.toggle-dot {
			@apply pointer-events-none inline-block transform rounded-full bg-white shadow-lg ring-0;
			@apply transition duration-200 ease-in-out;
		}

		.toggle-dot-sm { @apply h-4 w-4; }
		.toggle-dot-md { @apply h-5 w-5; }
		.toggle-dot-lg { @apply h-6 w-6; }

		.toggle-dot.checked.toggle-dot-sm { @apply translate-x-4; }
		.toggle-dot.checked.toggle-dot-md { @apply translate-x-5; }
		.toggle-dot.checked.toggle-dot-lg { @apply translate-x-7; }
		.toggle-dot:not(.checked) { @apply translate-x-0; }

		.toggle-label {
			@apply text-base font-medium;
			color: var(--color-text);
		}

		.toggle-description {
			@apply text-sm;
			color: var(--color-text-muted);
		}
	}
</style>
