<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import Navbar from '$lib/components/Navbar.svelte';
	import { fly, fade } from 'svelte/transition';

	let isChecking = $state(true);

	// Check if we're on a game page - hide navbar for immersive experience
	const isGamePage = $derived(
		$page.url.pathname.startsWith('/game/') || $page.url.pathname.startsWith('/playtest')
	);

	// Check authentication on mount
	onMount(() => {
		// Check if already authenticated (e.g., from root layout load)
		// Only call loadAuthFromStorage if not already authenticated
		if (!$auth.isAuthenticated) {
			// Try to restore auth from storage
			const restored = auth.loadAuthFromStorage();
			if (!restored) {
				// Not authenticated, redirect to login
				const returnUrl = $page.url.pathname + $page.url.search;
				goto(`/login?returnUrl=${encodeURIComponent(returnUrl)}`);
				return;
			}
		}

		// Authenticated, show content
		isChecking = false;

		// Periodically check token validity (every 60 seconds)
		const interval = setInterval(() => {
			const isValid = auth.checkTokenValidity();
			if (!isValid) {
				// Token expired, redirect to login
				const returnUrl = $page.url.pathname + $page.url.search;
				goto(`/login?returnUrl=${encodeURIComponent(returnUrl)}`);
			}
		}, 60000);

		// Cleanup interval on unmount
		return () => clearInterval(interval);
	});
</script>

{#if isChecking}
	<div class="loading-container" transition:fade={{ duration: 200 }}>
		<div class="spinner"></div>
		<p>Gathering Mana...</p>
	</div>
{:else}
	<div class="app-container" class:game-mode={isGamePage} transition:fade={{ duration: 300 }}>
		{#if !isGamePage}
			<Navbar />
		{/if}
		<main class="main-content" class:full-screen={isGamePage}>
			<div class="content-wrapper" class:no-padding={isGamePage} in:fly={{ y: 20, duration: 300, delay: 100 }}>
				<slot />
			</div>
		</main>
	</div>
{/if}

<style>
	/* Loading State */
	.loading-container {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		gap: var(--space-4);
		background: var(--ci-blind-eternities);
	}

	.spinner {
		width: 48px;
		height: 48px;
		border: 3px solid rgba(59, 130, 246, 0.2);
		border-top: 3px solid var(--ci-jace-cloak);
		border-radius: var(--radius-full);
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	.loading-container p {
		color: var(--ci-scroll-parchment);
		font-size: var(--text-base);
		font-weight: var(--weight-medium);
		font-style: italic;
	}

	/* App Container */
	.app-container {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background: var(--bg-void);
	}

	/* Main Content */
	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.content-wrapper {
		flex: 1;
		width: 100%;
		margin: 0 auto;
	}

	/* Game Mode - Full screen immersive experience */
	.app-container.game-mode {
		background: #0a0d12;
	}

	.main-content.full-screen {
		position: fixed;
		inset: 0;
		padding: 0;
	}

	.content-wrapper.no-padding {
		padding: 0;
		height: 100%;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.content-wrapper {
			padding: 0;
		}
	}
</style>
