<script lang="ts">
	import { validateEmail, validatePassword, validatePasswordConfirmation, validateName, type PasswordRequirements } from '$lib/utils/validation';
	import { auth, authError, authLoading } from '$lib/stores';
	import { tick } from 'svelte';

	interface Props {
		plan?: string;
		onSuccess?: (checkoutUrl?: string) => void;
		onLoginClick?: () => void;
	}

	let { plan = 'free', onSuccess, onLoginClick }: Props = $props();

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);

	let nameError = $state('');
	let emailError = $state('');
	let passwordError = $state('');
	let confirmError = $state('');

	let touched = $state({ name: false, email: false, password: false, confirm: false });

	let nameInputEl: HTMLInputElement | undefined = $state();
	let emailInputEl: HTMLInputElement | undefined = $state();
	let passwordInputEl: HTMLInputElement | undefined = $state();
	let confirmInputEl: HTMLInputElement | undefined = $state();

	// Generate unique IDs for accessibility
	const formId = $derived(`register-form-${Math.random().toString(36).substr(2, 9)}`);
	const nameId = $derived(`${formId}-name`);
	const emailId = $derived(`${formId}-email`);
	const passwordId = $derived(`${formId}-password`);
	const confirmId = $derived(`${formId}-confirm`);
	const nameErrorId = $derived(`${formId}-name-error`);
	const emailErrorId = $derived(`${formId}-email-error`);
	const passwordErrorId = $derived(`${formId}-password-error`);
	const confirmErrorId = $derived(`${formId}-confirm-error`);
	const generalErrorId = $derived(`${formId}-general-error`);
	const passwordHintId = $derived(`${formId}-password-hint`);

	// Reactive validations
	const nameValidation = $derived.by(() => {
		if (!touched.name && !name) return { isValid: true };
		return validateName(name);
	});

	const emailValidation = $derived.by(() => {
		if (!touched.email && !email) return { isValid: true };
		return validateEmail(email);
	});

	const passwordValidation = $derived.by(() => {
		if (!touched.password && !password) return { isValid: true, requirements: undefined as PasswordRequirements | undefined };
		return validatePassword(password);
	});

	const confirmValidation = $derived.by(() => {
		if (!touched.confirm && !confirmPassword) return { isValid: true };
		return validatePasswordConfirmation(password, confirmPassword);
	});

	function handleNameInput() {
		touched.name = true;
		nameError = '';
		auth.clearError();
	}

	function handleEmailInput() {
		touched.email = true;
		emailError = '';
		auth.clearError();
	}

	function handlePasswordInput() {
		touched.password = true;
		passwordError = '';
		auth.clearError();
		// Re-validate confirm if already touched
		if (touched.confirm && confirmPassword) {
			const result = validatePasswordConfirmation(password, confirmPassword);
			confirmError = result.isValid ? '' : (result.error || '');
		}
	}

	function handleConfirmInput() {
		touched.confirm = true;
		confirmError = '';
		auth.clearError();
	}

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		touched = { name: true, email: true, password: true, confirm: true };

		// Validate all fields
		const nameResult = validateName(name);
		if (!nameResult.isValid) {
			nameError = nameResult.error || 'Invalid name';
			await tick();
			nameInputEl?.focus();
			return;
		}

		const emailResult = validateEmail(email);
		if (!emailResult.isValid) {
			emailError = emailResult.error || 'Invalid email';
			await tick();
			emailInputEl?.focus();
			return;
		}

		const passwordResult = validatePassword(password);
		if (!passwordResult.isValid) {
			passwordError = passwordResult.error || 'Invalid password';
			await tick();
			passwordInputEl?.focus();
			return;
		}

		const confirmResult = validatePasswordConfirmation(password, confirmPassword);
		if (!confirmResult.isValid) {
			confirmError = confirmResult.error || 'Passwords do not match';
			await tick();
			confirmInputEl?.focus();
			return;
		}

		nameError = '';
		emailError = '';
		passwordError = '';
		confirmError = '';

		const result = await auth.register({
			name: name.trim(),
			email: email.trim(),
			password,
			plan
		});

		if (result.success) {
			onSuccess?.(result.checkoutUrl);
		}
	}
</script>

<form class="register-form" onsubmit={handleSubmit} aria-label="Registration form">
	<div class="form-header">
		<h2 class="form-title">Create Account</h2>
		<p class="form-subtitle">Create your account to get started</p>
	</div>

	{#if $authError}
		<div id={generalErrorId} class="error-banner" role="alert" aria-live="assertive">
			<span class="error-icon" aria-hidden="true">!</span>
			{$authError}
		</div>
	{/if}

	<div class="form-fields">
		<!-- Name Field -->
		<div class="field-group">
			<label for={nameId} class="field-label">
				Full Name
				<span class="required" aria-hidden="true">*</span>
			</label>
			<input
				id={nameId}
				type="text"
				bind:value={name}
				bind:this={nameInputEl}
				oninput={handleNameInput}
				placeholder="John Doe"
				class="input-field"
				class:error={(!nameValidation.isValid && touched.name) || nameError}
				disabled={$authLoading}
				aria-describedby={(!nameValidation.isValid && touched.name) || nameError ? nameErrorId : undefined}
				aria-invalid={(!nameValidation.isValid && touched.name) || !!nameError}
				aria-required="true"
				autocomplete="name"
			/>
			{#if (!nameValidation.isValid && touched.name) || nameError}
				<div id={nameErrorId} class="field-error" role="alert">
					{nameError || nameValidation.error}
				</div>
			{/if}
		</div>

		<!-- Email Field -->
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

		<!-- Password Field -->
		<div class="field-group">
			<label for={passwordId} class="field-label">
				Password
				<span class="required" aria-hidden="true">*</span>
			</label>
			<div class="password-input-wrapper">
				<input
					id={passwordId}
					type={showPassword ? 'text' : 'password'}
					bind:value={password}
					bind:this={passwordInputEl}
					oninput={handlePasswordInput}
					placeholder="Create a strong password"
					class="input-field"
					class:error={(!passwordValidation.isValid && touched.password) || passwordError}
					disabled={$authLoading}
					aria-describedby={[
						passwordHintId,
						(!passwordValidation.isValid && touched.password) || passwordError ? passwordErrorId : ''
					].filter(Boolean).join(' ')}
					aria-invalid={(!passwordValidation.isValid && touched.password) || !!passwordError}
					aria-required="true"
					autocomplete="new-password"
				/>
				<button
					type="button"
					class="password-toggle"
					onclick={togglePasswordVisibility}
					aria-label={showPassword ? 'Hide password' : 'Show password'}
					disabled={$authLoading}
				>
					{#if showPassword}
						<svg class="toggle-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path>
							<line x1="1" y1="1" x2="23" y2="23"></line>
						</svg>
					{:else}
						<svg class="toggle-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
							<circle cx="12" cy="12" r="3"></circle>
						</svg>
					{/if}
				</button>
			</div>
			{#if (!passwordValidation.isValid && touched.password) || passwordError}
				<div id={passwordErrorId} class="field-error" role="alert">
					{passwordError || passwordValidation.error}
				</div>
			{/if}

			<!-- Password Requirements -->
			{#if touched.password && passwordValidation.requirements}
				<div id={passwordHintId} class="password-requirements" aria-label="Password requirements">
					<ul class="requirements-list">
						<li class:met={passwordValidation.requirements.minLength}>
							<span class="check-icon" aria-hidden="true">{passwordValidation.requirements.minLength ? '✓' : '○'}</span>
							At least 8 characters
						</li>
						<li class:met={passwordValidation.requirements.hasUppercase}>
							<span class="check-icon" aria-hidden="true">{passwordValidation.requirements.hasUppercase ? '✓' : '○'}</span>
							One uppercase letter
						</li>
						<li class:met={passwordValidation.requirements.hasLowercase}>
							<span class="check-icon" aria-hidden="true">{passwordValidation.requirements.hasLowercase ? '✓' : '○'}</span>
							One lowercase letter
						</li>
						<li class:met={passwordValidation.requirements.hasNumber}>
							<span class="check-icon" aria-hidden="true">{passwordValidation.requirements.hasNumber ? '✓' : '○'}</span>
							One number
						</li>
						<li class:met={passwordValidation.requirements.hasSpecialChar} class="optional">
							<span class="check-icon" aria-hidden="true">{passwordValidation.requirements.hasSpecialChar ? '✓' : '○'}</span>
							Special character (recommended)
						</li>
					</ul>
				</div>
			{:else}
				<div id={passwordHintId} class="password-hint">
					Password must be at least 8 characters with uppercase, lowercase, and numbers
				</div>
			{/if}
		</div>

		<!-- Confirm Password Field -->
		<div class="field-group">
			<label for={confirmId} class="field-label">
				Confirm Password
				<span class="required" aria-hidden="true">*</span>
			</label>
			<input
				id={confirmId}
				type={showPassword ? 'text' : 'password'}
				bind:value={confirmPassword}
				bind:this={confirmInputEl}
				oninput={handleConfirmInput}
				placeholder="Confirm your password"
				class="input-field"
				class:error={(!confirmValidation.isValid && touched.confirm) || confirmError}
				disabled={$authLoading}
				aria-describedby={(!confirmValidation.isValid && touched.confirm) || confirmError ? confirmErrorId : undefined}
				aria-invalid={(!confirmValidation.isValid && touched.confirm) || !!confirmError}
				aria-required="true"
				autocomplete="new-password"
			/>
			{#if (!confirmValidation.isValid && touched.confirm) || confirmError}
				<div id={confirmErrorId} class="field-error" role="alert">
					{confirmError || confirmValidation.error}
				</div>
			{/if}
		</div>
	</div>

	<button
		type="submit"
		class="submit-button"
		disabled={$authLoading}
		aria-busy={$authLoading}
	>
		{#if $authLoading}
			<span class="spinner" aria-hidden="true"></span>
			<span class="visually-hidden">Please wait, </span>Creating account...
		{:else}
			Create Account
		{/if}
	</button>

	<p class="terms-text">
		By creating an account, you agree to our
		<a href="/terms" class="link">Terms of Service</a>
		and
		<a href="/privacy" class="link">Privacy Policy</a>
	</p>

	{#if onLoginClick}
		<div class="form-footer">
			<p class="footer-text">
				Already have an account?
				<button
					type="button"
					class="link-button"
					onclick={() => {
						auth.clearError();
						onLoginClick?.();
					}}
					disabled={$authLoading}
				>
					Sign in
				</button>
			</p>
		</div>
	{/if}
</form>

<style>
	/* ===== Form Container - Mobile First ===== */
	.register-form {
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

	/* ===== Password Input ===== */
	.password-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.password-input-wrapper .input-field {
		padding-right: 3rem;
	}

	.password-toggle {
		position: absolute;
		right: var(--space-2, 0.5rem);
		display: flex;
		align-items: center;
		justify-content: center;
		width: 2.5rem;
		height: 2.5rem;
		background: transparent;
		border: none;
		cursor: pointer;
		color: #94a3b8;
		border-radius: var(--radius-sm, 4px);
		transition: color 0.2s ease;
	}

	.password-toggle:hover:not(:disabled) {
		color: #ffffff;
	}

	.password-toggle:focus-visible {
		outline: 2px solid #3b82f6;
		outline-offset: 2px;
	}

	.password-toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.toggle-icon {
		width: 1.25rem;
		height: 1.25rem;
	}

	/* ===== Password Requirements ===== */
	.password-hint {
		font-size: var(--font-size-xs, 0.75rem);
		color: #94a3b8;
		margin-top: var(--space-1, 0.25rem);
	}

	.password-requirements {
		margin-top: var(--space-2, 0.5rem);
		padding: var(--space-2, 0.5rem) var(--space-3, 0.75rem);
		background: rgba(255, 255, 255, 0.05);
		border-radius: var(--radius-sm, 4px);
	}

	.requirements-list {
		list-style: none;
		margin: 0;
		padding: 0;
		font-size: var(--font-size-xs, 0.75rem);
		display: grid;
		gap: var(--space-1, 0.25rem);
	}

	.requirements-list li {
		display: flex;
		align-items: center;
		gap: var(--space-2, 0.5rem);
		color: #94a3b8;
	}

	.requirements-list li.met {
		color: #34d399;
	}

	.requirements-list li.optional {
		font-style: italic;
	}

	.check-icon {
		font-size: 0.75rem;
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

	/* ===== Terms Text ===== */
	.terms-text {
		font-size: var(--font-size-xs, 0.75rem);
		color: #94a3b8;
		text-align: center;
		margin: var(--space-4, 1rem) 0 0;
	}

	/* ===== Form Footer ===== */
	.form-footer {
		text-align: center;
		margin-top: var(--space-4, 1rem);
		padding-top: var(--space-4, 1rem);
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}

	.footer-text {
		font-size: var(--font-size-sm, 0.875rem);
		color: #94a3b8;
		margin: 0;
	}

	/* ===== Links ===== */
	.link,
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

	.link:hover,
	.link-button:hover:not(:disabled) {
		color: #a5b4fc;
		text-decoration: underline;
	}

	.link:focus-visible,
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
		.register-form {
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
		.register-form {
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
