<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { RegisterData } from '$lib/types/auth';

	// Form state
	let username = '';
	let email = '';
	let password = '';
	let confirmPassword = '';
	let isLoading = false;
	let errorMessage = '';
	let successMessage = '';

	// Validation errors
	let usernameError = '';
	let emailError = '';
	let passwordError = '';
	let confirmPasswordError = '';

	// Get return URL from query params
	$: returnUrl = $page.url.searchParams.get('returnUrl') || '/lobby';

	// Redirect if already authenticated
	onMount(() => {
		const restored = auth.loadAuthFromStorage();
		if (restored) {
			goto(returnUrl);
		}
	});

	function validateUsername(value: string): string {
		if (!value) {
			return 'Username is required';
		}
		if (value.length < 3) {
			return 'Username must be at least 3 characters';
		}
		if (value.length > 20) {
			return 'Username must be no more than 20 characters';
		}
		if (!/^[a-zA-Z0-9_]+$/.test(value)) {
			return 'Username can only contain letters, numbers, and underscores';
		}
		return '';
	}

	function validateEmail(value: string): string {
		if (!value) {
			return 'Email is required';
		}
		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		if (!emailRegex.test(value)) {
			return 'Please enter a valid email address';
		}
		return '';
	}

	function validatePassword(value: string): string {
		if (!value) {
			return 'Password is required';
		}
		if (value.length < 8) {
			return 'Password must be at least 8 characters';
		}
		return '';
	}

	function validateConfirmPassword(value: string, passwordValue: string): string {
		if (!value) {
			return 'Please confirm your password';
		}
		if (value !== passwordValue) {
			return 'Passwords do not match';
		}
		return '';
	}

	function validateForm(): boolean {
		usernameError = validateUsername(username);
		emailError = validateEmail(email);
		passwordError = validatePassword(password);
		confirmPasswordError = validateConfirmPassword(confirmPassword, password);

		return !usernameError && !emailError && !passwordError && !confirmPasswordError;
	}

	// Real-time validation on blur
	function handleUsernameBlur() {
		if (username) {
			usernameError = validateUsername(username);
		}
	}

	function handleEmailBlur() {
		if (email) {
			emailError = validateEmail(email);
		}
	}

	function handlePasswordBlur() {
		if (password) {
			passwordError = validatePassword(password);
		}
	}

	function handleConfirmPasswordBlur() {
		if (confirmPassword) {
			confirmPasswordError = validateConfirmPassword(confirmPassword, password);
		}
	}

	// Clear individual errors on input
	function handleUsernameInput() {
		if (usernameError) usernameError = '';
		if (errorMessage) errorMessage = '';
	}

	function handleEmailInput() {
		if (emailError) emailError = '';
		if (errorMessage) errorMessage = '';
	}

	function handlePasswordInput() {
		if (passwordError) passwordError = '';
		if (errorMessage) errorMessage = '';
	}

	function handleConfirmPasswordInput() {
		if (confirmPasswordError) confirmPasswordError = '';
		if (errorMessage) errorMessage = '';
	}

	async function handleRegister(event: Event) {
		event.preventDefault();

		if (!validateForm()) {
			return;
		}

		isLoading = true;
		errorMessage = '';
		successMessage = '';

		try {
			const registerData: RegisterData = {
				username,
				email,
				password
			};

			// TODO: Replace with actual API call when backend is ready
			await simulateRegister(registerData);

			// Show success message
			successMessage = 'Account created successfully! Logging you in...';

			// Auto-login after successful registration
			await new Promise((resolve) => setTimeout(resolve, 1000));

			// Redirect to original URL or lobby
			goto(returnUrl);
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
			} else {
				errorMessage = 'Registration failed. Please try again.';
			}
		} finally {
			isLoading = false;
		}
	}

	// Simulated registration function - replace with actual API call
	async function simulateRegister(data: RegisterData): Promise<void> {
		// Simulate network delay
		await new Promise((resolve) => setTimeout(resolve, 1500));

		// Simulate username already taken error (for demo)
		if (data.username.toLowerCase() === 'admin' || data.username.toLowerCase() === 'test') {
			throw new Error('Username already taken. Please choose another one.');
		}

		// Simulate email already registered error (for demo)
		if (data.email === 'taken@example.com') {
			throw new Error('Email already registered. Please use another email or login.');
		}

		// Create a mock JWT token with expiry in 24 hours
		const now = Math.floor(Date.now() / 1000);
		const exp = now + 86400; // 24 hours
		const payload = {
			sub: 'user-' + Math.random().toString(36).substr(2, 9),
			username: data.username,
			email: data.email,
			exp,
			iat: now
		};

		// Simple base64 encoding for demo (not secure, just for simulation)
		const payloadStr = JSON.stringify(payload);
		const encodedPayload = btoa(payloadStr);
		const mockToken = `mock.${encodedPayload}.signature`;

		// Auto-login after registration
		auth.login(mockToken, {
			id: payload.sub,
			username: payload.username,
			email: payload.email
		});

		console.log('Registration successful:', payload.username);
	}
</script>

<svelte:head>
	<title>Register - MAGE</title>
</svelte:head>

<div class="container">
	<div class="card">
		<h1>Create Account</h1>
		<p>Join MAGE and start playing Magic: The Gathering online</p>

		{#if errorMessage}
			<div class="error-message" role="alert" aria-live="polite">
				{errorMessage}
			</div>
		{/if}

		{#if successMessage}
			<div class="success-message" role="alert" aria-live="polite">
				{successMessage}
			</div>
		{/if}

		<form class="register-form" on:submit={handleRegister}>
			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					name="username"
					bind:value={username}
					on:input={handleUsernameInput}
					on:blur={handleUsernameBlur}
					placeholder="Choose a username (3-20 characters)"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={usernameError ? 'true' : 'false'}
					aria-describedby={usernameError ? 'username-error' : 'username-help'}
				/>
				<span class="field-help" id="username-help">Alphanumeric characters and underscores only</span>
				{#if usernameError}
					<span class="field-error" id="username-error" role="alert">{usernameError}</span>
				{/if}
			</div>

			<div class="form-group">
				<label for="email">Email</label>
				<input
					type="email"
					id="email"
					name="email"
					bind:value={email}
					on:input={handleEmailInput}
					on:blur={handleEmailBlur}
					placeholder="Enter your email address"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={emailError ? 'true' : 'false'}
					aria-describedby={emailError ? 'email-error' : undefined}
				/>
				{#if emailError}
					<span class="field-error" id="email-error" role="alert">{emailError}</span>
				{/if}
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					name="password"
					bind:value={password}
					on:input={handlePasswordInput}
					on:blur={handlePasswordBlur}
					placeholder="Create a password (min 8 characters)"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={passwordError ? 'true' : 'false'}
					aria-describedby={passwordError ? 'password-error' : 'password-help'}
				/>
				<span class="field-help" id="password-help">Use at least 8 characters</span>
				{#if passwordError}
					<span class="field-error" id="password-error" role="alert">{passwordError}</span>
				{/if}
			</div>

			<div class="form-group">
				<label for="confirm-password">Confirm Password</label>
				<input
					type="password"
					id="confirm-password"
					name="confirm-password"
					bind:value={confirmPassword}
					on:input={handleConfirmPasswordInput}
					on:blur={handleConfirmPasswordBlur}
					placeholder="Confirm your password"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={confirmPasswordError ? 'true' : 'false'}
					aria-describedby={confirmPasswordError ? 'confirm-password-error' : undefined}
				/>
				{#if confirmPasswordError}
					<span class="field-error" id="confirm-password-error" role="alert">
						{confirmPasswordError}
					</span>
				{/if}
			</div>

			<button type="submit" class="btn-primary" disabled={isLoading}>
				{#if isLoading}
					<span class="spinner" aria-hidden="true"></span>
					Creating Account...
				{:else}
					Create Account
				{/if}
			</button>
		</form>

		<div class="links">
			<a href="/login">Already have an account? Login</a>
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		padding: 2rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.card {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 400px;
	}

	h1 {
		margin: 0 0 0.5rem 0;
		font-size: 2rem;
		color: #333;
	}

	p {
		margin: 0 0 2rem 0;
		color: #666;
	}

	.error-message {
		background: #fee;
		border: 1px solid #fcc;
		color: #c33;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.success-message {
		background: #efe;
		border: 1px solid #cfc;
		color: #3c3;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.register-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	label {
		font-weight: 500;
		color: #333;
	}

	input {
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	input[aria-invalid='true'] {
		border-color: #c33;
	}

	input:disabled {
		background: #f5f5f5;
		cursor: not-allowed;
	}

	.field-help {
		color: #666;
		font-size: 0.875rem;
	}

	.field-error {
		color: #c33;
		font-size: 0.875rem;
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.btn-primary:hover:not(:disabled) {
		background: #5568d3;
	}

	.btn-primary:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.spinner {
		width: 16px;
		height: 16px;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.links {
		margin-top: 1rem;
		text-align: center;
	}

	.links a {
		color: #667eea;
		text-decoration: none;
	}

	.links a:hover {
		text-decoration: underline;
	}
</style>
