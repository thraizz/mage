<script lang="ts">
	import type { CardView } from '$lib/generated/mage/v1/models';
	import Card from './Card.svelte';

	interface Props {
		cards: CardView[];
		onComplete: (keepOnTop: CardView[], putToBottom: CardView[]) => void;
		onCancel: () => void;
	}

	let { cards, onComplete, onCancel }: Props = $props();

	// State: cards in each zone
	let unassigned = $state<CardView[]>([...cards]);
	let keepOnTop = $state<CardView[]>([]);
	let putToBottom = $state<CardView[]>([]);

	// Drag state
	let draggedCard: CardView | null = null;
	let dragSource: 'unassigned' | 'keepOnTop' | 'putToBottom' | null = null;

	function handleDragStart(card: CardView, source: typeof dragSource) {
		draggedCard = card;
		dragSource = source;
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
	}

	function handleDrop(e: DragEvent, target: 'unassigned' | 'keepOnTop' | 'putToBottom') {
		e.preventDefault();
		if (!draggedCard || !dragSource) return;

		// Remove from source
		if (dragSource === 'unassigned') {
			unassigned = unassigned.filter((c) => c.id !== draggedCard!.id);
		} else if (dragSource === 'keepOnTop') {
			keepOnTop = keepOnTop.filter((c) => c.id !== draggedCard!.id);
		} else if (dragSource === 'putToBottom') {
			putToBottom = putToBottom.filter((c) => c.id !== draggedCard!.id);
		}

		// Add to target
		if (target === 'unassigned') {
			unassigned = [...unassigned, draggedCard];
		} else if (target === 'keepOnTop') {
			keepOnTop = [...keepOnTop, draggedCard];
		} else if (target === 'putToBottom') {
			putToBottom = [...putToBottom, draggedCard];
		}

		draggedCard = null;
		dragSource = null;
	}

	function moveToKeepOnTop(card: CardView) {
		unassigned = unassigned.filter((c) => c.id !== card.id);
		putToBottom = putToBottom.filter((c) => c.id !== card.id);
		keepOnTop = [...keepOnTop, card];
	}

	function moveToPutToBottom(card: CardView) {
		unassigned = unassigned.filter((c) => c.id !== card.id);
		keepOnTop = keepOnTop.filter((c) => c.id !== card.id);
		putToBottom = [...putToBottom, card];
	}

	function moveToUnassigned(card: CardView) {
		keepOnTop = keepOnTop.filter((c) => c.id !== card.id);
		putToBottom = putToBottom.filter((c) => c.id !== card.id);
		unassigned = [...unassigned, card];
	}

	function handleComplete() {
		// All cards must be assigned
		if (unassigned.length > 0) {
			alert('Please assign all cards to either keep on top or put to bottom.');
			return;
		}
		onComplete(keepOnTop, putToBottom);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onCancel();
		} else if (e.key === 'Enter' && unassigned.length === 0) {
			handleComplete();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="overlay" role="dialog" aria-labelledby="scry-dialog-title">
	<div class="dialog">
		<div class="dialog-header">
			<h2 id="scry-dialog-title">Scry {cards.length}</h2>
			<button class="close-button" onclick={onCancel} aria-label="Close dialog">×</button>
		</div>

		<div class="instructions">
			<p>Organize the cards. Keep some on top, put the rest to the bottom.</p>
			{#if unassigned.length > 0}
				<p class="warning">⚠️ {unassigned.length} card(s) not yet assigned</p>
			{/if}
		</div>

		<div class="zones">
			<!-- Keep on Top Zone -->
			<div
				class="zone keep-zone"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'keepOnTop')}
			>
				<div class="zone-header">
					<h3>Keep on Top</h3>
					<span class="zone-count">{keepOnTop.length} card(s)</span>
				</div>
				<div class="card-list">
					{#if keepOnTop.length === 0}
						<p class="empty-message">Drag cards here to keep on top</p>
					{:else}
						{#each keepOnTop as card, index (card.id)}
							<div
								class="card-item"
								draggable="true"
								ondragstart={() => handleDragStart(card, 'keepOnTop')}
							>
								<span class="card-order">{index + 1}.</span>
								<div class="card-preview">
									<Card {card} interactive={false} />
								</div>
								<button class="btn-remove" onclick={() => moveToUnassigned(card)}>
									↩ Return
								</button>
							</div>
						{/each}
					{/if}
				</div>
			</div>

			<!-- Unassigned Zone -->
			<div
				class="zone unassigned-zone"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'unassigned')}
			>
				<div class="zone-header">
					<h3>Scrying</h3>
					<span class="zone-count">{unassigned.length} card(s)</span>
				</div>
				<div class="card-list">
					{#if unassigned.length === 0}
						<p class="empty-message">All cards assigned!</p>
					{:else}
						{#each unassigned as card (card.id)}
							<div
								class="card-item"
								draggable="true"
								ondragstart={() => handleDragStart(card, 'unassigned')}
							>
								<div class="card-preview">
									<Card {card} interactive={false} />
								</div>
								<div class="card-actions">
									<button class="btn-top" onclick={() => moveToKeepOnTop(card)}>
										↑ Keep
									</button>
									<button class="btn-bottom" onclick={() => moveToPutToBottom(card)}>
										↓ Bottom
									</button>
								</div>
							</div>
						{/each}
					{/if}
				</div>
			</div>

			<!-- Put to Bottom Zone -->
			<div
				class="zone bottom-zone"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'putToBottom')}
			>
				<div class="zone-header">
					<h3>Put to Bottom</h3>
					<span class="zone-count">{putToBottom.length} card(s)</span>
				</div>
				<div class="card-list">
					{#if putToBottom.length === 0}
						<p class="empty-message">Drag cards here to put to bottom</p>
					{:else}
						{#each putToBottom as card, index (card.id)}
							<div
								class="card-item"
								draggable="true"
								ondragstart={() => handleDragStart(card, 'putToBottom')}
							>
								<span class="card-order">{index + 1}.</span>
								<div class="card-preview">
									<Card {card} interactive={false} />
								</div>
								<button class="btn-remove" onclick={() => moveToUnassigned(card)}>
									↩ Return
								</button>
							</div>
						{/each}
					{/if}
				</div>
			</div>
		</div>

		<div class="dialog-footer">
			<button class="btn-secondary" onclick={onCancel}>Cancel</button>
			<button class="btn-primary" onclick={handleComplete} disabled={unassigned.length > 0}>
				Complete Scry
			</button>
		</div>
	</div>
</div>

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		backdrop-filter: blur(4px);
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

	.dialog {
		background: #1a1f2e;
		border: 2px solid #3a4451;
		border-radius: 12px;
		padding: 1.5rem;
		max-width: 1200px;
		width: 95%;
		max-height: 90vh;
		overflow-y: auto;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
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

	.dialog-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid #3a4451;
	}

	.dialog-header h2 {
		margin: 0;
		color: #fff;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.close-button {
		background: transparent;
		border: none;
		color: #9ca3af;
		font-size: 2rem;
		line-height: 1;
		cursor: pointer;
		padding: 0;
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: color 0.2s;
	}

	.close-button:hover {
		color: #fff;
	}

	.instructions {
		margin-bottom: 1.5rem;
		padding: 1rem;
		background: rgba(102, 126, 234, 0.1);
		border: 1px solid #667eea;
		border-radius: 8px;
	}

	.instructions p {
		margin: 0;
		color: #ddd;
		font-size: 0.875rem;
	}

	.instructions .warning {
		margin-top: 0.5rem;
		color: #fbbf24;
		font-weight: 600;
	}

	.zones {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1rem;
		margin-bottom: 1.5rem;
	}

	.zone {
		border: 2px solid #3a4451;
		border-radius: 8px;
		padding: 1rem;
		min-height: 300px;
		transition: border-color 0.2s;
	}

	.zone.keep-zone {
		border-color: #10b981;
		background: rgba(16, 185, 129, 0.05);
	}

	.zone.bottom-zone {
		border-color: #ef4444;
		background: rgba(239, 68, 68, 0.05);
	}

	.zone.unassigned-zone {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.05);
	}

	.zone-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
		padding-bottom: 0.5rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.zone-header h3 {
		margin: 0;
		color: #fff;
		font-size: 0.875rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.zone-count {
		color: #9ca3af;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.card-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.empty-message {
		color: #6b7280;
		font-style: italic;
		text-align: center;
		padding: 2rem 1rem;
		font-size: 0.875rem;
	}

	.card-item {
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid #3a4451;
		border-radius: 6px;
		padding: 0.75rem;
		display: flex;
		align-items: center;
		gap: 0.75rem;
		cursor: grab;
		transition: all 0.2s;
	}

	.card-item:hover {
		background: rgba(255, 255, 255, 0.08);
		transform: translateY(-1px);
	}

	.card-item:active {
		cursor: grabbing;
	}

	.card-order {
		color: #9ca3af;
		font-weight: 700;
		font-size: 0.875rem;
		min-width: 1.5rem;
	}

	.card-preview {
		flex: 1;
		display: flex;
		align-items: center;
	}

	.card-preview :global(.card) {
		max-width: 120px;
	}

	.card-actions {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.btn-top,
	.btn-bottom,
	.btn-remove {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
		white-space: nowrap;
	}

	.btn-top {
		background: rgba(16, 185, 129, 0.2);
		color: #10b981;
		border: 1px solid #10b981;
	}

	.btn-top:hover {
		background: rgba(16, 185, 129, 0.3);
	}

	.btn-bottom {
		background: rgba(239, 68, 68, 0.2);
		color: #ef4444;
		border: 1px solid #ef4444;
	}

	.btn-bottom:hover {
		background: rgba(239, 68, 68, 0.3);
	}

	.btn-remove {
		background: rgba(156, 163, 175, 0.2);
		color: #9ca3af;
		border: 1px solid #6b7280;
	}

	.btn-remove:hover {
		background: rgba(156, 163, 175, 0.3);
	}

	.dialog-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid #3a4451;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn-primary {
		background: #667eea;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background: #5568d3;
		transform: translateY(-1px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: rgba(255, 255, 255, 0.1);
		color: #fff;
		border: 1px solid #3a4451;
	}

	.btn-secondary:hover {
		background: rgba(255, 255, 255, 0.15);
	}
</style>
