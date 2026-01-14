<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Deck } from '$lib/types/deck';
	import { fetchUserDecks } from '$lib/api/decks';
	import { playtestGameStore, type PlaytestSessionMeta } from '$lib/stores/playtest-game';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import Play from '@lucide/svelte/icons/play';
	import Clock from '@lucide/svelte/icons/clock';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import Users from '@lucide/svelte/icons/users';
	import RotateCcw from '@lucide/svelte/icons/rotate-ccw';

	// Props
	let {
		open = $bindable(false),
		onClose
	}: {
		open?: boolean;
		onClose: () => void;
	} = $props();

	// View mode: 'sessions' or 'new'
	let viewMode = $state<'sessions' | 'new'>('sessions');

	// Session state
	let availableSessions = $state<PlaytestSessionMeta[]>([]);

	// Deck state
	let availableDecks = $state<Deck[]>([]);
	let loadingDecks = $state(false);
	let deckError = $state<string | null>(null);

	// Selected decks for each player
	let selectedDecks = $state<(string | null)[]>([null, null, null, null]);
	let playerCount = $state(2);

	// UI state
	let starting = $state(false);
	let continuing = $state(false);

	/**
	 * Load available playtest sessions
	 */
	function loadSessions(): void {
		availableSessions = playtestGameStore.listSessions();
	}

	/**
	 * Continue a playtest session
	 */
	async function continueSession(sessionId: string): Promise<void> {
		continuing = true;
		try {
			await goto(`/playtest?playtestId=${sessionId}`);
			onClose();
		} catch (err) {
			console.error('Failed to continue playtest:', err);
			deckError = 'Failed to continue playtest session';
		} finally {
			continuing = false;
		}
	}

	/**
	 * Delete a playtest session
	 */
	function deleteSession(sessionId: string): void {
		playtestGameStore.deleteSession(sessionId);
		loadSessions();
	}

	/**
	 * Format date for display
	 */
	function formatDate(timestamp: number): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		const diffHours = Math.floor(diffMs / 3600000);
		const diffDays = Math.floor(diffMs / 86400000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		if (diffDays < 7) return `${diffDays}d ago`;
		return date.toLocaleDateString();
	}

	/**
	 * Get an unchosen deck from available decks
	 */
	function getUnchosenDeck(): string | null {
		const chosen = new Set(selectedDecks.filter((id): id is string => id !== null));
		const unchosen = availableDecks.find((deck) => !chosen.has(deck.id));
		return unchosen?.id ?? null;
	}

	/**
	 * Prepopulate deck slots with unchosen decks
	 */
	function prepopulateDeckSlots(): void {
		for (let i = 0; i < playerCount; i++) {
			if (selectedDecks[i] === null) {
				selectedDecks[i] = getUnchosenDeck();
			}
		}
	}

	/**
	 * Load all available decks
	 */
	async function loadDecks(): Promise<void> {
		loadingDecks = true;
		deckError = null;
		try {
			availableDecks = await fetchUserDecks();

			// Auto-select first available decks if we have any
			prepopulateDeckSlots();
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
			.every((deckId) => deckId !== null);

		return hasMinimumPlayers && requiredDecksSelected;
	});

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
		continuing = false;
		deckError = null;
	}

	/**
	 * Load data when modal opens
	 */
	$effect(() => {
		if (open) {
			loadSessions();
			loadDecks();
		} else {
			resetForm();
			viewMode = 'sessions';
		}
	});

	/**
	 * Prepopulate deck slots when player count changes
	 */
	$effect(() => {
		// React to playerCount changes
		void playerCount;
		if (open && availableDecks.length > 0 && viewMode === 'new') {
			prepopulateDeckSlots();
		}
	});
</script>

<Modal bind:open title="Playtest" {onClose} size="medium">
	<div class="playtest-modal">
		<!-- Tab switcher -->
		<div class="tabs">
			<button
				class="tab-btn"
				class:active={viewMode === 'sessions'}
				onclick={() => (viewMode = 'sessions')}
			>
				<Clock size={16} />
				Continue Session
			</button>
			<button class="tab-btn" class:active={viewMode === 'new'} onclick={() => (viewMode = 'new')}>
				<Play size={16} />
				Start New
			</button>
		</div>

		{#if viewMode === 'sessions'}
			<!-- Sessions view -->
			<div class="sessions-view">
				{#if availableSessions.length === 0}
					<div class="empty-state">
						<Clock size={48} class="empty-icon" />
						<p>No saved sessions</p>
						<p class="hint">Start a new playtest to create your first session.</p>
					</div>
				{:else}
					<div class="sessions-list">
						{#each availableSessions as session (session.id)}
							<div class="session-card">
								<div class="session-info">
									<div class="session-label">{session.label}</div>
									<div class="session-meta">
										<span class="meta-item">
											<Users size={14} />
											{session.playerCount} players
										</span>
										<span class="meta-item">
											<RotateCcw size={14} />
											Turn {session.turn}
										</span>
										<span class="meta-item">
											<Clock size={14} />
											{formatDate(session.savedAt)}
										</span>
									</div>
								</div>
								<div class="session-actions">
									<button
										class="btn-continue"
										onclick={() => continueSession(session.id)}
										disabled={continuing}
									>
										<Play size={16} />
										Continue
									</button>
									<button
										class="btn-delete"
										onclick={() => deleteSession(session.id)}
										title="Delete session"
									>
										<Trash2 size={16} />
									</button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		{:else if viewMode === 'new'}
			<!-- New playtest view -->
			{#if loadingDecks}
				<div class="loading-container">
					<LoadingSpinner />
					<p>Loading your decks...</p>
				</div>
			{:else if deckError}
				<div class="error-container">
					<p class="error-message">{deckError}</p>
					<button class="btn-retry" onclick={loadDecks}> Retry </button>
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
									onclick={() => (playerCount = count)}
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
								<select id="deck-{i + 1}" bind:value={selectedDecks[i]} class="deck-select">
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
							<strong>Playtest Mode</strong> is a client-side sandbox for testing decks. No rules enforcement,
							no server connection required.
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
						<button class="btn-secondary" onclick={onClose} disabled={starting}> Cancel </button>
						<button class="btn-primary" onclick={startPlaytest} disabled={!isValid() || starting}>
							{starting ? 'Starting...' : 'Start Playtest'}
						</button>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</Modal>

<style>
	.playtest-modal {
		width: 100%;
		min-height: 300px;
	}

	.tabs {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 1.5rem;
		border-bottom: 1px solid #2a3441;
	}

	.tab-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: transparent;
		border: none;
		border-bottom: 2px solid transparent;
		color: #94a3b8;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		margin-bottom: -1px;
	}

	.tab-btn:hover {
		color: #f8fafc;
	}

	.tab-btn.active {
		color: #667eea;
		border-bottom-color: #667eea;
	}

	.sessions-view {
		min-height: 300px;
	}

	.sessions-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		max-height: 400px;
		overflow-y: auto;
		padding: 0.5rem;
		margin: -0.5rem;
	}

	.session-card {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 8px;
		transition: all 0.2s;
	}

	.session-card:hover {
		border-color: #667eea;
		background: #1f2533;
	}

	.session-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.session-label {
		font-weight: 600;
		color: #f8fafc;
		font-size: 0.9375rem;
	}

	.session-meta {
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.meta-item {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		color: #94a3b8;
		font-size: 0.8125rem;
	}

	.meta-item :global(svg) {
		opacity: 0.7;
	}

	.session-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-continue {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-continue:hover:not(:disabled) {
		background: #5568d3;
	}

	.btn-continue:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-delete {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.5rem;
		background: transparent;
		border: 1px solid #374151;
		border-radius: 6px;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-delete:hover {
		background: #1f2937;
		border-color: #ef4444;
		color: #ef4444;
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

	.empty-icon {
		color: #94a3b8;
		opacity: 0.5;
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
