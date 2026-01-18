<script lang="ts">
	import type { CardView } from '$lib/generated/mage/v1/models';
	import ManaCost from '$lib/components/mtg/ManaCost.svelte';
	import { getScryfallImageUrl } from '$lib/utils/scryfall';

	let {
		cards,
		playerName = 'You',
		onMove,
		onShuffle,
		onClose
	}: {
		cards: CardView[];
		playerName?: string;
		onMove: (_cardId: string, _zone: 'HAND' | 'BATTLEFIELD' | 'GRAVEYARD' | 'EXILE') => void;
		onShuffle?: () => void;
		onClose: () => void;
	} = $props();

	let searchQuery = $state('');
	let filterType = $state<
		| 'all'
		| 'creature'
		| 'instant'
		| 'sorcery'
		| 'artifact'
		| 'enchantment'
		| 'land'
		| 'planeswalker'
	>('all');
	let selectedDestination = $state<'HAND' | 'BATTLEFIELD' | 'GRAVEYARD' | 'EXILE'>('HAND');
	let shuffleAfter = $state(true);

	const cardTypes = [
		'all',
		'creature',
		'instant',
		'sorcery',
		'artifact',
		'enchantment',
		'land',
		'planeswalker'
	] as const;

	const filteredCards = $derived(() => {
		let result = [...cards];

		if (searchQuery.trim()) {
			const q = searchQuery.toLowerCase();
			result = result.filter((c) => {
				const name = (c.name || '').toLowerCase();
				const type = (c.type || '').toLowerCase();
				// rulesText isn't always present in CardView, so keep it optional
				const rulesText = (c as unknown as { rulesText?: string }).rulesText?.toLowerCase() || '';
				return name.includes(q) || type.includes(q) || rulesText.includes(q);
			});
		}

		if (filterType !== 'all') {
			result = result.filter((c) => (c.type || '').toLowerCase().includes(filterType));
		}

		return result;
	});

	function imageUrlFor(cardName: string): string {
		return getScryfallImageUrl(cardName, 'small');
	}

	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	function handleBackdropClick(event: MouseEvent): void {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	function sendCard(cardId: string): void {
		onMove(cardId, selectedDestination);
		if (shuffleAfter) {
			onShuffle?.();
		}
		onClose();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="overlay"
	role="dialog"
	aria-modal="true"
	aria-label="Search library"
	tabindex="-1"
	onclick={handleBackdropClick}
>
	<div class="modal">
		<header class="header">
			<div class="title">
				<h2>📚 {playerName === 'You' ? 'Your' : `${playerName}'s`} Deck</h2>
				<span class="count">{cards.length} cards</span>
			</div>
			<button class="close" onclick={onClose} title="Close">✕</button>
		</header>

		<div class="content">
			<div class="controls">
				<input
					class="search"
					type="text"
					placeholder="Search by name or type..."
					bind:value={searchQuery}
				/>
				<select class="type" bind:value={filterType}>
					{#each cardTypes as t}
						<option value={t}
							>{t === 'all' ? 'All Types' : t.charAt(0).toUpperCase() + t.slice(1)}</option
						>
					{/each}
				</select>
			</div>

			<div class="send-row">
				<label class="send-label">
					Send to:
					<select class="dest" bind:value={selectedDestination}>
						<option value="HAND">🖐️ Hand</option>
						<option value="BATTLEFIELD">⚔️ Battlefield</option>
						<option value="GRAVEYARD">🪦 Graveyard</option>
						<option value="EXILE">✨ Exile</option>
					</select>
				</label>

				<label class="shuffle">
					<input type="checkbox" bind:checked={shuffleAfter} />
					Shuffle after
				</label>
			</div>

			<div class="list">
				{#if filteredCards().length === 0}
					<div class="empty">
						{#if cards.length === 0}
							Deck is empty
						{:else}
							No cards match your search
						{/if}
					</div>
				{:else}
					{#each filteredCards() as card (card.id)}
						<button
							class="row"
							type="button"
							onclick={() => sendCard(card.id)}
							title="Click to send"
						>
							<div class="thumb">
								{#if imageUrlFor(card.name)}
									<img src={imageUrlFor(card.name)} alt={card.name} draggable="false" />
								{:else}
									<span class="placeholder">🃏</span>
								{/if}
							</div>
							<div class="info">
								<div class="name-row">
									<span class="name">{card.name}</span>
									{#if card.manaCost}
										<ManaCost cost={card.manaCost} size="sm" />
									{/if}
								</div>
								<span class="type-text">{card.type}</span>
							</div>
						</button>
					{/each}
				{/if}
			</div>
		</div>
	</div>
</div>

<style>
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 2000;
	}

	.modal {
		width: min(900px, 95vw);
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		background: #1a1f2e;
		border: 2px solid rgba(34, 197, 94, 0.45);
		border-radius: 12px;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.55);
		overflow: hidden;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem 1.25rem;
		background: linear-gradient(135deg, rgba(34, 197, 94, 0.25), rgba(102, 126, 234, 0.18));
		border-bottom: 1px solid rgba(42, 52, 65, 1);
	}

	.title {
		display: flex;
		align-items: baseline;
		gap: 0.75rem;
	}

	.title h2 {
		margin: 0;
		font-size: 1.15rem;
		color: #e2e8f0;
	}

	.count {
		font-size: 0.85rem;
		color: #9ca3af;
	}

	.close {
		width: 32px;
		height: 32px;
		border-radius: 6px;
		border: 1px solid rgba(55, 65, 81, 0.7);
		background: rgba(17, 24, 39, 0.7);
		color: #e2e8f0;
		cursor: pointer;
	}

	.close:hover {
		background: rgba(55, 65, 81, 0.8);
	}

	.content {
		padding: 1rem 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		overflow: hidden;
	}

	.controls {
		display: flex;
		gap: 0.75rem;
	}

	.search {
		flex: 1;
		padding: 0.75rem 1rem;
		border-radius: 8px;
		border: 1px solid #2a3441;
		background: #141821;
		color: #e2e8f0;
	}

	.search:focus {
		outline: none;
		border-color: rgba(34, 197, 94, 0.8);
	}

	.type {
		padding: 0.75rem 1rem;
		border-radius: 8px;
		border: 1px solid #2a3441;
		background: #141821;
		color: #e2e8f0;
	}

	.send-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.send-label {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		color: #9ca3af;
		font-size: 0.875rem;
	}

	.dest {
		padding: 0.4rem 0.75rem;
		border-radius: 8px;
		border: 1px solid #2a3441;
		background: #141821;
		color: #e2e8f0;
	}

	.shuffle {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		color: #9ca3af;
		font-size: 0.875rem;
	}

	.list {
		flex: 1;
		overflow: auto;
		border: 1px solid #2a3441;
		border-radius: 8px;
		background: #141821;
	}

	.empty {
		padding: 2rem;
		text-align: center;
		color: #6b7280;
	}

	.row {
		width: 100%;
		display: flex;
		gap: 0.75rem;
		align-items: center;
		padding: 0.5rem 0.75rem;
		border: none;
		border-bottom: 1px solid #2a3441;
		background: transparent;
		color: inherit;
		cursor: pointer;
		text-align: left;
	}

	.row:hover {
		background: rgba(34, 197, 94, 0.1);
	}

	.row:last-child {
		border-bottom: none;
	}

	.thumb {
		width: 50px;
		height: 70px;
		border-radius: 4px;
		overflow: hidden;
		background: #0d1117;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.thumb img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.placeholder {
		opacity: 0.6;
	}

	.info {
		display: flex;
		flex-direction: column;
		min-width: 0;
		gap: 0.2rem;
	}

	.name-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.name {
		font-weight: 600;
		color: #e2e8f0;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.type-text {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	@media (max-width: 768px) {
		.controls {
			flex-direction: column;
		}
	}
</style>
