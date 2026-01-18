<script lang="ts">
	import Card from './Card.svelte';
	import type { CardView } from '$lib/generated/mage/v1/models';
	import type { Action } from 'svelte/action';

	interface Props {
		battlefieldNonlands: CardView[];
		battlefieldLands: CardView[];
		commandCards: CardView[];
		isCommanderGame: boolean;
		isDragging: boolean;
		isOverValidDrop: boolean;
		dropZone: string | null;
		hoveredCardId: string | null;
		onCardClick: (cardId: string) => void;
		onCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
		onCardContextMenu: (cardId: string, cardName: string) => void;
		onCommandCardMouseDown: (cardId: string, cardName: string, e: MouseEvent) => void;
		onCardHover: (cardId: string | null) => void;
		battlefieldDropZoneRef?: (el: HTMLDivElement | null) => void;
		commandDropZoneRef?: (el: HTMLDivElement | null) => void;
	}

	let {
		battlefieldNonlands,
		battlefieldLands,
		commandCards,
		isCommanderGame,
		isDragging,
		isOverValidDrop,
		dropZone,
		hoveredCardId,
		onCardClick,
		onCardMouseDown,
		onCardContextMenu,
		onCommandCardMouseDown,
		onCardHover,
		battlefieldDropZoneRef,
		commandDropZoneRef
	}: Props = $props();

	const totalCards = $derived(battlefieldNonlands.length + battlefieldLands.length);

	// Create action for battlefield drop zone
	const battlefieldDropZone: Action<HTMLDivElement> = (node) => {
		if (battlefieldDropZoneRef) battlefieldDropZoneRef(node);
		return {
			destroy() {
				if (battlefieldDropZoneRef) battlefieldDropZoneRef(null);
			}
		};
	};

	// Create action for command drop zone
	const commandDropZone: Action<HTMLDivElement> = (node) => {
		if (commandDropZoneRef) commandDropZoneRef(node);
		return {
			destroy() {
				if (commandDropZoneRef) commandDropZoneRef(null);
			}
		};
	};
</script>

<div
	class="battlefield-area my-battlefield"
	class:drag-active={isDragging}
	class:drag-valid={isDragging && isOverValidDrop && dropZone === 'battlefield'}
	use:battlefieldDropZone
>
	<div class="battlefield-content-wrapper">
		<!-- Main Battlefield (left side) -->
		<div class="battlefield-main">
			<span class="zone-label">Your Battlefield</span>
			<div class="battlefield-rows">
				<div class="battlefield-row battlefield-row--nonlands">
					{#if battlefieldNonlands.length > 0}
						{#each battlefieldNonlands as card (card.id)}
							<div
								class="battlefield-card-wrapper"
								class:is-hovered={hoveredCardId === card.id}
								role="button"
								tabindex="0"
								aria-label={card.name}
								onmousedown={(e) => onCardMouseDown(card.id, card.name, e)}
								onkeydown={(e) => {
									if (e.key === 'Enter' || e.key === ' ') {
										e.preventDefault();
										onCardClick(card.id);
									}
								}}
								onmouseenter={() => onCardHover(card.id)}
								onmouseleave={() => {
									if (hoveredCardId === card.id) onCardHover(null);
								}}
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
									onclick={() => onCardClick(card.id)}
									oncontextmenu={(e) => {
										e.preventDefault();
										onCardContextMenu(card.id, card.name);
									}}
								/>
							</div>
						{/each}
					{/if}
					{#if battlefieldNonlands.length === 0}
						<div class="empty-battlefield">No permanents</div>
					{/if}
				</div>

				<div class="battlefield-row battlefield-row--lands">
					{#if battlefieldLands.length > 0}
						{#each battlefieldLands as card (card.id)}
							<div
								class="battlefield-card-wrapper"
								class:is-hovered={hoveredCardId === card.id}
								role="button"
								tabindex="0"
								aria-label={card.name}
								onmousedown={(e) => onCardMouseDown(card.id, card.name, e)}
								onkeydown={(e) => {
									if (e.key === 'Enter' || e.key === ' ') {
										e.preventDefault();
										onCardClick(card.id);
									}
								}}
								onmouseenter={() => onCardHover(card.id)}
								onmouseleave={() => {
									if (hoveredCardId === card.id) onCardHover(null);
								}}
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
									onclick={() => onCardClick(card.id)}
									oncontextmenu={(e) => {
										e.preventDefault();
										onCardContextMenu(card.id, card.name);
									}}
								/>
							</div>
						{/each}
					{/if}
					{#if battlefieldLands.length === 0}
						<div class="empty-battlefield">No lands</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Command Zone (right side) -->
		{#if isCommanderGame}
			<div
				class="command-zone"
				class:drag-valid={isDragging && isOverValidDrop && dropZone === 'command'}
				use:commandDropZone
			>
				<span class="zone-label">Command Zone</span>
				<div class="command-cards">
					{#if commandCards.length === 0}
						<div class="command-zone-empty">
							<span class="zone-empty-text">Empty</span>
							<span class="zone-empty-hint">Drag commander here</span>
						</div>
					{/if}
					{#each commandCards as card (card.id)}
						<div
							class="command-card-wrapper"
							title={card.name}
							role="button"
							tabindex="0"
							aria-label={card.name}
							onmousedown={(e) => onCommandCardMouseDown(card.id, card.name, e)}
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
								onclick={() => {}}
							/>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>
