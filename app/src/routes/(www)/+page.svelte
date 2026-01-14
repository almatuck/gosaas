<script lang="ts">
	import { setSEO, getFAQSchema } from '$lib/utils/seo';
	import { Button } from '$lib/components/ui';
	import {
		Zap,
		Shield,
		BarChart3,
		ArrowRight,
		Check,
		Linkedin,
		Quote
	} from 'lucide-svelte';

	const seo = setSEO({
		title: 'Your SaaS Name - Build Something Amazing',
		description:
			'The modern platform for building and scaling your SaaS business. Ship faster with authentication, billing, and email built-in.',
		keywords: [
			'saas boilerplate',
			'saas starter kit',
			'sveltekit template',
			'stripe integration',
			'authentication'
		],
		url: '/'
	});

	const features = [
		{
			icon: Zap,
			title: 'Launch Fast',
			description: 'Ship your MVP in days, not months. Authentication, billing, and email are already wired up.',
			color: 'var(--color-primary)'
		},
		{
			icon: Shield,
			title: 'Secure by Default',
			description: 'Enterprise-grade security with JWT auth, CORS protection, and rate limiting built-in.',
			color: 'var(--color-secondary)'
		},
		{
			icon: BarChart3,
			title: 'Scale Effortlessly',
			description: 'Built on proven technologies that grow with you from 0 to millions of users.',
			color: 'var(--color-tertiary)'
		}
	];

	const steps = [
		{
			step: '1',
			title: 'Clone',
			description: 'Fork the repo and run the setup script'
		},
		{
			step: '2',
			title: 'Configure',
			description: 'Add your API keys and customize the theme'
		},
		{
			step: '3',
			title: 'Launch',
			description: 'Deploy to your favorite platform and start selling'
		}
	];

	const faqs = [
		{
			question: 'What is included in this starter?',
			answer:
				'Everything you need to launch a SaaS: user authentication, Stripe billing integration, email sequences, a beautiful landing page, dashboard, and account management.'
		},
		{
			question: 'What technologies does this use?',
			answer:
				'The frontend is built with SvelteKit 5 and Tailwind CSS v4. The backend uses Go with the go-zero framework. Levee handles auth, billing, and email.'
		},
		{
			question: 'Can I use this for commercial projects?',
			answer:
				'Yes! This is a boilerplate designed to be renamed and customized for your own SaaS product. Use it however you like.'
		},
		{
			question: 'How do I customize the design?',
			answer:
				'All styling is centralized in app.css using CSS custom properties. Change the color palette, fonts, and component styles in one place.'
		}
	];

	const faqSchema = JSON.stringify(getFAQSchema(faqs));
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
	{@html `<script type="application/ld+json">${faqSchema}</script>`}
</svelte:head>

<!-- Hero Section -->
<section class="hero-section">
	<div class="animate-fade-in">
		<span class="hero-badge">
			<span class="hero-badge-dot"></span>
			Now in public beta
		</span>

		<h1 class="hero-title">
			Build your SaaS<br />
			<span class="text-gradient">faster than ever</span>
		</h1>

		<p class="hero-subtitle">
			Stop reinventing the wheel. This boilerplate gives you authentication, billing, email,
			and a polished UI so you can focus on what makes your product unique.
		</p>

		<div class="hero-cta">
			<Button type="primary" size="lg" href="/auth/register">
				Get Started Free
				<ArrowRight class="w-5 h-5" />
			</Button>
			<Button type="secondary" size="lg" href="/pricing">View Pricing</Button>
		</div>
	</div>
</section>

<!-- Features Grid -->
<section class="section" id="features">
	<div class="section-header">
		<h2 class="section-title">Everything you need to launch</h2>
		<p class="section-subtitle">
			Built with modern technologies and best practices so you can ship with confidence.
		</p>
	</div>

	<div class="features-grid">
		{#each features as feature, i}
			{@const FeatureIcon = feature.icon}
			{@const delay = i + 1}
			<div class="feature-card animate-slide-up opacity-0 stagger-delay" style:--delay={delay} style:--theme-color={feature.color}>
				<div class="feature-icon themed-bg">
					<FeatureIcon class="w-7 h-7 themed-color" />
				</div>
				<h3 class="feature-title">{feature.title}</h3>
				<p class="feature-description">{feature.description}</p>
			</div>
		{/each}
	</div>
</section>

<!-- How It Works -->
<div class="section-alt">
	<section class="section">
		<div class="section-header">
			<h2 class="section-title">Get started in minutes</h2>
			<p class="section-subtitle">Three simple steps to your production-ready SaaS.</p>
		</div>

		<div class="steps-container">
			<div class="steps-line"></div>
			<div class="steps-grid">
				{#each steps as item, i}
					{@const delay = i + 1}
					<div class="step-item animate-slide-up opacity-0 stagger-delay-150" style:--delay={delay}>
						<div class="step-number">
							<span>{item.step}</span>
						</div>
						<h3 class="step-title">{item.title}</h3>
						<p class="step-description">{item.description}</p>
					</div>
				{/each}
			</div>
		</div>
	</section>
</div>

<!-- What's Included -->
<section class="section">
	<div class="section-header">
		<h2 class="section-title">What's included</h2>
		<p class="section-subtitle">
			A complete foundation for your SaaS business, ready to customize.
		</p>
	</div>

	<div class="features-grid">
		{#each [
			{ title: 'User Authentication', items: ['Email/password login', 'Password reset flow', 'Email verification', 'Session management'] },
			{ title: 'Stripe Billing', items: ['Subscription management', 'Customer portal', 'Usage-based pricing', 'Webhook handling'] },
			{ title: 'Email Marketing', items: ['Transactional emails', 'Drip campaigns', 'List management', 'Template system'] },
			{ title: 'Beautiful UI', items: ['50+ UI components', 'Dark mode ready', 'Responsive design', 'Accessibility built-in'] },
			{ title: 'Developer Experience', items: ['TypeScript throughout', 'API client generation', 'Hot reload', 'Type-safe routing'] },
			{ title: 'Production Ready', items: ['Security headers', 'Rate limiting', 'Error handling', 'Logging & monitoring'] }
		] as section, i}
			{@const delay = i + 1}
			<div class="feature-card animate-slide-up opacity-0 stagger-delay-80" style:--delay={delay}>
				<h3 class="feature-title">{section.title}</h3>
				<ul class="mt-4 space-y-2">
					{#each section.items as item}
						<li class="flex items-center gap-2 text-base-300 text-sm">
							<Check class="w-4 h-4 text-success flex-shrink-0" />
							<span>{item}</span>
						</li>
					{/each}
				</ul>
			</div>
		{/each}
	</div>
</section>

<!-- Testimonial -->
<div class="section-alt">
	<section class="section">
		<div class="max-w-3xl mx-auto">
			<div class="quote-card">
				<Quote class="w-8 h-8 text-primary mb-4" />
				<p class="quote-text">
					"This boilerplate saved us months of development time. We launched our MVP in under two weeks
					and have been iterating on customer feedback ever since."
				</p>
				<div class="quote-author">
					<div class="quote-avatar">JD</div>
					<div>
						<p class="quote-name">Jane Doe</p>
						<p class="quote-role">Founder, Example Startup</p>
					</div>
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
			{#each faqs as faq}
				<div class="faq-item">
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
		<h2 class="cta-title">
			Ready to build your SaaS?
		</h2>
		<p class="cta-subtitle">
			Start with a solid foundation and focus on what makes your product unique.
		</p>
		<Button type="primary" size="lg" href="/auth/register">
			Get Started Free
			<ArrowRight class="w-5 h-5" />
		</Button>
	</section>
</div>
