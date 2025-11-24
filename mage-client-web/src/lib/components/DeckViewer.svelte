<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Deck, DeckCard } from '$lib/types/deck';
	import LoadingSpinner from './LoadingSpinner.svelte';

	export let deck: Deck;
	export let loading = false;

	const dispatch = createEventDispatcher();

	// Card type grouping
	interface CardGroup {
		type: string;
		cards: Array<{ name: string; quantity: number }>;
		totalCount: number;
	}

	$: cardGroups = groupCardsByType(deck.mainDeck);
	$: sideboardCards = groupCards(deck.sideboard);
	$: deckStats = calculateDeckStats(deck.mainDeck);

	function groupCardsByType(cards: DeckCard[]): CardGroup[] {
		// Group cards by type - in a real implementation, we'd need card data with types
		// For now, we'll create a simple grouped list by card name
		const groups = new Map<string, Array<{ name: string; quantity: number }>>();
		const typeCounts = new Map<string, number>();

		// For this MVP, we'll group all cards under "Cards" since we don't have type data
		// In production, this would query card database for types
		const allCards = new Map<string, number>();

		for (const card of cards) {
			const current = allCards.get(card.cardName) || 0;
			allCards.set(card.cardName, current + card.quantity);
		}

		const cardList = Array.from(allCards.entries())
			.map(([name, quantity]) => ({ name, quantity }))
			.sort((a, b) => a.name.localeCompare(b.name));

		return [
			{
				type: 'Main Deck',
				cards: cardList,
				totalCount: cards.reduce((sum, c) => sum + c.quantity, 0)
			}
		];
	}

	function groupCards(cards: DeckCard[]): Array<{ name: string; quantity: number }> {
		const grouped = new Map<string, number>();
		for (const card of cards) {
			const current = grouped.get(card.cardName) || 0;
			grouped.set(card.cardName, current + card.quantity);
		}
		return Array.from(grouped.entries())
			.map(([name, quantity]) => ({ name, quantity }))
			.sort((a, b) => a.name.localeCompare(b.name));
	}

	function calculateDeckStats(cards: DeckCard[]) {
		// Calculate mana curve and color distribution
		// For MVP, we'll use simple placeholder stats
		// In production, this would query card database for mana costs and colors
		const totalCards = cards.reduce((sum, c) => sum + c.quantity, 0);
		const uniqueCards = new Set(cards.map(c => c.cardName)).size;

		// Placeholder mana curve (0-7+)
		const manaCurve = [0, 0, 0, 0, 0, 0, 0, 0];
		// For now, distribute randomly as placeholder
		const avgCost = 3;
		for (const card of cards) {
			// Placeholder: assume average card cost is around 3
			const cost = Math.min(7, Math.floor(Math.random() * 6) + 1);
			manaCurve[cost] += card.quantity;
		}

		// Placeholder color distribution
		const colors = {
			White: 0,
			Blue: 0,
			Black: 0,
			Red: 0,
			Green: 0,
			Colorless: 0
		};

		return {
			totalCards,
			uniqueCards,
			manaCurve,
			colors
		};
	}

	function handleExport() {
		// Export deck as text
		let deckText = `# ${deck.name}\n`;
		deckText += `# Format: ${deck.format}\n`;
		deckText += `# Cards: ${deck.cardCount}\n\n`;

		// Main deck
		for (const card of deck.mainDeck) {
			if (card.quantity > 0) {
				deckText += `${card.quantity} ${card.cardName}\n`;
			}
		}

		// Sideboard
		if (deck.sideboard.length > 0) {
			deckText += `\nSideboard:\n`;
			for (const card of deck.sideboard) {
				if (card.quantity > 0) {
					deckText += `${card.quantity} ${card.cardName}\n`;
				}
			}
		}

		// Create download
		const blob = new Blob([deckText], { type: 'text/plain' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${deck.name.replace(/[^a-z0-9]/gi, '_').toLowerCase()}.txt`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	function handleDelete() {
		dispatch('delete');
	}

	function handleClose() {
		dispatch('close');
	}
</script>

<div class="deck-viewer">
	{#if loading}
		<div class="loading-container">
			<LoadingSpinner size="large" />
			<p>Loading deck details...</p>
		</div>
	{:else}
		<!-- Header -->
		<div class="viewer-header">
			<div class="header-left">
				<button class="btn-back" on:click={handleClose} title="Back to deck list">
					<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 19l-7-7m0 0l7-7m-7 7h18"
						/>
					</svg>
				</button>
				<div class="header-text">
					<h1>{deck.name}</h1>
					<div class="deck-meta">
						<span class="format-badge">{deck.format}</span>
						<span class="card-count">{deck.cardCount} cards</span>
					</div>
				</div>
			</div>
			<div class="header-actions">
				<button class="btn-export" on:click={handleExport}>
					<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
						/>
					</svg>
					Export
				</button>
				<button class="btn-delete" on:click={handleDelete}>
					<svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
						/>
					</svg>
					Delete
				</button>
			</div>
		</div>

		<div class="viewer-content">
			<!-- Stats Section -->
			<div class="stats-section">
				<div class="stat-card">
					<div class="stat-label">Total Cards</div>
					<div class="stat-value">{deckStats.totalCards}</div>
				</div>
				<div class="stat-card">
					<div class="stat-label">Unique Cards</div>
					<div class="stat-value">{deckStats.uniqueCards}</div>
				</div>
				<div class="stat-card">
					<div class="stat-label">Main Deck</div>
					<div class="stat-value">{deck.mainDeck.reduce((sum, c) => sum + c.quantity, 0)}</div>
				</div>
				{#if deck.sideboard.length > 0}
					<div class="stat-card">
						<div class="stat-label">Sideboard</div>
						<div class="stat-value">{deck.sideboard.reduce((sum, c) => sum + c.quantity, 0)}</div>
					</div>
				{/if}
			</div>

			<!-- Mana Curve -->
			<div class="visualization-section">
				<h2>Mana Curve</h2>
				<div class="mana-curve">
					{#each deckStats.manaCurve as count, cost}
						<div class="curve-bar">
							<div class="bar-fill" style="height: {count > 0 ? (count / Math.max(...deckStats.manaCurve)) * 100 : 0}%">
								<span class="bar-count">{count}</span>
							</div>
							<div class="bar-label">{cost === 7 ? '7+' : cost}</div>
						</div>
					{/each}
				</div>
				<p class="chart-note">Note: Mana costs are placeholder data in this MVP version</p>
			</div>

			<!-- Color Distribution -->
			<div class="visualization-section">
				<h2>Color Distribution</h2>
				<div class="color-distribution">
					{#each Object.entries(deckStats.colors) as [color, count]}
						<div class="color-bar">
							<div class="color-label">{color}</div>
							<div class="color-bar-track">
								<div
									class="color-bar-fill color-{color.toLowerCase()}"
									style="width: {count > 0 ? (count / deckStats.totalCards) * 100 : 0}%"
								></div>
							</div>
							<div class="color-count">{count}</div>
						</div>
					{/each}
				</div>
				<p class="chart-note">Note: Color data is placeholder in this MVP version</p>
			</div>

			<!-- Card List -->
			<div class="card-list-section">
				<h2>Card List</h2>

				{#each cardGroups as group}
					<div class="card-group">
						<h3 class="group-header">
							{group.type} ({group.totalCount})
						</h3>
						<div class="card-items">
							{#each group.cards as card}
								<div class="card-item">
									<span class="card-quantity">{card.quantity}x</span>
									<span class="card-name">{card.name}</span>
								</div>
							{/each}
						</div>
					</div>
				{/each}

				{#if deck.sideboard.length > 0}
					<div class="card-group">
						<h3 class="group-header">
							Sideboard ({deck.sideboard.reduce((sum, c) => sum + c.quantity, 0)})
						</h3>
						<div class="card-items">
							{#each sideboardCards as card}
								<div class="card-item">
									<span class="card-quantity">{card.quantity}x</span>
									<span class="card-name">{card.name}</span>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.deck-viewer {
		background-color: #f9fafb;
		min-height: 100vh;
		padding: 2rem 0;
	}

	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 0;
	}

	.loading-container p {
		margin-top: 1rem;
		color: #4b5563;
	}

	/* Header */
	.viewer-header {
		max-width: 1280px;
		margin: 0 auto 2rem;
		padding: 0 1rem;
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.header-left {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		flex: 1;
	}

	.btn-back {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		background-color: white;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-back:hover {
		background-color: #f9fafb;
	}

	.btn-back .icon {
		width: 1.25rem;
		height: 1.25rem;
	}

	.header-text h1 {
		font-size: 1.875rem;
		font-weight: 700;
		color: #111827;
		margin: 0 0 0.5rem 0;
	}

	.deck-meta {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.format-badge {
		display: inline-block;
		padding: 0.25rem 0.75rem;
		font-size: 0.875rem;
		font-weight: 500;
		background-color: #dbeafe;
		color: #1e40af;
		border-radius: 0.375rem;
	}

	.card-count {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.header-actions {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.btn-export,
	.btn-delete {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.btn-export {
		background-color: #3b82f6;
		color: white;
	}

	.btn-export:hover {
		background-color: #2563eb;
	}

	.btn-delete {
		background-color: #ef4444;
		color: white;
	}

	.btn-delete:hover {
		background-color: #dc2626;
	}

	.btn-export .icon,
	.btn-delete .icon {
		width: 1.25rem;
		height: 1.25rem;
	}

	/* Content */
	.viewer-content {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 1rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	/* Stats Section */
	.stats-section {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 1rem;
	}

	.stat-card {
		background-color: white;
		padding: 1.5rem;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
		text-align: center;
	}

	.stat-label {
		font-size: 0.875rem;
		color: #6b7280;
		margin-bottom: 0.5rem;
	}

	.stat-value {
		font-size: 2rem;
		font-weight: 700;
		color: #111827;
	}

	/* Visualizations */
	.visualization-section {
		background-color: white;
		padding: 1.5rem;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	.visualization-section h2 {
		font-size: 1.25rem;
		font-weight: 600;
		color: #111827;
		margin: 0 0 1.5rem 0;
	}

	.chart-note {
		margin-top: 1rem;
		font-size: 0.75rem;
		color: #9ca3af;
		font-style: italic;
	}

	/* Mana Curve */
	.mana-curve {
		display: flex;
		align-items: flex-end;
		justify-content: space-around;
		height: 200px;
		gap: 0.5rem;
		padding: 0 1rem;
	}

	.curve-bar {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.bar-fill {
		width: 100%;
		background: linear-gradient(to top, #3b82f6, #60a5fa);
		border-radius: 0.25rem 0.25rem 0 0;
		display: flex;
		align-items: flex-start;
		justify-content: center;
		padding-top: 0.25rem;
		min-height: 0;
		transition: height 0.3s ease;
	}

	.bar-count {
		font-size: 0.75rem;
		font-weight: 600;
		color: white;
	}

	.bar-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #4b5563;
	}

	/* Color Distribution */
	.color-distribution {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.color-bar {
		display: grid;
		grid-template-columns: 100px 1fr 60px;
		align-items: center;
		gap: 1rem;
	}

	.color-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	.color-bar-track {
		height: 24px;
		background-color: #f3f4f6;
		border-radius: 0.25rem;
		overflow: hidden;
	}

	.color-bar-fill {
		height: 100%;
		transition: width 0.3s ease;
	}

	.color-white {
		background-color: #fef3c7;
	}

	.color-blue {
		background-color: #3b82f6;
	}

	.color-black {
		background-color: #1f2937;
	}

	.color-red {
		background-color: #ef4444;
	}

	.color-green {
		background-color: #10b981;
	}

	.color-colorless {
		background-color: #9ca3af;
	}

	.color-count {
		font-size: 0.875rem;
		font-weight: 600;
		color: #111827;
		text-align: right;
	}

	/* Card List */
	.card-list-section {
		background-color: white;
		padding: 1.5rem;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	.card-list-section > h2 {
		font-size: 1.25rem;
		font-weight: 600;
		color: #111827;
		margin: 0 0 1.5rem 0;
	}

	.card-group {
		margin-bottom: 2rem;
	}

	.card-group:last-child {
		margin-bottom: 0;
	}

	.group-header {
		font-size: 1rem;
		font-weight: 600;
		color: #374151;
		margin: 0 0 1rem 0;
		padding-bottom: 0.5rem;
		border-bottom: 2px solid #e5e7eb;
	}

	.card-items {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: 0.5rem;
	}

	.card-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0.75rem;
		background-color: #f9fafb;
		border-radius: 0.25rem;
		transition: background-color 0.2s;
	}

	.card-item:hover {
		background-color: #f3f4f6;
	}

	.card-quantity {
		font-size: 0.875rem;
		font-weight: 600;
		color: #6b7280;
		min-width: 2rem;
	}

	.card-name {
		font-size: 0.875rem;
		color: #111827;
		flex: 1;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.viewer-header {
			flex-direction: column;
		}

		.header-actions {
			width: 100%;
		}

		.btn-export,
		.btn-delete {
			flex: 1;
		}

		.card-items {
			grid-template-columns: 1fr;
		}

		.mana-curve {
			padding: 0 0.5rem;
		}
	}
</style>
