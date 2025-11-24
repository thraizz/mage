<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import Navbar from '$lib/components/Navbar.svelte';
	import { fly, fade } from 'svelte/transition';

	let isChecking = true;

	// Check authentication on mount
	onMount(() => {
		// Try to restore auth from storage
		const restored = auth.loadAuthFromStorage();

		if (!restored) {
			// Not authenticated, redirect to login
			const returnUrl = $page.url.pathname + $page.url.search;
			goto(`/login?returnUrl=${encodeURIComponent(returnUrl)}`);
			return;
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
		<p>Loading...</p>
	</div>
{:else}
	<div class="app-container" transition:fade={{ duration: 300 }}>
		<Navbar />
		<main class="main-content">
			<div class="content-wrapper" in:fly={{ y: 20, duration: 300, delay: 100 }}>
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
		gap: 1rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(255, 255, 255, 0.3);
		border-top: 4px solid white;
		border-radius: 50%;
		animation: spin 1s linear infinite;
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
		color: white;
		font-size: 1rem;
		font-weight: 500;
	}

	/* App Container */
	.app-container {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background-color: #f9fafb;
	}

	/* Main Content */
	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.content-wrapper {
		flex: 1;
		max-width: 1280px;
		width: 100%;
		margin: 0 auto;
		padding: 2rem 1rem;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.content-wrapper {
			padding: 1rem;
		}
	}
</style>
