<script lang="ts">
	import { SparklineChart } from '$lib/components/charts';

	interface Props {
		title?: string;
		value: string | number;
		change?: number;
		changeType?: 'increase' | 'decrease' | 'neutral';
		subtitle?: string;
		icon?: string;
		loading?: boolean;
		error?: string;
		format?: 'number' | 'currency' | 'percentage' | 'ratio';
		precision?: number;
		sparklineData?: number[]; // Historical data points for mini chart
	}

	let {
		title,
		value,
		change,
		changeType = 'neutral',
		subtitle,
		icon,
		loading = false,
		error,
		format = 'number',
		precision = 2,
		sparklineData = []
	}: Props = $props();

	function formatValue(val: string | number): string {
		if (typeof val === 'string') return val;
		if (val === undefined || val === null || typeof val !== 'number') return '0';

		switch (format) {
			case 'currency':
				return new Intl.NumberFormat('en-US', {
					style: 'currency',
					currency: 'USD',
					minimumFractionDigits: 0,
					maximumFractionDigits: precision
				}).format(val);
			case 'percentage':
				return `${(val * 100).toFixed(precision)}%`;
			case 'ratio':
				return `${val.toFixed(precision)}x`;
			default:
				return new Intl.NumberFormat('en-US', {
					minimumFractionDigits: 0,
					maximumFractionDigits: precision
				}).format(val);
		}
	}

	function formatChange(val: number | undefined): string {
		if (val === undefined || val === null || typeof val !== 'number') {
			return '0.0%';
		}
		const sign = val >= 0 ? '+' : '';
		return `${sign}${val.toFixed(1)}%`;
	}

	// Determine card styling based on performance
	let cardClasses = $derived.by(() => {
		const base = 'relative overflow-hidden rounded-lg p-5 transition-all duration-200';

		if (changeType === 'increase') {
			return `${base} bg-gradient-to-br from-green-500/10 via-slate-800/50 to-slate-800 border border-green-500/30 hover:border-green-500/50`;
		} else if (changeType === 'decrease') {
			return `${base} bg-gradient-to-br from-red-500/10 via-slate-800/50 to-slate-800 border border-red-500/30 hover:border-red-500/50`;
		} else {
			return `${base} bg-gradient-to-br from-slate-700/50 via-slate-800/50 to-slate-800 border border-slate-600/50 hover:border-slate-500/50`;
		}
	});

	let sparklineColor = $derived.by(() => {
		if (changeType === 'increase') return 'var(--color-accent-success)';
		if (changeType === 'decrease') return 'var(--color-accent-error)';
		return 'var(--color-base-400)';
	});
</script>

<div class="w-full h-full">
	{#if loading}
		<div class="{cardClasses} h-full">
			<div class="animate-pulse space-y-3">
				{#if title}
					<div class="h-3 w-24 bg-slate-700 rounded"></div>
				{/if}
				<div class="h-8 w-24 bg-slate-700 rounded"></div>
				<div class="h-2 w-20 bg-slate-700 rounded"></div>
			</div>
		</div>
	{:else if error}
		<div class="relative overflow-hidden rounded-lg p-4 bg-gradient-to-br from-red-500/10 via-slate-800/50 to-slate-800 border border-red-500/30 h-full flex flex-col justify-center">
			<div class="text-red-400 text-center">
				{#if title}
					<div class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-2">{title}</div>
				{/if}
				<i class="fas fa-exclamation-triangle text-3xl mb-2"></i>
				<div class="text-xs">{error}</div>
			</div>
		</div>
	{:else}
		<div class="{cardClasses} h-full flex flex-col">
			<!-- Background accent gradient -->
			<div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-white/5 to-transparent rounded-full blur-2xl -z-10"></div>

			<!-- Content -->
			<div class="relative z-10 flex-1 flex flex-col">
				<!-- Title Row -->
				{#if title}
					<div class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-3">{title}</div>
				{/if}

				<!-- Main Value & Change -->
				<div class="mb-4">
					<div class="flex items-end gap-3 mb-1">
						<div class="text-4xl font-bold text-white tracking-tight leading-none">
							{formatValue(value)}
						</div>
						{#if change !== undefined}
							<div class="flex items-center gap-1 px-2.5 py-1 rounded-full mb-1 {changeType === 'increase' ? 'bg-green-500/20 text-green-400' : changeType === 'decrease' ? 'bg-red-500/20 text-red-400' : 'bg-slate-500/20 text-slate-400'}">
								{#if changeType === 'increase'}
									<svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
									</svg>
								{:else if changeType === 'decrease'}
									<svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z" clip-rule="evenodd" />
									</svg>
								{/if}
								<span class="text-xs font-bold">{formatChange(change)}</span>
							</div>
						{/if}
					</div>
					{#if subtitle}
						<div class="text-xs text-slate-400">{subtitle}</div>
					{/if}
				</div>

				<!-- Large Sparkline Chart (takes remaining space) -->
				{#if sparklineData && sparklineData.length > 0}
					<div class="flex-1 min-h-[80px] -mb-2">
						<SparklineChart
							data={sparklineData}
							width={300}
							height={80}
							color={sparklineColor}
							fill={true}
						/>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.metric-card {
		/* Metric card uses utility classes */
	}
</style>
