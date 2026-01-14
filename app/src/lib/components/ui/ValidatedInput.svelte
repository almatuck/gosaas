<!--
  ValidatedInput Component
  Input field with real-time validation feedback
-->

<script lang="ts">
	interface FieldValidation<T> {
		validate: (value: T, context?: any) => ValidationResult;
		showOn?: 'change' | 'blur';
		debounceMs?: number;
	}

	interface ValidationResult {
		isValid: boolean;
		hint?: string;
		warning?: string;
		severity?: 'warning' | 'error';
	}

	interface Props {
		value?: string;
		validation?: FieldValidation<string>;
		context?: any;
		showValidation?: boolean;
		forceValidate?: boolean;
		type?: string;
		placeholder?: string;
		disabled?: boolean;
		class?: string;
		onBlur?: (e: FocusEvent & { currentTarget: HTMLInputElement }) => void;
	}

	let {
		value = $bindable(''),
		validation,
		context,
		showValidation = true,
		forceValidate = false,
		type = 'text',
		placeholder = '',
		disabled = false,
		class: className = '',
		onBlur
	}: Props = $props();

	let validationResult = $state<ValidationResult | null>(null);
	let shouldShow = $state(false);
	let debounceTimer: ReturnType<typeof setTimeout> | undefined = $state();
	let showTimer: ReturnType<typeof setTimeout> | undefined = $state();

	// Run validation when value changes
	$effect(() => {
		if (!validation) return;

		// Clear any existing show timer when value changes
		if (showTimer) {
			clearTimeout(showTimer);
			shouldShow = false;
		}

		// Clear previous validation timer
		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		// Run validation immediately to get result, but don't show yet
		const result = validation.validate(value, context);
		validationResult = result;

		// If force validate, show immediately
		if (forceValidate) {
			shouldShow = true;
		}

		// Set timer to show validation after debounce
		if (validation.showOn === 'change' || !validation.showOn) {
			showTimer = setTimeout(() => {
				shouldShow = true;
			}, validation.debounceMs || 500);
		}

		return () => {
			if (debounceTimer) {
				clearTimeout(debounceTimer);
			}
			if (showTimer) {
				clearTimeout(showTimer);
			}
		};
	});

	function handleBlur(e: FocusEvent & { currentTarget: HTMLInputElement }) {
		// Show validation on blur if we have a result
		if (validation && validationResult && (validation.showOn === 'blur' || validation.showOn === 'change' || !validation.showOn)) {
			shouldShow = true;
		}
		onBlur?.(e);
	}

	// Determine validation state class
	const validationClass = $derived(() => {
		if (!shouldShow || !showValidation || !validationResult) {
			return '';
		}

		if (validationResult.severity === 'warning') {
			return 'warning';
		}

		if (!validationResult.isValid) {
			return 'invalid';
		}

		return '';
	});

	// Get the hint text to display
	const hintText = $derived(() => {
		if (!shouldShow || !showValidation || !validationResult) return null;
		return validationResult.hint || validationResult.warning;
	});

	// Get hint class
	const hintClass = $derived(() => {
		if (!validationResult) return '';

		if (validationResult.severity === 'warning') {
			return 'hint-warning';
		}

		return 'hint-info';
	});
</script>

<div class="w-full">
	<input
		{type}
		{placeholder}
		{disabled}
		bind:value
		onblur={handleBlur}
		class="validated-input {validationClass()} {className}"
	/>

	<!-- Validation hint -->
	{#if hintText()}
		<div
			class="validated-hint {hintClass()}"
			role="status"
			aria-live="polite"
		>
			{hintText()}
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.validated-input {
		.validated-input {
			@apply w-full px-4 py-2 rounded-lg text-sm border transition-colors duration-150;
			@apply focus:outline-none;
			background: var(--color-bg);
			color: var(--color-text);
			border-color: var(--color-border);
		}

		.validated-input:focus {
			border-color: var(--color-primary);
		}

		.validated-input.warning {
			border-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
		}

		.validated-input.warning:focus {
			border-color: color-mix(in srgb, var(--color-warning) 50%, transparent);
		}

		.validated-input.invalid {
			border-color: color-mix(in srgb, var(--color-warning) 20%, transparent);
		}

		.validated-input.invalid:focus {
			border-color: color-mix(in srgb, var(--color-warning) 40%, transparent);
		}

		.validated-hint {
			@apply text-xs mt-1.5;
		}

		.hint-warning {
			color: var(--color-warning);
		}

		.hint-info {
			color: var(--color-text-muted);
		}
	}
</style>
