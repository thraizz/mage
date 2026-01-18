<script lang="ts">
	import { sendPlayerUUID } from '$lib/api/game';
	import { moveCard } from '$lib/api/direct-actions';
	import type { LibrarySearchView } from '$lib/generated/mage/v1/models';
	import ManaCost from '$lib/components/mtg/ManaCost.svelte';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';

	// Props
	let {
		gameId,
		librarySearchData,
		onComplete = () => {},
		onCancel = () => {}
	}: {
		gameId: string;
		librarySearchData: LibrarySearchView;
		onComplete?: () => void;
		onCancel?: () => void;
	} = $props();

	// State
	let selectedCardId = $state<string | null>(null);
	let searchQuery = $state('');
	let filterType = $state<string>('all');
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	// Drag state
	let draggedCard = $state<{ id: string; name: string } | null>(null);
	let dragPosition = $state({ x: 0, y: 0 });
	let isDragging = $state(false);
	let hoveredDropZone = $state<string | null>(null);

	// Selected destination for click-to-send
	let selectedDestination = $state<'hand' | 'battlefield' | 'graveyard' | 'exile'>('hand');

	// Derived
	const message = $derived(librarySearchData.message || 'Search your library for a card');
	const canCancel = $derived(librarySearchData.canCancel ?? true);
	const cards = $derived(librarySearchData.cards || []);

	// Initialize selected destination from server if provided
	$effect(() => {
		const serverDest = librarySearchData.destination;
		if (serverDest && ['hand', 'battlefield', 'graveyard', 'exile'].includes(serverDest)) {
			selectedDestination = serverDest as 'hand' | 'battlefield' | 'graveyard' | 'exile';
		}
	});

	// Card type filters
	const cardTypes = [
		'all',
		'creature',
		'instant',
		'sorcery',
		'artifact',
		'enchantment',
		'land',
		'planeswalker'
	];

	// Filtered and sorted cards
	const filteredCards = $derived(() => {
		let result = [...cards];

		// Filter by search query
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			result = result.filter(
				(card) =>
					card.name?.toLowerCase().includes(query) ||
					card.type?.toLowerCase().includes(query) ||
					card.rulesText?.toLowerCase().includes(query)
			);
		}

		// Filter by card type
		if (filterType !== 'all') {
			result = result.filter((card) => card.type?.toLowerCase().includes(filterType));
		}

		// Sort by name
		result.sort((a, b) => (a.name || '').localeCompare(b.name || ''));

		return result;
	});

	/**
	 * Handle confirm selection - moves card to selected destination
	 */
	async function handleConfirm() {
		if (isLoading || !selectedCardId) return;
		await moveCardToZone(selectedCardId, selectedDestination);
	}

	/**
	 * Handle cancel
	 */
	async function handleCancel() {
		if (isLoading || !canCancel) return;

		isLoading = true;
		error = null;

		try {
			await sendPlayerUUID(gameId, 'CANCEL');
			onCancel();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to cancel';
		} finally {
			isLoading = false;
		}
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && canCancel) {
			handleCancel();
		} else if (event.key === 'Enter' && selectedCardId) {
			handleConfirm();
		}
	}

	/**
	 * Get selected card
	 */
	const selectedCard = $derived(selectedCardId ? cards.find((c) => c.id === selectedCardId) : null);

	// Drop zones configuration
	const dropZones = [
		{ id: 'hand', label: '🖐️ Hand', color: '#3b82f6' },
		{ id: 'battlefield', label: '⚔️ Battlefield', color: '#22c55e' },
		{ id: 'graveyard', label: '💀 Graveyard', color: '#6b7280' },
		{ id: 'exile', label: '✨ Exile', color: '#a855f7' }
	] as const;

	/**
	 * Handle drag start on a card
	 */
	function handleDragStart(cardId: string, cardName: string, event: MouseEvent) {
		event.preventDefault();
		event.stopPropagation();

		draggedCard = { id: cardId, name: cardName };
		dragPosition = { x: event.clientX, y: event.clientY };
		isDragging = true;

		document.addEventListener('mousemove', handleDragMove);
		document.addEventListener('mouseup', handleDragEnd);
		document.body.style.userSelect = 'none';
	}

	/**
	 * Handle drag move
	 */
	function handleDragMove(event: MouseEvent) {
		if (!isDragging) return;
		dragPosition = { x: event.clientX, y: event.clientY };
	}

	/**
	 * Handle drag end
	 */
	async function handleDragEnd() {
		document.removeEventListener('mousemove', handleDragMove);
		document.removeEventListener('mouseup', handleDragEnd);
		document.body.style.userSelect = '';

		if (draggedCard && hoveredDropZone) {
			await moveCardToZone(draggedCard.id, hoveredDropZone);
		}

		draggedCard = null;
		isDragging = false;
		hoveredDropZone = null;
	}

	/**
	 * Handle drop zone hover
	 */
	function handleDropZoneEnter(zoneId: string) {
		if (isDragging) {
			hoveredDropZone = zoneId;
		}
	}

	/**
	 * Handle drop zone leave
	 */
	function handleDropZoneLeave() {
		hoveredDropZone = null;
	}

	/**
	 * Move card to a specific zone using direct actions
	 */
	async function moveCardToZone(cardId: string, zoneId: string) {
		if (isLoading) return;

		isLoading = true;
		error = null;

		try {
			// Map zone id to API format
			const zoneMap: Record<string, string> = {
				hand: 'HAND',
				battlefield: 'BATTLEFIELD',
				graveyard: 'GRAVEYARD',
				exile: 'EXILE'
			};

			await moveCard(gameId, cardId, zoneMap[zoneId] || 'HAND');
			onComplete();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to move card';
		} finally {
			isLoading = false;
		}
	}

	/**
	 * Handle quick-send to selected destination
	 */
	async function handleQuickSend(cardId: string) {
		await moveCardToZone(cardId, selectedDestination);
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="library-search-overlay" role="dialog" aria-labelledby="search-title" aria-modal="true">
	<div class="library-search-modal">
		<div class="modal-header">
			<h2 id="search-title">📚 Your Library</h2>
			<span class="card-count">{cards.length} cards</span>
		</div>

		<div class="modal-content">
			<p class="message">{message}</p>
			<p class="destination-info">
				<span class="hint">💡 Drag cards to drop zones or double-click to send to:</span>
				<select class="destination-select" bind:value={selectedDestination}>
					<option value="hand">🖐️ Hand</option>
					<option value="battlefield">⚔️ Battlefield</option>
					<option value="graveyard">💀 Graveyard</option>
					<option value="exile">✨ Exile</option>
				</select>
			</p>

			{#if error}
				<div class="error-message" role="alert">
					<span class="error-icon">⚠️</span>
					{error}
				</div>
			{/if}

			<!-- Search and filter controls -->
			<div class="search-controls">
				<input
					type="text"
					class="search-input"
					placeholder="Search by name, type, or text..."
					bind:value={searchQuery}
				/>
				<select class="type-filter" bind:value={filterType}>
					{#each cardTypes as type}
						<option value={type}
							>{type === 'all' ? 'All Types' : type.charAt(0).toUpperCase() + type.slice(1)}</option
						>
					{/each}
				</select>
			</div>

			<!-- Main content area with cards and drop zones -->
			<div class="content-area">
				<!-- Card list -->
				<div class="card-list">
					{#if filteredCards().length === 0}
						<div class="no-results">
							{#if cards.length === 0}
								Your library is empty
							{:else}
								No cards match your search
							{/if}
						</div>
					{:else}
						{#each filteredCards() as card (card.id)}
							<!-- svelte-ignore a11y_no_static_element_interactions -->
							<div
								class="card-item"
								class:selected={selectedCardId === card.id}
								class:dragging={draggedCard?.id === card.id}
								onmousedown={(e) => handleDragStart(card.id, card.name, e)}
								ondblclick={() => handleQuickSend(card.id)}
								role="button"
								tabindex="0"
							>
								<div class="card-thumbnail">
									{#if getScryfallImageUrl(card.name, 'small')}
										<img
											src={getScryfallImageUrl(card.name, 'small')}
											alt={card.name}
											class="card-image"
											draggable="false"
										/>
									{:else}
										<div class="card-placeholder">🃏</div>
									{/if}
								</div>
								<div class="card-info">
									<div class="card-header">
										<span class="card-name">{card.name}</span>
										{#if card.manaCost}
											<ManaCost cost={card.manaCost} size="sm" />
										{/if}
									</div>
									<span class="card-type">{card.type}</span>
									{#if card.rulesText}
										<p class="card-rules">{card.rulesText}</p>
									{/if}
									{#if card.power && card.toughness}
										<span class="card-pt">{card.power}/{card.toughness}</span>
									{/if}
								</div>
								<div class="drag-hint">⋮⋮</div>
							</div>
						{/each}
					{/if}
				</div>

				<!-- Drop zones sidebar -->
				<div class="drop-zones-sidebar" class:active={isDragging}>
					<div class="drop-zones-label">Drop here:</div>
					{#each dropZones as zone}
						<!-- svelte-ignore a11y_no_static_element_interactions -->
						<div
							class="drop-zone"
							class:hovered={hoveredDropZone === zone.id}
							style="--zone-color: {zone.color}"
							onmouseenter={() => handleDropZoneEnter(zone.id)}
							onmouseleave={handleDropZoneLeave}
						>
							<span class="zone-label">{zone.label}</span>
							{#if hoveredDropZone === zone.id && draggedCard}
								<span class="drop-hint-text">Release to send here</span>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<!-- Selected card preview -->
			{#if selectedCard}
				<div class="selected-preview">
					<div class="preview-header">Selected:</div>
					<div class="preview-card">
						<strong>{selectedCard.name}</strong>
						{#if selectedCard.manaCost}
							<ManaCost cost={selectedCard.manaCost} size="sm" />
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<div class="modal-footer">
			{#if canCancel}
				<button class="btn-cancel" onclick={handleCancel} disabled={isLoading}>
					Cancel (Close)
				</button>
			{/if}
			<button class="btn-confirm" onclick={handleConfirm} disabled={isLoading || !selectedCardId}>
				{isLoading
					? 'Moving...'
					: selectedCard
						? `Send ${selectedCard.name} to ${selectedDestination}`
						: 'Select a Card'}
			</button>
		</div>
	</div>

	<!-- Drag ghost -->
	{#if isDragging && draggedCard}
		{@const dragImageUrl = getScryfallImageUrl(draggedCard.name, 'small')}
		<div class="drag-ghost" style="left: {dragPosition.x}px; top: {dragPosition.y}px;">
			<div class="drag-ghost-card" class:over-zone={hoveredDropZone !== null}>
				{#if dragImageUrl}
					<img src={dragImageUrl} alt={draggedCard.name} class="drag-image" draggable="false" />
				{:else}
					<span class="drag-name">{draggedCard.name}</span>
				{/if}
			</div>
			{#if hoveredDropZone}
				<div class="drag-destination">
					→ {dropZones.find((z) => z.id === hoveredDropZone)?.label}
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.library-search-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.9);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.library-search-modal {
		background: #1a1f2e;
		border: 2px solid #667eea;
		border-radius: 12px;
		width: 95%;
		max-width: 900px;
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		box-shadow:
			0 25px 50px -12px rgba(0, 0, 0, 0.5),
			0 0 40px rgba(102, 126, 234, 0.15);
		animation: slideUp 0.3s ease-out;
	}

	@keyframes slideUp {
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
		padding: 1rem 1.5rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 10px 10px 0 0;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
		color: white;
		letter-spacing: 0.5px;
	}

	.card-count {
		background: rgba(255, 255, 255, 0.2);
		padding: 0.25rem 0.75rem;
		border-radius: 1rem;
		font-size: 0.875rem;
		color: white;
	}

	.modal-content {
		padding: 1rem 1.5rem;
		overflow-y: auto;
		flex: 1;
	}

	.message {
		margin: 0 0 0.5rem;
		font-size: 1rem;
		color: #e2e8f0;
		text-align: center;
	}

	.destination-info {
		margin: 0 0 1rem;
		font-size: 0.875rem;
		color: #9ca3af;
		text-align: center;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.hint {
		color: #6b7280;
	}

	.destination-select {
		padding: 0.5rem 1rem;
		background: #141821;
		border: 1px solid #3b82f6;
		border-radius: 6px;
		color: #e2e8f0;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.destination-select:focus {
		outline: none;
		border-color: #667eea;
	}

	.error-message {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: rgba(239, 68, 68, 0.15);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 8px;
		color: #f87171;
		margin-bottom: 1rem;
	}

	.error-icon {
		font-size: 1rem;
	}

	/* Search controls */
	.search-controls {
		display: flex;
		gap: 0.75rem;
		margin-bottom: 1rem;
	}

	.search-input {
		flex: 1;
		padding: 0.75rem 1rem;
		background: #141821;
		border: 1px solid #2a3441;
		border-radius: 8px;
		color: #e2e8f0;
		font-size: 0.875rem;
	}

	.search-input:focus {
		outline: none;
		border-color: #667eea;
	}

	.search-input::placeholder {
		color: #6b7280;
	}

	.type-filter {
		padding: 0.75rem 1rem;
		background: #141821;
		border: 1px solid #2a3441;
		border-radius: 8px;
		color: #e2e8f0;
		font-size: 0.875rem;
		cursor: pointer;
	}

	.type-filter:focus {
		outline: none;
		border-color: #667eea;
	}

	/* Content area with cards and drop zones */
	.content-area {
		display: flex;
		gap: 1rem;
	}

	/* Card list */
	.card-list {
		flex: 1;
		max-height: 400px;
		overflow-y: auto;
		border: 1px solid #2a3441;
		border-radius: 8px;
		background: #141821;
	}

	.no-results {
		padding: 2rem;
		text-align: center;
		color: #6b7280;
	}

	.card-item {
		width: 100%;
		padding: 0.5rem 0.75rem;
		background: transparent;
		border: none;
		border-bottom: 1px solid #2a3441;
		cursor: grab;
		text-align: left;
		transition: all 0.15s ease;
		display: flex;
		align-items: center;
		gap: 0.75rem;
		user-select: none;
	}

	.card-item:last-child {
		border-bottom: none;
	}

	.card-item:hover:not(.dragging) {
		background: rgba(102, 126, 234, 0.1);
	}

	.card-item:active {
		cursor: grabbing;
	}

	.card-item.selected {
		background: rgba(102, 126, 234, 0.2);
		border-left: 3px solid #667eea;
	}

	.card-item.dragging {
		opacity: 0.4;
		background: rgba(102, 126, 234, 0.1);
	}

	.card-thumbnail {
		width: 50px;
		height: 70px;
		border-radius: 4px;
		overflow: hidden;
		flex-shrink: 0;
		background: #0d1117;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.card-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.card-placeholder {
		font-size: 1.5rem;
		opacity: 0.5;
	}

	.card-info {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
		flex: 1;
		min-width: 0;
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
	}

	.card-name {
		font-weight: 600;
		color: #e2e8f0;
		font-size: 0.9rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.card-type {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.card-rules {
		margin: 0;
		font-size: 0.7rem;
		color: #6b7280;
		line-height: 1.3;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.card-pt {
		font-size: 0.8rem;
		font-weight: 600;
		color: #fbbf24;
	}

	.drag-hint {
		color: #4b5563;
		font-size: 1rem;
		padding: 0 0.25rem;
		opacity: 0;
		transition: opacity 0.2s;
	}

	.card-item:hover .drag-hint {
		opacity: 1;
	}

	/* Drop zones sidebar */
	.drop-zones-sidebar {
		width: 140px;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		opacity: 0.5;
		transition: opacity 0.2s;
	}

	.drop-zones-sidebar.active {
		opacity: 1;
	}

	.drop-zones-label {
		font-size: 0.75rem;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		padding: 0.25rem 0;
	}

	.drop-zone {
		padding: 1rem;
		background: rgba(255, 255, 255, 0.03);
		border: 2px dashed #2a3441;
		border-radius: 8px;
		text-align: center;
		transition: all 0.2s ease;
		cursor: default;
	}

	.drop-zone.hovered {
		border-color: var(--zone-color);
		background: rgba(var(--zone-color), 0.1);
		box-shadow: 0 0 20px rgba(var(--zone-color), 0.2);
		transform: scale(1.02);
	}

	.zone-label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #9ca3af;
	}

	.drop-zone.hovered .zone-label {
		color: #e2e8f0;
	}

	.drop-hint-text {
		display: block;
		font-size: 0.7rem;
		color: #6b7280;
		margin-top: 0.25rem;
	}

	/* Selected preview */
	.selected-preview {
		margin-top: 1rem;
		padding: 0.75rem 1rem;
		background: rgba(102, 126, 234, 0.15);
		border: 1px solid rgba(102, 126, 234, 0.3);
		border-radius: 8px;
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.preview-header {
		color: #9ca3af;
		font-size: 0.875rem;
	}

	.preview-card {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #e2e8f0;
	}

	/* Footer */
	.modal-footer {
		padding: 1rem 1.5rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		border-radius: 0 0 10px 10px;
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
	}

	.btn-cancel {
		padding: 0.75rem 1.5rem;
		background: transparent;
		color: #9ca3af;
		border: 1px solid #374151;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-cancel:hover:not(:disabled) {
		background: #374151;
		color: #e2e8f0;
		border-color: #4b5563;
	}

	.btn-confirm {
		padding: 0.75rem 1.5rem;
		background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.btn-confirm:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3);
	}

	.btn-confirm:disabled,
	.btn-cancel:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Drag ghost */
	.drag-ghost {
		position: fixed;
		pointer-events: none;
		z-index: 10000;
		transform: translate(-50%, -50%);
	}

	.drag-ghost-card {
		width: 70px;
		height: 98px;
		background: linear-gradient(135deg, #1a1f2e, #0d1117);
		border: 2px solid #667eea;
		border-radius: 6px;
		display: flex;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		box-shadow:
			0 15px 40px rgba(0, 0, 0, 0.6),
			0 0 30px rgba(102, 126, 234, 0.3);
		opacity: 0.95;
		transform: rotate(-5deg) scale(1.1);
		transition: all 0.15s ease;
	}

	.drag-ghost-card.over-zone {
		border-color: #22c55e;
		box-shadow:
			0 15px 40px rgba(0, 0, 0, 0.6),
			0 0 30px rgba(34, 197, 94, 0.5);
		transform: rotate(0deg) scale(1.15);
	}

	.drag-image {
		width: 100%;
		height: 100%;
		object-fit: cover;
		border-radius: 4px;
	}

	.drag-name {
		font-size: 0.6rem;
		font-weight: 600;
		color: white;
		text-align: center;
		padding: 0.25rem;
		line-height: 1.2;
	}

	.drag-destination {
		position: absolute;
		top: 100%;
		left: 50%;
		transform: translateX(-50%);
		margin-top: 0.5rem;
		padding: 0.25rem 0.5rem;
		background: rgba(34, 197, 94, 0.9);
		color: white;
		font-size: 0.7rem;
		font-weight: 600;
		border-radius: 4px;
		white-space: nowrap;
	}

	/* Scrollbar styling */
	.card-list::-webkit-scrollbar {
		width: 8px;
	}

	.card-list::-webkit-scrollbar-track {
		background: #141821;
		border-radius: 4px;
	}

	.card-list::-webkit-scrollbar-thumb {
		background: #374151;
		border-radius: 4px;
	}

	.card-list::-webkit-scrollbar-thumb:hover {
		background: #4b5563;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.library-search-modal {
			width: 100%;
			max-height: 100vh;
			border-radius: 0;
			max-width: none;
		}

		.modal-header {
			border-radius: 0;
		}

		.modal-footer {
			border-radius: 0;
		}

		.search-controls {
			flex-direction: column;
		}

		.content-area {
			flex-direction: column;
		}

		.drop-zones-sidebar {
			width: 100%;
			flex-direction: row;
			flex-wrap: wrap;
			justify-content: center;
		}

		.drop-zone {
			flex: 1;
			min-width: 80px;
			padding: 0.75rem;
		}

		.drop-zones-label {
			width: 100%;
			text-align: center;
		}

		.card-list {
			max-height: 300px;
		}

		.card-thumbnail {
			width: 40px;
			height: 56px;
		}
	}
</style>
