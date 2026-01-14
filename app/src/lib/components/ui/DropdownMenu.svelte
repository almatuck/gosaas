<!--
  DropdownMenu Component
  Dropdown menu with items, checkboxes, radio items, etc.
-->

<script lang="ts">
	interface MenuItem {
		label: string;
		icon?: any;
		onClick?: () => void;
		separator?: boolean;
		disabled?: boolean;
		type?: 'item' | 'checkbox' | 'radio';
		checked?: boolean;
	}

	interface Props {
		trigger?: any;
		items: MenuItem[];
		align?: 'start' | 'end';
		sideOffset?: number;
	}

	let {
		trigger,
		items = [],
		align = 'end',
		sideOffset = 4
	}: Props = $props();

	let isOpen = $state(false);

	function toggleMenu() {
		isOpen = !isOpen;
	}

	function handleItemClick(item: MenuItem) {
		if (!item.disabled) {
			item.onClick?.();
			if (item.type !== 'checkbox' && item.type !== 'radio') {
				isOpen = false;
			}
		}
	}

	function handleClickOutside(e: MouseEvent) {
		if (isOpen && !(e.target as HTMLElement).closest('.dropdown-menu-container')) {
			isOpen = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<div class="relative dropdown-menu-container">
	{#if trigger}
		<div onclick={toggleMenu} role="button" tabindex="0" onkeydown={(e) => e.key === 'Enter' && toggleMenu()}>
			{@render trigger()}
		</div>
	{/if}

	{#if isOpen}
		<div
			class="ddmenu {align === 'end' ? 'ddmenu-end' : 'ddmenu-start'}"
			role="menu"
		>
			{#each items as item}
				{#if item.separator}
					<div class="ddmenu-separator"></div>
				{:else}
					<button
						type="button"
						onclick={() => handleItemClick(item)}
						disabled={item.disabled}
						class="ddmenu-item"
						role="menuitem"
					>
						{#if item.type === 'checkbox' || item.type === 'radio'}
							<span class="ddmenu-indicator">
								{#if item.checked}
									{#if item.type === 'checkbox'}
										<svg class="ddmenu-check" fill="none" viewBox="0 0 24 24" stroke="currentColor">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
										</svg>
									{:else}
										<svg class="ddmenu-radio" viewBox="0 0 8 8">
											<circle cx="4" cy="4" r="4" />
										</svg>
									{/if}
								{/if}
							</span>
							<span class="ddmenu-label-with-indicator">{item.label}</span>
						{:else}
							{#if item.icon}
								<span class="ddmenu-icon">
									{@render item.icon()}
								</span>
							{/if}
							<span>{item.label}</span>
						{/if}
					</button>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.dropdown-menu {
		.ddmenu {
			@apply absolute z-50 mt-2 min-w-[12rem] rounded-xl p-1 shadow-xl;
			background: var(--color-base-800);
			border: 1px solid color-mix(in srgb, var(--color-text) 10%, transparent);
		}

		.ddmenu-start {
			@apply left-0;
		}

		.ddmenu-end {
			@apply right-0;
		}

		.ddmenu-separator {
			@apply my-1 h-px;
			background: var(--color-base-700);
		}

		.ddmenu-item {
			@apply relative flex w-full items-center rounded-lg px-2 py-1.5 text-sm;
			@apply transition-colors cursor-pointer text-left;
			color: var(--color-text);

			&:hover:not(:disabled) {
				background: color-mix(in srgb, var(--color-base-700) 50%, transparent);
			}

			&:disabled {
				@apply opacity-50 cursor-not-allowed;
			}
		}

		.ddmenu-indicator {
			@apply absolute left-2 flex h-3.5 w-3.5 items-center justify-center;
		}

		.ddmenu-check {
			@apply h-4 w-4;
		}

		.ddmenu-radio {
			@apply h-2 w-2 fill-current;
		}

		.ddmenu-label-with-indicator {
			@apply pl-6;
		}

		.ddmenu-icon {
			@apply mr-2 h-4 w-4;
		}
	}
</style>
