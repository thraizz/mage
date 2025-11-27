<script lang="ts">
	/**
	 * OpponentPanel - Collapsible panel for opponent's game state
	 * Can be expanded/collapsed to save space in 4-player games
	 */
	import type { CardView, Player } from '$lib/generated/mage/v1/models';
	import Card from './Card.svelte';
	import Graveyard from './Graveyard.svelte';
	import type { GameCard } from '$lib/types/game';

	// Props
	let {
		opponent,
		battlefieldCards = [],
		selectedCardIds = [],
		expanded = $bindable(true),
		position = 'top',
		onCardClick = () => {}
	}: {
		opponent: Player;
		battlefieldCards?: CardView[];
		selectedCardIds?: string[];
		expanded?: boolean;
		position?: 'top' | 'left' | 'right';
		onCardClick?: (cardId: string) => void;
	} = $props();

	/**
	 * Convert CardView to GameCard for Graveyard component
	 */
	function toGameCard(card: CardView): GameCard {
		return {
			id: card.id,
			name: card.name,
			manaCost: card.manaCost,
			cardType: card.type,
			power: card.power,
			toughness: card.toughness,
			imageUrl: '',
			isTapped: card.tapped,
			isSelected: false,
			ownerId: card.ownerId,
			controllerId: card.controllerId
		};
	}

	/**
	 * Toggle expanded state
	 */
	function toggle(): void {
		expanded = !expanded;
	}

	/**
	 * Format life total
	 */
	function formatLife(life: number): string {
		return life.toString();
	}
</script>

<div class="opponent-panel" class:expanded class:collapsed={!expanded} class:position-left={position === 'left'} class:position-right={position === 'right'} class:position-top={position === 'top'}>
	<!-- Header - Always visible -->
	<button class="panel-header" onclick={toggle}>
		<div class="player-info">
			<span class="player-name" class:has-priority={opponent.hasPriority}>
				{opponent.name}
				{#if opponent.hasPriority}
					<span class="priority-dot"></span>
				{/if}
			</span>
			<div class="player-stats">
				<span class="life" title="Life">
					<span class="stat-icon">❤️</span>
					{formatLife(opponent.life)}
				</span>
				{#if opponent.poison > 0}
					<span class="poison" title="Poison">
						<span class="stat-icon">☠️</span>
						{opponent.poison}
					</span>
				{/if}
				<span class="library" title="Library">
					<span class="stat-icon">📚</span>
					{opponent.libraryCount}
				</span>
				<span class="hand-count" title="Cards in hand">
					<span class="stat-icon">🃏</span>
					{opponent.handCount}
				</span>
			</div>
		</div>
		<span class="expand-icon">{expanded ? '▼' : '▶'}</span>
	</button>

	<!-- Expandable Content -->
	{#if expanded}
		<div class="panel-content">
			<!-- Battlefield -->
			{#if battlefieldCards.length > 0}
				<div class="battlefield-section">
					<div class="section-label">Battlefield ({battlefieldCards.length})</div>
					<div class="battlefield-cards">
						{#each battlefieldCards as card (card.id)}
							<div class="card-wrapper">
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost}
									cardType={card.type}
									power={card.power}
									toughness={card.toughness}
									imageUrl=""
									isTapped={card.tapped}
									isSelected={selectedCardIds.includes(card.id)}
									size="small"
									onclick={() => onCardClick(card.id)}
								/>
							</div>
						{/each}
					</div>
				</div>
			{:else}
				<div class="empty-battlefield">
					<span class="empty-text">No permanents</span>
				</div>
			{/if}

			<!-- Graveyard -->
			{#if opponent.graveyard && opponent.graveyard.length > 0}
				<div class="graveyard-section">
					<Graveyard
						cards={opponent.graveyard.map(toGameCard)}
						playerName={opponent.name}
						isOpponent={true}
						onCardClick={onCardClick}
					/>
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.opponent-panel {
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 8px;
		overflow: hidden;
		transition: all 0.3s ease;
		display: flex;
		flex-direction: column;
	}

	.opponent-panel.collapsed {
		max-height: 52px;
	}

	.opponent-panel.expanded {
		max-height: none;
	}

	/* Position variants for 4-player layout */
	.opponent-panel.position-left,
	.opponent-panel.position-right {
		min-width: 180px;
	}

	.opponent-panel.position-top {
		min-height: auto;
	}

	/* Panel Header */
	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.625rem 0.875rem;
		background: #141821;
		border: none;
		border-bottom: 1px solid #2a3441;
		cursor: pointer;
		transition: background 0.2s;
		width: 100%;
		text-align: left;
		color: inherit;
	}

	.panel-header:hover {
		background: #1a1f2e;
	}

	.player-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		flex: 1;
		min-width: 0;
	}

	.player-name {
		font-weight: 600;
		font-size: 0.9375rem;
		color: #fff;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.player-name.has-priority {
		color: #22c55e;
	}

	.priority-dot {
		width: 8px;
		height: 8px;
		background: #22c55e;
		border-radius: 50%;
		animation: pulse 1.5s infinite;
		flex-shrink: 0;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; transform: scale(1); }
		50% { opacity: 0.6; transform: scale(1.2); }
	}

	.player-stats {
		display: flex;
		gap: 0.75rem;
		font-size: 0.8125rem;
		flex-wrap: wrap;
	}

	.player-stats span {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.stat-icon {
		font-size: 0.75rem;
	}

	.life { color: #ef4444; font-weight: 700; }
	.poison { color: #a855f7; }
	.library { color: #3b82f6; }
	.hand-count { color: #fbbf24; }

	.expand-icon {
		font-size: 0.75rem;
		color: #6b7280;
		margin-left: 0.5rem;
		flex-shrink: 0;
	}

	/* Panel Content */
	.panel-content {
		padding: 0.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		overflow: hidden;
		animation: slide-down 0.2s ease;
	}

	@keyframes slide-down {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* Battlefield Section */
	.battlefield-section {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}

	.section-label {
		font-size: 0.6875rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
	}

	.battlefield-cards {
		display: flex;
		flex-wrap: wrap;
		gap: 0.25rem;
		max-height: 200px;
		overflow-y: auto;
	}

	.card-wrapper {
		flex-shrink: 0;
	}

	.empty-battlefield {
		padding: 1rem;
		text-align: center;
	}

	.empty-text {
		font-size: 0.8125rem;
		color: #4b5563;
		font-style: italic;
	}

	/* Graveyard Section */
	.graveyard-section {
		border-top: 1px solid #2a3441;
		padding-top: 0.5rem;
	}

	/* Scrollbar */
	.battlefield-cards::-webkit-scrollbar {
		width: 4px;
	}

	.battlefield-cards::-webkit-scrollbar-track {
		background: #0f1419;
	}

	.battlefield-cards::-webkit-scrollbar-thumb {
		background: #2a3441;
		border-radius: 2px;
	}

	/* Responsive - Horizontal layout for side panels */
	@media (min-width: 1200px) {
		.opponent-panel.position-left .panel-content,
		.opponent-panel.position-right .panel-content {
			max-height: 400px;
			overflow-y: auto;
		}
	}
</style>

