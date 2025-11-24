<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { fetchUserDecks } from '$lib/api/decks';
	import type { Deck } from '$lib/types/deck';
	import DeckCard from '$lib/components/DeckCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import DeckUploadModal from '$lib/components/DeckUploadModal.svelte';

	let decks: Deck[] = [];
	let loading = true;
	let error: string | null = null;
	let selectedFormat = '';
	let showUploadModal = false;

	// Load decks on mount
	onMount(async () => {
		await loadDecks();
	});

	async function loadDecks() {
		loading = true;
		error = null;
		try {
			decks = await fetchUserDecks(selectedFormat || undefined);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load decks';
			console.error('Failed to load decks:', err);
		} finally {
			loading = false;
		}
	}

	function handleDeckClick(deckId: string) {
		goto(`/decks/${deckId}`);
	}

	function handleUploadNewDeck() {
		showUploadModal = true;
	}

	function handleUploadSuccess() {
		showUploadModal = false;
		// Reload decks to show the newly uploaded deck
		loadDecks();
	}

	const formats = ['Standard', 'Modern', 'Commander', 'Legacy', 'Vintage', 'Pioneer', 'Pauper', 'Historic'];
</script>

<svelte:head>
	<title>My Decks - MAGE</title>
</svelte:head>

<div class="decks-page">
	<div class="decks-container">
		<!-- Header -->
		<div class="header">
			<div class="header-content">
				<div class="header-text">
					<h1>My Decks</h1>
					<p>Manage your deck collection</p>
				</div>
				<button class="btn-upload" on:click={handleUploadNewDeck}>
					<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
					</svg>
					Upload New Deck
				</button>
			</div>
		</div>

		<!-- Filters -->
		<div class="filters">
			<label for="format" class="filter-label">
				Filter by format:
			</label>
			<select
				id="format"
				bind:value={selectedFormat}
				on:change={loadDecks}
				class="format-select"
			>
				<option value="">All Formats</option>
				{#each formats as format}
					<option value={format}>{format}</option>
				{/each}
			</select>
			<button
				on:click={loadDecks}
				class="btn-refresh"
				title="Refresh deck list"
			>
				<svg class="refresh-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
			</button>
		</div>

		<!-- Content -->
		{#if loading}
			<div class="loading-state">
				<LoadingSpinner size="large" />
				<p>Loading your decks...</p>
			</div>
		{:else if error}
			<div class="error-state">
				<div class="error-content">
					<div class="error-icon-wrapper">
						<svg class="error-icon" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="error-text">
						<h3>Error loading decks</h3>
						<p>{error}</p>
						<button class="btn-retry" on:click={loadDecks}>
							Try Again
						</button>
					</div>
				</div>
			</div>
		{:else if decks.length === 0}
			<div class="empty-state">
				<svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
				</svg>
				<h3>No decks found</h3>
				<p>
					{selectedFormat ? `No ${selectedFormat} decks in your collection.` : 'Get started by uploading your first deck.'}
				</p>
				<div class="empty-actions">
					<button class="btn-upload-empty" on:click={handleUploadNewDeck}>
						<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
						Upload Deck
					</button>
				</div>
			</div>
		{:else}
			<!-- Deck Grid -->
			<div class="deck-grid">
				{#each decks as deck (deck.id)}
					<DeckCard {deck} on:click={() => handleDeckClick(deck.id)} />
				{/each}
			</div>

			<!-- Deck count -->
			<div class="deck-count">
				Showing {decks.length} deck{decks.length !== 1 ? 's' : ''}
				{#if selectedFormat}
					in {selectedFormat}
				{/if}
			</div>
		{/if}
	</div>
</div>

<!-- Deck Upload Modal -->
<DeckUploadModal
	open={showUploadModal}
	on:close={() => showUploadModal = false}
	on:success={handleUploadSuccess}
/>

<style>
	.decks-page {
		min-height: 100vh;
		background-color: #f9fafb;
		padding: 2rem 0;
	}

	.decks-container {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 1rem;
	}

	/* Header */
	.header {
		margin-bottom: 2rem;
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.header-text h1 {
		font-size: 1.875rem;
		font-weight: 700;
		color: #111827;
		margin: 0;
	}

	.header-text p {
		margin: 0.5rem 0 0 0;
		font-size: 0.875rem;
		color: #4b5563;
	}

	.btn-upload {
		display: inline-flex;
		align-items: center;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: white;
		background-color: #3b82f6;
		cursor: pointer;
		transition: background-color 0.2s;
		box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
	}

	.btn-upload:hover {
		background-color: #2563eb;
	}

	.btn-upload:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);
	}

	.btn-upload .icon {
		width: 1.25rem;
		height: 1.25rem;
		margin-right: 0.5rem;
	}

	/* Filters */
	.filters {
		margin-bottom: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.filter-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	.format-select {
		display: block;
		width: 12rem;
		padding: 0.5rem 0.75rem;
		padding-right: 2.5rem;
		font-size: 0.875rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		background-color: white;
		box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
		cursor: pointer;
	}

	.format-select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 1px #3b82f6;
	}

	.btn-refresh {
		display: inline-flex;
		align-items: center;
		padding: 0.5rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
		background-color: white;
		cursor: pointer;
		transition: background-color 0.2s;
		box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
	}

	.btn-refresh:hover {
		background-color: #f9fafb;
	}

	.btn-refresh:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);
	}

	.refresh-icon {
		width: 1rem;
		height: 1rem;
	}

	/* Loading State */
	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 0;
	}

	.loading-state p {
		margin-top: 1rem;
		color: #4b5563;
	}

	/* Error State */
	.error-state {
		background-color: #fef2f2;
		border-left: 4px solid #f87171;
		padding: 1rem;
		border-radius: 0.375rem;
	}

	.error-content {
		display: flex;
	}

	.error-icon-wrapper {
		flex-shrink: 0;
	}

	.error-icon {
		width: 1.25rem;
		height: 1.25rem;
		color: #f87171;
	}

	.error-text {
		margin-left: 0.75rem;
	}

	.error-text h3 {
		font-size: 0.875rem;
		font-weight: 500;
		color: #991b1b;
		margin: 0;
	}

	.error-text p {
		margin: 0.25rem 0 0 0;
		font-size: 0.875rem;
		color: #b91c1c;
	}

	.btn-retry {
		margin-top: 0.75rem;
		display: inline-flex;
		align-items: center;
		padding: 0.5rem 0.75rem;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: #b91c1c;
		background-color: #fee2e2;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-retry:hover {
		background-color: #fecaca;
	}

	.btn-retry:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.5);
	}

	/* Empty State */
	.empty-state {
		text-align: center;
		padding: 3rem;
		background-color: white;
		border-radius: 0.5rem;
		border: 2px dashed #d1d5db;
	}

	.empty-icon {
		margin: 0 auto;
		width: 3rem;
		height: 3rem;
		color: #9ca3af;
	}

	.empty-state h3 {
		margin: 0.5rem 0 0 0;
		font-size: 0.875rem;
		font-weight: 500;
		color: #111827;
	}

	.empty-state p {
		margin: 0.25rem 0 0 0;
		font-size: 0.875rem;
		color: #6b7280;
	}

	.empty-actions {
		margin-top: 1.5rem;
	}

	.btn-upload-empty {
		display: inline-flex;
		align-items: center;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: white;
		background-color: #3b82f6;
		cursor: pointer;
		transition: background-color 0.2s;
		box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
	}

	.btn-upload-empty:hover {
		background-color: #2563eb;
	}

	.btn-upload-empty:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5);
	}

	.btn-upload-empty .icon {
		width: 1.25rem;
		height: 1.25rem;
		margin-right: 0.5rem;
	}

	/* Deck Grid */
	.deck-grid {
		display: grid;
		grid-template-columns: repeat(1, 1fr);
		gap: 1.5rem;
	}

	.deck-count {
		margin-top: 1.5rem;
		text-align: center;
		font-size: 0.875rem;
		color: #4b5563;
	}

	/* Responsive */
	@media (min-width: 640px) {
		.deck-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (min-width: 1024px) {
		.deck-grid {
			grid-template-columns: repeat(3, 1fr);
		}
	}

	@media (min-width: 1280px) {
		.deck-grid {
			grid-template-columns: repeat(4, 1fr);
		}
	}

	@media (max-width: 640px) {
		.header-content {
			flex-direction: column;
			align-items: flex-start;
			gap: 1rem;
		}

		.filters {
			flex-wrap: wrap;
		}

		.format-select {
			width: 100%;
		}
	}
</style>
