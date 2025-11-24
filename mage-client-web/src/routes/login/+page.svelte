<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { getMageClient } from '$lib/grpc/client';
	import { createSessionToken } from '$lib/utils/jwt';
	import Modal from '$lib/components/Modal.svelte';

	// Form state
	let username = '';
	let password = '';
	let rememberMe = false;
	let isLoading = false;
	let errorMessage = '';

	// Validation errors
	let usernameError = '';
	let passwordError = '';

	// Guest registration modal state
	let showPasswordModal = false;
	let guestUsername = '';
	let guestPassword = '';

	// Get return URL from query params
	$: returnUrl = $page.url.searchParams.get('returnUrl') || '/lobby';

	// Redirect if already authenticated
	onMount(() => {
		const restored = auth.loadAuthFromStorage();
		if (restored) {
			goto(returnUrl);
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

	async function performLogin(loginUsername: string, loginPassword: string) {
		isLoading = true;
		errorMessage = '';

		try {
			// Connect to server using real gRPC client
			const client = getMageClient();
			const response = await client.connectUser(loginUsername, loginPassword);

			if (!response.success) {
				throw new Error(response.error || 'Login failed');
			}

			// Debug: log the response to see what we're getting
			if (import.meta.env.DEV) {
				console.log('Login response:', {
					success: response.success,
					sessionId: response.sessionId,
					userId: response.userId,
					error: response.error
				});
			}

			// Check if sessionId is in the response or already set in client
			let sessionId = response.sessionId || client.getSessionId();
			
			// If still no sessionId, this is an error
			if (!sessionId || sessionId.trim() === '') {
				console.error('Login response missing sessionId:', response);
				throw new Error('Login succeeded but no session ID received from server');
			}

			// Ensure sessionId is set in client
			if (response.sessionId && response.sessionId !== client.getSessionId()) {
				client.setSessionId(response.sessionId);
			} else if (!client.getSessionId()) {
				client.setSessionId(sessionId);
			}

			// Create a session-based token from server response
			// The server returns sessionId and userId, which we use to create a token
			const token = createSessionToken(
				sessionId,
				response.userId,
				loginUsername,
				`${loginUsername}@example.com`
			);

			// Store in auth store (this will also ensure sessionId is set)
			auth.login(token, {
				id: response.userId,
				username: loginUsername,
				email: `${loginUsername}@example.com`
			});

			// Show success toast
			toast.success(`Welcome back, ${loginUsername}!`);

			// Small delay to ensure everything is set before navigation
			await new Promise(resolve => setTimeout(resolve, 50));

			// Redirect to original URL or lobby on successful login
			goto(returnUrl);
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
				toast.error(error.message);
			} else {
				errorMessage = 'Login failed. Please try again.';
				toast.error('Login failed. Please try again.');
			}
		} finally {
			isLoading = false;
		}
	}

	async function handleLogin(event: Event) {
		event.preventDefault();

		if (!validateForm()) {
			return;
		}

		await performLogin(username, password);
	}

	/**
	 * Generate a random password with at least 8 characters
	 */
	function generateRandomPassword(): string {
		const length = 12; // Generate 12 character password
		const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
		let password = '';
		for (let i = 0; i < length; i++) {
			password += charset.charAt(Math.floor(Math.random() * charset.length));
		}
		return password;
	}

	/**
	 * Generate a random guest username
	 */
	function generateGuestUsername(): string {
		return 'Guest_' + Math.random().toString(36).substring(2, 9);
	}

	async function handleGuestLogin() {
		isLoading = true;
		errorMessage = '';

		try {
			// Generate random credentials
			guestUsername = generateGuestUsername();
			guestPassword = generateRandomPassword();

			const client = getMageClient();

			// Register the guest user
			const registerResponse = await client.register(guestUsername, guestPassword);

			if (!registerResponse.success) {
				throw new Error(registerResponse.error || 'Guest registration failed');
			}

			// Show password modal
			showPasswordModal = true;
			isLoading = false;
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
				toast.error(error.message);
			} else {
				errorMessage = 'Guest registration failed. Please try again.';
				toast.error('Guest registration failed. Please try again.');
			}
			isLoading = false;
		}
	}

	async function handleGuestLoginAfterModal() {
		// Close modal and proceed with login
		showPasswordModal = false;
		isLoading = true;
		errorMessage = '';

		try {
			const client = getMageClient();
			const response = await client.connectUser(guestUsername, guestPassword);

			if (!response.success) {
				throw new Error(response.error || 'Guest login failed');
			}

			// Verify sessionId is set in client
			if (!response.sessionId) {
				throw new Error('Guest login succeeded but no session ID received');
			}

			// Ensure sessionId is in client
			const currentSessionId = client.getSessionId();
			if (!currentSessionId) {
				client.setSessionId(response.sessionId);
			}

			// Create a session-based token from server response
			const token = createSessionToken(
				response.sessionId,
				response.userId,
				guestUsername,
				'guest@example.com'
			);

			// Store in auth store
			auth.login(token, {
				id: response.userId,
				username: guestUsername,
				email: 'guest@example.com'
			});

			// Show success toast
			toast.success(`Welcome, ${guestUsername}!`);

			// Redirect to original URL or lobby on successful login
			goto(returnUrl);
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = error.message;
				toast.error(error.message);
			} else {
				errorMessage = 'Guest login failed. Please try again.';
				toast.error('Guest login failed. Please try again.');
			}
		} finally {
			isLoading = false;
		}
	}

	function copyPasswordToClipboard() {
		if (typeof navigator !== 'undefined' && navigator.clipboard) {
			navigator.clipboard.writeText(guestPassword).then(() => {
				toast.success('Password copied to clipboard!');
			}).catch(() => {
				toast.error('Failed to copy password');
			});
		}
	}

	async function handleDevLogin() {
		// Auto-fill credentials and trigger login
		username = 'thraizz';
		password = 'Test123!';
		
		// Perform login directly with dev credentials
		await performLogin(username, password);
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

		{#if import.meta.env.DEV}
			<div class="divider">
				<span>DEV MODE</span>
			</div>

			<button type="button" class="btn-dev" on:click={handleDevLogin} disabled={isLoading}>
				{#if isLoading}
					<span class="spinner" aria-hidden="true"></span>
					Logging in...
				{:else}
					Log in as thraizz
				{/if}
			</button>
		{/if}

		<div class="links">
			<a href="/register">Don't have an account? Register</a>
		</div>
	</div>
</div>

<!-- Guest Password Modal -->
<Modal
	open={showPasswordModal}
	title="Save Your Guest Account Details"
	onClose={() => {
		showPasswordModal = false;
		handleGuestLoginAfterModal();
	}}
	closeOnBackdrop={false}
>
	<div class="password-modal-content">
		<p class="password-warning">
			A guest account has been created for you. <strong>Save your password</strong> if you want to
			play again with this account.
		</p>

		<div class="credentials-box">
			<div class="credential-item">
				<label>Username:</label>
				<div class="credential-value">
					<code>{guestUsername}</code>
					<button
						type="button"
						class="copy-button"
						on:click={() => {
							if (typeof navigator !== 'undefined' && navigator.clipboard) {
								navigator.clipboard.writeText(guestUsername);
								toast.success('Username copied!');
							}
						}}
						title="Copy username"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
							<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
						</svg>
					</button>
				</div>
			</div>

			<div class="credential-item">
				<label>Password:</label>
				<div class="credential-value">
					<code class="password-display">{guestPassword}</code>
					<button
						type="button"
						class="copy-button"
						on:click={copyPasswordToClipboard}
						title="Copy password"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="16"
							height="16"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
							<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
						</svg>
					</button>
				</div>
			</div>
		</div>

		<div class="modal-actions">
			<button type="button" class="btn-primary" on:click={handleGuestLoginAfterModal}>
				Continue to Game
			</button>
		</div>
	</div>
</Modal>

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
	.btn-secondary,
	.btn-dev {
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

	.btn-dev {
		background: #ff6b6b;
		color: white;
	}

	.btn-dev:hover:not(:disabled) {
		background: #ee5a5a;
	}

	.btn-primary:disabled,
	.btn-secondary:disabled,
	.btn-dev:disabled {
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

	/* Password Modal Styles */
	:global(.password-modal-content) {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	:global(.password-warning) {
		color: #666;
		line-height: 1.6;
		margin: 0;
	}

	:global(.password-warning strong) {
		color: #c33;
	}

	:global(.credentials-box) {
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	:global(.credential-item) {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	:global(.credential-item label) {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	:global(.credential-value) {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background: white;
		padding: 0.75rem;
		border-radius: 4px;
		border: 1px solid #d1d5db;
	}

	:global(.credential-value code) {
		flex: 1;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-size: 0.9375rem;
		color: #111827;
		background: transparent;
		padding: 0;
		border: none;
		word-break: break-all;
	}

	:global(.password-display) {
		font-weight: 600;
		letter-spacing: 0.05em;
	}

	:global(.copy-button) {
		background: #f3f4f6;
		border: 1px solid #d1d5db;
		border-radius: 4px;
		padding: 0.5rem;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #6b7280;
		transition: all 0.2s;
		flex-shrink: 0;
	}

	:global(.copy-button:hover) {
		background: #e5e7eb;
		color: #374151;
	}

	:global(.modal-actions) {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		margin-top: 0.5rem;
	}
</style>
