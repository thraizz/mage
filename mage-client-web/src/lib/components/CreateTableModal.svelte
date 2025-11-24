<script lang="ts">
	import type { GameFormat } from '$lib/types/table';
	import type { Deck } from '$lib/types/deck';
	import { getGameFormats, createTable } from '$lib/api/lobby';
	import { fetchUserDecks } from '$lib/api/decks';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';

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
	const isValid = $derived(() => {
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
		if (!isValid()) {
			return;
		}

		loading = true;
		error = null;

		try {
			const table = await createTable({
				name: tableName || undefined,
				format: selectedFormat,
				maxPlayers,
				password: password || undefined
			});

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
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="12" y1="8" x2="12" y2="12"></line>
						<line x1="12" y1="16" x2="12.01" y2="16"></line>
					</svg>
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
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="18"
							height="18"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path
								d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
							></path>
							<line x1="1" y1="1" x2="23" y2="23"></line>
						</svg>
					{:else}
						<!-- Eye Icon -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="18"
							height="18"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
							<circle cx="12" cy="12" r="3"></circle>
						</svg>
					{/if}
				</button>
			</div>
			<p class="help-text">Password-protected tables will show a lock icon</p>
		</div>

		<!-- Error Message -->
		{#if error}
			<div class="error-banner">
				<svg
					class="error-icon"
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="12" r="10"></circle>
					<line x1="12" y1="8" x2="12" y2="12"></line>
					<line x1="12" y1="16" x2="12.01" y2="16"></line>
				</svg>
				<span>{error}</span>
			</div>
		{/if}

		<!-- Form Actions -->
		<div class="form-actions">
			<button type="button" class="cancel-button" onclick={handleClose} disabled={loading}>
				Cancel
			</button>
			<button type="submit" class="submit-button" disabled={loading || !isValid()}>
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
		gap: 1.5rem;
	}

	/* Form Groups */
	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #374151;
	}

	.required-text {
		color: #ef4444;
	}

	.optional-text {
		color: #9ca3af;
		font-weight: 400;
	}

	/* Form Inputs */
	.form-input,
	.form-select {
		width: 100%;
		padding: 0.625rem 0.875rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		transition: all 0.2s;
		background-color: white;
	}

	.form-input:focus,
	.form-select:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.form-input:disabled,
	.form-select:disabled {
		background-color: #f3f4f6;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.form-select {
		cursor: pointer;
		appearance: none;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
		background-position: right 0.5rem center;
		background-repeat: no-repeat;
		background-size: 1.5em 1.5em;
		padding-right: 2.5rem;
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
		right: 0.625rem;
		top: 50%;
		transform: translateY(-50%);
		background: transparent;
		border: none;
		color: #6b7280;
		cursor: pointer;
		padding: 0.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	.toggle-password:hover:not(:disabled) {
		color: #374151;
		background-color: #f3f4f6;
	}

	.toggle-password:disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	/* Help Text */
	.help-text {
		font-size: 0.75rem;
		color: #6b7280;
		margin: 0;
	}

	/* Error Text */
	.error-text {
		font-size: 0.75rem;
		color: #ef4444;
		margin: 0;
	}

	/* Error Banner */
	.error-banner {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem 1rem;
		background-color: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 0.5rem;
		color: #991b1b;
		font-size: 0.875rem;
	}

	.error-icon {
		flex-shrink: 0;
		color: #ef4444;
	}

	/* Loading Decks State */
	.loading-decks {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem 1rem;
		background-color: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		color: #6b7280;
		font-size: 0.875rem;
	}

	/* No Decks Message */
	.no-decks-message {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.875rem 1rem;
		background-color: #fffbeb;
		border: 1px solid #fde68a;
		border-radius: 0.5rem;
		color: #92400e;
		font-size: 0.875rem;
		line-height: 1.5;
	}

	.no-decks-message svg {
		flex-shrink: 0;
		color: #f59e0b;
	}

	/* Form Actions */
	.form-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		padding-top: 0.5rem;
	}

	.cancel-button,
	.submit-button {
		padding: 0.625rem 1.5rem;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.cancel-button {
		background-color: white;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.cancel-button:hover:not(:disabled) {
		background-color: #f9fafb;
		border-color: #9ca3af;
	}

	.cancel-button:disabled {
		cursor: not-allowed;
		opacity: 0.6;
	}

	.submit-button {
		background-color: #667eea;
		color: white;
	}

	.submit-button:hover:not(:disabled) {
		background-color: #5568d3;
	}

	.submit-button:disabled {
		background-color: #9ca3af;
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
