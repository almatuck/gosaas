<!--
  Reusable Button Component
-->

<script lang="ts">
	let {
		type = 'primary',
		size = 'md',
		disabled = false,
		active = false,
		onclick,
		href,
		external = false,
		htmlType = 'button',
		children,
		class: extraClass = ''
	}: {
		type?: 'primary' | 'secondary' | 'link' | 'danger' | 'ghost';
		size?: 'sm' | 'md' | 'lg' | 'icon';
		disabled?: boolean;
		active?: boolean;
		onclick?: () => void;
		href?: string;
		external?: boolean;
		htmlType?: 'button' | 'submit' | 'reset';
		children: any;
		class?: string;
	} = $props();

	const typeClasses = {
		primary: 'btn-primary',
		secondary: 'btn-secondary',
		link: 'btn-link',
		danger: 'btn-danger',
		ghost: 'btn-ghost'
	};

	const sizeClasses = {
		sm: 'btn-sm',
		md: 'btn-md',
		lg: 'btn-lg',
		icon: 'btn-icon'
	};

	const className = $derived(`${typeClasses[type]} ${sizeClasses[size]} ${active ? 'active' : ''} ${extraClass}`.trim());
</script>

{#if href}
	<a {href} class={className} target={external ? '_blank' : undefined} rel={external ? 'noopener noreferrer' : undefined}>
		{@render children()}
	</a>
{:else}
	<button class={className} type={htmlType} {disabled} {onclick}>
		{@render children()}
	</button>
{/if}

<style>
	@reference "$src/app.css";
	@layer components.button {
		/* Button styles are defined in app.css */
	}
</style>
