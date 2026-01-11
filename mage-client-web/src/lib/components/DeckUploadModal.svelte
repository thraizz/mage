<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { uploadDeck } from '$lib/api/decks';
	import type { DeckUploadRequest } from '$lib/types/deck';
	import { toast } from '$lib/stores/toast';
	import CircleAlert from '@lucide/svelte/icons/circle-alert';
	import Plus from '@lucide/svelte/icons/plus';
	import X from '@lucide/svelte/icons/x';
	import ArrowUp from '@lucide/svelte/icons/arrow-up';
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import Crown from '@lucide/svelte/icons/crown';
	import TextCursorInput from '@lucide/svelte/icons/text-cursor-input';
	import Rows3 from '@lucide/svelte/icons/rows-3';
	import BookOpen from '@lucide/svelte/icons/book-open';
	import Eraser from '@lucide/svelte/icons/eraser';
	import Save from '@lucide/svelte/icons/save';
	import {
		parseStructuredCards,
		structuredCardsToText,
		parseDeckList,
		validateDeck as validateDeckUtil,
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
	$: hasCommander =
		selectedFormat === 'Commander' && structuredCards.some((c) => c.section === 'commander' && c.quantity > 0);

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

	function removeCard(_index: number) {
		structuredCards.splice(_index, 1);
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
							{#if viewMode === 'text'}
								<Rows3 class="icon" size={14} aria-hidden="true" />
								Structured
							{:else}
								<TextCursorInput class="icon" size={14} aria-hidden="true" />
								Text
							{/if}
						</button>
						<button type="button" class="example-link" on:click={loadExample} disabled={loading}>
							<BookOpen class="icon" size={14} aria-hidden="true" />
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
										<Plus class="icon" size={14} aria-hidden="true" />
										Add
									</button>
								</div>
								<div class="card-list">
									{#each structuredCards.filter((c) => c.section === 'commander') as card}
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
												disabled={loading}
												aria-label="Remove card"
												title="Remove"
											>
												<X class="icon" size={16} aria-hidden="true" />
											</button
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
									<Plus class="icon" size={14} aria-hidden="true" />
									Add
								</button>
							</div>
							<div class="card-list">
								{#each structuredCards.filter((c) => c.section === 'main') as card}
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
											{#if !hasCommander}
												<button
													type="button"
													class="btn-move btn-move-commander"
													on:click={() => moveCard(globalIndex, 'commander')}
													disabled={loading}
													title="Set as Commander"
													aria-label="Set as Commander"
												>
													<Crown class="icon crown-icon" size={16} aria-hidden="true" />
												</button>
											{/if}
										{/if}
										{#if selectedFormat !== 'Commander'}
											<button
												type="button"
												class="btn-move"
												on:click={() => moveCard(globalIndex, 'sideboard')}
												disabled={loading}
												title="Move to Sideboard"
												aria-label="Move to Sideboard"
											>
												<ArrowDown class="icon" size={16} aria-hidden="true" />
											</button
											>
										{/if}
										<button
											type="button"
											class="btn-remove"
											on:click={() => removeCard(globalIndex)}
											disabled={loading}
											aria-label="Remove card"
											title="Remove"
										>
											<X class="icon" size={16} aria-hidden="true" />
										</button
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
										<Plus class="icon" size={14} aria-hidden="true" />
										Add
									</button>
								</div>
								<div class="card-list">
									{#each structuredCards.filter((c) => c.section === 'sideboard') as card}
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
												title="Move to Main Deck"
												aria-label="Move to Main Deck"
											>
												<ArrowUp class="icon" size={16} aria-hidden="true" />
											</button
											>
											<button
												type="button"
												class="btn-remove"
												on:click={() => removeCard(globalIndex)}
												disabled={loading}
												aria-label="Remove card"
												title="Remove"
											>
												<X class="icon" size={16} aria-hidden="true" />
											</button
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
						<CircleAlert class="error-icon" size={20} aria-hidden="true" />
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
					<Eraser class="icon" size={16} aria-hidden="true" />
					Clear
				</button>

				<div class="actions-right">
					<button type="button" class="btn-secondary" on:click={handleClose} disabled={loading}>
						<X class="icon" size={16} aria-hidden="true" />
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
							<Save class="icon" size={16} aria-hidden="true" />
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
		padding: var(--space-6);
	}

	.modal-header {
		margin-bottom: 1.5rem;
	}

	.modal-header h2 {
		margin: 0 0 0.5rem 0;
		font-size: var(--text-2xl);
		font-weight: var(--weight-bold);
		color: var(--text-bright);
	}

	.subtitle {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--text-muted);
	}

	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.label-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.label-row-actions {
		display: flex;
		gap: var(--space-3);
		align-items: center;
	}

	label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	.required {
		color: var(--status-error);
	}

	input,
	select,
	textarea {
		padding: var(--space-2) var(--space-3);
		border: 1px solid var(--input-border);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-family: var(--font-body);
		color: var(--text-bright);
		background: var(--input-bg);
		transition: all var(--transition-fast);
	}

	input:focus,
	select:focus,
	textarea:focus {
		outline: none;
		border-color: var(--input-focus-border);
		box-shadow: 0 0 0 3px var(--input-focus-ring);
	}

	input:disabled,
	select:disabled,
	textarea:disabled {
		background: var(--bg-slate);
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
		color: var(--accent-gold);
		font-size: var(--text-xs);
		cursor: pointer;
		padding: 0;
		text-decoration: underline;
	}

	.example-link :global(svg.icon) {
		margin-right: var(--space-1);
	}

	.example-link:hover:not(:disabled) {
		color: var(--accent-gold-bright);
	}

	.example-link:disabled {
		color: var(--text-ghost);
		cursor: not-allowed;
		text-decoration: none;
	}

	.view-toggle {
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		color: var(--text-muted);
		font-size: var(--text-xs);
		cursor: pointer;
		padding: var(--space-1) var(--space-2);
		transition: all var(--transition-fast);
	}

	.view-toggle :global(svg.icon) {
		color: currentColor;
	}

	.view-toggle:hover:not(:disabled) {
		background: var(--bg-steel);
		border-color: var(--border-strong);
		color: var(--text-bright);
	}

	.view-toggle:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.structured-editor {
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		background: var(--bg-obsidian);
		max-height: 24rem;
		overflow-y: auto;
	}

	.card-section {
		border-bottom: 1px solid var(--border-subtle);
	}

	.card-section:last-child {
		border-bottom: none;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-3) var(--space-4);
		background: var(--bg-slate);
		border-bottom: 1px solid var(--border-subtle);
		position: sticky;
		top: 0;
		z-index: 1;
	}

	.section-header h4 {
		margin: 0;
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	.btn-add {
		background: var(--accent-gold);
		color: var(--bg-void);
		border: none;
		border-radius: var(--radius-sm);
		padding: var(--space-1) var(--space-2);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-add :global(svg.icon) {
		color: currentColor;
	}

	.btn-add:hover:not(:disabled) {
		background: var(--accent-gold-bright);
		box-shadow: var(--shadow-glow);
	}

	.btn-add:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.card-list {
		padding: var(--space-2);
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.card-item {
		display: flex;
		gap: var(--space-2);
		align-items: center;
		padding: var(--space-2);
		background: var(--bg-iron);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-sm);
		transition: background var(--transition-fast), border-color var(--transition-fast);
	}

	.card-item:hover {
		background: var(--bg-steel);
		border-color: var(--border-default);
	}

	.card-quantity {
		width: 3.5rem;
		padding: var(--space-2);
		border: 1px solid var(--input-border);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
		text-align: center;
	}

	.card-name {
		flex: 1;
		padding: var(--space-2) var(--space-2);
		border: 1px solid var(--input-border);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
		font-family: var(--font-body);
	}

	.card-name:focus,
	.card-quantity:focus {
		outline: none;
		border-color: var(--accent-gold);
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
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

	.btn-move :global(svg.icon),
	.btn-remove :global(svg.icon) {
		color: currentColor;
	}

	.btn-move-commander :global(svg.crown-icon) {
		color: var(--accent-gold);
	}

	.btn-move {
		background: var(--bg-steel);
		font-size: var(--text-sm);
	}

	.btn-move:hover:not(:disabled) {
		background: var(--border-strong);
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
		font-size: var(--text-xs);
		color: var(--text-dim);
	}

	.hint code {
		background: var(--bg-iron);
		color: var(--text-bright);
		padding: 0.125rem var(--space-1);
		border-radius: var(--radius-sm);
		font-family: monospace;
		font-size: var(--text-xs);
	}

	.stats-box {
		background: var(--bg-obsidian);
		border: 1px solid color-mix(in srgb, var(--status-info) 35%, var(--border-subtle));
		border-radius: var(--radius-lg);
		padding: var(--space-4);
	}

	.stats-box h3 {
		margin: 0 0 0.75rem 0;
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--status-info);
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: var(--space-4);
		font-size: var(--text-sm);
	}

	@media (min-width: 640px) {
		.stats-grid {
			grid-template-columns: repeat(3, 1fr);
		}
	}

	.stat {
		display: flex;
		gap: var(--space-2);
	}

	.stat-label {
		color: var(--text-muted);
	}

	.stat-value {
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	.error-box {
		background: var(--status-error-dim);
		border: 1px solid color-mix(in srgb, var(--status-error) 45%, var(--border-subtle));
		border-radius: var(--radius-lg);
		padding: var(--space-4);
	}

	.error-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 0.75rem;
	}

	:global(svg.error-icon) {
		color: var(--status-error);
		flex-shrink: 0;
	}

	.error-box h3 {
		margin: 0;
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--status-error);
	}

	.error-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 15rem;
		overflow-y: auto;
	}

	.error-item {
		padding: var(--space-2);
		background: color-mix(in srgb, var(--status-error) 12%, transparent);
		border-left: 3px solid var(--status-error);
		border-radius: var(--radius-sm);
		font-size: var(--text-sm);
		color: var(--text-bright);
		font-family: 'Monaco', 'Courier New', monospace;
	}

	.error-hint {
		margin: 0.75rem 0 0 0;
		font-size: var(--text-xs);
		color: var(--text-muted);
		font-style: italic;
	}

	.actions {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 1rem;
		border-top: 1px solid var(--border-subtle);
	}

	.actions-right {
		display: flex;
		gap: 0.75rem;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		cursor: pointer;
		transition: all var(--transition-fast);
		border: 1px solid;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-primary {
		background: var(--accent-gold);
		color: var(--bg-void);
		border-color: var(--accent-gold);
	}

	.btn-primary:hover:not(:disabled) {
		background: var(--accent-gold-bright);
		box-shadow: var(--shadow-glow);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: var(--bg-iron);
		color: var(--text-bright);
		border-color: var(--border-default);
	}

	.btn-secondary:hover:not(:disabled) {
		background: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.btn-secondary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
