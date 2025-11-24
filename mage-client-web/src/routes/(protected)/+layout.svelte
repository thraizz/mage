<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';

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
	<div class="loading-container">
		<div class="spinner"></div>
		<p>Loading...</p>
	</div>
{:else}
	<slot />
{/if}

<style>
	.loading-container {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		gap: 1rem;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #f3f3f3;
		border-top: 4px solid #667eea;
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

	p {
		color: #666;
		font-size: 1rem;
	}
</style>
