<!--
  ThemeToggle Component
  Toggle between light, dark, and system theme
-->

<script lang="ts">
	import { writable } from 'svelte/store';

	// Simple theme store (you may want to connect this to a global theme store)
	const theme = writable<'light' | 'dark' | 'system'>('system');

	let isOpen = $state(false);

	function setTheme(newTheme: 'light' | 'dark' | 'system') {
		theme.set(newTheme);
		isOpen = false;

		// Apply theme to document
		if (newTheme === 'system') {
			const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
			document.documentElement.classList.toggle('dark', isDark);
		} else {
			document.documentElement.classList.toggle('dark', newTheme === 'dark');
		}
	}

	function toggleDropdown() {
		isOpen = !isOpen;
	}

	function handleClickOutside(e: MouseEvent) {
		if (isOpen && !(e.target as HTMLElement).closest('.theme-toggle-container')) {
			isOpen = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<div class="relative theme-toggle-container">
	<button
		type="button"
		onclick={toggleDropdown}
		class="rounded-full border border-gray-200 dark:border-gray-700 p-2 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
		aria-label="Toggle theme"
	>
		<svg class="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
		</svg>
		<svg class="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100 top-2 left-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
		</svg>
		<span class="sr-only">Toggle theme</span>
	</button>

	{#if isOpen}
		<div class="absolute right-0 mt-2 w-40 rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 z-50">
			<div class="py-1" role="menu" aria-orientation="vertical">
				<button
					type="button"
					onclick={() => setTheme('light')}
					class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2"
					role="menuitem"
				>
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
					<span>Light</span>
				</button>
				<button
					type="button"
					onclick={() => setTheme('dark')}
					class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2"
					role="menuitem"
				>
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
					</svg>
					<span>Dark</span>
				</button>
				<button
					type="button"
					onclick={() => setTheme('system')}
					class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2"
					role="menuitem"
				>
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
					</svg>
					<span>System</span>
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
@reference "$src/app.css";

@layer components.theme-toggle {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
