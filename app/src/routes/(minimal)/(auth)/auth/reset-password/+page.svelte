<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { PasswordResetForm } from '$lib/components/auth';

	const token = $derived($page.url.searchParams.get('token') || '');
	const mode = $derived(token ? 'reset' : 'request') as 'request' | 'reset';

	function handleSuccess() {
		if (mode === 'reset') {
			goto('/auth/login');
		}
	}

	function handleBackToLogin() {
		goto('/auth/login');
	}
</script>

<svelte:head>
	<title>Reset Password - GoSaaS</title>
	<meta name="description" content="Reset your GoSaaS account password." />
</svelte:head>

<h1 class="font-display text-2xl font-bold text-white text-center mb-2">
	{mode === 'reset' ? 'Set new password' : 'Reset your password'}
</h1>
<p class="text-center text-[var(--color-base-400)] text-sm mb-8">
	{mode === 'reset'
		? 'Enter your new password below'
		: "Enter your email and we'll send you a reset link"}
</p>

<PasswordResetForm {mode} {token} onSuccess={handleSuccess} onBackToLogin={handleBackToLogin} />
