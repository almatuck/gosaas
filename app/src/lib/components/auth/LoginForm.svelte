<script lang="ts">
	import { validateEmail } from '$lib/utils/validation';
	import { auth, authError, authLoading } from '$lib/stores';
	import { tick } from 'svelte';

	interface Props {
		onSuccess?: () => void;
		onRegisterClick?: () => void;
		onForgotPasswordClick?: () => void;
	}

	let { onSuccess, onRegisterClick, onForgotPasswordClick }: Props = $props();

	let email = $state('');
	let password = $state('');
	let emailError = $state('');
	let passwordError = $state('');
	let touched = $state({ email: false, password: false });
	let emailInputEl: HTMLInputElement | undefined = $state();
	let passwordInputEl: HTMLInputElement | undefined = $state();

	// Generate unique IDs for accessibility
	const formId = $derived(`login-form-${Math.random().toString(36).substr(2, 9)}`);
	const emailId = $derived(`${formId}-email`);
	const passwordId = $derived(`${formId}-password`);
	const emailErrorId = $derived(`${formId}-email-error`);
	const passwordErrorId = $derived(`${formId}-password-error`);
	const generalErrorId = $derived(`${formId}-general-error`);

	// Reactive validation
	const emailValidation = $derived.by(() => {
		if (!touched.email && !email) return { isValid: true };
		return validateEmail(email);
	});

	const hasEmailError = $derived((!emailValidation.isValid && touched.email) || !!emailError);
	const hasPasswordError = $derived(!!passwordError);

	function handleEmailInput() {
		touched.email = true;
		emailError = '';
		auth.clearError();
	}

	function handlePasswordInput() {
		touched.password = true;
		passwordError = '';
		auth.clearError();
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		touched = { email: true, password: true };

		// Validate email
		const emailResult = validateEmail(email);
		if (!emailResult.isValid) {
			emailError = emailResult.error || 'Invalid email';
			await tick();
			emailInputEl?.focus();
			return;
		}

		// Validate password
		if (!password) {
			passwordError = 'Please enter your password';
			await tick();
			passwordInputEl?.focus();
			return;
		}

		emailError = '';
		passwordError = '';

		const success = await auth.login({ email: email.trim(), password });

		if (success) {
			onSuccess?.();
		}
	}
</script>

<form class="card bg-base-200 w-full max-w-md mx-auto" onsubmit={handleSubmit} aria-label="Login form">
	<div class="card-body">
		<div class="text-center mb-4">
			<h2 class="card-title justify-center text-2xl">Welcome Back</h2>
			<p class="text-base-content/60 text-sm">Sign in to your account to continue</p>
		</div>

		{#if $authError}
			<div id={generalErrorId} class="alert alert-error mb-4" role="alert" aria-live="assertive">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
				</svg>
				<span>{$authError}</span>
			</div>
		{/if}

		<div class="space-y-4">
			<div class="form-control w-full">
				<label for={emailId} class="label">
					<span class="label-text">Email Address <span class="text-error">*</span></span>
				</label>
				<input
					id={emailId}
					type="email"
					bind:value={email}
					bind:this={emailInputEl}
					oninput={handleEmailInput}
					placeholder="you@example.com"
					class="input input-bordered w-full"
					class:input-error={hasEmailError}
					disabled={$authLoading}
					aria-describedby={hasEmailError ? emailErrorId : undefined}
					aria-invalid={hasEmailError}
					aria-required="true"
					autocomplete="email"
				/>
				{#if hasEmailError}
					<label class="label" id={emailErrorId}>
						<span class="label-text-alt text-error">{emailError || emailValidation.error}</span>
					</label>
				{/if}
			</div>

			<div class="form-control w-full">
				<label for={passwordId} class="label">
					<span class="label-text">Password <span class="text-error">*</span></span>
				</label>
				<input
					id={passwordId}
					type="password"
					bind:value={password}
					bind:this={passwordInputEl}
					oninput={handlePasswordInput}
					placeholder="Enter your password"
					class="input input-bordered w-full"
					class:input-error={hasPasswordError}
					disabled={$authLoading}
					aria-describedby={hasPasswordError ? passwordErrorId : undefined}
					aria-invalid={hasPasswordError}
					aria-required="true"
					autocomplete="current-password"
				/>
				{#if hasPasswordError}
					<label class="label" id={passwordErrorId}>
						<span class="label-text-alt text-error">{passwordError}</span>
					</label>
				{/if}
			</div>

			{#if onForgotPasswordClick}
				<div class="text-right">
					<button
						type="button"
						class="link link-primary text-sm"
						onclick={onForgotPasswordClick}
						disabled={$authLoading}
					>
						Forgot your password?
					</button>
				</div>
			{/if}
		</div>

		<div class="form-control mt-6">
			<button
				type="submit"
				class="btn btn-primary w-full"
				disabled={$authLoading}
				aria-busy={$authLoading}
			>
				{#if $authLoading}
					<span class="loading loading-spinner loading-sm"></span>
					Signing in...
				{:else}
					Sign In
				{/if}
			</button>
		</div>

		{#if onRegisterClick}
			<div class="divider"></div>
			<p class="text-center text-sm text-base-content/60">
				Don't have an account?
				<button
					type="button"
					class="link link-primary"
					onclick={() => {
						auth.clearError();
						onRegisterClick?.();
					}}
					disabled={$authLoading}
				>
					Create one
				</button>
			</p>
		{/if}
	</div>
</form>
