<!--
  TagInput Component
  Input field for entering multiple tags/items

  Features:
  - Add tags by pressing Enter or comma
  - Remove tags by clicking X or pressing Backspace on empty input
  - Visual feedback for dark/light themes
  - Keyboard accessible
-->

<script lang="ts">
	interface Props {
		value?: string[];
		placeholder?: string;
		class?: string;
		disabled?: boolean;
		transformInput?: (value: string) => string;
	}

	let {
		value = $bindable([]),
		placeholder = 'Type and press Enter...',
		class: className = '',
		disabled = false,
		transformInput
	}: Props = $props();

	// Ensure we always have a valid array to work with
	let tags = $derived(value ?? []);

	let inputValue = $state('');

	function handleKeyDown(e: KeyboardEvent) {
		// Add tag on Enter or comma
		if (e.key === 'Enter' || e.key === ',') {
			e.preventDefault();
			addTag();
		}
		// Remove last tag on Backspace if input is empty
		else if (e.key === 'Backspace' && !inputValue && tags.length > 0) {
			removeTag(tags.length - 1);
		}
	}

	function addTag() {
		let trimmedValue = inputValue.trim();
		if (trimmedValue) {
			// Apply transformation if provided (e.g., uppercase for Jira keys)
			if (transformInput) {
				trimmedValue = transformInput(trimmedValue);
			}
			// Check for duplicates after transformation
			if (!tags.includes(trimmedValue)) {
				value = [...tags, trimmedValue];
				inputValue = '';
			} else {
				// Still clear input if it's a duplicate
				inputValue = '';
			}
		}
	}

	function removeTag(indexToRemove: number) {
		value = tags.filter((_, index) => index !== indexToRemove);
	}

	function handleInputBlur() {
		// Add tag on blur if there's content
		if (inputValue.trim()) {
			addTag();
		}
	}
</script>

<div class="tag-input-container {disabled ? 'disabled' : ''} {className}">
	<!-- Render existing tags -->
	{#each tags as tag, index}
		<div class="tag">
			<span>{tag}</span>
			{#if !disabled}
				<button
					type="button"
					onclick={() => removeTag(index)}
					class="tag-remove"
					aria-label="Remove {tag}"
				>
					<svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			{/if}
		</div>
	{/each}

	<!-- Input for new tags -->
	{#if !disabled}
		<input
			type="text"
			bind:value={inputValue}
			onkeydown={handleKeyDown}
			onblur={handleInputBlur}
			placeholder={tags.length === 0 ? placeholder : ''}
			class="tag-input"
		/>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.tag-input {
		.tag-input-container {
			@apply flex flex-wrap gap-2 min-h-[42px] w-full rounded-lg border p-2 bg-transparent;
			border-color: var(--color-border);
		}

		.tag-input-container:focus-within {
			@apply outline-none;
			border-color: var(--color-text-muted);
		}

		.tag-input-container.disabled {
			@apply opacity-50 cursor-not-allowed;
		}

		.tag {
			@apply flex items-center gap-1.5 px-2.5 py-1 rounded-md text-sm font-medium border transition-colors;
			background: color-mix(in srgb, var(--color-info) 10%, transparent);
			color: var(--color-info);
			border-color: color-mix(in srgb, var(--color-info) 30%, transparent);
		}

		.tag-remove {
			@apply rounded-sm p-0.5 focus:outline-none focus:ring-1 transition-colors;
		}

		.tag-remove:hover {
			background: color-mix(in srgb, var(--color-info) 20%, transparent);
		}

		.tag-remove:focus {
			ring-color: var(--color-info);
		}

		.tag-input {
			@apply flex-1 min-w-[120px] bg-transparent outline-none text-sm;
			color: var(--color-text);
		}

		.tag-input::placeholder {
			color: var(--color-text-muted);
		}
	}
</style>
