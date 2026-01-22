<script lang="ts">
	import Card from './Card.svelte';
	import type { GameCard } from '$lib/types/game';
	import {
		dragDropStore,
		isDragging as isDraggingStore,
		draggedCardId,
		getAllValidDropZones,
		type SourceZone
	} from '$lib/utils/drag-drop';
	import Sparkles from '@lucide/svelte/icons/sparkles';

	// Props
	let {
		cards = [],
		playerName = 'Player',
		isOpponent = false,
		canDrag = false,
		compact = false,
		// eslint-disable-next-line no-unused-vars
		onCardClick = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardDragStart = (cardId: string) => {},
		// eslint-disable-next-line no-unused-vars
		onCardDragEnd = (cardId: string, dropped: boolean) => {}
	}: {
		cards?: GameCard[];
		playerName?: string;
		isOpponent?: boolean;
		canDrag?: boolean;
		compact?: boolean;
		// eslint-disable-next-line no-unused-vars
		onCardClick?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardDragStart?: (cardId: string) => void;
		// eslint-disable-next-line no-unused-vars
		onCardDragEnd?: (cardId: string, dropped: boolean) => void;
	} = $props();

	// State
	let showModal = $state(false);
	let selectedCardId = $state<string | null>(null);
	let searchQuery = $state('');
	let dragStartPosition = $state<{ x: number; y: number } | null>(null);
	let isDragPending = $state(false);
	const DRAG_THRESHOLD = 5;

	// Drag state from store
	const isDraggingGlobal = $derived($isDraggingStore);
	const draggedId = $derived($draggedCardId);

	/**
	 * Toggle exile modal
	 */
	function toggleModal(): void {
		showModal = !showModal;
		if (showModal) searchQuery = '';
	}

	/**
	 * Close modal
	 */
	function closeModal(): void {
		showModal = false;
		searchQuery = '';
	}

	/**
	 * Handle card click in modal
	 */
	function handleCardClick(_cardId: string): void {
		selectedCardId = selectedCardId === _cardId ? null : _cardId;
		onCardClick(_cardId);
	}

	/**
	 * Handle backdrop click
	 */
	function handleBackdropClick(event: MouseEvent): void {
		if (event.target === event.currentTarget) {
			closeModal();
		}
	}

	/**
	 * Handle mouse down - start drag tracking
	 */
	function handleMouseDown(cardId: string, cardName: string, event: MouseEvent): void {
		if (event.button !== 0) return; // Only left click
		if (!canDrag || isOpponent) return;

		event.preventDefault();
		event.stopPropagation();

		dragStartPosition = { x: event.clientX, y: event.clientY };
		isDragPending = true;

		const handleMouseMove = (moveEvent: MouseEvent) => {
			if (!dragStartPosition || !isDragPending) return;

			const dx = moveEvent.clientX - dragStartPosition.x;
			const dy = moveEvent.clientY - dragStartPosition.y;
			const distance = Math.sqrt(dx * dx + dy * dy);

			if (distance >= DRAG_THRESHOLD) {
				isDragPending = false;
				const validZones = getAllValidDropZones('exile' as SourceZone);
				dragDropStore.startDrag(
					cardId,
					cardName,
					'exile' as SourceZone,
					moveEvent.clientX,
					moveEvent.clientY,
					validZones
				);
				onCardDragStart(cardId);
				closeModal(); // Close modal when starting drag

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
	 * Prevent native drag events
	 */
	function handleDragStart(event: DragEvent): void {
		event.preventDefault();
	}

	/**
	 * Check if a card is being dragged
	 */
	function isCardDragging(cardId: string): boolean {
		return isDraggingGlobal && draggedId === cardId;
	}

	// Derived values
	const cardCount = $derived(cards.length);
	const isEmpty = $derived(cardCount === 0);
	const filteredCards = $derived.by(() => {
		if (!searchQuery.trim()) return cards;
		const q = searchQuery.toLowerCase();
		return cards.filter((c) => {
			const name = (c.name || '').toLowerCase();
			const type = (c.cardType || '').toLowerCase();
			return name.includes(q) || type.includes(q);
		});
	});
</script>

{#if compact}
	<!-- Compact badge mode -->
	<button
		class="exile-compact"
		class:has-cards={!isEmpty}
		class:opponent={isOpponent}
		onclick={toggleModal}
		title="{playerName}'s Exile ({cardCount} cards){isEmpty ? '' : ' - Click to view'}"
	>
		<span class="exile-icon-compact">
			<Sparkles size={14} />
		</span>
		<span class="exile-label-compact">Exile</span>
		<span class="card-count-compact" class:zero={isEmpty}>{cardCount}</span>
	</button>
{:else}
	<!-- Full card-shaped mode -->
	<div class="exile-zone" class:opponent={isOpponent}>
		<button
			class="exile-button"
			class:empty={isEmpty}
			onclick={toggleModal}
			disabled={isEmpty}
			title="{playerName}'s Exile ({cardCount} cards)"
		>
			<div class="exile-icon">
				{#if isEmpty}
					<span class="icon-empty">🌌</span>
				{:else}
					<span class="icon-filled">✨</span>
				{/if}
			</div>

			<div class="card-count-badge" class:zero={cardCount === 0}>
				{cardCount}
			</div>
		</button>

		<div class="exile-label">Exile</div>
	</div>
{/if}

<!-- Exile Modal -->
{#if showModal}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="exile-modal-backdrop" onclick={handleBackdropClick}>
		<div class="exile-modal">
			<div class="modal-header">
				<h3>{playerName}'s Exile</h3>
				<span class="card-count-text">{cardCount} card{cardCount !== 1 ? 's' : ''}</span>
				<button class="close-button" onclick={closeModal} title="Close">✕</button>
			</div>

			<div class="modal-content">
				{#if isEmpty}
					<div class="empty-state">
						<p>No cards in exile</p>
					</div>
				{:else}
					<div class="search-row">
						<input
							class="search-input"
							type="text"
							placeholder="Search exile..."
							bind:value={searchQuery}
						/>
					</div>

					<div class="card-grid">
						{#each filteredCards as card (card.id)}
							<!-- svelte-ignore a11y_no_static_element_interactions -->
							<div
								class="card-grid-item"
								class:draggable={canDrag && !isOpponent}
								class:is-dragging={isCardDragging(card.id)}
								onmousedown={(e) => handleMouseDown(card.id, card.name, e)}
								ondragstart={handleDragStart}
							>
								<Card
									cardId={card.id}
									cardName={card.name}
									manaCost={card.manaCost || ''}
									cardType={card.cardType || ''}
									power={card.power || ''}
									toughness={card.toughness || ''}
									imageUrl={card.imageUrl || ''}
									isSelected={selectedCardId === card.id}
									size="normal"
									onclick={() => handleCardClick(card.id)}
									isDragging={isCardDragging(card.id)}
								/>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<style>
	/* Exile Zone */
	.exile-zone {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
	}

	.exile-button {
		position: relative;
		width: 80px;
		height: 112px;
		border: 2px solid #7c3aed;
		border-radius: 8px;
		background: linear-gradient(135deg, #1a1f2e 0%, #2a1a3e 100%);
		cursor: pointer;
		transition:
			transform 0.2s,
			border-color 0.2s,
			box-shadow 0.2s;
		overflow: hidden;
	}

	.exile-button:not(.empty):hover {
		transform: scale(1.05);
		border-color: #a78bfa;
		box-shadow: 0 4px 12px rgba(124, 58, 237, 0.4);
	}

	.exile-button:disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	.exile-button.empty {
		border-style: dashed;
		border-color: #4c1d95;
		background: #1a1f2e;
	}

	.exile-icon {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 3rem;
	}

	.icon-empty {
		opacity: 0.3;
	}

	.icon-filled {
		opacity: 0.8;
		animation: sparkle 2s ease-in-out infinite;
	}

	@keyframes sparkle {
		0%,
		100% {
			opacity: 0.6;
			transform: scale(1);
		}
		50% {
			opacity: 1;
			transform: scale(1.1);
		}
	}

	.card-count-badge {
		position: absolute;
		top: 0.25rem;
		right: 0.25rem;
		padding: 0.25rem 0.5rem;
		background: #7c3aed;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 700;
		color: white;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
	}

	.card-count-badge.zero {
		background: #4c1d95;
		color: #9ca3af;
	}

	.exile-label {
		font-size: 0.75rem;
		color: #a78bfa;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	/* Modal Backdrop */
	.exile-modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.75);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.2s;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	/* Modal */
	.exile-modal {
		background: linear-gradient(135deg, #1a1f2e 0%, #1f1a2e 100%);
		border-radius: 12px;
		border: 2px solid #7c3aed;
		max-width: 90vw;
		max-height: 85vh;
		width: 800px;
		display: flex;
		flex-direction: column;
		box-shadow:
			0 20px 25px -5px rgba(124, 58, 237, 0.3),
			0 10px 10px -5px rgba(0, 0, 0, 0.5);
		animation: slideIn 0.3s;
	}

	@keyframes slideIn {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.modal-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1.5rem;
		border-bottom: 1px solid #4c1d95;
		background: linear-gradient(135deg, #141821 0%, #1a1428 100%);
	}

	.modal-header h3 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: #a78bfa;
	}

	.card-count-text {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.close-button {
		margin-left: auto;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: none;
		border-radius: 4px;
		font-size: 1.5rem;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
	}

	.close-button:hover {
		background: #4c1d95;
		color: white;
	}

	.modal-content {
		flex: 1;
		overflow-y: auto;
		padding: 1.5rem;
	}

	.search-row {
		margin-bottom: 1rem;
	}

	.search-input {
		width: 100%;
		padding: 0.75rem 1rem;
		background: #141821;
		border: 1px solid #2a3441;
		border-radius: 8px;
		color: #e2e8f0;
		font-size: 0.875rem;
	}

	.search-input:focus {
		outline: none;
		border-color: rgba(124, 58, 237, 0.7);
	}

	.empty-state {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 200px;
	}

	.empty-state p {
		color: #6b7280;
		font-style: italic;
		margin: 0;
	}

	.card-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 1rem;
	}

	.card-grid-item {
		display: flex;
		justify-content: center;
		user-select: none;
		-webkit-user-select: none;
		transition:
			transform 0.2s ease,
			opacity 0.2s ease;
	}

	.card-grid-item.draggable {
		cursor: grab;
	}

	.card-grid-item.draggable:active {
		cursor: grabbing;
	}

	.card-grid-item.is-dragging {
		opacity: 0.4;
		transform: scale(0.95);
	}

	/* Scrollbar Styling */
	.modal-content::-webkit-scrollbar {
		width: 8px;
	}

	.modal-content::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.modal-content::-webkit-scrollbar-thumb {
		background: #7c3aed;
		border-radius: 4px;
	}

	.modal-content::-webkit-scrollbar-thumb:hover {
		background: #a78bfa;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.exile-modal {
			width: 95vw;
			max-height: 90vh;
		}

		.modal-header {
			padding: 1rem;
		}

		.modal-content {
			padding: 1rem;
		}

		.card-grid {
			grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
			gap: 0.75rem;
		}
	}

	/* Compact Badge Mode */
	.exile-compact {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background: linear-gradient(135deg, rgba(26, 31, 46, 0.6) 0%, rgba(17, 24, 39, 0.6) 100%);
		border: 1px solid rgba(167, 139, 250, 0.3);
		border-radius: 8px;
		min-height: 36px;
		cursor: default;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		color: inherit;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
		/* Allow parent drop zone to receive drag events */
		pointer-events: auto;
	}

	.exile-compact.has-cards {
		cursor: pointer;
		background: linear-gradient(135deg, rgba(167, 139, 250, 0.15) 0%, rgba(26, 31, 46, 0.9) 100%);
		border-color: rgba(167, 139, 250, 0.5);
	}

	.exile-compact.has-cards:hover {
		background: linear-gradient(135deg, rgba(167, 139, 250, 0.25) 0%, rgba(42, 52, 65, 0.9) 100%);
		border-color: rgba(167, 139, 250, 0.7);
		box-shadow:
			0 2px 6px rgba(0, 0, 0, 0.3),
			0 1px 3px rgba(167, 139, 250, 0.2);
		transform: translateY(-1px);
	}

	.exile-icon-compact {
		font-size: 0.875rem;
		opacity: 0.7;
	}

	.exile-compact.has-cards .exile-icon-compact {
		opacity: 1;
	}

	.exile-label-compact {
		font-size: 0.6875rem;
		color: #a78bfa;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.card-count-compact {
		font-size: 0.75rem;
		font-weight: 700;
		color: #a78bfa;
		background: rgba(124, 58, 237, 0.2);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		min-width: 1.25rem;
		text-align: center;
	}

	.card-count-compact.zero {
		color: #4b5563;
		background: transparent;
	}
</style>
