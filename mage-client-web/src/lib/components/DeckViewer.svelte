<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { Deck, DeckCard } from '$lib/types/deck';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { toast } from '$lib/stores/toast';
	import ArrowLeft from '@lucide/svelte/icons/arrow-left';
	import Copy from '@lucide/svelte/icons/copy';
	import Download from '@lucide/svelte/icons/download';
	import Trash2 from '@lucide/svelte/icons/trash-2';

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
	$: commanderCards = groupCards(deck.commanders || []);
	$: deckStats = calculateDeckStats(deck.mainDeck);

	function groupCardsByType(cards: DeckCard[]): CardGroup[] {
		// Group cards by type - in a real implementation, we'd need card data with types
		// For now, we'll create a simple grouped list by card name

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
		const totalCards = cards.reduce((sum, c) => sum + c.quantity, 0);
		const uniqueCards = new Set(cards.map((c) => c.cardName)).size;

		// Calculate mana curve from actual card data
		const manaCurve = [0, 0, 0, 0, 0, 0, 0, 0];
		for (const card of cards) {
			if (card.manaCost) {
				const cmc = calculateCMC(card.manaCost);
				const index = Math.min(7, cmc);
				manaCurve[index] += card.quantity;
			} else {
				// No mana cost (likely a land), count as 0
				manaCurve[0] += card.quantity;
			}
		}

		// Calculate color distribution from actual card data
		const colors = {
			White: 0,
			Blue: 0,
			Black: 0,
			Red: 0,
			Green: 0,
			Colorless: 0
		};

		for (const card of cards) {
			if (card.colors && card.colors.length > 0) {
				for (const color of card.colors) {
					switch (color) {
						case 'W':
							colors.White += card.quantity;
							break;
						case 'U':
							colors.Blue += card.quantity;
							break;
						case 'B':
							colors.Black += card.quantity;
							break;
						case 'R':
							colors.Red += card.quantity;
							break;
						case 'G':
							colors.Green += card.quantity;
							break;
					}
				}
			} else {
				colors.Colorless += card.quantity;
			}
		}

		return {
			totalCards,
			uniqueCards,
			manaCurve,
			colors
		};
	}

	// Calculate Converted Mana Cost (CMC) from mana cost string
	function calculateCMC(manaCost: string): number {
		if (!manaCost) return 0;

		let cmc = 0;
		// Match all numbers in braces {2}, {3}, etc.
		const genericMatches = manaCost.match(/\{(\d+)\}/g);
		if (genericMatches) {
			for (const match of genericMatches) {
				const num = match.match(/\d+/);
				if (num) {
					cmc += parseInt(num[0]);
				}
			}
		}

		// Count colored mana symbols
		const coloredSymbols = manaCost.match(/\{[WUBRG]\}/g);
		if (coloredSymbols) {
			cmc += coloredSymbols.length;
		}

		// Count hybrid mana symbols (e.g., {W/U}, {2/W})
		const hybridSymbols = manaCost.match(/\{[WUBRG0-9]\/[WUBRG0-9]\}/g);
		if (hybridSymbols) {
			cmc += hybridSymbols.length;
		}

		// Count Phyrexian mana symbols (e.g., {W/P})
		const phyrexianSymbols = manaCost.match(/\{[WUBRG]\/P\}/g);
		if (phyrexianSymbols) {
			cmc += phyrexianSymbols.length;
		}

		return cmc;
	}

	function buildDeckText(): string {
		let deckText = `# ${deck.name}\n`;
		deckText += `# Format: ${deck.format}\n`;
		deckText += `# Cards: ${deck.cardCount}\n\n`;

		// Commander(s)
		if (deck.commanders && deck.commanders.length > 0) {
			deckText += `Commander:\n`;
			for (const card of deck.commanders) {
				if (card.quantity > 0) {
					deckText += `${card.quantity} ${card.cardName}\n`;
				}
			}
			deckText += `\n`;
		}

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

		return deckText;
	}

	function handleExport() {
		// Export deck as text
		const deckText = buildDeckText();

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

	async function handleCopy() {
		const deckText = buildDeckText();

		try {
			if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
				await navigator.clipboard.writeText(deckText);
				toast.success('Deck copied to clipboard!');
				return;
			}

			// Fallback for older browsers / non-secure contexts
			const textarea = document.createElement('textarea');
			textarea.value = deckText;
			textarea.style.position = 'fixed';
			textarea.style.top = '0';
			textarea.style.left = '0';
			textarea.style.opacity = '0';
			document.body.appendChild(textarea);
			textarea.focus();
			textarea.select();
			const ok = document.execCommand('copy');
			document.body.removeChild(textarea);

			if (!ok) throw new Error('execCommand(copy) returned false');
			toast.success('Deck copied to clipboard!');
		} catch (err) {
			console.error('Failed to copy deck to clipboard:', err);
			toast.error('Failed to copy deck');
		}
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
				<button class="btn-back" onclick={handleClose} title="Back to deck list">
					<ArrowLeft class="icon" size={20} aria-hidden="true" />
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
				<button class="btn-copy" onclick={handleCopy} title="Copy deck to clipboard">
					<Copy class="icon" size={20} aria-hidden="true" />
					Copy
				</button>
				<button class="btn-export" onclick={handleExport}>
					<Download class="icon" size={20} aria-hidden="true" />
					Export
				</button>
				<button class="btn-delete" onclick={handleDelete}>
					<Trash2 class="icon" size={20} aria-hidden="true" />
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
				{#if deck.commanders && deck.commanders.length > 0}
					<div class="stat-card">
						<div class="stat-label">Commander{deck.commanders.length > 1 ? 's' : ''}</div>
						<div class="stat-value">{deck.commanders.reduce((sum, c) => sum + c.quantity, 0)}</div>
					</div>
				{/if}
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
							<div
								class="bar-fill"
								style="height: {count > 0 ? (count / Math.max(...deckStats.manaCurve)) * 100 : 0}%"
							>
								<span class="bar-count">{count}</span>
							</div>
							<div class="bar-label">{cost === 7 ? '7+' : cost}</div>
						</div>
					{/each}
				</div>
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
			</div>

			<!-- Card List -->
			<div class="card-list-section">
				<h2>Card List</h2>

				{#if deck.commanders && deck.commanders.length > 0}
					<div class="card-group">
						<h3 class="group-header commander-header">
							Commander{deck.commanders.length > 1 ? 's' : ''} ({deck.commanders.reduce(
								(sum, c) => sum + c.quantity,
								0
							)})
						</h3>
						<div class="card-items">
							{#each commanderCards as card}
								<div class="card-item commander-card">
									<span class="card-quantity">{card.quantity}x</span>
									<span class="card-name">{card.name}</span>
								</div>
							{/each}
						</div>
					</div>
				{/if}

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
		background: var(--bg-void);
		min-height: 100vh;
		padding: var(--space-8) 0;
		color: var(--text-bright);
	}

	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: var(--space-16) 0;
	}

	.loading-container p {
		margin-top: var(--space-4);
		color: var(--text-muted);
	}

	/* Header */
	.viewer-header {
		max-width: 1280px;
		margin: 0 auto var(--space-8);
		padding: 0 var(--space-4);
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: var(--space-4);
		flex-wrap: wrap;
	}

	.header-left {
		display: flex;
		align-items: flex-start;
		gap: var(--space-4);
		flex: 1;
	}

	.btn-back {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-2);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		background: var(--bg-iron);
		color: var(--text-muted);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-back:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
		color: var(--text-bright);
	}

	.btn-back:focus-visible {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
		border-color: var(--accent-gold);
	}

	.header-text h1 {
		font-family: var(--font-display);
		font-size: var(--text-3xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
		margin: 0 0 var(--space-2) 0;
	}

	.deck-meta {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		flex-wrap: wrap;
	}

	.format-badge {
		display: inline-block;
		padding: var(--space-1) var(--space-3);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		background: var(--bg-iron);
		color: var(--text-bright);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-full);
	}

	.card-count {
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	.header-actions {
		display: flex;
		gap: var(--space-3);
		flex-wrap: wrap;
	}

	.btn-export,
	.btn-delete,
	.btn-copy {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-4);
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-export {
		background: var(--accent-gold);
		color: var(--bg-void);
		border-color: var(--accent-gold);
	}

	.btn-export:hover {
		background: var(--accent-gold-bright);
		box-shadow: var(--shadow-glow);
	}

	.btn-copy {
		background: var(--bg-iron);
		color: var(--text-bright);
		border-color: var(--border-default);
	}

	.btn-copy:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.btn-delete {
		background: var(--status-error);
		color: white;
		border-color: var(--status-error);
	}

	.btn-delete:hover {
		background-color: #dc2626;
	}

	.btn-export:focus-visible,
	.btn-copy:focus-visible,
	.btn-delete:focus-visible {
		outline: none;
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	:global(svg.icon) {
		display: block;
		flex-shrink: 0;
	}

	/* Content */
	.viewer-content {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 var(--space-4);
		display: flex;
		flex-direction: column;
		gap: var(--space-8);
	}

	/* Stats Section */
	.stats-section {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: var(--space-4);
	}

	.stat-card {
		background: var(--bg-obsidian);
		padding: var(--space-6);
		border-radius: var(--radius-lg);
		border: 1px solid var(--border-subtle);
		box-shadow: var(--shadow-sm);
		text-align: center;
	}

	.stat-label {
		font-size: var(--text-sm);
		color: var(--text-muted);
		margin-bottom: var(--space-2);
	}

	.stat-value {
		font-family: var(--font-display);
		font-size: var(--text-3xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
	}

	/* Visualizations */
	.visualization-section {
		background: var(--bg-obsidian);
		padding: var(--space-6);
		border-radius: var(--radius-lg);
		border: 1px solid var(--border-subtle);
		box-shadow: var(--shadow-sm);
	}

	.visualization-section h2 {
		font-family: var(--font-display);
		font-size: var(--text-2xl);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0 0 var(--space-6) 0;
	}

	.chart-note {
		margin-top: var(--space-4);
		font-size: var(--text-xs);
		color: var(--text-dim);
		font-style: italic;
	}

	/* Mana Curve */
	.mana-curve {
		display: flex;
		align-items: flex-end;
		justify-content: space-around;
		height: 200px;
		gap: var(--space-2);
		padding: 0 var(--space-4);
	}

	.curve-bar {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-2);
	}

	.bar-fill {
		width: 100%;
		background: linear-gradient(to top, var(--accent-gold-dim), var(--accent-gold-bright));
		border-radius: var(--radius-sm) var(--radius-sm) 0 0;
		display: flex;
		align-items: flex-start;
		justify-content: center;
		padding-top: var(--space-1);
		min-height: 0;
		transition: height var(--transition-slow);
	}

	.bar-count {
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		color: var(--bg-void);
	}

	.bar-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	/* Color Distribution */
	.color-distribution {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.color-bar {
		display: grid;
		grid-template-columns: 100px 1fr 60px;
		align-items: center;
		gap: var(--space-4);
	}

	.color-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	.color-bar-track {
		height: 24px;
		background: var(--bg-iron);
		border-radius: var(--radius-sm);
		border: 1px solid var(--border-subtle);
		overflow: hidden;
	}

	.color-bar-fill {
		height: 100%;
		transition: width var(--transition-slow);
	}

	.color-white {
		background: var(--mana-white);
	}

	.color-blue {
		background: var(--mana-blue);
	}

	.color-black {
		background: var(--mana-black);
	}

	.color-red {
		background: var(--mana-red);
	}

	.color-green {
		background: var(--mana-green);
	}

	.color-colorless {
		background: var(--mana-colorless);
	}

	.color-count {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		text-align: right;
	}

	/* Card List */
	.card-list-section {
		background: var(--bg-obsidian);
		padding: var(--space-6);
		border-radius: var(--radius-lg);
		border: 1px solid var(--border-subtle);
		box-shadow: var(--shadow-sm);
	}

	.card-list-section > h2 {
		font-family: var(--font-display);
		font-size: var(--text-2xl);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0 0 var(--space-6) 0;
	}

	.card-group {
		margin-bottom: var(--space-8);
	}

	.card-group:last-child {
		margin-bottom: 0;
	}

	.group-header {
		font-family: var(--font-display);
		font-size: var(--text-lg);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		margin: 0 0 var(--space-4) 0;
		padding-bottom: var(--space-2);
		border-bottom: 1px solid var(--border-subtle);
	}

	.commander-header {
		color: var(--accent-gold);
		border-bottom-color: var(--border-accent);
	}

	.card-items {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
		gap: var(--space-2);
	}

	.card-item {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-slate);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-sm);
		transition: all var(--transition-fast);
	}

	.card-item:hover {
		background: var(--bg-iron);
		border-color: var(--border-default);
	}

	.commander-card {
		background: var(--bg-slate);
		border-left: 3px solid var(--accent-gold);
	}

	.commander-card:hover {
		background: var(--bg-iron);
	}

	.card-quantity {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--text-dim);
		min-width: 2rem;
	}

	.card-name {
		font-size: var(--text-sm);
		color: var(--text-bright);
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
		.btn-delete,
		.btn-copy {
			flex: 1;
		}

		.card-items {
			grid-template-columns: 1fr;
		}

		.mana-curve {
			padding: 0 var(--space-2);
		}
	}
</style>
