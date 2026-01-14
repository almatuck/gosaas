<!--
  Button Group Component
  Segmented control buttons with shared borders
-->

<script lang="ts">
	interface Button {
		value: string;
		label: string;
		icon?: string;
	}

	interface Props {
		buttons: Button[];
		selected?: string;
		onselect?: (value: string) => void;
		size?: 'sm' | 'md' | 'lg';
	}

	let {
		buttons,
		selected = $bindable(''),
		onselect,
		size = 'md'
	}: Props = $props();

	function handleSelect(value: string) {
		selected = value;
		if (onselect) {
			onselect(value);
		}
	}

</script>

<span class="button-group">
	{#each buttons as button, index}
		<button
			type="button"
			onclick={() => handleSelect(button.value)}
			class="button-group-item button-group-{size}
				{index === 0 ? 'button-group-first' : ''}
				{index === buttons.length - 1 ? 'button-group-last' : ''}
				{index !== 0 ? 'button-group-middle' : ''}
				{selected === button.value ? 'button-group-selected' : ''}"
		>
			{#if button.icon}
				<i class="fas fa-{button.icon}"></i>
			{/if}
			{button.label}
		</button>
	{/each}
</span>

<style>
	@reference "$src/app.css";
	@layer components.button-group {
		.button-group {
			@apply isolate inline-flex rounded-xl shadow-lg;
			box-shadow: 0 10px 25px color-mix(in srgb, var(--color-base-900) 20%, transparent);
		}

		.button-group-item {
			@apply relative inline-flex items-center gap-2 font-semibold transition-all;
		}

		.button-group-item:focus {
			@apply z-10;
		}

		.button-group-sm { @apply px-3 py-1.5 text-sm; }
		.button-group-md { @apply px-4 py-2.5 text-base; }
		.button-group-lg { @apply px-6 py-3 text-lg; }

		.button-group-first { @apply rounded-l-xl; }
		.button-group-last { @apply rounded-r-xl; }
		.button-group-middle { @apply -ml-px; }

		.button-group-selected {
			@apply text-white shadow-md z-10 border-2 border-transparent;
			background: linear-gradient(to right, var(--color-primary), var(--color-secondary));
		}

		.button-group-item:not(.button-group-selected) {
			@apply border-2;
			background: color-mix(in srgb, var(--color-base-800) 60%, transparent);
			color: var(--color-base-300);
			border-color: color-mix(in srgb, var(--color-base-600) 60%, transparent);
		}

		.button-group-item:not(.button-group-selected):hover {
			@apply text-white;
			background: color-mix(in srgb, var(--color-base-700) 60%, transparent);
		}
	}
</style>
