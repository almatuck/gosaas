<!--
  Breadcrumb Component
  Navigation breadcrumbs for hierarchy
-->

<script lang="ts">
	interface BreadcrumbItem {
		label: string;
		href?: string;
		icon?: string;
	}

	interface Props {
		items: BreadcrumbItem[];
		separator?: 'slash' | 'chevron' | 'arrow';
	}

	let {
		items,
		separator = 'chevron'
	}: Props = $props();

	const separatorIcons = {
		slash: '/',
		chevron: 'fa-chevron-right',
		arrow: 'fa-arrow-right'
	};
</script>

<nav class="breadcrumb" aria-label="Breadcrumb">
	<ol class="breadcrumb-list">
		{#each items as item, index}
			<li class="breadcrumb-item">
				{#if index > 0}
					<span class="breadcrumb-separator">
						{#if separator === 'slash'}
							/
						{:else}
							<i class="fas {separatorIcons[separator]} text-xs"></i>
						{/if}
					</span>
				{/if}

				{#if item.href && index < items.length - 1}
					<a
						href={item.href}
						class="breadcrumb-link"
					>
						{#if item.icon}
							<i class="fas fa-{item.icon}"></i>
						{/if}
						{item.label}
					</a>
				{:else}
					<span class="breadcrumb-current">
						{#if item.icon}
							<i class="fas fa-{item.icon}"></i>
						{/if}
						{item.label}
					</span>
				{/if}
			</li>
		{/each}
	</ol>
</nav>

<style>
	@reference "$src/app.css";
	@layer components.breadcrumb {
		.breadcrumb {
			@apply flex;
		}

		.breadcrumb-list {
			@apply flex items-center space-x-2;
		}

		.breadcrumb-item {
			@apply flex items-center;
		}

		.breadcrumb-separator {
			@apply mx-2;
			color: var(--color-base-500);
		}

		.breadcrumb-link {
			@apply flex items-center gap-2 text-sm font-medium transition-colors;
			color: var(--color-base-400);
		}

		.breadcrumb-link:hover {
			@apply text-white;
		}

		.breadcrumb-current {
			@apply flex items-center gap-2 text-sm font-medium text-white;
		}
	}
</style>
