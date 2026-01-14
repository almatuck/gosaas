<!--
  Reusable Loading Spinner Component
-->

<script lang="ts">
	let { size = 'normal' }: { size?: 'normal' | 'large' } = $props();

	const containerClasses = {
		normal: 'loading-container',
		large: 'loading-container-large'
	};

	const spinnerClasses = {
		normal: 'spinner',
		large: 'spinner-large'
	};
</script>

<div class={containerClasses[size]}>
	<div class={spinnerClasses[size]}></div>
</div>

<style>
@reference "$src/app.css";

@layer components.loading-spinner {
	/* No custom styles needed - using Tailwind utilities */
}
</style>
