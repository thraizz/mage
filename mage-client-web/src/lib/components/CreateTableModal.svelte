<script lang="ts">
	import type { GameFormat } from '$lib/types/table';
	import type { Deck } from '$lib/types/deck';
	import { getGameFormats, createTable } from '$lib/api/lobby';
	import { submitDeck } from '$lib/api/table';
	import { fetchUserDecks, getDeckDetails } from '$lib/api/decks';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import CircleAlert from '@lucide/svelte/icons/circle-alert';
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';

	// Props
	let {
		open = $bindable(false),
		onClose,
		onSuccess
	}: {
		open?: boolean;
		onClose: () => void;
		// eslint-disable-next-line no-unused-vars
		onSuccess: (tableId: string) => void;
	} = $props();

	// Available formats and player counts
	const formats = getGameFormats();
	const playerCounts = [2, 3, 4, 5, 6, 7, 8];

	// Form state
	let tableName = $state('');
	let selectedFormat = $state<GameFormat>('Commander');
	let maxPlayers = $state(4);
	let password = $state('');
	let showPassword = $state(false);
	let selectedDeck = $state<string | null>(null);

	// Deck state
	let availableDecks = $state<Deck[]>([]);
	let loadingDecks = $state(false);

	// UI state
	let loading = $state(false);
	let error = $state<string | null>(null);

	// Validation state
	let touched = $state({
		format: false,
		maxPlayers: false,
		deck: false
	});

	/**
	 * Validate form
	 */
	const isValid = $derived.by(() => {
		return selectedFormat !== null && maxPlayers >= 2 && maxPlayers <= 8 && selectedDeck !== null;
	});

	/**
	 * Load decks for selected format
	 */
	async function loadDecksForFormat(): Promise<void> {
		loadingDecks = true;
		try {
			availableDecks = await fetchUserDecks(selectedFormat);
			// Reset deck selection when format changes
			selectedDeck = availableDecks.length > 0 ? availableDecks[0].id : null;
		} catch (err) {
			console.error('Failed to load decks:', err);
			availableDecks = [];
			selectedDeck = null;
		} finally {
			loadingDecks = false;
		}
	}

	/**
	 * Reset form to defaults
	 */
	function resetForm(): void {
		tableName = '';
		selectedFormat = 'Commander';
		maxPlayers = 4;
		password = '';
		showPassword = false;
		selectedDeck = null;
		availableDecks = [];
		loading = false;
		loadingDecks = false;
		error = null;
		touched = {
			format: false,
			maxPlayers: false,
			deck: false
		};
	}

	/**
	 * Load decks when modal opens or format changes
	 */
	$effect(() => {
		if (open) {
			loadDecksForFormat();
		}
	});

	/**
	 * Reload decks when format changes
	 */
	$effect(() => {
		if (open && selectedFormat) {
			loadDecksForFormat();
		}
	});

	/**
	 * Handle form submission
	 */
	async function handleSubmit(e: Event): Promise<void> {
		e.preventDefault();

		// Mark all fields as touched
		touched = {
			format: true,
			maxPlayers: true,
			deck: true
		};

		// Validate
		if (!isValid) {
			return;
		}

		loading = true;
		error = null;

		try {
			// Find the selected deck
			const deck = availableDecks.find((d) => d.id === selectedDeck);
			if (!deck) {
				error = 'Please select a valid deck';
				loading = false;
				return;
			}

			console.log('[CreateTable] Selected deck summary:', {
				deckId: deck.id,
				deckName: deck.name
			});

			// Fetch full deck details (including card lists)
			// fetchUserDecks only returns summary info without card details
			const fullDeck = await getDeckDetails(deck.id);

			console.log('[CreateTable] Full deck loaded:', {
				deckId: fullDeck.id,
				deckName: fullDeck.name,
				mainDeckCount: fullDeck.mainDeck.length,
				sideboardCount: fullDeck.sideboard.length,
				commanderCount: fullDeck.commanders.length
			});

			const table = await createTable({
				name: tableName || undefined,
				format: selectedFormat,
				maxPlayers,
				password: password || undefined
			});

			console.log('[CreateTable] Table created:', table.id, '- Now submitting deck...');

			// Submit the deck for the creator
			try {
				await submitDeck(table.id, {
					mainDeck: fullDeck.mainDeck.map((c) => ({
						name: c.cardName,
						quantity: c.quantity
					})),
					sideboard: fullDeck.sideboard.map((c) => ({
						name: c.cardName,
						quantity: c.quantity
					})),
					commanders: fullDeck.commanders.map((c) => ({
						name: c.cardName,
						quantity: c.quantity
					}))
				});
				console.log('[CreateTable] Deck submitted successfully!');
			} catch (deckErr) {
				console.error('[CreateTable] Failed to submit deck:', deckErr);
				// Don't fail table creation if deck submit fails - user can resubmit
			}

			// Success - notify parent and close modal
			onSuccess(table.id);
			handleClose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create table';
			console.error('Failed to create table:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Handle modal close
	 */
	function handleClose(): void {
		if (!loading) {
			resetForm();
			onClose();
		}
	}

	/**
	 * Toggle password visibility
	 */
	function togglePasswordVisibility(): void {
		showPassword = !showPassword;
	}
</script>

<Modal {open} onClose={handleClose} title="Create New Table" size="medium">
	<form class="create-table-form" onsubmit={handleSubmit}>
		<!-- Table Name (Optional) -->
		<div class="form-group">
			<label for="table-name" class="form-label">
				Table Name <span class="optional-text">(optional)</span>
			</label>
			<input
				id="table-name"
				type="text"
				class="form-input"
				placeholder="e.g., Friday Night Magic"
				bind:value={tableName}
				disabled={loading}
			/>
		</div>

		<!-- Format Selector -->
		<div class="form-group">
			<label for="format" class="form-label">
				Game Format <span class="required-text">*</span>
			</label>
			<select
				id="format"
				class="form-select"
				bind:value={selectedFormat}
				disabled={loading}
				onfocus={() => (touched.format = true)}
			>
				{#each formats as format}
					<option value={format}>{format}</option>
				{/each}
			</select>
			{#if touched.format && !selectedFormat}
				<p class="error-text">Please select a format</p>
			{/if}
		</div>

		<!-- Deck Selection -->
		<div class="form-group">
			<label for="deck" class="form-label">
				Select Deck <span class="required-text">*</span>
			</label>
			{#if loadingDecks}
				<div class="loading-decks">
					<LoadingSpinner size="small" />
					<span>Loading decks...</span>
				</div>
			{:else if availableDecks.length === 0}
				<div class="no-decks-message">
					<CircleAlert size={18} aria-hidden="true" />
					<span>No {selectedFormat} decks found. Please create a deck first.</span>
				</div>
			{:else}
				<select
					id="deck"
					class="form-select"
					bind:value={selectedDeck}
					disabled={loading}
					onfocus={() => (touched.deck = true)}
				>
					{#each availableDecks as deck}
						<option value={deck.id}>
							{deck.name} ({deck.cardCount} cards)
							{#if !deck.isValid}
								- Invalid
							{/if}
						</option>
					{/each}
				</select>
				{#if touched.deck && !selectedDeck}
					<p class="error-text">Please select a deck</p>
				{/if}
			{/if}
		</div>

		<!-- Player Count -->
		<div class="form-group">
			<label for="max-players" class="form-label">
				Max Players <span class="required-text">*</span>
			</label>
			<select
				id="max-players"
				class="form-select"
				bind:value={maxPlayers}
				disabled={loading}
				onfocus={() => (touched.maxPlayers = true)}
			>
				{#each playerCounts as count}
					<option value={count}>{count} Players</option>
				{/each}
			</select>
			{#if touched.maxPlayers && (maxPlayers < 2 || maxPlayers > 8)}
				<p class="error-text">Player count must be between 2 and 8</p>
			{/if}
		</div>

		<!-- Password (Optional) -->
		<div class="form-group">
			<label for="password" class="form-label">
				Password <span class="optional-text">(optional)</span>
			</label>
			<div class="password-input-wrapper">
				<input
					id="password"
					type={showPassword ? 'text' : 'password'}
					class="form-input password-input"
					placeholder="Leave empty for public table"
					bind:value={password}
					disabled={loading}
				/>
				<button
					type="button"
					class="toggle-password"
					onclick={togglePasswordVisibility}
					disabled={loading}
					title={showPassword ? 'Hide password' : 'Show password'}
				>
					{#if showPassword}
						<!-- Eye Off Icon -->
						<EyeOff size={18} aria-hidden="true" />
					{:else}
						<!-- Eye Icon -->
						<Eye size={18} aria-hidden="true" />
					{/if}
				</button>
			</div>
			<p class="help-text">Password-protected tables will show a lock icon</p>
		</div>

		<!-- Error Message -->
		{#if error}
			<div class="error-banner">
				<CircleAlert class="error-icon" size={20} aria-hidden="true" />
				<span>{error}</span>
			</div>
		{/if}

		<!-- Form Actions -->
		<div class="form-actions">
			<button type="button" class="cancel-button" onclick={handleClose} disabled={loading}>
				Cancel
			</button>
			<button type="submit" class="submit-button" disabled={loading || !isValid}>
				{#if loading}
					<LoadingSpinner size="small" color="white" />
					<span>Creating...</span>
				{:else}
					<span>Create & Join</span>
				{/if}
			</button>
		</div>
	</form>
</Modal>

<style>
	.create-table-form {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}

	/* Form Groups */
	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.form-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--ci-scroll-parchment);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.required-text {
		color: var(--ci-mountain-ember);
	}

	.optional-text {
		color: var(--ci-swamp-obsidian);
		font-weight: var(--weight-normal);
		text-transform: none;
	}

	/* Form Inputs */
	.form-input,
	.form-select {
		width: 100%;
		padding: var(--space-3) var(--space-4);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		transition: all var(--transition-fast);
		background-color: var(--bg-iron);
		color: var(--ci-scroll-parchment);
	}

	.form-input::placeholder {
		color: var(--text-ghost);
	}

	.form-input:focus,
	.form-select:focus {
		outline: none;
		border-color: var(--ci-jace-cloak);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
		background-color: var(--bg-obsidian);
	}

	.form-input:disabled,
	.form-select:disabled {
		background-color: var(--bg-slate);
		cursor: not-allowed;
		opacity: 0.5;
	}

	.form-select {
		cursor: pointer;
		appearance: none;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%23A69AA8' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
		background-position: right 0.5rem center;
		background-repeat: no-repeat;
		background-size: 1.5em 1.5em;
		padding-right: 2.5rem;
	}

	.form-select option {
		background: var(--bg-slate);
		color: var(--ci-scroll-parchment);
	}

	/* Password Input */
	.password-input-wrapper {
		position: relative;
	}

	.password-input {
		padding-right: 2.75rem;
	}

	.toggle-password {
		position: absolute;
		right: var(--space-2);
		top: 50%;
		transform: translateY(-50%);
		background: transparent;
		border: none;
		color: var(--ci-swamp-obsidian);
		cursor: pointer;
		padding: var(--space-1);
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-sm);
		transition: all var(--transition-fast);
	}

	.toggle-password:hover:not(:disabled) {
		color: var(--ci-scroll-parchment);
		background-color: var(--bg-steel);
	}

	.toggle-password:disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	/* Help Text */
	.help-text {
		font-size: var(--text-xs);
		color: var(--ci-swamp-obsidian);
		font-style: italic;
		margin: 0;
	}

	/* Error Text */
	.error-text {
		font-size: var(--text-xs);
		color: var(--ci-mountain-ember);
		margin: 0;
	}

	/* Error Banner */
	.error-banner {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		background-color: rgba(255, 77, 77, 0.1);
		border: 1px solid var(--ci-mountain-ember);
		border-radius: var(--radius-md);
		color: var(--ci-mountain-ember);
		font-size: var(--text-sm);
	}

	:global(svg.error-icon) {
		flex-shrink: 0;
		color: var(--ci-mountain-ember);
	}

	/* Loading Decks State */
	.loading-decks {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		background-color: var(--bg-slate);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-md);
		color: var(--ci-swamp-obsidian);
		font-size: var(--text-sm);
		font-style: italic;
	}

	/* No Decks Message */
	.no-decks-message {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		background-color: rgba(245, 158, 11, 0.1);
		border: 1px solid var(--status-warning);
		border-radius: var(--radius-md);
		color: var(--status-warning);
		font-size: var(--text-sm);
		line-height: var(--leading-relaxed);
	}

	.no-decks-message :global(svg) {
		flex-shrink: 0;
		color: var(--status-warning);
	}

	/* Form Actions */
	.form-actions {
		display: flex;
		gap: var(--space-3);
		justify-content: flex-end;
		padding-top: var(--space-2);
	}

	.cancel-button,
	.submit-button {
		padding: var(--space-3) var(--space-6);
		border-radius: var(--radius-md);
		font-weight: var(--weight-semibold);
		font-size: var(--text-sm);
		cursor: pointer;
		transition: all var(--transition-base);
		border: none;
		display: flex;
		align-items: center;
		gap: var(--space-2);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.cancel-button {
		background-color: var(--bg-iron);
		color: var(--ci-scroll-parchment);
		border: 1px solid var(--border-default);
	}

	.cancel-button:hover:not(:disabled) {
		background-color: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.cancel-button:disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	.submit-button {
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
		color: var(--ci-scroll-parchment);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
	}

	.submit-button:hover:not(:disabled) {
		background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
		box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
		transform: translateY(-1px);
	}

	.submit-button:disabled {
		background: var(--bg-steel);
		box-shadow: none;
		cursor: not-allowed;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.form-actions {
			flex-direction: column-reverse;
		}

		.cancel-button,
		.submit-button {
			width: 100%;
			justify-content: center;
		}
	}
</style>
