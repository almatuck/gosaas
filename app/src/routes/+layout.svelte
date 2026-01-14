<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { initMonitoring, logger } from '$lib/monitoring';
	import { getMonitoringConfig } from '$lib/monitoring/config';
	import '$src/app.css';

	let { children } = $props();

	onMount(() => {
		const config = getMonitoringConfig();
		initMonitoring(config).then(() => {
			logger.info('GoSaaS application started');
		});

		const handleError = (event: ErrorEvent) => {
			logger.error('Unhandled error', event.error, {
				message: event.message,
				filename: event.filename,
				lineno: event.lineno,
				colno: event.colno
			});
		};

		const handleRejection = (event: PromiseRejectionEvent) => {
			logger.error('Unhandled promise rejection', event.reason);
		};

		window.addEventListener('error', handleError);
		window.addEventListener('unhandledrejection', handleRejection);

		return () => {
			window.removeEventListener('error', handleError);
			window.removeEventListener('unhandledrejection', handleRejection);
		};
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</svelte:head>

<a href="#main-content" class="skip-link">Skip to main content</a>
{@render children()}
