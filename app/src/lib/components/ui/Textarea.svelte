<!--
  Textarea Component
  Multi-line text input
-->

<script lang="ts">
	interface Props {
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		rows?: number;
		class?: string;
	}

	let {
		value = $bindable(''),
		placeholder = '',
		disabled = false,
		rows = 3,
		class: className = ''
	}: Props = $props();
</script>

<textarea
	bind:value
	{placeholder}
	{disabled}
	{rows}
	class="textarea {className}"
></textarea>

<style>
	@reference "$src/app.css";
	@layer components.textarea {
		.textarea {
			@apply flex min-h-[80px] w-full rounded-md border px-3 py-2 text-sm;
			@apply focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2;
			@apply disabled:cursor-not-allowed disabled:opacity-50;
			background: var(--color-bg);
			color: var(--color-text);
			border-color: var(--color-border);
		}

		.textarea::placeholder {
			color: var(--color-text-muted);
		}

		.textarea:focus-visible {
			@apply ring-primary;
			ring-offset-color: var(--color-bg);
		}
	}
</style>
