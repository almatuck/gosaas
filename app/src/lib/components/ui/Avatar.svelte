<!--
  Avatar Component
  User profile images with fallback initials
-->

<script lang="ts">
	interface Props {
		src?: string;
		alt?: string;
		initials?: string;
		size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
		status?: 'online' | 'offline' | 'away' | 'busy';
		rounded?: boolean;
	}

	let {
		src,
		alt = '',
		initials = '',
		size = 'md',
		status,
		rounded = true
	}: Props = $props();

	let imageError = $state(false);

	const sizeClasses = {
		xs: 'h-6 w-6 text-xs',
		sm: 'h-8 w-8 text-sm',
		md: 'h-10 w-10 text-base',
		lg: 'h-12 w-12 text-lg',
		xl: 'h-16 w-16 text-xl'
	};

	const statusClasses = {
		online: 'bg-green-500 ring-green-500/20',
		offline: 'bg-slate-500 ring-slate-500/20',
		away: 'bg-yellow-500 ring-yellow-500/20',
		busy: 'bg-red-500 ring-red-500/20'
	};

	const statusSizes = {
		xs: 'h-1.5 w-1.5',
		sm: 'h-2 w-2',
		md: 'h-2.5 w-2.5',
		lg: 'h-3 w-3',
		xl: 'h-4 w-4'
	};
</script>

<div class="avatar-wrapper">
	<div
		class="avatar avatar-{size} {rounded ? 'avatar-rounded' : 'avatar-square'}"
	>
		{#if src && !imageError}
			<img
				{src}
				{alt}
				onerror={() => imageError = true}
				class="avatar-image"
			/>
		{:else if initials}
			{initials}
		{:else}
			<i class="fas fa-user"></i>
		{/if}
	</div>

	{#if status}
		<span
			class="avatar-status avatar-status-{size} avatar-status-{status} {rounded ? 'avatar-status-rounded' : 'avatar-status-square'}"
		></span>
	{/if}
</div>

<style>
	@reference "$src/app.css";
	@layer components.avatar {
		.avatar-wrapper {
			@apply relative inline-block;
		}

		.avatar {
			@apply inline-flex items-center justify-center overflow-hidden font-semibold text-white;
			background: linear-gradient(to bottom right, var(--color-primary), var(--color-secondary));
		}

		.avatar-xs { @apply h-6 w-6 text-xs; }
		.avatar-sm { @apply h-8 w-8 text-sm; }
		.avatar-md { @apply h-10 w-10 text-base; }
		.avatar-lg { @apply h-12 w-12 text-lg; }
		.avatar-xl { @apply h-16 w-16 text-xl; }

		.avatar-rounded { @apply rounded-full; }
		.avatar-square { @apply rounded-lg; }

		.avatar-image {
			@apply h-full w-full object-cover;
		}

		.avatar-status {
			@apply absolute bottom-0 right-0 block ring-2;
			ring-color: var(--color-base-900);
		}

		.avatar-status-xs { @apply h-1.5 w-1.5; }
		.avatar-status-sm { @apply h-2 w-2; }
		.avatar-status-md { @apply h-2.5 w-2.5; }
		.avatar-status-lg { @apply h-3 w-3; }
		.avatar-status-xl { @apply h-4 w-4; }

		.avatar-status-rounded { @apply rounded-full; }
		.avatar-status-square { @apply rounded; }

		.avatar-status-online {
			background-color: var(--color-success);
		}

		.avatar-status-offline {
			@apply bg-base-500;
		}

		.avatar-status-away {
			background-color: var(--color-warning);
		}

		.avatar-status-busy {
			background-color: var(--color-error);
		}
	}
</style>
