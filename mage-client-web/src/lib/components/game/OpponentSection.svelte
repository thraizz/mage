<script lang="ts">
	import Heart from '@lucide/svelte/icons/heart';
	import Skull from '@lucide/svelte/icons/skull';
	import Hand from '@lucide/svelte/icons/hand';
	import BookOpen from '@lucide/svelte/icons/book-open';
	import GalleryVertical from '@lucide/svelte/icons/gallery-vertical';
	import Card from './Card.svelte';
	import type { CardView } from '$lib/generated/mage/v1/models';

	interface Opponent {
		playerId: string;
		name: string;
		life: number;
		poison: number;
		handCount: number;
		libraryCount: number;
		graveyard: CardView[];
	}

	interface Props {
		opponent: Opponent;
		otherPlayers: Opponent[];
		battlefieldNonlands: CardView[];
		battlefieldLands: CardView[];
		commandCards: CardView[];
		isCommanderGame: boolean;
		showLifeMenu: boolean;
		onSelectOpponent?: (playerId: string) => void;
		onLifeChange: (delta: number, playerId: string) => void;
		onPoisonChange: (delta: number, playerId: string) => void;
		onToggleLifeMenu: () => void;
		onCardContextMenu: (cardId: string, cardName: string) => void;
	}

	let {
		opponent,
		otherPlayers,
		battlefieldNonlands,
		battlefieldLands,
		commandCards,
		isCommanderGame,
		showLifeMenu,
		onSelectOpponent,
		onLifeChange,
		onPoisonChange,
		onToggleLifeMenu,
		onCardContextMenu
	}: Props = $props();
</script>

<div class="opponent-section">
	<div class="battlefield-area opponent-battlefield">
		<!-- Opponent Info Overlay -->
		<div class="opponent-info-overlay">
			<div class="opponent-identity">
				{#if otherPlayers.length > 1 && onSelectOpponent}
					<select
						class="opponent-select"
						value={opponent.playerId}
						onchange={(e) => onSelectOpponent(e.currentTarget.value)}
					>
						{#each otherPlayers as opp}
							<option value={opp.playerId}>{opp.name}</option>
						{/each}
					</select>
				{:else}
					<span class="opponent-name-label">{opponent.name}</span>
				{/if}
			</div>

			<div class="opponent-stats-compact">
				<div class="life-group">
					<button
						class="stat-btn minus"
						onclick={() => onLifeChange(-1, opponent.playerId)}
						title="Decrease life"
					>
						−
					</button>
					<button class="stat-display life" onclick={onToggleLifeMenu} title="Life total">
						<span class="stat-icon"><Heart size={14} /></span>
						<span class="stat-value">{opponent.life}</span>
					</button>
					<button
						class="stat-btn plus"
						onclick={() => onLifeChange(1, opponent.playerId)}
						title="Increase life"
					>
						+
					</button>
				</div>

				{#if opponent.poison > 0}
					<div class="stat-display poison" title="Poison counters">
						<span class="stat-icon"><Skull size={14} /></span>
						<span class="stat-value">{opponent.poison}</span>
					</div>
				{/if}

				<div class="opponent-counts">
					<span class="opponent-count" title="Hand cards">
						<Hand size={12} />
						{opponent.handCount}
					</span>
					<span class="opponent-count" title="Library cards">
						<BookOpen size={12} />
						{opponent.libraryCount}
					</span>
					<span class="opponent-count" title="Graveyard cards">
						<GalleryVertical size={12} />
						{opponent.graveyard.length}
					</span>
				</div>

				{#if showLifeMenu}
					<div class="quick-menu opponent-menu">
						<div class="menu-section">
							<span class="menu-label">Life</span>
							<div class="menu-row">
								<button onclick={() => onLifeChange(-5, opponent.playerId)}>−5</button>
								<button onclick={() => onLifeChange(-1, opponent.playerId)}>−1</button>
								<button onclick={() => onLifeChange(1, opponent.playerId)}>+1</button>
								<button onclick={() => onLifeChange(5, opponent.playerId)}>+5</button>
							</div>
						</div>
						<div class="menu-section">
							<span class="menu-label">Poison</span>
							<div class="menu-row">
								<button onclick={() => onPoisonChange(-1, opponent.playerId)}>−1</button>
								<span class="menu-value">{opponent.poison}</span>
								<button onclick={() => onPoisonChange(1, opponent.playerId)}>+1</button>
							</div>
						</div>
						<button class="menu-close" onclick={onToggleLifeMenu}> ✕ </button>
					</div>
				{/if}
			</div>
		</div>

		<!-- Opponent Battlefield Content Wrapper -->
		<div class="battlefield-content-wrapper">
			<!-- Opponent Battlefield Main -->
			<div class="battlefield-main">
				<div class="battlefield-rows">
					{#if battlefieldNonlands.length > 0}
						<div class="battlefield-row battlefield-row--nonlands">
							{#each battlefieldNonlands as card (card.id)}
								<div
									class="battlefield-card-wrapper readonly"
									title="{card.name} (controlled by {opponent.name})"
								>
									<Card
										cardId={card.id}
										cardName={card.name}
										manaCost={card.manaCost}
										cardType={card.type}
										power={card.power}
										toughness={card.toughness}
										color={card.color}
										imageUrl=""
										isTapped={card.tapped}
										isSelected={false}
										counters={card.counters}
										size="normal"
										onclick={() => {}}
										oncontextmenu={(e) => {
											e.preventDefault();
											onCardContextMenu(card.id, card.name);
										}}
									/>
								</div>
							{/each}
						</div>
					{/if}

					{#if battlefieldLands.length > 0}
						<div class="battlefield-row battlefield-row--lands">
							{#each battlefieldLands as card (card.id)}
								<div
									class="battlefield-card-wrapper readonly"
									title="{card.name} (controlled by {opponent.name})"
								>
									<Card
										cardId={card.id}
										cardName={card.name}
										manaCost={card.manaCost}
										cardType={card.type}
										power={card.power}
										toughness={card.toughness}
										color={card.color}
										imageUrl=""
										isTapped={card.tapped}
										isSelected={false}
										counters={card.counters}
										size="normal"
										onclick={() => {}}
										oncontextmenu={(e) => {
											e.preventDefault();
											onCardContextMenu(card.id, card.name);
										}}
									/>
								</div>
							{/each}
						</div>
					{/if}

					{#if battlefieldNonlands.length === 0 && battlefieldLands.length === 0}
						<div class="empty-battlefield">No permanents</div>
					{/if}
				</div>
			</div>

			<!-- Opponent Command Zone (right side) -->
			{#if isCommanderGame}
				<div class="command-zone opponent-command-zone">
					<span class="zone-label">Command Zone</span>
					<div class="command-cards">
						{#if commandCards.length === 0}
							<div class="command-zone-empty">
								<span class="zone-empty-text">Empty</span>
							</div>
						{/if}
						{#each commandCards as card (card.id)}
							<div class="command-card-wrapper readonly" title={card.name}>
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost}
									cardType={card.type}
									power={card.power}
									toughness={card.toughness}
									color={card.color}
									imageUrl=""
									isTapped={card.tapped}
									isSelected={false}
									size="small"
									onclick={() => {}}
								/>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
