<script lang="ts">
	import type { Table } from '$lib/types/table';
	import type { Deck } from '$lib/types/deck';
	import { joinTable } from '$lib/api/table';
	import { fetchUserDecks } from '$lib/api/decks';
	import { toast } from '$lib/stores/toast';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { structuredCardsToText, type CardEntry } from '$lib/utils/deck-parser';

	// Props
	let {
		open = $bindable(false),
		table,
		// eslint-disable-next-line no-unused-vars
		onSuccess
	}: {
		open: boolean;
		table: Table | null;
		// eslint-disable-next-line no-unused-vars
		onSuccess: (tableId: string) => void;
	} = $props();

	// State
	let decks = $state<Deck[]>([]);
	let selectedDeckId = $state<string>('');
	let password = $state<string>('');
	let loading = $state(false);
	let loadingDecks = $state(false);
	let error = $state<string | null>(null);

	// Load decks when modal opens
	$effect(() => {
		if (open && table) {
			loadDecks();
		}
	});

	/**
	 * Load user's decks for this format
	 */
	async function loadDecks(): Promise<void> {
		if (!table) return;

		loadingDecks = true;
		error = null;
		selectedDeckId = '';

		try {
			const allDecks = await fetchUserDecks();
			// Filter decks by format
			decks = allDecks.filter((d: Deck) => d.format === table.format);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load decks';
			console.error('Failed to load decks:', err);
		} finally {
			loadingDecks = false;
		}
	}

	/**
	 * Handle join
	 */
	async function handleJoin(): Promise<void> {
		if (!table) return;

		if (!selectedDeckId) {
			error = 'Please select a deck';
			return;
		}

		if (table.hasPassword && !password.trim()) {
			error = 'Password is required';
			return;
		}

		loading = true;
		error = null;

		try {
			// Find selected deck
			const deck = decks.find((d) => d.id === selectedDeckId);
			if (!deck) {
				throw new Error('Selected deck not found');
			}

			// Convert deck to text format
			const deckCards: CardEntry[] = [
				...deck.commanders.map((c) => ({
					name: c.cardName,
					quantity: c.quantity,
					section: 'commander' as const
				})),
				...deck.mainDeck.map((c) => ({
					name: c.cardName,
					quantity: c.quantity,
					section: 'main' as const
				})),
				...deck.sideboard.map((c) => ({
					name: c.cardName,
					quantity: c.quantity,
					section: 'sideboard' as const
				}))
			];
			const deckList = structuredCardsToText(deckCards);

			// Join table with deck
			await joinTable(table.id, deckList, table.hasPassword ? password : undefined);

			toast.success(`Joined table: ${table.name}`);
			onSuccess(table.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to join table';
			console.error('Failed to join table:', err);
			toast.error(error);
		} finally {
			loading = false;
		}
	}

	/**
	 * Handle close
	 */
	function handleClose(): void {
		if (loading) return;
		open = false;
		password = '';
		error = null;
	}
</script>

<Modal bind:open size="medium" closeOnBackdrop={!loading}>
	<div class="join-modal">
		<h2 class="modal-title">Join Table</h2>

		{#if table}
			<div class="table-info">
				<div class="info-row">
					<span class="info-label">Table:</span>
					<span class="info-value">{table.name}</span>
				</div>
				<div class="info-row">
					<span class="info-label">Format:</span>
					<span class="info-value format-badge">{table.format}</span>
				</div>
				<div class="info-row">
					<span class="info-label">Host:</span>
					<span class="info-value">{table.hostUsername}</span>
				</div>
				<div class="info-row">
					<span class="info-label">Players:</span>
					<span class="info-value">{table.players.length}/{table.maxPlayers}</span>
				</div>
			</div>

			{#if loadingDecks}
				<div class="loading-container">
					<LoadingSpinner size="medium" />
					<p class="loading-text">Loading your decks...</p>
				</div>
			{:else if decks.length === 0}
				<div class="no-decks">
					<p class="no-decks-message">
						You don't have any {table.format} decks. Please create one in the Deck Manager first.
					</p>
					<a href="/decks" class="create-deck-link">Go to Deck Manager</a>
				</div>
			{:else}
				<div class="form-group">
					<label for="deck-select" class="form-label">Select Deck</label>
					<select
						id="deck-select"
						class="form-select"
						bind:value={selectedDeckId}
						disabled={loading}
					>
						<option value="">-- Choose a deck --</option>
						{#each decks as deck}
							<option value={deck.id}>{deck.name} ({deck.cardCount} cards)</option>
						{/each}
					</select>
				</div>

				{#if table.hasPassword}
					<div class="form-group">
						<label for="password-input" class="form-label">Password</label>
						<input
							id="password-input"
							type="password"
							class="form-input"
							placeholder="Enter table password"
							bind:value={password}
							disabled={loading}
						/>
					</div>
				{/if}

				{#if error}
					<div class="error-message">{error}</div>
				{/if}

				<div class="modal-actions">
					<button class="btn btn-secondary" onclick={handleClose} disabled={loading}>
						Cancel
					</button>
					<button
						class="btn btn-primary"
						onclick={handleJoin}
						disabled={loading || !selectedDeckId}
					>
						{#if loading}
							<LoadingSpinner size="small" />
							<span>Joining...</span>
						{:else}
							Join Table
						{/if}
					</button>
				</div>
			{/if}
		{/if}
	</div>
</Modal>

<style>
	.join-modal {
		padding: 1.5rem;
	}

	.modal-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #111827;
		margin: 0 0 1.5rem 0;
	}

	.table-info {
		background-color: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 1rem;
		margin-bottom: 1.5rem;
	}

	.info-row {
		display: flex;
		justify-content: space-between;
		padding: 0.5rem 0;
	}

	.info-row:not(:last-child) {
		border-bottom: 1px solid #e5e7eb;
	}

	.info-label {
		font-weight: 600;
		color: #6b7280;
	}

	.info-value {
		color: #111827;
	}

	.format-badge {
		background-color: #667eea;
		color: white;
		padding: 0.125rem 0.5rem;
		border-radius: 0.25rem;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		padding: 2rem 0;
	}

	.loading-text {
		color: #6b7280;
		margin: 0;
	}

	.no-decks {
		text-align: center;
		padding: 2rem 0;
	}

	.no-decks-message {
		color: #6b7280;
		margin: 0 0 1rem 0;
	}

	.create-deck-link {
		display: inline-block;
		padding: 0.5rem 1rem;
		background-color: #667eea;
		color: white;
		border-radius: 0.375rem;
		text-decoration: none;
		font-weight: 600;
		transition: background-color 0.2s;
	}

	.create-deck-link:hover {
		background-color: #5568d3;
	}

	.form-group {
		margin-bottom: 1.25rem;
	}

	.form-label {
		display: block;
		font-weight: 600;
		color: #374151;
		margin-bottom: 0.5rem;
		font-size: 0.875rem;
	}

	.form-select,
	.form-input {
		width: 100%;
		padding: 0.625rem 0.875rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		transition: all 0.2s;
	}

	.form-select:focus,
	.form-input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.form-select:disabled,
	.form-input:disabled {
		background-color: #f9fafb;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.error-message {
		background-color: #fef2f2;
		border: 1px solid #fecaca;
		color: #dc2626;
		padding: 0.75rem 1rem;
		border-radius: 0.375rem;
		margin-bottom: 1.25rem;
		font-size: 0.875rem;
	}

	.modal-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		margin-top: 1.5rem;
	}

	.btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 0.375rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.btn-secondary {
		background-color: white;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.btn-secondary:hover:not(:disabled) {
		background-color: #f9fafb;
	}

	.btn-primary {
		background-color: #667eea;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background-color: #5568d3;
	}
</style>
