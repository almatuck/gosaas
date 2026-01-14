<script lang="ts">
	import { RefreshCw, Download, Maximize2, Minimize2, X, AlertCircle } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		subtitle?: string;
		description?: string;
		loading?: boolean;
		error?: string;
		lastUpdated?: Date;
		onRefresh?: () => void | Promise<void>;
		onExport?: (format: 'png' | 'csv' | 'json') => void | Promise<void>;
		children: Snippet;
		class?: string;
	}

	let {
		title,
		subtitle,
		description,
		loading = false,
		error,
		lastUpdated,
		onRefresh,
		onExport,
		children,
		class: className = ''
	}: Props = $props();

	let isFullscreen = $state(false);
	let isRefreshing = $state(false);
	let showExportMenu = $state(false);

	async function handleRefresh() {
		if (!onRefresh || isRefreshing) return;
		isRefreshing = true;
		try {
			await onRefresh();
		} finally {
			isRefreshing = false;
		}
	}

	async function handleExport(format: 'png' | 'csv' | 'json') {
		if (!onExport) return;
		showExportMenu = false;
		await onExport(format);
	}

	function toggleFullscreen() {
		isFullscreen = !isFullscreen;
	}

	function formatLastUpdated(date: Date): string {
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const minutes = Math.floor(diff / 60000);
		const hours = Math.floor(diff / 3600000);
		const days = Math.floor(diff / 86400000);

		if (minutes < 1) return 'Just now';
		if (minutes < 60) return `${minutes}m ago`;
		if (hours < 24) return `${hours}h ago`;
		return `${days}d ago`;
	}
</script>

{#if isFullscreen}
	<!-- Backdrop - rendered first so it's behind the card -->
	<button
		class="fixed inset-0 bg-black/80 backdrop-blur-sm z-40"
		onclick={toggleFullscreen}
		aria-label="Close fullscreen"
	></button>
{/if}

<div
	class="chart-card-wrapper card {className} {isFullscreen ? 'fullscreen-card' : ''}"
>
	<!-- Header -->
	<div class="card-header">
		<div class="flex items-start justify-between gap-4">
			<div class="flex-1 min-w-0">
				<h3 class="card-title truncate">{title}</h3>
				{#if subtitle}
					<p class="card-subtitle truncate">{subtitle}</p>
				{/if}
				{#if description}
					<p class="text-xs text-slate-500 mt-1">{description}</p>
				{/if}
			</div>

			<div class="flex items-center gap-2 flex-shrink-0">
				{#if lastUpdated}
					<span class="text-xs text-slate-500 hidden sm:inline">
						{formatLastUpdated(lastUpdated)}
					</span>
				{/if}

				{#if onRefresh}
					<button
						class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/50 transition-colors disabled:opacity-50"
						onclick={handleRefresh}
						disabled={isRefreshing || loading}
						title="Refresh"
						aria-label="Refresh chart"
					>
						<RefreshCw class="w-4 h-4 {isRefreshing ? 'animate-spin' : ''}" />
					</button>
				{/if}

				{#if onExport}
					<div class="relative">
						<button
							class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/50 transition-colors"
							onclick={() => showExportMenu = !showExportMenu}
							title="Export"
							aria-label="Export chart"
						>
							<Download class="w-4 h-4" />
						</button>

						{#if showExportMenu}
							<div
								class="absolute right-0 mt-2 w-36 rounded-lg bg-slate-800 border border-slate-700 shadow-xl z-10"
								onmouseleave={() => showExportMenu = false}
								role="menu"
								tabindex="-1"
							>
								<button
									class="w-full px-4 py-2 text-left text-sm text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors first:rounded-t-lg"
									onclick={() => handleExport('png')}
								>
									Export as PNG
								</button>
								<button
									class="w-full px-4 py-2 text-left text-sm text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors"
									onclick={() => handleExport('csv')}
								>
									Export as CSV
								</button>
								<button
									class="w-full px-4 py-2 text-left text-sm text-slate-300 hover:bg-slate-700/50 hover:text-white transition-colors last:rounded-b-lg"
									onclick={() => handleExport('json')}
								>
									Export as JSON
								</button>
							</div>
						{/if}
					</div>
				{/if}

				<button
					class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/50 transition-colors"
					onclick={toggleFullscreen}
					title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
					aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'}
				>
					{#if isFullscreen}
						<Minimize2 class="w-4 h-4" />
					{:else}
						<Maximize2 class="w-4 h-4" />
					{/if}
				</button>

				{#if isFullscreen}
					<button
						class="p-2 rounded-lg text-slate-400 hover:text-white hover:bg-slate-700/50 transition-colors"
						onclick={toggleFullscreen}
						title="Close"
						aria-label="Close fullscreen"
					>
						<X class="w-4 h-4" />
					</button>
				{/if}
			</div>
		</div>
	</div>

	<!-- Body -->
	<div class="card-body {isFullscreen ? 'flex-1 overflow-auto' : ''}">
		{#if error}
			<div class="alert alert-error flex items-start gap-3">
				<AlertCircle class="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" />
				<div class="flex-1">
					<p class="alert-title text-red-300">Error loading chart</p>
					<p class="alert-body text-red-400/80">{error}</p>
				</div>
				{#if onRefresh}
					<button
						class="btn-secondary text-xs py-1"
						onclick={handleRefresh}
					>
						Retry
					</button>
				{/if}
			</div>
		{:else if loading}
			<div class="flex flex-col items-center justify-center py-12">
				<div class="spinner-large mb-4"></div>
				<p class="text-sm text-slate-400">Loading chart data...</p>
			</div>
		{:else}
			{@render children()}
		{/if}
	</div>
</div>

<style>
	.chart-card-wrapper {
		position: relative;
		transition: all 0.3s ease;
	}

	.fullscreen-card {
		position: fixed !important;
		top: 1rem !important;
		right: 1rem !important;
		bottom: 1rem !important;
		left: 1rem !important;
		z-index: 50 !important;
		display: flex !important;
		flex-direction: column !important;
		box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25) !important;
		animation: slideIn 0.2s ease;
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}
</style>
