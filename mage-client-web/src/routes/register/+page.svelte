<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { toast } from '$lib/stores/toast';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { getMageClient } from '$lib/grpc/client';
	import { createSessionToken } from '$lib/utils/jwt';

	// Form state
	let username = '';
	let password = '';
	let confirmPassword = '';
	let isLoading = false;
	let errorMessage = '';
	let successMessage = '';

	// Validation errors
	let usernameError = '';
	let passwordError = '';
	let confirmPasswordError = '';

	// Available background images
	const backgroundImages = ['Boros.jpg', 'Golgari.jpg', 'Gruul.jpg', 'Izzet.jpg', 'Rakdos.jpg'];
	// Randomly select a background image
	let selectedBackground = backgroundImages[Math.floor(Math.random() * backgroundImages.length)];

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
		passwordError = validatePassword(password);
		confirmPasswordError = validateConfirmPassword(confirmPassword, password);

		return !usernameError && !passwordError && !confirmPasswordError;
	}

	// Real-time validation on blur
	function handleUsernameBlur() {
		if (username) {
			usernameError = validateUsername(username);
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
			// Register user using real gRPC client
			const client = getMageClient();
			const registerResponse = await client.register(username, password);

			if (!registerResponse.success) {
				throw new Error(registerResponse.error || 'Spark ignition failed');
			}

			// Show success message
			successMessage = 'Spark Ignited! Drawing your destiny...';
			toast.success('Spark Ignited! Welcome, Planeswalker.');

			// Auto-login after successful registration
			await new Promise((resolve) => setTimeout(resolve, 500));

			const loginResponse = await client.connectUser(username, password);

			if (!loginResponse.success) {
				// Registration succeeded but login failed - redirect to login page
				toast.warning('Please enter the Blind Eternities with your new account');
				goto('/login');
				return;
			}

			// Create a session-based token from server response
			const token = createSessionToken(
				loginResponse.sessionId,
				loginResponse.userId,
				username,
				`${username}@example.com`
			);

			// Store in auth store
			auth.login(token, {
				id: loginResponse.userId,
				username: username
			});

			// Redirect to original URL or lobby
			goto(returnUrl);
		} catch (error) {
			if (error instanceof Error) {
				errorMessage = `Spell Fizzled. ${error.message}`;
				toast.error(`Spell Fizzled. ${error.message}`);
			} else {
				errorMessage = 'Spell Fizzled. Unable to ignite your spark.';
				toast.error('Spell Fizzled. Unable to ignite your spark.');
			}
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Ignite Your Spark - MAGE</title>
</svelte:head>

<div class="container" style="background-image: url('/images/{selectedBackground}')">
	<div class="card">
		<h1>Ignite Your Spark</h1>
		<p class="flavor-text">
			Welcome, aspiring Planeswalker. Forge your identity and begin your journey.
		</p>

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

		<form class="register-form" onsubmit={handleRegister}>
			<div class="form-group">
				<label for="username">Choose Your Name</label>
				<input
					type="text"
					id="username"
					name="username"
					bind:value={username}
					oninput={handleUsernameInput}
					onblur={handleUsernameBlur}
					placeholder="Enter your planeswalker name"
					disabled={isLoading}
					aria-required="true"
					aria-invalid={usernameError ? 'true' : 'false'}
					aria-describedby={usernameError ? 'username-error' : 'username-help'}
				/>
				<span class="field-help" id="username-help"
					>3-20 characters: letters, numbers, and underscores only</span
				>
				{#if usernameError}
					<span class="field-error" id="username-error" role="alert">{usernameError}</span>
				{/if}
			</div>

			<div class="form-group">
				<label for="password">Create Your Password</label>
				<input
					type="password"
					id="password"
					name="password"
					bind:value={password}
					oninput={handlePasswordInput}
					onblur={handlePasswordBlur}
					placeholder="Minimum 8 characters"
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
				<label for="confirm-password">Confirm Your Password</label>
				<input
					type="password"
					id="confirm-password"
					name="confirm-password"
					bind:value={confirmPassword}
					oninput={handleConfirmPasswordInput}
					onblur={handleConfirmPasswordBlur}
					placeholder="Re-enter your password"
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
					Igniting Spark...
				{:else}
					Ignite My Spark
				{/if}
			</button>
		</form>

		<div class="links">
			<a href="/login">Spark already ignited? <strong>Enter the Blind Eternities</strong></a>
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
		position: relative;
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
	}

	.container::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(11, 12, 16, 0.7);
		z-index: 0;
	}

	.card {
		width: 100%;
		max-width: 420px;
		position: relative;
		z-index: 1;
	}

	h1 {
		font-size: 2.25rem;
		text-align: center;
		margin: 0 0 0.75rem 0;
	}

	.flavor-text {
		margin: 0 0 2rem 0;
		text-align: center;
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

	.field-help {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.875rem;
	}

	.field-error {
		color: #ff6b6b;
		font-size: 0.875rem;
	}

	.error-message {
		background: rgba(255, 107, 107, 0.1);
		border: 1px solid rgba(255, 107, 107, 0.3);
		color: #ff6b6b;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.success-message {
		background: rgba(81, 207, 102, 0.1);
		border: 1px solid rgba(81, 207, 102, 0.3);
		color: #51cf66;
		padding: 0.75rem;
		border-radius: 4px;
		margin-bottom: 1rem;
	}
</style>
