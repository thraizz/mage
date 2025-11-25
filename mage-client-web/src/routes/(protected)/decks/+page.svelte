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

	const formats = [
		'Standard',
		'Modern',
		'Commander',
		'Legacy',
		'Vintage',
		'Pioneer',
		'Pauper',
		'Historic'
	];
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
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 4v16m8-8H4"
						/>
					</svg>
					Upload New Deck
				</button>
			</div>
		</div>

		<!-- Filters -->
		<div class="filters">
			<label for="format" class="filter-label"> Filter by format: </label>
			<select id="format" bind:value={selectedFormat} on:change={loadDecks} class="format-select">
				<option value="">All Formats</option>
				{#each formats as format}
					<option value={format}>{format}</option>
				{/each}
			</select>
			<button on:click={loadDecks} class="btn-refresh" title="Refresh deck list">
				<svg class="refresh-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
					/>
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
							<path
								fill-rule="evenodd"
								d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
								clip-rule="evenodd"
							/>
						</svg>
					</div>
					<div class="error-text">
						<h3>Error loading decks</h3>
						<p>{error}</p>
						<button class="btn-retry" on:click={loadDecks}> Try Again </button>
					</div>
				</div>
			</div>
		{:else if decks.length === 0}
			<div class="empty-state">
				<svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
					/>
				</svg>
				<h3>No decks found</h3>
				<p>
					{selectedFormat
						? `No ${selectedFormat} decks in your collection.`
						: 'Get started by uploading your first deck.'}
				</p>
				<div class="empty-actions">
					<button class="btn-upload-empty" on:click={handleUploadNewDeck}>
						<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 4v16m8-8H4"
							/>
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
	on:close={() => (showUploadModal = false)}
	on:success={handleUploadSuccess}
/>

<style>
	.decks-page {
		min-height: 100vh;
		background: var(--bg-void);
		padding: var(--space-8) 0;
	}

	.decks-container {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 var(--space-4);
	}

	/* Header */
	.header {
		margin-bottom: var(--space-8);
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.header-text h1 {
		font-family: var(--font-display);
		font-size: var(--text-3xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
		margin: 0;
	}

	.header-text p {
		margin: var(--space-2) 0 0 0;
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	.btn-upload {
		display: inline-flex;
		align-items: center;
		padding: var(--space-2) var(--space-4);
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--bg-void);
		background: var(--accent-gold);
		cursor: pointer;
		transition: background var(--transition-fast);
	}

	.btn-upload:hover {
		background: var(--accent-gold-bright);
	}

	.btn-upload:focus {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	.btn-upload .icon {
		width: 1.25rem;
		height: 1.25rem;
		margin-right: var(--space-2);
	}

	/* Filters */
	.filters {
		margin-bottom: var(--space-6);
		display: flex;
		align-items: center;
		gap: var(--space-4);
	}

	.filter-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	.format-select {
		display: block;
		width: 12rem;
		padding: var(--space-2) var(--space-3);
		padding-right: var(--space-10);
		font-size: var(--text-sm);
		font-family: var(--font-body);
		color: var(--text-bright);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.format-select:focus {
		outline: none;
		border-color: var(--accent-gold);
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	.format-select option {
		background: var(--bg-slate);
		color: var(--text-bright);
	}

	.btn-refresh {
		display: inline-flex;
		align-items: center;
		padding: var(--space-2) var(--space-3);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
		background: var(--bg-iron);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-refresh:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
		color: var(--text-bright);
	}

	.btn-refresh:focus {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
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
		padding: var(--space-12) 0;
	}

	.loading-state p {
		margin-top: var(--space-4);
		color: var(--text-muted);
	}

	/* Error State */
	.error-state {
		background: var(--status-error-dim);
		border-left: 4px solid var(--status-error);
		padding: var(--space-4);
		border-radius: var(--radius-md);
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
		color: var(--status-error);
	}

	.error-text {
		margin-left: var(--space-3);
	}

	.error-text h3 {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--status-error);
		margin: 0;
	}

	.error-text p {
		margin: var(--space-1) 0 0 0;
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	.btn-retry {
		margin-top: var(--space-3);
		display: inline-flex;
		align-items: center;
		padding: var(--space-2) var(--space-3);
		border: 1px solid var(--status-error);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--status-error);
		background: transparent;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-retry:hover {
		background: var(--status-error);
		color: var(--text-bright);
	}

	.btn-retry:focus {
		outline: none;
		box-shadow: 0 0 0 3px var(--status-error-dim);
	}

	/* Empty State */
	.empty-state {
		text-align: center;
		padding: var(--space-12);
		background: var(--bg-obsidian);
		border-radius: var(--radius-lg);
		border: 2px dashed var(--border-default);
	}

	.empty-icon {
		margin: 0 auto;
		width: 3rem;
		height: 3rem;
		color: var(--text-ghost);
	}

	.empty-state h3 {
		margin: var(--space-2) 0 0 0;
		font-family: var(--font-display);
		font-size: var(--text-lg);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	.empty-state p {
		margin: var(--space-1) 0 0 0;
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	.empty-actions {
		margin-top: var(--space-6);
	}

	.btn-upload-empty {
		display: inline-flex;
		align-items: center;
		padding: var(--space-2) var(--space-4);
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--bg-void);
		background: var(--accent-gold);
		cursor: pointer;
		transition: background var(--transition-fast);
	}

	.btn-upload-empty:hover {
		background: var(--accent-gold-bright);
	}

	.btn-upload-empty:focus {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	.btn-upload-empty .icon {
		width: 1.25rem;
		height: 1.25rem;
		margin-right: var(--space-2);
	}

	/* Deck Grid */
	.deck-grid {
		display: grid;
		grid-template-columns: repeat(1, 1fr);
		gap: var(--space-6);
	}

	.deck-count {
		margin-top: var(--space-6);
		text-align: center;
		font-size: var(--text-sm);
		color: var(--text-muted);
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
			gap: var(--space-4);
		}

		.filters {
			flex-wrap: wrap;
		}

		.format-select {
			width: 100%;
		}
	}
</style>
