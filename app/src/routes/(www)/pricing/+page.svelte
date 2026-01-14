<script lang="ts">
	import { setSEO } from '$lib/utils/seo';
	import { Button } from '$lib/components/ui';
	import { Check, ArrowRight, Zap, Building2, Rocket } from 'lucide-svelte';
	import type { PricingPageData } from './+page.server';

	let { data }: { data: PricingPageData } = $props();

	const seo = setSEO({
		title: 'Pricing - Simple, Transparent Pricing',
		description: 'Choose the plan that fits your needs. Start free and scale as you grow. No hidden fees, no surprises.',
		keywords: ['pricing', 'plans', 'subscription', 'saas pricing'],
		url: '/pricing'
	});

	// Map product index to icon
	const icons = [Zap, Rocket, Building2];
	const colors = ['var(--color-base-500)', 'var(--color-accent-primary)', 'var(--color-accent-secondary)'];

	// Format price from cents
	function formatPrice(cents: number): string {
		if (cents === 0) return '$0';
		return `$${(cents / 100).toFixed(0)}`;
	}

	// Get period text
	function getPeriod(price: { interval?: string; unitAmountCents: number }): string {
		if (price.unitAmountCents === 0) return 'forever';
		if (price.interval === 'month') return '/month';
		if (price.interval === 'year') return '/year';
		return '';
	}

	// Determine if product is popular (highlighted price or middle product)
	function isPopular(productIndex: number): boolean {
		const product = data.products[productIndex];
		return product.prices?.some(p => p.highlighted) || productIndex === 1;
	}

	const faqs = [
		{
			question: 'Can I switch plans anytime?',
			answer: 'Yes! You can upgrade or downgrade your plan at any time. Changes take effect immediately, and we\'ll prorate any billing differences.'
		},
		{
			question: 'Is there a free trial for paid plans?',
			answer: 'Yes, Pro and Team plans come with a 14-day free trial. No credit card required to start.'
		},
		{
			question: 'Do you offer annual billing?',
			answer: 'Yes! Annual billing saves you 20%. Pay yearly and get 2 months free on any paid plan.'
		},
		{
			question: 'What payment methods do you accept?',
			answer: 'We accept all major credit cards (Visa, Mastercard, American Express) through our secure payment processor, Stripe.'
		}
	];
</script>

<svelte:head>
	<title>{seo.title}</title>
	<meta name="description" content={seo.description} />
	<meta name="keywords" content={seo.keywords} />
	<link rel="canonical" href={seo.canonical} />

	<meta property="og:type" content={seo.type} />
	<meta property="og:title" content={seo.title} />
	<meta property="og:description" content={seo.description} />
	<meta property="og:url" content={seo.url} />
	<meta property="og:image" content={seo.image} />
	<meta property="og:site_name" content={seo.siteName} />
	<meta property="og:locale" content={seo.locale} />

	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:site" content={seo.twitterHandle} />
	<meta name="twitter:title" content={seo.title} />
	<meta name="twitter:description" content={seo.description} />
	<meta name="twitter:image" content={seo.image} />
</svelte:head>

<!-- Hero Section -->
<section class="hero-section">
	<div class="animate-fade-in">
		<h1 class="hero-title">
			Simple, Transparent<br />
			<span class="text-gradient">Pricing</span>
		</h1>

		<p class="hero-subtitle">
			Start free and scale as you grow. No hidden fees, no surprises.
			Every plan includes our core features.
		</p>
	</div>
</section>

<!-- Pricing Cards -->
<section class="section">
	<div class="grid md:grid-cols-3 gap-6 lg:gap-8">
		{#each data.products as product, i}
			{@const price = product.prices?.[0]}
			{@const popular = isPopular(i)}
			{@const Icon = icons[i] || Zap}
			{@const color = colors[i] || colors[0]}
			{@const delay = i + 1}
			<div
				class="relative animate-slide-up opacity-0 stagger-delay rounded-2xl border transition-all duration-300 {popular ? 'bg-gradient-to-b from-[var(--color-accent-primary)]/20 to-transparent border-[var(--color-accent-primary)]/40' : 'bg-base-800/80 border-base-600'}"
				style:--delay={delay}
				style:--theme-color={color}
			>
				{#if popular}
					<div class="absolute -top-4 left-1/2 -translate-x-1/2">
						<span class="px-4 py-1.5 rounded-full text-xs font-bold bg-[var(--color-accent-primary)] text-base-900 uppercase tracking-wider">
							Most Popular
						</span>
					</div>
				{/if}

				<div class="p-8">
					<div class="flex items-center gap-3 mb-4">
						<div class="w-10 h-10 rounded-xl flex items-center justify-center themed-bg">
							<Icon class="w-5 h-5 themed-color" />
						</div>
						<h3 class="font-display text-xl font-bold text-white">{product.name}</h3>
					</div>

					<p class="text-base-400 text-sm mb-6">{product.description}</p>

					<div class="mb-8">
						{#if price}
							<span class="font-display text-4xl font-black text-white">{formatPrice(price.unitAmountCents)}</span>
							<span class="text-base-400 text-sm">{getPeriod(price)}</span>
						{/if}
					</div>

					<Button
						type={popular ? 'primary' : 'secondary'}
						size="md"
						href="/auth/register?plan={product.id}"
						class="w-full justify-center mb-8"
					>
						{price?.unitAmountCents === 0 ? 'Get Started' : `Choose ${product.name}`}
						<ArrowRight class="w-4 h-4" />
					</Button>

					{#if product.features && product.features.length > 0}
						<ul class="space-y-3">
							{#each product.features as feature}
								<li class="flex items-start gap-3 text-sm">
									<Check class="w-4 h-4 mt-0.5 flex-shrink-0 themed-color" />
									<span class="text-base-200">{feature}</span>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</div>
		{/each}
	</div>
</section>

<!-- Enterprise Section -->
<div class="section-alt">
	<section class="section">
		<div class="quote-card">
			<div class="grid lg:grid-cols-2 gap-8 items-center">
				<div>
					<span class="inline-flex items-center gap-2 rounded-full border border-[var(--color-accent-secondary)]/20 bg-[var(--color-accent-secondary)]/10 px-3 py-1.5 text-xs font-medium text-[var(--color-accent-secondary)] mb-4">
						<Building2 class="w-3.5 h-3.5" />
						ENTERPRISE
					</span>
					<h2 class="font-display text-2xl sm:text-3xl font-bold text-white mb-4">
						Need a custom solution?
					</h2>
					<p class="text-base-300 leading-relaxed mb-6">
						For large organizations with specific requirements, we offer custom plans with
						dedicated support, SLA guarantees, and tailored integrations.
					</p>
					<ul class="space-y-2 mb-8">
						{#each ['Custom user limits', 'SSO & SAML authentication', 'Dedicated account manager', 'Custom contracts & invoicing'] as item}
							<li class="flex items-center gap-2 text-sm text-base-200">
								<Check class="w-4 h-4 text-[var(--color-accent-secondary)]" />
								{item}
							</li>
						{/each}
					</ul>
				</div>
				<div class="text-center lg:text-right">
					<Button type="secondary" size="lg" href="mailto:sales@example.com">
						Contact Sales
					</Button>
				</div>
			</div>
		</div>
	</section>
</div>

<!-- FAQ Section -->
<section class="section" id="faq">
	<div class="section-header">
		<h2 class="section-title">Frequently Asked Questions</h2>
	</div>

	<div class="max-w-3xl mx-auto">
		<div class="faq-list">
			{#each faqs as faq, i}
				{@const delay = i + 1}
				<div class="faq-item animate-slide-up opacity-0 stagger-delay" style:--delay={delay}>
					<h3 class="faq-question">{faq.question}</h3>
					<p class="faq-answer">{faq.answer}</p>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- CTA Section -->
<div class="section-alt">
	<section class="cta-section">
		<h2 class="cta-title">Ready to get started?</h2>
		<p class="cta-subtitle">
			No credit card required. Start building today.
		</p>
		<Button type="primary" size="lg" href="/auth/register">
			Get Started Free
			<ArrowRight class="w-5 h-5" />
		</Button>
	</section>
</div>
