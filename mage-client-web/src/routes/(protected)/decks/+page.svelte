<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchUserDecks, deleteDeck } from '$lib/api/decks';
	import type { Deck } from '$lib/types/deck';

	let decks: Deck[] = [];
	let loading = true;
	let error: string | null = null;
	let selectedFormat = '';

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

	async function handleDelete(deckId: string, deckName: string) {
		if (!confirm(`Are you sure you want to delete "${deckName}"?`)) {
			return;
		}

		try {
			await deleteDeck(deckId);
			// Remove from local list
			decks = decks.filter((d) => d.id !== deckId);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete deck');
			console.error('Failed to delete deck:', err);
		}
	}

	function formatDate(timestamp: number): string {
		const now = Date.now();
		const diffMs = now - timestamp;
		const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

		if (diffDays === 0) return 'today';
		if (diffDays === 1) return 'yesterday';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
		if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
		return `${Math.floor(diffDays / 365)} years ago`;
	}

	// Extract mana colors from deck (simplified - would need actual card data)
	// eslint-disable-next-line no-unused-vars
	function getManaColors(_deck: Deck): string[] {
		// TODO: This would need card database integration to extract actual mana colors
		// For now, return empty array
		return [];
	}

	const formats = ['Standard', 'Modern', 'Commander', 'Legacy', 'Vintage', 'Pioneer'];
</script>

<svelte:head>
	<title>My Decks - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>My Decks</h1>
		<button class="btn-primary" on:click={() => alert('Upload deck feature coming soon!')}>
			Upload New Deck
		</button>
	</header>

	<div class="filters">
		<label for="format">Filter by format:</label>
		<select id="format" bind:value={selectedFormat} on:change={loadDecks}>
			<option value="">All Formats</option>
			{#each formats as format}
				<option value={format}>{format}</option>
			{/each}
		</select>
	</div>

	{#if loading}
		<div class="loading">
			<p>Loading decks...</p>
		</div>
	{:else if error}
		<div class="error">
			<p>❌ {error}</p>
			<button class="btn-primary" on:click={loadDecks}>Retry</button>
		</div>
	{:else if decks.length === 0}
		<div class="empty-state">
			<p>📦 No decks found{selectedFormat ? ` in ${selectedFormat}` : ''}</p>
			<p class="hint">Upload your first deck to get started!</p>
			<button class="btn-primary" on:click={() => alert('Upload deck feature coming soon!')}>
				Upload Deck
			</button>
		</div>
	{:else}
		<div class="decks-grid">
			{#each decks as deck (deck.id)}
				<div class="deck-card">
					<div class="deck-header">
						<h3>{deck.name}</h3>
						<span class="format-badge">{deck.format}</span>
					</div>
					<div class="deck-stats">
						<p>{deck.cardCount} cards</p>
						<p>Last modified: {formatDate(deck.updatedAt)}</p>
					</div>
					{#if getManaColors(deck).length > 0}
						<div class="deck-colors">
							{#each getManaColors(deck) as color}
								<span class="mana-symbol {color}">{color}</span>
							{/each}
						</div>
					{/if}
					<div class="deck-actions">
						<button
							class="btn-secondary"
							on:click={() => alert(`View deck ${deck.id} - coming soon!`)}
						>
							View
						</button>
						<button
							class="btn-secondary"
							on:click={() => alert(`Edit deck ${deck.id} - coming soon!`)}
						>
							Edit
						</button>
						<button class="btn-danger" on:click={() => handleDelete(deck.id, deck.name)}>
							Delete
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h1 {
		margin: 0;
		font-size: 2.5rem;
		color: #333;
	}

	.filters {
		margin-bottom: 2rem;
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.filters label {
		font-weight: 500;
		color: #555;
	}

	.filters select {
		padding: 0.5rem 1rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		background: white;
		font-size: 1rem;
		cursor: pointer;
	}

	.loading,
	.error {
		text-align: center;
		padding: 3rem;
		background: #f9fafb;
		border-radius: 8px;
	}

	.error {
		border: 2px solid #ef4444;
	}

	.error p {
		color: #ef4444;
		font-size: 1.125rem;
		margin: 0 0 1rem 0;
	}

	.loading p {
		color: #666;
		font-size: 1.125rem;
		margin: 0;
	}

	.decks-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.deck-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		transition: box-shadow 0.2s;
	}

	.deck-card:hover {
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.deck-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
	}

	.deck-header h3 {
		margin: 0;
		font-size: 1.25rem;
		color: #333;
		flex: 1;
	}

	.format-badge {
		padding: 0.25rem 0.75rem;
		background: #667eea;
		color: white;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 500;
		white-space: nowrap;
	}

	.deck-stats {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.deck-stats p {
		margin: 0;
		font-size: 0.875rem;
		color: #666;
	}

	.deck-colors {
		display: flex;
		gap: 0.5rem;
		font-size: 1.5rem;
	}

	.deck-actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	.btn-secondary {
		padding: 0.5rem 1rem;
		background: white;
		color: #667eea;
		border: 1px solid #667eea;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		flex: 1;
	}

	.btn-secondary:hover {
		background: #667eea;
		color: white;
	}

	.btn-danger {
		padding: 0.5rem 1rem;
		background: white;
		color: #ef4444;
		border: 1px solid #ef4444;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-danger:hover {
		background: #ef4444;
		color: white;
	}

	.empty-state {
		text-align: center;
		padding: 3rem;
		background: #f9fafb;
		border: 2px dashed #ddd;
		border-radius: 8px;
	}

	.empty-state p {
		margin: 0 0 1rem 0;
		font-size: 1.125rem;
		color: #666;
	}

	.empty-state .hint {
		font-size: 0.875rem;
		color: #888;
	}
</style>
