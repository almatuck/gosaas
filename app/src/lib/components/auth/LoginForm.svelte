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

<form class="login-form" onsubmit={handleSubmit} aria-label="Login form">
	<div class="form-header">
		<h2 class="form-title">Welcome Back</h2>
		<p class="form-subtitle">Sign in to your account to continue</p>
	</div>

	{#if $authError}
		<div id={generalErrorId} class="error-banner" role="alert" aria-live="assertive">
			<span class="error-icon" aria-hidden="true">!</span>
			{$authError}
		</div>
	{/if}

	<div class="form-fields">
		<div class="field-group">
			<label for={emailId} class="field-label">
				Email Address
				<span class="required" aria-hidden="true">*</span>
			</label>
			<input
				id={emailId}
				type="email"
				bind:value={email}
				bind:this={emailInputEl}
				oninput={handleEmailInput}
				placeholder="you@example.com"
				class="input-field"
				class:error={(!emailValidation.isValid && touched.email) || emailError}
				disabled={$authLoading}
				aria-describedby={(!emailValidation.isValid && touched.email) || emailError ? emailErrorId : undefined}
				aria-invalid={(!emailValidation.isValid && touched.email) || !!emailError}
				aria-required="true"
				autocomplete="email"
			/>
			{#if (!emailValidation.isValid && touched.email) || emailError}
				<div id={emailErrorId} class="field-error" role="alert">
					{emailError || emailValidation.error}
				</div>
			{/if}
		</div>

		<div class="field-group">
			<label for={passwordId} class="field-label">
				Password
				<span class="required" aria-hidden="true">*</span>
			</label>
			<input
				id={passwordId}
				type="password"
				bind:value={password}
				bind:this={passwordInputEl}
				oninput={handlePasswordInput}
				placeholder="Enter your password"
				class="input-field"
				class:error={passwordError}
				disabled={$authLoading}
				aria-describedby={passwordError ? passwordErrorId : undefined}
				aria-invalid={!!passwordError}
				aria-required="true"
				autocomplete="current-password"
			/>
			{#if passwordError}
				<div id={passwordErrorId} class="field-error" role="alert">
					{passwordError}
				</div>
			{/if}
		</div>

		{#if onForgotPasswordClick}
			<div class="forgot-password-wrapper">
				<button
					type="button"
					class="link-button"
					onclick={onForgotPasswordClick}
					disabled={$authLoading}
				>
					Forgot your password?
				</button>
			</div>
		{/if}
	</div>

	<button
		type="submit"
		class="submit-button"
		disabled={$authLoading}
		aria-busy={$authLoading}
	>
		{#if $authLoading}
			<span class="spinner" aria-hidden="true"></span>
			<span class="visually-hidden">Please wait, </span>Signing in...
		{:else}
			Sign In
		{/if}
	</button>

	{#if onRegisterClick}
		<div class="form-footer">
			<p class="footer-text">
				Don't have an account?
				<button
					type="button"
					class="link-button"
					onclick={() => {
						auth.clearError();
						onRegisterClick?.();
					}}
					disabled={$authLoading}
				>
					Create one
				</button>
			</p>
		</div>
	{/if}
</form>

<style>
	/* ===== Form Container - Mobile First ===== */
	.login-form {
		max-width: var(--max-width-sm, 480px);
		margin: 0 auto;
		padding: var(--space-6, 1.5rem);
		background: var(--color-base-800, #1e293b);
		border-radius: var(--radius-lg, 12px);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	/* ===== Form Header ===== */
	.form-header {
		text-align: center;
		margin-bottom: var(--space-6, 1.5rem);
	}

	.form-title {
		font-size: var(--font-size-2xl, 1.5rem);
		font-weight: 700;
		color: #ffffff;
		margin: 0 0 var(--space-2, 0.5rem);
	}

	.form-subtitle {
		font-size: var(--font-size-sm, 0.875rem);
		color: #94a3b8;
		margin: 0;
	}

	/* ===== Error Banner ===== */
	.error-banner {
		display: flex;
		align-items: flex-start;
		gap: var(--space-2, 0.5rem);
		padding: var(--space-3, 0.75rem);
		margin-bottom: var(--space-4, 1rem);
		font-size: var(--font-size-sm, 0.875rem);
		color: #b91c1c;
		background: #fef2f2;
		border-radius: var(--radius-md, 8px);
		border-left: 4px solid #ef4444;
	}

	.error-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1.25rem;
		height: 1.25rem;
		background: #ef4444;
		color: #ffffff;
		font-size: var(--font-size-xs, 0.75rem);
		font-weight: 700;
		border-radius: 50%;
		flex-shrink: 0;
	}

	/* ===== Form Fields ===== */
	.form-fields {
		display: flex;
		flex-direction: column;
		gap: var(--space-4, 1rem);
		margin-bottom: var(--space-6, 1.5rem);
	}

	.field-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-1, 0.25rem);
	}

	.field-label {
		font-size: var(--font-size-sm, 0.875rem);
		font-weight: 500;
		color: #e2e8f0;
	}

	.required {
		color: #f87171;
		margin-left: 2px;
	}

	.input-field {
		width: 100%;
		padding: var(--space-3, 0.75rem) var(--space-4, 1rem);
		font-size: var(--font-size-base, 1rem);
		color: #ffffff;
		background: var(--color-base-700, #334155);
		border: 1px solid rgba(255, 255, 255, 0.1);
		border-radius: var(--radius-md, 8px);
		transition: all 0.2s ease;
		box-sizing: border-box;
		min-height: var(--touch-target-min, 44px);
	}

	.input-field::placeholder {
		color: #64748b;
	}

	.input-field:focus {
		outline: none;
		border-color: var(--color-accent-primary, #6366f1);
		background: var(--color-base-700, #334155);
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2);
	}

	.input-field:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.input-field.error {
		border-color: #ef4444;
		background: rgba(239, 68, 68, 0.1);
	}

	.input-field.error:focus {
		box-shadow: 0 0 0 3px rgb(239 68 68 / 0.2);
	}

	.field-error {
		font-size: var(--font-size-xs, 0.75rem);
		color: #f87171;
		margin-top: var(--space-1, 0.25rem);
	}

	/* ===== Forgot Password ===== */
	.forgot-password-wrapper {
		display: flex;
		justify-content: flex-end;
	}

	/* ===== Submit Button ===== */
	.submit-button {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2, 0.5rem);
		width: 100%;
		padding: var(--space-3, 0.75rem) var(--space-4, 1rem);
		font-size: var(--font-size-base, 1rem);
		font-weight: 600;
		color: #ffffff;
		background: linear-gradient(135deg, var(--color-accent-primary, #6366f1), var(--color-accent-primary-dark, #4f46e5));
		border: none;
		border-radius: var(--radius-md, 8px);
		cursor: pointer;
		transition: all 0.2s ease;
		min-height: var(--touch-target-comfortable, 48px);
	}

	.submit-button:hover:not(:disabled) {
		background: linear-gradient(135deg, var(--color-accent-primary-dark, #4f46e5), #4338ca);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
	}

	.submit-button:active:not(:disabled) {
		transform: translateY(0);
	}

	.submit-button:focus-visible {
		outline: 3px solid #1d4ed8;
		outline-offset: 2px;
	}

	.submit-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	.spinner {
		width: 1rem;
		height: 1rem;
		border: 2px solid rgb(255 255 255 / 0.3);
		border-top-color: #ffffff;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	/* ===== Form Footer ===== */
	.form-footer {
		text-align: center;
		margin-top: var(--space-6, 1.5rem);
		padding-top: var(--space-4, 1rem);
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}

	.footer-text {
		font-size: var(--font-size-sm, 0.875rem);
		color: #94a3b8;
		margin: 0;
	}

	/* ===== Link Button ===== */
	.link-button {
		font-size: inherit;
		font-weight: 500;
		color: var(--color-accent-primary-light, #818cf8);
		background: none;
		border: none;
		padding: 0;
		cursor: pointer;
		text-decoration: none;
		transition: color 0.2s ease;
	}

	.link-button:hover:not(:disabled) {
		color: #a5b4fc;
		text-decoration: underline;
	}

	.link-button:focus-visible {
		outline: 2px solid #3b82f6;
		outline-offset: 2px;
		border-radius: 2px;
	}

	.link-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	/* ===== Accessibility ===== */
	.visually-hidden {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}

	/* ===== Tablets (640px+) ===== */
	@media (min-width: 640px) {
		.login-form {
			padding: var(--space-8, 2rem);
		}

		.form-title {
			font-size: var(--font-size-3xl, 1.875rem);
		}
	}

	/* ===== Touch Device Optimizations ===== */
	@media (hover: none) and (pointer: coarse) {
		.input-field {
			min-height: var(--touch-target-comfortable, 48px);
		}

		.submit-button:hover:not(:disabled) {
			transform: none;
			box-shadow: none;
		}
	}

	/* ===== Reduced Motion ===== */
	@media (prefers-reduced-motion: reduce) {
		.input-field,
		.submit-button,
		.link-button,
		.spinner {
			transition: none;
			animation: none;
		}
	}

	/* ===== High Contrast Mode ===== */
	@media (prefers-contrast: high) {
		.login-form {
			border: 2px solid #ffffff;
		}

		.input-field {
			border-width: 2px;
			border-color: #ffffff;
		}

		.submit-button {
			border: 2px solid var(--color-accent-primary-light, #818cf8);
		}
	}
</style>
