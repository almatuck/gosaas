<script lang="ts">
	import { validateEmail, validatePassword, validatePasswordConfirmation, type PasswordRequirements } from '$lib/utils/validation';
	import { passwordReset } from '$lib/stores';
	import { tick } from 'svelte';

	type Mode = 'request' | 'reset';

	interface Props {
		mode?: Mode;
		token?: string;
		onSuccess?: () => void;
		onBackToLogin?: () => void;
	}

	let { mode = 'request', token = '', onSuccess, onBackToLogin }: Props = $props();

	// Request mode state
	let email = $state('');
	let emailError = $state('');
	let emailTouched = $state(false);

	// Reset mode state
	let newPassword = $state('');
	let confirmPassword = $state('');
	let passwordError = $state('');
	let confirmError = $state('');
	let showPassword = $state(false);
	let passwordTouched = $state(false);
	let confirmTouched = $state(false);

	let emailInputEl: HTMLInputElement | undefined = $state();
	let passwordInputEl: HTMLInputElement | undefined = $state();
	let confirmInputEl: HTMLInputElement | undefined = $state();

	// Get store values
	let isLoading = $derived($passwordReset.isLoading);
	let isSuccess = $derived($passwordReset.isSuccess);
	let storeError = $derived($passwordReset.error);

	// Generate unique IDs for accessibility
	const formId = $derived(`password-reset-form-${Math.random().toString(36).substr(2, 9)}`);
	const emailId = $derived(`${formId}-email`);
	const passwordId = $derived(`${formId}-password`);
	const confirmId = $derived(`${formId}-confirm`);
	const emailErrorId = $derived(`${formId}-email-error`);
	const passwordErrorId = $derived(`${formId}-password-error`);
	const confirmErrorId = $derived(`${formId}-confirm-error`);
	const generalErrorId = $derived(`${formId}-general-error`);
	const passwordHintId = $derived(`${formId}-password-hint`);

	// Reactive validations
	const emailValidation = $derived.by(() => {
		if (!emailTouched && !email) return { isValid: true };
		return validateEmail(email);
	});

	const passwordValidation = $derived.by(() => {
		if (!passwordTouched && !newPassword) return { isValid: true, requirements: undefined as PasswordRequirements | undefined };
		return validatePassword(newPassword);
	});

	const confirmValidation = $derived.by(() => {
		if (!confirmTouched && !confirmPassword) return { isValid: true };
		return validatePasswordConfirmation(newPassword, confirmPassword);
	});

	function handleEmailInput() {
		emailTouched = true;
		emailError = '';
		passwordReset.reset();
	}

	function handlePasswordInput() {
		passwordTouched = true;
		passwordError = '';
		passwordReset.reset();
		// Re-validate confirm if already touched
		if (confirmTouched && confirmPassword) {
			const result = validatePasswordConfirmation(newPassword, confirmPassword);
			confirmError = result.isValid ? '' : (result.error || '');
		}
	}

	function handleConfirmInput() {
		confirmTouched = true;
		confirmError = '';
		passwordReset.reset();
	}

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	async function handleRequestSubmit(event: SubmitEvent) {
		event.preventDefault();
		emailTouched = true;

		const emailResult = validateEmail(email);
		if (!emailResult.isValid) {
			emailError = emailResult.error || 'Invalid email';
			await tick();
			emailInputEl?.focus();
			return;
		}

		emailError = '';
		const success = await passwordReset.requestReset(email.trim());

		if (success) {
			onSuccess?.();
		}
	}

	async function handleResetSubmit(event: SubmitEvent) {
		event.preventDefault();
		passwordTouched = true;
		confirmTouched = true;

		const passwordResult = validatePassword(newPassword);
		if (!passwordResult.isValid) {
			passwordError = passwordResult.error || 'Invalid password';
			await tick();
			passwordInputEl?.focus();
			return;
		}

		const confirmResult = validatePasswordConfirmation(newPassword, confirmPassword);
		if (!confirmResult.isValid) {
			confirmError = confirmResult.error || 'Passwords do not match';
			await tick();
			confirmInputEl?.focus();
			return;
		}

		passwordError = '';
		confirmError = '';

		const success = await passwordReset.resetPassword(token, newPassword);

		if (success) {
			onSuccess?.();
		}
	}
</script>

{#if mode === 'request'}
	<!-- Request Password Reset Form -->
	{#if isSuccess}
		<div class="success-card" role="status" aria-live="polite">
			<div class="success-icon" aria-hidden="true">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
					<polyline points="22 4 12 14.01 9 11.01"></polyline>
				</svg>
			</div>
			<h2 class="success-title">Check Your Email</h2>
			<p class="success-message">
				We've sent password reset instructions to <strong>{email}</strong>.
				Please check your inbox and follow the link to reset your password.
			</p>
			<p class="success-note">
				Didn't receive the email? Check your spam folder or
				<button
					type="button"
					class="link-button"
					onclick={() => passwordReset.reset()}
				>
					try again
				</button>
			</p>
			{#if onBackToLogin}
				<button
					type="button"
					class="secondary-button"
					onclick={onBackToLogin}
				>
					Back to Sign In
				</button>
			{/if}
		</div>
	{:else}
		<form class="reset-form" onsubmit={handleRequestSubmit} aria-label="Password reset request form">
			<div class="form-header">
				<h2 class="form-title">Reset Password</h2>
				<p class="form-subtitle">
					Enter your email address and we'll send you instructions to reset your password.
				</p>
			</div>

			{#if storeError}
				<div id={generalErrorId} class="error-banner" role="alert" aria-live="assertive">
					<span class="error-icon" aria-hidden="true">!</span>
					{storeError}
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
						class:error={(!emailValidation.isValid && emailTouched) || emailError}
						disabled={isLoading}
						aria-describedby={(!emailValidation.isValid && emailTouched) || emailError ? emailErrorId : undefined}
						aria-invalid={(!emailValidation.isValid && emailTouched) || !!emailError}
						aria-required="true"
						autocomplete="email"
					/>
					{#if (!emailValidation.isValid && emailTouched) || emailError}
						<div id={emailErrorId} class="field-error" role="alert">
							{emailError || emailValidation.error}
						</div>
					{/if}
				</div>
			</div>

			<button
				type="submit"
				class="submit-button"
				disabled={isLoading}
				aria-busy={isLoading}
			>
				{#if isLoading}
					<span class="spinner" aria-hidden="true"></span>
					<span class="visually-hidden">Please wait, </span>Sending...
				{:else}
					Send Reset Link
				{/if}
			</button>

			{#if onBackToLogin}
				<div class="form-footer">
					<p class="footer-text">
						Remember your password?
						<button
							type="button"
							class="link-button"
							onclick={onBackToLogin}
							disabled={isLoading}
						>
							Sign in
						</button>
					</p>
				</div>
			{/if}
		</form>
	{/if}

{:else}
	<!-- Reset Password Form -->
	{#if isSuccess}
		<div class="success-card" role="status" aria-live="polite">
			<div class="success-icon" aria-hidden="true">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
					<polyline points="22 4 12 14.01 9 11.01"></polyline>
				</svg>
			</div>
			<h2 class="success-title">Password Reset Complete</h2>
			<p class="success-message">
				Your password has been successfully reset. You can now sign in with your new password.
			</p>
			{#if onBackToLogin}
				<button
					type="button"
					class="submit-button"
					onclick={onBackToLogin}
				>
					Sign In
				</button>
			{/if}
		</div>
	{:else}
		<form class="reset-form" onsubmit={handleResetSubmit} aria-label="Create new password form">
			<div class="form-header">
				<h2 class="form-title">Create New Password</h2>
				<p class="form-subtitle">
					Enter your new password below.
				</p>
			</div>

			{#if storeError}
				<div id={generalErrorId} class="error-banner" role="alert" aria-live="assertive">
					<span class="error-icon" aria-hidden="true">!</span>
					{storeError}
				</div>
			{/if}

			<div class="form-fields">
				<!-- New Password Field -->
				<div class="field-group">
					<label for={passwordId} class="field-label">
						New Password
						<span class="required" aria-hidden="true">*</span>
					</label>
					<div class="password-input-wrapper">
						<input
							id={passwordId}
							type={showPassword ? 'text' : 'password'}
							bind:value={newPassword}
							bind:this={passwordInputEl}
							oninput={handlePasswordInput}
							placeholder="Enter your new password"
							class="input-field"
							class:error={(!passwordValidation.isValid && passwordTouched) || passwordError}
							disabled={isLoading}
							aria-describedby={[
								passwordHintId,
								(!passwordValidation.isValid && passwordTouched) || passwordError ? passwordErrorId : ''
							].filter(Boolean).join(' ')}
							aria-invalid={(!passwordValidation.isValid && passwordTouched) || !!passwordError}
							aria-required="true"
							autocomplete="new-password"
						/>
						<button
							type="button"
							class="password-toggle"
							onclick={togglePasswordVisibility}
							aria-label={showPassword ? 'Hide password' : 'Show password'}
							disabled={isLoading}
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
					{#if (!passwordValidation.isValid && passwordTouched) || passwordError}
						<div id={passwordErrorId} class="field-error" role="alert">
							{passwordError || passwordValidation.error}
						</div>
					{/if}

					<!-- Password Requirements -->
					{#if passwordTouched && passwordValidation.requirements}
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
						Confirm New Password
						<span class="required" aria-hidden="true">*</span>
					</label>
					<input
						id={confirmId}
						type={showPassword ? 'text' : 'password'}
						bind:value={confirmPassword}
						bind:this={confirmInputEl}
						oninput={handleConfirmInput}
						placeholder="Confirm your new password"
						class="input-field"
						class:error={(!confirmValidation.isValid && confirmTouched) || confirmError}
						disabled={isLoading}
						aria-describedby={(!confirmValidation.isValid && confirmTouched) || confirmError ? confirmErrorId : undefined}
						aria-invalid={(!confirmValidation.isValid && confirmTouched) || !!confirmError}
						aria-required="true"
						autocomplete="new-password"
					/>
					{#if (!confirmValidation.isValid && confirmTouched) || confirmError}
						<div id={confirmErrorId} class="field-error" role="alert">
							{confirmError || confirmValidation.error}
						</div>
					{/if}
				</div>
			</div>

			<button
				type="submit"
				class="submit-button"
				disabled={isLoading}
				aria-busy={isLoading}
			>
				{#if isLoading}
					<span class="spinner" aria-hidden="true"></span>
					<span class="visually-hidden">Please wait, </span>Resetting...
				{:else}
					Reset Password
				{/if}
			</button>

			{#if onBackToLogin}
				<div class="form-footer">
					<p class="footer-text">
						<button
							type="button"
							class="link-button"
							onclick={onBackToLogin}
							disabled={isLoading}
						>
							Back to Sign In
						</button>
					</p>
				</div>
			{/if}
		</form>
	{/if}
{/if}

<style>
	/* ===== Form Container - Mobile First ===== */
	.reset-form,
	.success-card {
		max-width: var(--max-width-sm, 480px);
		margin: 0 auto;
		padding: var(--space-6, 1.5rem);
		background: var(--color-base-800, #1e293b);
		border-radius: var(--radius-lg, 12px);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	/* ===== Success Card ===== */
	.success-card {
		text-align: center;
	}

	.success-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 4rem;
		height: 4rem;
		margin-bottom: var(--space-4, 1rem);
		color: #34d399;
		background: rgba(52, 211, 153, 0.2);
		border-radius: 50%;
	}

	.success-icon svg {
		width: 2rem;
		height: 2rem;
	}

	.success-title {
		font-size: var(--font-size-xl, 1.25rem);
		font-weight: 700;
		color: #ffffff;
		margin: 0 0 var(--space-2, 0.5rem);
	}

	.success-message {
		font-size: var(--font-size-sm, 0.875rem);
		color: #94a3b8;
		margin: 0 0 var(--space-4, 1rem);
		line-height: 1.6;
	}

	.success-note {
		font-size: var(--font-size-xs, 0.75rem);
		color: #64748b;
		margin: 0 0 var(--space-6, 1.5rem);
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
		line-height: 1.5;
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

	.check-icon {
		font-size: 0.75rem;
	}

	/* ===== Buttons ===== */
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

	.secondary-button {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		padding: var(--space-3, 0.75rem) var(--space-4, 1rem);
		font-size: var(--font-size-base, 1rem);
		font-weight: 500;
		color: #ffffff;
		background: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: var(--radius-md, 8px);
		cursor: pointer;
		transition: all 0.2s ease;
		min-height: var(--touch-target-comfortable, 48px);
	}

	.secondary-button:hover {
		background: rgba(255, 255, 255, 0.15);
		border-color: rgba(255, 255, 255, 0.3);
	}

	.secondary-button:focus-visible {
		outline: 3px solid #3b82f6;
		outline-offset: 2px;
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
		.reset-form,
		.success-card {
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
		.secondary-button,
		.link-button,
		.spinner {
			transition: none;
			animation: none;
		}
	}

	/* ===== High Contrast Mode ===== */
	@media (prefers-contrast: high) {
		.reset-form,
		.success-card {
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
