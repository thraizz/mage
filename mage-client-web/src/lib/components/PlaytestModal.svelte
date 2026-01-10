<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { Deck } from '$lib/types/deck';
	import { fetchUserDecks } from '$lib/api/decks';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';

	// Props
	let {
		open = $bindable(false),
		onClose
	}: {
		open?: boolean;
		onClose: () => void;
	} = $props();

	// Deck state
	let availableDecks = $state<Deck[]>([]);
	let loadingDecks = $state(false);
	let deckError = $state<string | null>(null);

	// Selected decks for each player
	let selectedDecks = $state<(string | null)[]>([null, null, null, null]);
	let playerCount = $state(2);

	// UI state
	let starting = $state(false);

	/**
	 * Load all available decks
	 */
	async function loadDecks(): Promise<void> {
		loadingDecks = true;
		deckError = null;
		try {
			availableDecks = await fetchUserDecks();
			
			// Auto-select first available decks if we have any
			if (availableDecks.length > 0) {
				selectedDecks[0] = availableDecks[0].id;
				if (availableDecks.length > 1) {
					selectedDecks[1] = availableDecks[1].id;
				}
			}
		} catch (err) {
			console.error('Failed to load decks:', err);
			deckError = err instanceof Error ? err.message : 'Failed to load decks';
			availableDecks = [];
		} finally {
			loadingDecks = false;
		}
	}

	/**
	 * Check if form is valid
	 */
	const isValid = $derived(() => {
		// Need at least 2 players with selected decks
		const hasMinimumPlayers = playerCount >= 2;
		const requiredDecksSelected = selectedDecks
			.slice(0, playerCount)
			.every(deckId => deckId !== null);
		
		return hasMinimumPlayers && requiredDecksSelected;
	});

	/**
	 * Get deck name by ID
	 */
	function getDeckName(deckId: string | null): string {
		if (!deckId) return 'Select a deck';
		const deck = availableDecks.find(d => d.id === deckId);
		return deck ? `${deck.name} (${deck.format})` : 'Unknown deck';
	}

	/**
	 * Start playtest session
	 */
	async function startPlaytest(): Promise<void> {
		if (!isValid()) return;

		starting = true;
		try {
			// Build URL params with selected deck IDs
			const params = new URLSearchParams();
			for (let i = 0; i < playerCount; i++) {
				const deckId = selectedDecks[i];
				if (deckId) {
					params.append(`d${i + 1}`, deckId);
				}
			}

			// Navigate to playtest page
			await goto(`/playtest?${params.toString()}`);
			
			// Close modal
			onClose();
		} catch (err) {
			console.error('Failed to start playtest:', err);
			deckError = 'Failed to start playtest session';
		} finally {
			starting = false;
		}
	}

	/**
	 * Reset form
	 */
	function resetForm(): void {
		selectedDecks = [null, null, null, null];
		playerCount = 2;
		starting = false;
		deckError = null;
	}

	/**
	 * Load decks when modal opens
	 */
	$effect(() => {
		if (open) {
			loadDecks();
		} else {
			resetForm();
		}
	});
</script>

<Modal bind:open title="Configure Playtest" {onClose} size="medium">
	<div class="playtest-modal">
		{#if loadingDecks}
			<div class="loading-container">
				<LoadingSpinner />
				<p>Loading your decks...</p>
			</div>
		{:else if deckError}
			<div class="error-container">
				<p class="error-message">{deckError}</p>
				<button class="btn-retry" onclick={loadDecks}>
					Retry
				</button>
			</div>
		{:else if availableDecks.length === 0}
			<div class="empty-state">
				<p>You don't have any decks yet.</p>
				<p class="hint">Create some decks first to use playtest mode.</p>
			</div>
		{:else}
			<div class="form-content">
				<!-- Player count selector -->
				<div class="form-group">
					<label for="player-count">Number of Players</label>
					<div class="player-count-buttons">
						{#each [2, 3, 4] as count}
							<button
								class="player-count-btn"
								class:active={playerCount === count}
								onclick={() => playerCount = count}
							>
								{count}
							</button>
						{/each}
					</div>
				</div>

				<!-- Deck selectors for each player -->
				<div class="deck-selectors">
					{#each Array(playerCount) as _, i}
						<div class="form-group">
							<label for="deck-{i + 1}">Player {i + 1} Deck</label>
							<select
								id="deck-{i + 1}"
								bind:value={selectedDecks[i]}
								class="deck-select"
							>
								<option value={null}>Select a deck...</option>
								{#each availableDecks as deck}
									<option value={deck.id}>
										{deck.name} ({deck.format}, {deck.cardCount} cards)
									</option>
								{/each}
							</select>
						</div>
					{/each}
				</div>

				<!-- Info section -->
				<div class="info-box">
					<p class="info-text">
						<strong>Playtest Mode</strong> is a client-side sandbox for testing decks.
						No rules enforcement, no server connection required.
					</p>
					<ul class="features-list">
						<li>Switch between player perspectives</li>
						<li>Drag cards between zones</li>
						<li>Modify life totals freely</li>
						<li>Draw opening hands automatically</li>
					</ul>
				</div>

				<!-- Action buttons -->
				<div class="modal-actions">
					<button class="btn-secondary" onclick={onClose} disabled={starting}>
						Cancel
					</button>
					<button
						class="btn-primary"
						onclick={startPlaytest}
						disabled={!isValid() || starting}
					>
						{starting ? 'Starting...' : 'Start Playtest'}
					</button>
				</div>
			</div>
		{/if}
	</div>
</Modal>

<style>
	.playtest-modal {
		width: 100%;
		min-height: 300px;
	}

	.loading-container,
	.error-container,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 1rem;
		text-align: center;
		gap: 1rem;
	}

	.error-message {
		color: #ef4444;
		margin: 0;
	}

	.btn-retry {
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-retry:hover {
		background: #5568d3;
	}

	.empty-state p {
		margin: 0.5rem 0;
	}

	.hint {
		color: #94a3b8;
		font-size: 0.875rem;
	}

	.form-content {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-weight: 600;
		color: #f8fafc;
		font-size: 0.875rem;
	}

	.player-count-buttons {
		display: flex;
		gap: 0.5rem;
	}

	.player-count-btn {
		flex: 1;
		padding: 0.75rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 6px;
		color: #94a3b8;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.player-count-btn:hover {
		border-color: #667eea;
		color: #f8fafc;
	}

	.player-count-btn.active {
		background: #667eea;
		border-color: #667eea;
		color: white;
	}

	.deck-selectors {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		max-height: 300px;
		overflow-y: auto;
		padding: 0.5rem;
		margin: -0.5rem;
	}

	.deck-select {
		width: 100%;
		padding: 0.75rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 6px;
		color: #f8fafc;
		font-size: 0.875rem;
		cursor: pointer;
		transition: border-color 0.2s;
	}

	.deck-select:hover {
		border-color: #667eea;
	}

	.deck-select:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.info-box {
		padding: 1rem;
		background: rgba(102, 126, 234, 0.1);
		border: 1px solid rgba(102, 126, 234, 0.2);
		border-radius: 8px;
	}

	.info-text {
		margin: 0 0 0.75rem;
		color: #cbd5e1;
		font-size: 0.875rem;
		line-height: 1.5;
	}

	.features-list {
		margin: 0;
		padding-left: 1.5rem;
		color: #94a3b8;
		font-size: 0.8125rem;
	}

	.features-list li {
		margin: 0.25rem 0;
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		padding-top: 1rem;
		border-top: 1px solid #2a3441;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.75rem 1.5rem;
		border-radius: 6px;
		font-weight: 600;
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
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: transparent;
		border: 1px solid #374151;
		color: #9ca3af;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #1f2937;
		border-color: #4b5563;
		color: #f8fafc;
	}

	.btn-secondary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
