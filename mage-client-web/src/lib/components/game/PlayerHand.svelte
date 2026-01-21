<script lang="ts">
	import Card from './Card.svelte';
	import type { GameCard } from '$lib/types/game';
	import type { CardView } from '$lib/generated/mage/v1/models';
	import {
		myHand,
		selectedCards,
		gameStore,
		cardsBeingPlayed,
		hasPriority
	} from '$lib/stores/game';
	import {
		dragDropStore,
		isDragging as isDraggingStore,
		draggedCardId,
		getValidDropZonesForCard
	} from '$lib/utils/drag-drop';

	// Props (optional callbacks for additional handling like target selection)
	let {
		// eslint-disable-next-line no-unused-vars
		onCardClick = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardHover = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardDragStart = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardDragEnd = (cardId: string, dropped: boolean) => {},
		size = 'normal',
		currentPhase = '',
		canDrag = true,
		showHeader = false
	}: {
		// eslint-disable-next-line no-unused-vars
		onCardClick?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardHover?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardDragStart?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardDragEnd?: (cardId: string, dropped: boolean) => void;
		size?: 'small' | 'normal' | 'large';
		currentPhase?: string;
		canDrag?: boolean;
		showHeader?: boolean;
	} = $props();

	// Drag state from store
	const isDraggingGlobal = $derived($isDraggingStore);
	const draggedId = $derived($draggedCardId);
	const playingCards = $derived($cardsBeingPlayed);
	const playerHasPriority = $derived($hasPriority);

	// State
	let multiSelectMode = $state(false);
	let dragStartPosition = $state<{ x: number; y: number } | null>(null);
	let isDragPending = $state(false);
	const DRAG_THRESHOLD = 5; // Pixels before drag starts

	/**
	 * Convert CardView from proto to GameCard for components
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

	// Always use game store for cards and selection
	const handCards = $derived.by(() => ($myHand || []).map(toGameCard));
	const selectedCardIds = $derived($selectedCards || []);

	/**
	 * Handle card click
	 */
	function handleCardClick(cardId: string, event?: MouseEvent | KeyboardEvent): void {
		// Check for multi-select (Shift key)
		if (event?.shiftKey) {
			multiSelectMode = true;
			// Toggle selection in game store
			gameStore.toggleCardSelection(cardId);
		} else {
			// Single select - clear all and select this one
			multiSelectMode = false;
			const isSelected = selectedCardIds.includes(cardId);
			gameStore.clearSelection();
			if (!isSelected) {
				gameStore.toggleCardSelection(cardId);
			}
		}

		onCardClick(cardId);
	}

	/**
	 * Handle keyboard events for accessibility
	 */
	function handleKeydown(cardId: string, event: KeyboardEvent): void {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			handleCardClick(cardId, event);
		}
	}

	/**
	 * Handle card hover
	 */
	function handleCardHover(_cardId: string): void {
		onCardHover(_cardId);
	}

	/**
	 * Check if card is selected
	 */
	function isCardSelected(cardId: string): boolean {
		return selectedCardIds.includes(cardId);
	}

	/**
	 * Check if card can be dragged
	 */
	function canDragCard(cardType: string): boolean {
		if (!canDrag || !playerHasPriority) return false;
		const validZones = getValidDropZonesForCard(cardType, currentPhase, playerHasPriority);
		return validZones.length > 0;
	}

	/**
	 * Handle mouse down - start drag tracking
	 * We use pure mouse events instead of native drag API to avoid browser's drag image
	 */
	function handleMouseDown(cardId: string, cardType: string, event: MouseEvent): void {
		if (event.button !== 0) return; // Only left click
		if (!canDragCard(cardType)) return;

		// Prevent native drag and text selection
		event.preventDefault();

		dragStartPosition = { x: event.clientX, y: event.clientY };
		isDragPending = true;

		// Store the card info for potential drag
		const card = handCards.find((c) => c.id === cardId);
		if (!card) return;

		// Add mousemove listener to detect actual drag
		const handleMouseMove = (moveEvent: MouseEvent) => {
			if (!dragStartPosition || !isDragPending) return;

			const dx = moveEvent.clientX - dragStartPosition.x;
			const dy = moveEvent.clientY - dragStartPosition.y;
			const distance = Math.sqrt(dx * dx + dy * dy);

			if (distance >= DRAG_THRESHOLD) {
				// Start the drag
				isDragPending = false;
				const validZones = getValidDropZonesForCard(cardType, currentPhase, playerHasPriority);
				dragDropStore.startDrag(
					cardId,
					card.name,
					'hand',
					moveEvent.clientX,
					moveEvent.clientY,
					validZones
				);
				onCardDragStart(cardId);

				// Remove this listener
				document.removeEventListener('mousemove', handleMouseMove);
				document.removeEventListener('mouseup', handleMouseUp);
			}
		};

		const handleMouseUp = () => {
			isDragPending = false;
			dragStartPosition = null;
			document.removeEventListener('mousemove', handleMouseMove);
			document.removeEventListener('mouseup', handleMouseUp);
		};

		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
	}

	/**
	 * Prevent native drag events on cards
	 */
	function handleDragStart(event: DragEvent): void {
		event.preventDefault();
	}

	/**
	 * Handle touch start for mobile drag
	 */
	function handleTouchStart(cardId: string, cardType: string, event: TouchEvent): void {
		if (!canDragCard(cardType)) return;
		if (event.touches.length !== 1) return;

		const touch = event.touches[0];
		dragStartPosition = { x: touch.clientX, y: touch.clientY };
		isDragPending = true;

		const card = handCards.find((c) => c.id === cardId);
		if (!card) return;

		const handleTouchMove = (moveEvent: TouchEvent) => {
			if (!dragStartPosition || !isDragPending) return;
			if (moveEvent.touches.length !== 1) return;

			const touch = moveEvent.touches[0];
			const dx = touch.clientX - dragStartPosition.x;
			const dy = touch.clientY - dragStartPosition.y;
			const distance = Math.sqrt(dx * dx + dy * dy);

			if (distance >= DRAG_THRESHOLD) {
				moveEvent.preventDefault();
				isDragPending = false;
				const validZones = getValidDropZonesForCard(cardType, currentPhase, playerHasPriority);
				dragDropStore.startDrag(
					cardId,
					card.name,
					'hand',
					touch.clientX,
					touch.clientY,
					validZones
				);
				onCardDragStart(cardId);

				document.removeEventListener('touchmove', handleTouchMove);
				document.removeEventListener('touchend', handleTouchEnd);
			}
		};

		const handleTouchEnd = () => {
			isDragPending = false;
			dragStartPosition = null;
			document.removeEventListener('touchmove', handleTouchMove);
			document.removeEventListener('touchend', handleTouchEnd);
		};

		document.addEventListener('touchmove', handleTouchMove, { passive: false });
		document.addEventListener('touchend', handleTouchEnd);
	}

	/**
	 * Check if a card is currently being dragged
	 */
	function isCardDragging(cardId: string): boolean {
		return isDraggingGlobal && draggedId === cardId;
	}

	/**
	 * Check if a card is currently being played (animating)
	 */
	function isCardBeingPlayed(cardId: string): boolean {
		return playingCards.includes(cardId);
	}

	// Derived values
	const handCount = $derived(handCards.length);
	const isEmpty = $derived(handCount === 0);
</script>

<div class="player-hand" class:no-header={!showHeader}>
	{#if showHeader}
		<div class="hand-header">
			<span class="hand-label">Your Hand</span>
			<span class="hand-count">({handCount})</span>
			{#if selectedCardIds.length > 0}
				<span class="selected-count">{selectedCardIds.length} selected</span>
			{/if}
		</div>
	{/if}

	{#if isEmpty}
		<div class="empty-state">
			<p>No cards in hand</p>
		</div>
	{:else}
		<div class="hand-cards" class:dragging-active={isDraggingGlobal}>
			{#each handCards as card (card.id)}
				{@const cardIsDragging = isCardDragging(card.id)}
				{@const cardIsPlaying = isCardBeingPlayed(card.id)}
				{@const cardCanDrag = canDragCard(card.cardType || '')}
				<div
					class="card-wrapper"
					class:draggable={cardCanDrag}
					class:is-dragging={cardIsDragging}
					class:is-playing={cardIsPlaying}
					role="button"
					tabindex="0"
					draggable="false"
					onclick={(e) => handleCardClick(card.id, e)}
					onkeydown={(e) => handleKeydown(card.id, e)}
					onmouseenter={() => handleCardHover(card.id)}
					onmousedown={(e) => handleMouseDown(card.id, card.cardType || '', e)}
					ontouchstart={(e) => handleTouchStart(card.id, card.cardType || '', e)}
					ondragstart={handleDragStart}
					aria-label={`Card: ${card.name}${cardCanDrag ? ' - Drag to play' : ''}`}
				>
					<Card
						cardId={card.id}
						cardName={card.name}
						manaCost={card.manaCost || ''}
						cardType={card.cardType || ''}
						power={card.power || ''}
						toughness={card.toughness || ''}
						imageUrl={card.imageUrl || ''}
						isTapped={card.isTapped || false}
						isSelected={isCardSelected(card.id)}
						counters={card.counters?.map((c) => ({ name: c.type, count: c.count })) || []}
						{size}
						onclick={() => {}}
						onhover={() => {}}
						isDragging={cardIsDragging}
						isBeingPlayed={cardIsPlaying}
					/>
				</div>
			{/each}
		</div>
	{/if}

	{#if multiSelectMode}
		<div class="multi-select-hint">
			<span>💡 Shift+Click to select multiple cards</span>
		</div>
	{/if}
</div>

<style>
	.player-hand {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	/* When header is hidden, reduce top padding */
	.player-hand.no-header .hand-cards {
		padding-top: 0.5rem;
	}

	.hand-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: #141821;
		border-bottom: 1px solid #2a3441;
	}

	.hand-label {
		font-size: 0.875rem;
		color: #6b7280;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.hand-count {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.selected-count {
		margin-left: auto;
		padding: 0.25rem 0.5rem;
		background: #667eea;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 600;
		color: white;
	}

	.empty-state {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.empty-state p {
		color: #4b5563;
		font-style: italic;
		margin: 0;
	}

	.hand-cards {
		flex: 1;
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		overflow-x: auto;
		overflow-y: hidden;
		align-items: center;
	}

	.card-wrapper {
		flex-shrink: 0;
		transition:
			transform 0.2s ease,
			opacity 0.2s ease;
		user-select: none;
		-webkit-user-select: none;
		-webkit-user-drag: none;
	}

	.card-wrapper.draggable {
		cursor: grab;
	}

	.card-wrapper * {
		-webkit-user-drag: none;
	}

	.card-wrapper.draggable:active {
		cursor: grabbing;
	}

	.card-wrapper.is-dragging {
		opacity: 0.4;
		transform: scale(0.95);
	}

	.card-wrapper.is-playing {
		animation: card-exit 0.3s ease-out forwards;
		pointer-events: none;
	}

	@keyframes card-exit {
		0% {
			transform: translateY(0) scale(1);
			opacity: 1;
		}
		100% {
			transform: translateY(-50px) scale(0.9);
			opacity: 0;
		}
	}

	/* When dragging is active globally, dim non-dragged cards slightly */
	.hand-cards.dragging-active .card-wrapper:not(.is-dragging) {
		opacity: 0.6;
	}

	.multi-select-hint {
		padding: 0.5rem 1rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		text-align: center;
	}

	.multi-select-hint span {
		font-size: 0.75rem;
		color: #6b7280;
	}

	/* Scrollbar Styling */
	.hand-cards::-webkit-scrollbar {
		height: 8px;
	}

	.hand-cards::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.hand-cards::-webkit-scrollbar-thumb {
		background: #3a4451;
		border-radius: 4px;
	}

	.hand-cards::-webkit-scrollbar-thumb:hover {
		background: #4a5461;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.hand-cards {
			padding: 0.5rem;
			gap: 0.25rem;
		}

		.hand-header {
			padding: 0.5rem;
		}

		.multi-select-hint {
			display: none;
		}
	}
</style>
