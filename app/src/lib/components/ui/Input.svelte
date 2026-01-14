<!--
  Reusable Input Component
-->

<script lang="ts">
	let {
		type = 'text',
		size = 'md',
		placeholder = '',
		value = $bindable(),
		disabled = false,
		id,
		class: extraClass = '',
		oninput,
		onblur
	}: {
		type?: string;
		size?: 'sm' | 'md' | 'lg';
		placeholder?: string;
		value?: string;
		disabled?: boolean;
		id?: string;
		class?: string;
		oninput?: (e: Event) => void;
		onblur?: (e: Event) => void;
	} = $props();

	// Ensure value is never undefined
	let safeValue = $derived(value ?? '');

	const sizeClasses = {
		sm: 'input-sm',
		md: 'input-md',
		lg: 'input-lg'
	};

	const className = $derived(`input ${sizeClasses[size]} ${extraClass}`.trim());

	function handleInput(e: Event) {
		value = (e.currentTarget as HTMLInputElement).value;
		oninput?.(e);
	}
</script>

<input class={className} {type} {placeholder} value={safeValue} oninput={handleInput} onblur={onblur} {disabled} {id} />

<style>
	@reference "$src/app.css";
	@layer components.input {
		/* Size variants for input component */
		.input-sm {
			@apply h-8 px-3 py-1 text-sm;
		}

		.input-md {
			/* Default size - inherits from .input in app.css */
		}

		.input-lg {
			@apply h-12 px-5 py-4 text-base;
		}
	}
</style>
