<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getDeckDetails, deleteDeck } from '$lib/api/decks';
	import type { Deck } from '$lib/types/deck';
	import DeckViewer from '$lib/components/DeckViewer.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { toast } from '$lib/stores/toast';

	let deck: Deck | null = null;
	let loading = true;
	let error: string | null = null;
	let showDeleteConfirm = false;
	let isDeleting = false;

	$: deckId = $page.params.id;

	onMount(async () => {
		await loadDeck();
	});

	async function loadDeck() {
		if (!deckId) {
			error = 'Invalid deck ID';
			loading = false;
			return;
		}

		loading = true;
		error = null;
		try {
			deck = await getDeckDetails(deckId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load deck';
			console.error('Failed to load deck:', err);
			toast.error(error);
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		goto('/decks');
	}

	function handleDeleteClick() {
		showDeleteConfirm = true;
	}

	async function handleDeleteConfirm() {
		if (!deck) return;

		isDeleting = true;
		try {
			await deleteDeck(deck.id);
			toast.success('Deck deleted successfully');
			goto('/decks');
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to delete deck';
			console.error('Failed to delete deck:', err);
			toast.error(errorMsg);
			showDeleteConfirm = false;
		} finally {
			isDeleting = false;
		}
	}

	function handleDeleteCancel() {
		showDeleteConfirm = false;
	}
</script>

<svelte:head>
	<title>{deck ? `${deck.name} - Deck Viewer` : 'Loading Deck...'} - MAGE</title>
</svelte:head>

{#if error}
	<div class="error-page">
		<div class="error-content">
			<svg class="error-icon" fill="currentColor" viewBox="0 0 20 20">
				<path
					fill-rule="evenodd"
					d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
					clip-rule="evenodd"
				/>
			</svg>
			<h2>Error Loading Deck</h2>
			<p>{error}</p>
			<div class="error-actions">
				<button class="btn-retry" on:click={loadDeck}>Try Again</button>
				<button class="btn-back" on:click={handleClose}>Back to Decks</button>
			</div>
		</div>
	</div>
{:else if deck}
	<DeckViewer {deck} {loading} on:close={handleClose} on:delete={handleDeleteClick} />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
	bind:open={showDeleteConfirm}
	title="Delete Deck"
	message={deck ? `Are you sure you want to delete "${deck.name}"? This action cannot be undone.` : 'Are you sure?'}
	confirmText={isDeleting ? 'Deleting...' : 'Delete'}
	cancelText="Cancel"
	destructive={true}
	onConfirm={handleDeleteConfirm}
	onCancel={handleDeleteCancel}
/>

<style>
	.error-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #f9fafb;
		padding: 2rem;
	}

	.error-content {
		text-align: center;
		max-width: 400px;
	}

	.error-icon {
		width: 4rem;
		height: 4rem;
		color: #ef4444;
		margin: 0 auto 1rem;
	}

	.error-content h2 {
		font-size: 1.5rem;
		font-weight: 600;
		color: #111827;
		margin: 0 0 0.5rem 0;
	}

	.error-content p {
		font-size: 0.875rem;
		color: #6b7280;
		margin: 0 0 1.5rem 0;
	}

	.error-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: center;
	}

	.btn-retry,
	.btn-back {
		padding: 0.5rem 1rem;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-retry {
		background-color: #3b82f6;
		color: white;
		border: none;
	}

	.btn-retry:hover {
		background-color: #2563eb;
	}

	.btn-back {
		background-color: white;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.btn-back:hover {
		background-color: #f9fafb;
	}
</style>
