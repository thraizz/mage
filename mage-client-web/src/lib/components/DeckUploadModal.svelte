<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { uploadDeck } from '$lib/api/decks';
	import type { DeckUploadRequest } from '$lib/types/deck';
	import { toast } from '$lib/stores/toast';
	import {
		parseStructuredCards,
		structuredCardsToText,
		parseDeckList,
		validateDeck as validateDeckUtil,
		type DeckStats,
		type CardEntry,
		type DeckFormat
	} from '$lib/utils/deck-parser';

	export let open = false;

	const dispatch = createEventDispatcher<{
		close: void;
		success: { deckId: string };
	}>();

	// Form state
	let deckName = '';
	let selectedFormat = 'Standard';
	let deckList = '';
	let loading = false;
	let errors: string[] = [];
	let viewMode: 'text' | 'structured' = 'text';

	// Real-time stats
	$: stats = parseDeckList(deckList, selectedFormat as DeckFormat);

	const formats: DeckFormat[] = [
		'Standard',
		'Modern',
		'Commander',
		'Legacy',
		'Vintage',
		'Pioneer',
		'Pauper',
		'Historic'
	];

	// Structured card data
	let structuredCards: CardEntry[] = [];

	function validateDeck(): string[] {
		return validateDeckUtil(deckName, deckList, selectedFormat as DeckFormat, stats);
	}

	async function handleSubmit() {
		// Validate
		errors = validateDeck();
		if (errors.length > 0) {
			return;
		}

		loading = true;
		try {
			const request: DeckUploadRequest = {
				name: deckName.trim(),
				format: selectedFormat,
				deckList: deckList.trim()
			};

			const deck = await uploadDeck(request);
			toast.success(`Deck "${deckName}" uploaded successfully!`);
			dispatch('success', { deckId: deck.id });
			handleClose();
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Failed to upload deck';
			toast.error(errorMessage);
			errors = [errorMessage];
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		// Reset form
		deckName = '';
		selectedFormat = 'Standard';
		deckList = '';
		structuredCards = [];
		viewMode = 'text';
		errors = [];
		dispatch('close');
	}

	function handleClear() {
		deckList = '';
		structuredCards = [];
		errors = [];
	}

	function addCard(section: 'commander' | 'main' | 'sideboard') {
		const newCard: CardEntry = { name: '', quantity: 1, section };
		structuredCards.push(newCard);
		deckList = structuredCardsToText(structuredCards);
	}

	function removeCard(index: number) {
		structuredCards.splice(index, 1);
		deckList = structuredCardsToText(structuredCards);
	}

	function updateCard(
		index: number,
		field: 'name' | 'quantity' | 'section',
		value: string | number
	) {
		if (field === 'quantity') {
			const qty = typeof value === 'number' ? value : parseInt(value as string) || 1;
			structuredCards[index].quantity = Math.max(1, Math.min(100, qty));
		} else if (field === 'section') {
			structuredCards[index].section = value as 'commander' | 'main' | 'sideboard';
		} else {
			structuredCards[index].name = value as string;
		}
		deckList = structuredCardsToText(structuredCards);
	}

	function moveCard(index: number, newSection: 'commander' | 'main' | 'sideboard') {
		structuredCards[index].section = newSection;
		deckList = structuredCardsToText(structuredCards);
	}

	function switchToStructuredView() {
		viewMode = 'structured';
		structuredCards = parseStructuredCards(deckList);
	}

	function switchToTextView() {
		viewMode = 'text';
		// Text is already synced via deckList binding
	}

	function loadExample() {
		// Clear deck list before inserting example
		deckList = '';
		structuredCards = [];

		// Add format-specific example based on selected format
		let exampleText = '';

		if (selectedFormat === 'Commander') {
			exampleText = `Commander:
1 Hearthhull, the Worldseed

1 Aftermath Analyst
1 Augur of Autumn
1 Baloth Prime
1 Beast Within
1 Binding the Old Gods
1 Blasphemous Act
1 Bojuka Bog
1 Braids, Arisen Nightmare
1 Cabaretti Courtyard
1 Canyon Slough
1 Cinder Glade
1 Command Tower
1 Cultivate
1 Dakmor Salvage
1 Escape to the Wilds
1 Escape Tunnel
1 Eumidian Hatchery
1 Eumidian Wastewaker
1 Evendo Brushrazer
1 Evolving Wilds
1 Exploration Broodship
1 Fabled Passage
1 Farseek
1 Festering Thicket
4 Forest
4 Forest
1 Formless Genesis
1 Gaze of Granite
1 Greater Gargadon
1 Harrow
1 Horizon Explorer
1 Infernal Grasp
1 Juri, Master of the Revue
1 Karplusan Forest
1 Korvold, Fae-Cursed King
1 Llanowar Wastes
1 Lotus Cobra
1 Maestros Theater
1 Manifold Key
1 Mayhem Devil
1 Mazirek, Kraul Death Priest
1 Moraug, Fury of Akoum
1 Mountain
2 Mountain
1 Mountain Valley
1 Myriad Landscape
1 Nature's Lore
1 Night's Whisper
1 Omnath, Locus of Rage
1 Oracle of Mul Daya
1 Pest Infestation
1 Planetary Annihilation
1 Putrefy
1 Rakdos Charm
1 Rampaging Baloths
1 Ramunap Excavator
1 Riveteers Overlook
1 Rocky Tar Pit
1 Roiling Regrowth
1 Satyr Wayfinder
1 Scouring Swarm
1 Sheltered Thicket
1 Smoldering Marsh
1 Sol Ring
1 Soul of Windgrace
1 Splendid Reclamation
1 Springbloom Druid
1 Sprouting Goblin
1 Sulfurous Springs
3 Swamp
2 Swamp
1 Sylvan Safekeeper
1 Szarel, Genesis Shepherd
1 Tear Asunder
1 Terramorphic Expanse
1 The Gitrog Monster
1 Tiller Engine
1 Tireless Tracker
1 Titania, Protector of Argoth
1 Twilight Mire
1 Valakut Exploration
1 Vernal Fen
1 Viridescent Bog
1 Walk-In Closet/Forgotten Cellar
1 Wastes
1 Windgrace's Judgment
1 Worldsoul's Rage
1 Zask, Skittering Swarmlord
1 Zuran Orb`;
		} else if (selectedFormat === 'Modern') {
			exampleText = `# Creatures (12)
4 Dragon's Rage Channeler
4 Ragavan, Nimble Pilferer
4 Monastery Swiftspear

# Spells (28)
4 Lightning Bolt
4 Lava Dart
4 Unholy Heat
4 Light Up the Stage
4 Mishra's Bauble
4 Expressive Iteration
4 Rift Bolt

# Lands (20)
4 Wooded Foothills
4 Scalding Tarn
4 Arid Mesa
4 Steam Vents
2 Mountain
2 Spirebluff Canal

Sideboard:
2 Blood Moon
2 Alpine Moon
3 Surgical Extraction
2 Relic of Progenitus
2 Smash to Smithereens
2 Pyroclasm
2 Anger of the Gods`;
		} else if (selectedFormat === 'Standard') {
			exampleText = `# Creatures (16)
4 Monastery Swiftspear
4 Kumano Faces Kakkazan
4 Phoenix Chick
4 Bloodthirsty Adversary

# Spells (24)
4 Play with Fire
4 Lightning Strike
4 Shock
4 Strangle
4 Reckless Impulse
4 Burn Down the House

# Lands (20)
16 Mountain
4 Den of the Bugbear

Sideboard:
3 Abrade
3 Fry
3 Tibalt's Trickery
3 Roiling Vortex
3 End the Festivities`;
		} else if (selectedFormat === 'Pauper') {
			exampleText = `# Creatures (16)
4 Monastery Swiftspear
4 Kessig Flamebreather
4 Thermo-Alchemist
4 Voldaren Epicure

# Spells (24)
4 Lightning Bolt
4 Lava Spike
4 Chain Lightning
4 Skewer the Critics
4 Fireblast
4 Needle Drop

# Lands (20)
20 Mountain

Sideboard:
4 Pyroblast
3 Electrickery
3 Smash to Smithereens
2 Flaring Pain
3 End the Festivities`;
		} else {
			// Default: Red Deck Wins for Standard/Modern/etc
			exampleText = `# Aggressive Red Deck
4 Lightning Bolt
4 Monastery Swiftspear
4 Goblin Guide
4 Eidolon of the Great Revel
4 Lava Spike
4 Rift Bolt
4 Skewer the Critics
4 Light Up the Stage
20 Mountain
4 Sunbaked Canyon

Sideboard:
3 Smash to Smithereens
3 Pyroclasm
3 Roiling Vortex
3 Searing Blood
3 Ensnaring Bridge`;
		}

		// Set deck name and example text
		if (selectedFormat === 'Commander') {
			deckName = 'Hearthhull, the Worldseed';
		} else if (selectedFormat === 'Modern') {
			deckName = 'Burn';
		} else {
			deckName = 'Red Deck Wins';
		}
		deckList = exampleText;

		// If in structured view, update the structured cards
		if (viewMode === 'structured') {
			structuredCards = parseStructuredCards(deckList);
		}
	}
</script>

<Modal {open} size="large" on:close={handleClose}>
	<div class="modal-content">
		<!-- Header -->
		<div class="modal-header">
			<h2>Upload New Deck</h2>
			<p class="subtitle">Import your deck from a text list</p>
		</div>

		<!-- Form -->
		<form on:submit|preventDefault={handleSubmit}>
			<!-- Deck Name -->
			<div class="form-group">
				<label for="deck-name">
					Deck Name <span class="required">*</span>
				</label>
				<input
					id="deck-name"
					type="text"
					bind:value={deckName}
					placeholder="My Awesome Deck"
					disabled={loading}
				/>
			</div>

			<!-- Format Selector -->
			<div class="form-group">
				<label for="format">
					Format <span class="required">*</span>
				</label>
				<select id="format" bind:value={selectedFormat} disabled={loading}>
					{#each formats as format}
						<option value={format}>{format}</option>
					{/each}
				</select>
			</div>

			<!-- Deck List Editor -->
			<div class="form-group">
				<div class="label-row">
					<label for="deck-list">
						Deck List <span class="required">*</span>
					</label>
					<div class="label-row-actions">
						<button
							type="button"
							class="view-toggle"
							on:click={() => (viewMode === 'text' ? switchToStructuredView() : switchToTextView())}
							disabled={loading}
						>
							{viewMode === 'text' ? '📋 Structured View' : '📝 Text View'}
						</button>
						<button type="button" class="example-link" on:click={loadExample} disabled={loading}>
							Load Example
						</button>
					</div>
				</div>

				{#if viewMode === 'text'}
					<textarea
						id="deck-list"
						bind:value={deckList}
						placeholder="4 Lightning Bolt&#10;20 Mountain&#10;&#10;Sideboard:&#10;2 Dragon's Claw"
						rows="12"
						disabled={loading}
					></textarea>
					<p class="hint">
						Format: <code>4 Card Name</code> or <code>4x Card Name</code>
						{#if selectedFormat === 'Commander'}
							· Use "Commander:" for your commander card
						{:else}
							· Use "Sideboard:" to separate sideboard cards
						{/if}
					</p>
				{:else}
					<!-- Structured Card Editor -->
					<div class="structured-editor">
						{#if selectedFormat === 'Commander'}
							<div class="card-section">
								<div class="section-header">
									<h4>
										Commander ({structuredCards
											.filter((c) => c.section === 'commander')
											.reduce((sum, c) => sum + c.quantity, 0)})
									</h4>
									<button
										type="button"
										class="btn-add"
										on:click={() => addCard('commander')}
										disabled={loading}
									>
										+ Add
									</button>
								</div>
								<div class="card-list">
									{#each structuredCards.filter((c) => c.section === 'commander') as card, index}
										{@const globalIndex = structuredCards.findIndex((c) => c === card)}
										<div class="card-item">
											<input
												type="number"
												min="1"
												max="100"
												bind:value={card.quantity}
												on:change={() => updateCard(globalIndex, 'quantity', card.quantity)}
												class="card-quantity"
												disabled={loading}
											/>
											<input
												type="text"
												bind:value={card.name}
												on:input={(e) => updateCard(globalIndex, 'name', e.currentTarget.value)}
												placeholder="Card name"
												class="card-name"
												disabled={loading}
											/>
											<button
												type="button"
												class="btn-remove"
												on:click={() => removeCard(globalIndex)}
												disabled={loading}>×</button
											>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						<div class="card-section">
							<div class="section-header">
								<h4>
									Main Deck ({structuredCards
										.filter((c) => c.section === 'main')
										.reduce((sum, c) => sum + c.quantity, 0)})
								</h4>
								<button
									type="button"
									class="btn-add"
									on:click={() => addCard('main')}
									disabled={loading}
								>
									+ Add
								</button>
							</div>
							<div class="card-list">
								{#each structuredCards.filter((c) => c.section === 'main') as card, index}
									{@const globalIndex = structuredCards.findIndex((c) => c === card)}
									<div class="card-item">
										<input
											type="number"
											min="1"
											max="100"
											bind:value={card.quantity}
											on:change={() => updateCard(globalIndex, 'quantity', card.quantity)}
											class="card-quantity"
											disabled={loading}
										/>
										<input
											type="text"
											bind:value={card.name}
											on:input={(e) => updateCard(globalIndex, 'name', e.currentTarget.value)}
											placeholder="Card name"
											class="card-name"
											disabled={loading}
										/>
										{#if selectedFormat === 'Commander'}
											<button
												type="button"
												class="btn-move"
												on:click={() => moveCard(globalIndex, 'commander')}
												disabled={loading}
												title="Move to Commander">↑</button
											>
										{/if}
										{#if selectedFormat !== 'Commander'}
											<button
												type="button"
												class="btn-move"
												on:click={() => moveCard(globalIndex, 'sideboard')}
												disabled={loading}
												title="Move to Sideboard">↓</button
											>
										{/if}
										<button
											type="button"
											class="btn-remove"
											on:click={() => removeCard(globalIndex)}
											disabled={loading}>×</button
										>
									</div>
								{/each}
							</div>
						</div>

						{#if selectedFormat !== 'Commander'}
							<div class="card-section">
								<div class="section-header">
									<h4>
										Sideboard ({structuredCards
											.filter((c) => c.section === 'sideboard')
											.reduce((sum, c) => sum + c.quantity, 0)})
									</h4>
									<button
										type="button"
										class="btn-add"
										on:click={() => addCard('sideboard')}
										disabled={loading}
									>
										+ Add
									</button>
								</div>
								<div class="card-list">
									{#each structuredCards.filter((c) => c.section === 'sideboard') as card, index}
										{@const globalIndex = structuredCards.findIndex((c) => c === card)}
										<div class="card-item">
											<input
												type="number"
												min="1"
												max="100"
												bind:value={card.quantity}
												on:change={() => updateCard(globalIndex, 'quantity', card.quantity)}
												class="card-quantity"
												disabled={loading}
											/>
											<input
												type="text"
												bind:value={card.name}
												on:input={(e) => updateCard(globalIndex, 'name', e.currentTarget.value)}
												placeholder="Card name"
												class="card-name"
												disabled={loading}
											/>
											<button
												type="button"
												class="btn-move"
												on:click={() => moveCard(globalIndex, 'main')}
												disabled={loading}
												title="Move to Main Deck">↑</button
											>
											<button
												type="button"
												class="btn-remove"
												on:click={() => removeCard(globalIndex)}
												disabled={loading}>×</button
											>
										</div>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				{/if}
			</div>

			<!-- Real-time Stats -->
			{#if deckList.trim()}
				<div class="stats-box">
					<h3>Deck Statistics</h3>
					<div class="stats-grid">
						{#if selectedFormat === 'Commander' && stats.commanderCount > 0}
							<div class="stat">
								<span class="stat-label">Commander:</span>
								<span class="stat-value"
									>{stats.commanderCount} card{stats.commanderCount !== 1 ? 's' : ''}</span
								>
							</div>
						{/if}
						<div class="stat">
							<span class="stat-label">Main Deck:</span>
							<span class="stat-value">{stats.mainDeckCount} cards</span>
						</div>
						{#if selectedFormat !== 'Commander'}
							<div class="stat">
								<span class="stat-label">Sideboard:</span>
								<span class="stat-value">{stats.sideboardCount} cards</span>
							</div>
						{/if}
						<div class="stat">
							<span class="stat-label">Total:</span>
							<span class="stat-value">{stats.totalCount} cards</span>
						</div>
						<div class="stat">
							<span class="stat-label">Unique:</span>
							<span class="stat-value">{stats.uniqueCards} cards</span>
						</div>
					</div>
				</div>
			{/if}

			<!-- Validation Errors -->
			{#if errors.length > 0}
				<div class="error-box">
					<div class="error-header">
						<svg class="error-icon" fill="currentColor" viewBox="0 0 20 20">
							<path
								fill-rule="evenodd"
								d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
								clip-rule="evenodd"
							/>
						</svg>
						<h3>
							{errors.length === 1 ? '1 Validation Error' : `${errors.length} Validation Errors`}
						</h3>
					</div>
					<div class="error-list">
						{#each errors as error}
							<div class="error-item">
								{error}
							</div>
						{/each}
					</div>
					<p class="error-hint">Fix these issues before saving your deck.</p>
				</div>
			{/if}

			<!-- Actions -->
			<div class="actions">
				<button
					type="button"
					class="btn-secondary"
					on:click={handleClear}
					disabled={loading || !deckList.trim()}
				>
					Clear
				</button>

				<div class="actions-right">
					<button type="button" class="btn-secondary" on:click={handleClose} disabled={loading}>
						Cancel
					</button>
					<button
						type="submit"
						class="btn-primary"
						disabled={loading || !deckName.trim() || !deckList.trim()}
					>
						{#if loading}
							<LoadingSpinner size="small" />
							<span>Uploading...</span>
						{:else}
							Save Deck
						{/if}
					</button>
				</div>
			</div>
		</form>
	</div>
</Modal>

<style>
	.modal-content {
		padding: 1.5rem;
	}

	.modal-header {
		margin-bottom: 1.5rem;
	}

	.modal-header h2 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		font-weight: 700;
		color: #1f2937;
	}

	.subtitle {
		margin: 0;
		font-size: 0.875rem;
		color: #6b7280;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.label-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.label-row-actions {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	label {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	.required {
		color: #ef4444;
	}

	input,
	select,
	textarea {
		padding: 0.5rem 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-family: inherit;
		transition: all 0.2s;
	}

	input:focus,
	select:focus,
	textarea:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	input:disabled,
	select:disabled,
	textarea:disabled {
		background: #f3f4f6;
		cursor: not-allowed;
		opacity: 0.6;
	}

	textarea {
		font-family: 'Monaco', 'Courier New', monospace;
		resize: vertical;
	}

	.example-link {
		background: none;
		border: none;
		color: #3b82f6;
		font-size: 0.75rem;
		cursor: pointer;
		padding: 0;
		text-decoration: underline;
	}

	.example-link:hover:not(:disabled) {
		color: #2563eb;
	}

	.example-link:disabled {
		color: #9ca3af;
		cursor: not-allowed;
		text-decoration: none;
	}

	.view-toggle {
		background: #f3f4f6;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		color: #374151;
		font-size: 0.75rem;
		cursor: pointer;
		padding: 0.25rem 0.5rem;
		transition: all 0.2s;
	}

	.view-toggle:hover:not(:disabled) {
		background: #e5e7eb;
		border-color: #9ca3af;
	}

	.view-toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.structured-editor {
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		background: white;
		max-height: 24rem;
		overflow-y: auto;
	}

	.card-section {
		border-bottom: 1px solid #e5e7eb;
	}

	.card-section:last-child {
		border-bottom: none;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		background: #f9fafb;
		border-bottom: 1px solid #e5e7eb;
		position: sticky;
		top: 0;
		z-index: 1;
	}

	.section-header h4 {
		margin: 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #374151;
	}

	.btn-add {
		background: #3b82f6;
		color: white;
		border: none;
		border-radius: 0.25rem;
		padding: 0.25rem 0.5rem;
		font-size: 0.75rem;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-add:hover:not(:disabled) {
		background: #2563eb;
	}

	.btn-add:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.card-list {
		padding: 0.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.card-item {
		display: flex;
		gap: 0.5rem;
		align-items: center;
		padding: 0.5rem;
		background: white;
		border-radius: 0.25rem;
		transition: background 0.2s;
	}

	.card-item:hover {
		background: #f9fafb;
	}

	.card-quantity {
		width: 3.5rem;
		padding: 0.375rem;
		border: 1px solid #d1d5db;
		border-radius: 0.25rem;
		font-size: 0.875rem;
		text-align: center;
	}

	.card-name {
		flex: 1;
		padding: 0.375rem 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 0.25rem;
		font-size: 0.875rem;
		font-family: inherit;
	}

	.card-name:focus,
	.card-quantity:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
	}

	.btn-move,
	.btn-remove {
		background: #ef4444;
		color: white;
		border: none;
		border-radius: 0.25rem;
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		font-size: 1rem;
		transition: background 0.2s;
		flex-shrink: 0;
	}

	.btn-move {
		background: #6b7280;
		font-size: 0.875rem;
	}

	.btn-move:hover:not(:disabled) {
		background: #4b5563;
	}

	.btn-remove:hover:not(:disabled) {
		background: #dc2626;
	}

	.btn-move:disabled,
	.btn-remove:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.hint {
		margin: 0;
		font-size: 0.75rem;
		color: #6b7280;
	}

	.hint code {
		background: #f3f4f6;
		padding: 0.125rem 0.25rem;
		border-radius: 0.25rem;
		font-family: monospace;
		font-size: 0.75rem;
	}

	.stats-box {
		background: #eff6ff;
		border: 1px solid #bfdbfe;
		border-radius: 0.5rem;
		padding: 1rem;
	}

	.stats-box h3 {
		margin: 0 0 0.75rem 0;
		font-size: 0.875rem;
		font-weight: 500;
		color: #1e3a8a;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
		font-size: 0.875rem;
	}

	@media (min-width: 640px) {
		.stats-grid {
			grid-template-columns: repeat(3, 1fr);
		}
	}

	.stat {
		display: flex;
		gap: 0.5rem;
	}

	.stat-label {
		color: #1e40af;
	}

	.stat-value {
		font-weight: 600;
		color: #1e3a8a;
	}

	.error-box {
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 0.5rem;
		padding: 1rem;
	}

	.error-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 0.75rem;
	}

	.error-icon {
		width: 1.25rem;
		height: 1.25rem;
		color: #dc2626;
		flex-shrink: 0;
	}

	.error-box h3 {
		margin: 0;
		font-size: 0.875rem;
		font-weight: 600;
		color: #7f1d1d;
	}

	.error-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 15rem;
		overflow-y: auto;
	}

	.error-item {
		padding: 0.5rem;
		background: #fee2e2;
		border-left: 3px solid #dc2626;
		border-radius: 0.25rem;
		font-size: 0.875rem;
		color: #991b1b;
		font-family: 'Monaco', 'Courier New', monospace;
	}

	.error-hint {
		margin: 0.75rem 0 0 0;
		font-size: 0.75rem;
		color: #b91c1c;
		font-style: italic;
	}

	.actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 1rem;
		border-top: 1px solid #e5e7eb;
	}

	.actions-right {
		display: flex;
		gap: 0.75rem;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		border: 1px solid;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-primary {
		background: #3b82f6;
		color: white;
		border-color: transparent;
	}

	.btn-primary:hover:not(:disabled) {
		background: #2563eb;
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: white;
		color: #374151;
		border-color: #d1d5db;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #f9fafb;
	}

	.btn-secondary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
