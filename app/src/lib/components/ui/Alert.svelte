<!--
  Reusable Alert Component
  Ensures 100% consistency for all alerts/notices
  All styling controlled by theme tokens
-->

<script lang="ts">
	import {
		Info,
		CheckCircle,
		AlertTriangle,
		XCircle,
		ShieldAlert,
		Lightbulb,
		TrendingUp,
		Scale
	} from 'lucide-svelte';
	import type { ComponentType } from 'svelte';

	let {
		type = 'info',
		icon,
		title,
		children
	}: {
		type?: 'info' | 'success' | 'warning' | 'error' | 'security';
		icon?: ComponentType | string;
		title: string;
		children: any;
	} = $props();

	const defaultIcons = {
		info: Info,
		success: CheckCircle,
		warning: AlertTriangle,
		error: XCircle,
		security: ShieldAlert
	};

	const iconNameMap: Record<string, ComponentType> = {
		'exclamation-triangle': AlertTriangle,
		lightbulb: Lightbulb,
		'chart-line': TrendingUp,
		'info-circle': Info,
		'balance-scale': Scale
	};

	const Icon = $derived(
		!icon
			? defaultIcons[type]
			: typeof icon === 'string'
				? iconNameMap[icon] || defaultIcons[type]
				: icon
	);
</script>

<div class="alert alert-{type}">
	<div class="alert-content">
		<div class="alert-icon-container">
			<Icon class="alert-icon" />
		</div>
		<div class="alert-text">
			<h3 class="alert-title">{title}</h3>
			<div class="alert-body">
				{@render children()}
			</div>
		</div>
	</div>
</div>

<style>
	@reference "$src/app.css";
	@layer components.alert {
		/* ===== Alerts ===== */
		.alert {
			@apply rounded-xl p-4 border flex items-start gap-3;
		}

		.alert-info {
			background: color-mix(in srgb, var(--color-info) 10%, transparent);
			border-color: color-mix(in srgb, var(--color-info) 30%, transparent);
		}

		.alert-success {
			background: color-mix(in srgb, var(--color-success) 10%, transparent);
			border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
		}

		.alert-warning {
			background: color-mix(in srgb, var(--color-warning) 10%, transparent);
			border-color: color-mix(in srgb, var(--color-warning) 30%, transparent);
		}

		.alert-error {
			background: color-mix(in srgb, var(--color-error) 10%, transparent);
			border-color: color-mix(in srgb, var(--color-error) 30%, transparent);
		}

		.alert-security {
			background: color-mix(in srgb, var(--color-security) 10%, transparent);
			border-color: color-mix(in srgb, var(--color-security) 30%, transparent);
		}

		.alert-icon-container {
			@apply shrink-0 w-10 h-10 rounded-lg flex items-center justify-center;
		}

		.alert-info .alert-icon-container {
			background: color-mix(in srgb, var(--color-info) 20%, transparent);
		}
		.alert-success .alert-icon-container {
			background: color-mix(in srgb, var(--color-success) 20%, transparent);
		}
		.alert-warning .alert-icon-container {
			background: color-mix(in srgb, var(--color-warning) 20%, transparent);
		}
		.alert-error .alert-icon-container {
			background: color-mix(in srgb, var(--color-error) 20%, transparent);
		}
		.alert-security .alert-icon-container {
			background: color-mix(in srgb, var(--color-security) 20%, transparent);
		}

		.alert-content {
			@apply flex items-start gap-3;
		}

		.alert-text {
			@apply flex-1;
		}

		.alert-icon-container :global(.alert-icon) {
			@apply w-5 h-5;
		}

		.alert-info .alert-icon-container :global(.alert-icon) {
			@apply text-info;
		}
		.alert-success .alert-icon-container :global(.alert-icon) {
			@apply text-success;
		}
		.alert-warning .alert-icon-container :global(.alert-icon) {
			@apply text-warning;
		}
		.alert-error .alert-icon-container :global(.alert-icon) {
			@apply text-error;
		}
		.alert-security .alert-icon-container :global(.alert-icon) {
			@apply text-security;
		}

		.alert-title {
			@apply font-semibold text-text;
		}

		.alert-body {
			@apply text-text-secondary mt-1;
		}
	}
</style>
