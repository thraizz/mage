<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { getDeckDetails, deleteDeck } from '$lib/api/decks';
	import type { Deck } from '$lib/types/deck';
	import DeckViewer from '$lib/components/DeckViewer.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { toast } from '$lib/stores/toast';
	import CircleX from '@lucide/svelte/icons/circle-x';

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
			<CircleX class="error-icon" aria-hidden="true" />
			<h2>Error Loading Deck</h2>
			<p>{error}</p>
			<div class="error-actions">
				<button class="btn-retry" onclick={loadDeck}>Try Again</button>
				<button class="btn-back" onclick={handleClose}>Back to Decks</button>
			</div>
		</div>
	</div>
{:else if deck}
	<DeckViewer {deck} {loading} onclose={handleClose} ondelete={handleDeleteClick} />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
	bind:open={showDeleteConfirm}
	title="Delete Deck"
	message={deck
		? `Are you sure you want to delete "${deck.name}"? This action cannot be undone.`
		: 'Are you sure?'}
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
		background: var(--bg-void);
		padding: var(--space-8);
	}

	.error-content {
		text-align: center;
		max-width: 400px;
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-md);
		padding: var(--space-8);
	}

	/* Icon is rendered by a component, so style it globally */
	:global(svg.error-icon) {
		width: 4rem;
		height: 4rem;
		color: var(--status-error);
		margin: 0 auto var(--space-4);
	}

	.error-content h2 {
		font-family: var(--font-display);
		font-size: var(--text-2xl);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0 0 var(--space-2) 0;
	}

	.error-content p {
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin: 0 0 var(--space-6) 0;
	}

	.error-actions {
		display: flex;
		gap: var(--space-3);
		justify-content: center;
	}

	.btn-retry,
	.btn-back {
		padding: var(--space-2) var(--space-4);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-retry {
		background: var(--accent-gold);
		color: var(--bg-void);
		border: 1px solid var(--accent-gold);
	}

	.btn-retry:hover {
		background: var(--accent-gold-bright);
		box-shadow: var(--shadow-glow);
	}

	.btn-back {
		background: var(--bg-iron);
		color: var(--text-bright);
		border: 1px solid var(--border-default);
	}

	.btn-back:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.btn-retry:focus-visible,
	.btn-back:focus-visible {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}
</style>
