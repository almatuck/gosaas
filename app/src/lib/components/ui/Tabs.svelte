<!--
  Tabs Component
  Tab navigation with modern styling
-->

<script lang="ts">
	import type { ComponentType } from 'svelte';

	interface Tab {
		id: string;
		label: string;
		icon?: ComponentType;
		badge?: string | number;
	}

	interface Props {
		tabs: Tab[];
		activeTab?: string;
		onchange?: (tabId: string) => void;
		variant?: 'underline' | 'pills' | 'boxed';
	}

	let {
		tabs,
		activeTab = $bindable(''),
		onchange,
		variant = 'underline'
	}: Props = $props();

	function handleTabClick(tabId: string) {
		activeTab = tabId;
		if (onchange) {
			onchange(tabId);
		}
	}
</script>

{#if variant === 'underline'}
	<!-- Underline Tabs -->
	<div class="border-b-2 border-slate-700/50">
		<nav class="-mb-px flex space-x-8" aria-label="Tabs">
			{#each tabs as tab}
				<button
					type="button"
					onclick={() => handleTabClick(tab.id)}
					class="group inline-flex items-center gap-3 border-b-3 px-2 py-4 text-lg font-bold transition-all
						{activeTab === tab.id
							? 'border-indigo-400 text-white'
							: 'border-transparent text-slate-500 hover:border-slate-500 hover:text-slate-300'}"
					aria-current={activeTab === tab.id ? 'page' : undefined}
				>
					{#if tab.icon}
						{@const IconComponent = tab.icon}
				<IconComponent class="w-5 h-5 {activeTab === tab.id ? 'text-indigo-400' : 'text-slate-500'}" />
					{/if}
					{tab.label}
					{#if tab.badge}
						<span class="ml-2 rounded-full px-3 py-1 text-sm font-bold
							{activeTab === tab.id
								? 'bg-indigo-500/20 text-indigo-300'
								: 'bg-slate-700/50 text-slate-500'}">
							{tab.badge}
						</span>
					{/if}
				</button>
			{/each}
		</nav>
	</div>
{:else if variant === 'pills'}
	<!-- Pill Tabs -->
	<nav class="flex space-x-2" aria-label="Tabs">
		{#each tabs as tab}
			<button
				type="button"
				onclick={() => handleTabClick(tab.id)}
				class="inline-flex items-center gap-2 rounded-xl px-4 py-2.5 text-base font-semibold transition-all
					{activeTab === tab.id
						? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg shadow-indigo-500/20'
						: 'bg-slate-800/60 text-slate-300 hover:bg-slate-700/60 hover:text-white'}"
				aria-current={activeTab === tab.id ? 'page' : undefined}
			>
				{#if tab.icon}
					{@const IconComponent = tab.icon}
					<IconComponent class="w-4 h-4" />
				{/if}
				{tab.label}
				{#if tab.badge}
					<span class="ml-1 rounded-full bg-white/20 px-2 py-0.5 text-xs font-medium">
						{tab.badge}
					</span>
				{/if}
			</button>
		{/each}
	</nav>
{:else if variant === 'boxed'}
	<!-- Boxed Tabs -->
	<div class="rounded-xl bg-slate-800/40 p-1">
		<nav class="flex space-x-1" aria-label="Tabs">
			{#each tabs as tab}
				<button
					type="button"
					onclick={() => handleTabClick(tab.id)}
					class="flex-1 inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-base font-semibold transition-all
						{activeTab === tab.id
							? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white shadow-lg'
							: 'text-slate-400 hover:text-white hover:bg-slate-700/50'}"
					aria-current={activeTab === tab.id ? 'page' : undefined}
				>
					{#if tab.icon}
						{@const IconComponent = tab.icon}
					<IconComponent class="w-4 h-4" />
					{/if}
					{tab.label}
					{#if tab.badge}
						<span class="ml-1 rounded-full px-2 py-0.5 text-xs font-medium
							{activeTab === tab.id ? 'bg-white/20' : 'bg-slate-700'}">
							{tab.badge}
						</span>
					{/if}
				</button>
			{/each}
		</nav>
	</div>
{/if}

<style>
@reference "$src/app.css";

@layer components.tabs {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
