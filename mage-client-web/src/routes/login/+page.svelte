<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	// Form state
	let username = '';
	let password = '';
	let rememberMe = false;
	let isLoading = false;
	let errorMessage = '';

	// Validation errors
	let usernameError = '';
	let passwordError = '';

	// Redirect if already authenticated
	onMount(() => {
		const restored = auth.loadAuthFromStorage();
		if (restored) {
			goto('/lobby');
		}
	});

	function validateForm(): boolean {
		usernameError = '';
		passwordError = '';
		errorMessage = '';

		let isValid = true;

		if (!username) {
			usernameError = 'Username is required';
			isValid = false;
		} else if (username.length < 3) {
			usernameError = 'Username must be at least 3 characters';
			isValid = false;
		}

		if (!password) {
			passwordError = 'Password is required';
			isValid = false;
		} else if (password.length < 6) {
			passwordError = 'Password must be at least 6 characters';
			isValid = false;
		}

		return isValid;
	}

	async function handleLogin(event: Event) {
		event.preventDefault();

		if (!validateForm()) {
			return;
		}

		isLoading = true;
		errorMessage = '';

		try {
			// TODO: Replace with actual API call when backend is ready
			// For now, simulate API call with timeout
			await simulateLogin(username, password);

			// Redirect to lobby on successful login
			goto('/lobby');
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
			} else {
				errorMessage = 'Login failed. Please try again.';
			}
		} finally {
			isLoading = false;
		}
	}

	async function handleGuestLogin() {
		isLoading = true;
		errorMessage = '';

		try {
			// TODO: Replace with actual guest login API call
			await simulateLogin('Guest', '', true);

			// Redirect to lobby on successful login
			goto('/lobby');
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
			} else {
				errorMessage = 'Guest login failed. Please try again.';
			}
		} finally {
			isLoading = false;
		}
	}

	// Simulated login function - replace with actual API call
	async function simulateLogin(
		user: string,
		pass: string,
		isGuest: boolean = false
	): Promise<void> {
		// Simulate network delay
		await new Promise((resolve) => setTimeout(resolve, 1000));

		// For demo purposes, accept any credentials
		// In production, this would be replaced with an actual API call
		if (!isGuest && (!user || !pass)) {
			throw new Error('Invalid credentials');
		}

		// Create a mock JWT token with expiry in 24 hours
		const now = Math.floor(Date.now() / 1000);
		const exp = now + 86400; // 24 hours
		const payload = {
			sub: isGuest ? 'guest-' + Math.random().toString(36).substr(2, 9) : 'user-123',
			username: isGuest ? 'Guest' : user,
			email: isGuest ? 'guest@example.com' : `${user}@example.com`,
			exp,
			iat: now
		};

		// Simple base64 encoding for demo (not secure, just for simulation)
		const payloadStr = JSON.stringify(payload);
		const encodedPayload = btoa(payloadStr);
		const mockToken = `mock.${encodedPayload}.signature`;

		// Store in auth store
		auth.login(mockToken, {
			id: payload.sub,
			username: payload.username,
			email: payload.email
		});

		console.log('Login successful:', payload.username);
	}
</script>

<svelte:head>
	<title>Login - MAGE</title>
</svelte:head>

<div class="container">
	<div class="card">
		<h1>Login</h1>
		<p>Welcome to MAGE - Magic: The Gathering Online</p>

		{#if errorMessage}
			<div class="error-message" role="alert" aria-live="polite">
				{errorMessage}
			</div>
		{/if}

		<form class="login-form" on:submit={handleLogin}>
			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					name="username"
					bind:value={username}
					placeholder="Enter your username"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={usernameError ? 'true' : 'false'}
					aria-describedby={usernameError ? 'username-error' : undefined}
				/>
				{#if usernameError}
					<span class="field-error" id="username-error" role="alert">{usernameError}</span>
				{/if}
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					name="password"
					bind:value={password}
					placeholder="Enter your password"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={passwordError ? 'true' : 'false'}
					aria-describedby={passwordError ? 'password-error' : undefined}
				/>
				{#if passwordError}
					<span class="field-error" id="password-error" role="alert">{passwordError}</span>
				{/if}
			</div>

			<div class="form-group checkbox-group">
				<label for="remember-me">
					<input
						type="checkbox"
						id="remember-me"
						name="rememberMe"
						bind:checked={rememberMe}
						disabled={isLoading}
					/>
					<span>Remember me</span>
				</label>
			</div>

			<button type="submit" class="btn-primary" disabled={isLoading}>
				{#if isLoading}
					<span class="spinner" aria-hidden="true"></span>
					Logging in...
				{:else}
					Login
				{/if}
			</button>
		</form>

		<div class="divider">
			<span>OR</span>
		</div>

		<button type="button" class="btn-secondary" on:click={handleGuestLogin} disabled={isLoading}>
			{#if isLoading}
				<span class="spinner" aria-hidden="true"></span>
				Connecting...
			{:else}
				Continue as Guest
			{/if}
		</button>

		<div class="links">
			<a href="/register">Don't have an account? Register</a>
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

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.checkbox-group {
		flex-direction: row;
		align-items: center;
	}

	.checkbox-group label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}

	.checkbox-group input[type='checkbox'] {
		width: auto;
		cursor: pointer;
	}

	label {
		font-weight: 500;
		color: #333;
	}

	input[type='text'],
	input[type='password'] {
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	input[type='text']:focus,
	input[type='password']:focus {
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

	.field-error {
		color: #c33;
		font-size: 0.875rem;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.75rem 1.5rem;
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
	}

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background: #5568d3;
	}

	.btn-secondary {
		background: #f0f0f0;
		color: #333;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #e0e0e0;
	}

	.btn-primary:disabled,
	.btn-secondary:disabled {
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

	.divider {
		margin: 1.5rem 0;
		text-align: center;
		position: relative;
	}

	.divider::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 0;
		right: 0;
		height: 1px;
		background: #ddd;
	}

	.divider span {
		background: white;
		padding: 0 1rem;
		color: #999;
		font-size: 0.875rem;
		position: relative;
		z-index: 1;
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
