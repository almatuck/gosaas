<!--
  Reusable Table Component
-->

<script lang="ts">
	let { children }: { children: any } = $props();
</script>

<div class="table-container">
	<table class="table">
		{@render children()}
	</table>
</div>

<style>
	@reference "$src/app.css";
	@layer components.table {
		.table-container {
			@apply overflow-x-auto rounded-lg;
			background: color-mix(in srgb, var(--color-base-800) 50%, transparent);
			border: 1px solid var(--color-base-700);
		}

		.table {
			@apply w-full text-sm;
		}

		.table :global(thead) {
			@apply border-b border-base-700;
		}

		.table :global(th) {
			@apply px-4 py-3 text-left text-xs font-semibold uppercase tracking-wider text-base-400;
		}

		.table :global(td) {
			@apply px-4 py-3 text-base-200;
		}

		.table :global(tr) {
			@apply border-b border-base-700/50;
		}

		.table :global(tbody tr:hover) {
			background: color-mix(in srgb, var(--color-base-700) 30%, transparent);
		}
	}
</style>
