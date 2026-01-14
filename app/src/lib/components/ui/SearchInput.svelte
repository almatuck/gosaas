<!--
  Search Input Component
  Input with search icon and clear button
-->

<script lang="ts">
	interface Props {
		value?: string;
		placeholder?: string;
		disabled?: boolean;
		loading?: boolean;
		onsearch?: (value: string) => void;
		onclear?: () => void;
	}

	let {
		value = $bindable(''),
		placeholder = 'Search...',
		disabled = false,
		loading = false,
		onsearch,
		onclear
	}: Props = $props();

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
		if (onsearch) {
			onsearch(value);
		}
	}

	function handleClear() {
		value = '';
		if (onclear) {
			onclear();
		}
	}
</script>

<div class="relative">
	<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
		{#if loading}
			<div class="h-5 w-5 animate-spin rounded-full border-2 border-indigo-600 border-t-transparent"></div>
		{:else}
			<i class="fas fa-search text-slate-400"></i>
		{/if}
	</div>
	<input
		type="text"
		bind:value
		{placeholder}
		{disabled}
		oninput={handleInput}
		class="input w-full pl-11 {value ? 'pr-11' : ''} text-xs sm:text-sm md:text-base"
	/>
	{#if value && !disabled}
		<button
			type="button"
			onclick={handleClear}
			class="absolute inset-y-0 right-0 flex items-center pr-4 text-slate-400 hover:text-white transition-colors"
			aria-label="Clear search"
		>
			<i class="fas fa-times"></i>
		</button>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.search-input {
		/* Search input uses utility classes and input class from app.css */
	}
</style>
