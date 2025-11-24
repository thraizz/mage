<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Modal from './Modal.svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import { uploadDeck } from '$lib/api/decks';
	import type { DeckUploadRequest } from '$lib/types/deck';
	import { toast } from '$lib/stores/toast';

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

	// Real-time stats
	$: stats = parseDeckList(deckList);

	const formats = [
		'Standard',
		'Modern',
		'Commander',
		'Legacy',
		'Vintage',
		'Pioneer',
		'Pauper',
		'Historic'
	];

	interface DeckStats {
		mainDeckCount: number;
		sideboardCount: number;
		totalCount: number;
		uniqueCards: number;
		errors: string[];
	}

	function parseDeckList(text: string): DeckStats {
		const lines = text.split('\n').map((line) => line.trim());
		let mainDeckCount = 0;
		let sideboardCount = 0;
		const uniqueMainCards = new Set<string>();
		const uniqueSideCards = new Set<string>();
		const cardQuantities = new Map<string, number>();
		const parseErrors: string[] = [];
		let inSideboard = false;
		let inCommander = false;

		for (let i = 0; i < lines.length; i++) {
			const line = lines[i];

			// Skip empty lines and comments
			if (!line || line.startsWith('#') || line.startsWith('//')) {
				continue;
			}

			// Check for section markers
			if (line.toLowerCase().includes('commander') && line.toLowerCase().includes(':')) {
				inCommander = true;
				inSideboard = false;
				continue;
			}
			if (line.toLowerCase().includes('sideboard')) {
				inSideboard = true;
				inCommander = false;
				continue;
			}

			// Parse card line: "4 Lightning Bolt" or "Lightning Bolt" or "4x Lightning Bolt"
			const match = line.match(/^(\d+)x?\s+(.+)$/i);
			if (match) {
				const quantity = parseInt(match[1]);
				const cardName = match[2].trim();

				if (quantity <= 0 || quantity > 100) {
					parseErrors.push(`Line ${i + 1}: Invalid quantity (${quantity})`);
					continue;
				}

				if (!cardName) {
					parseErrors.push(`Line ${i + 1}: Missing card name`);
					continue;
				}

				// Track quantities for 4-of validation (skip basic lands and commander zone)
				if (!inCommander && !inSideboard) {
					const lowerName = cardName.toLowerCase();
					const isBasicLand = ['plains', 'island', 'swamp', 'mountain', 'forest', 'wastes', 'snow-covered plains', 'snow-covered island', 'snow-covered swamp', 'snow-covered mountain', 'snow-covered forest'].includes(lowerName);

					if (!isBasicLand) {
						const currentQty = cardQuantities.get(lowerName) || 0;
						cardQuantities.set(lowerName, currentQty + quantity);
					}
				}

				if (inSideboard) {
					sideboardCount += quantity;
					uniqueSideCards.add(cardName.toLowerCase());
				} else if (inCommander) {
					mainDeckCount += quantity;
					uniqueMainCards.add(cardName.toLowerCase());
				} else {
					mainDeckCount += quantity;
					uniqueMainCards.add(cardName.toLowerCase());
				}
			} else if (line) {
				// Single card without quantity prefix
				const cardName = line.trim();
				const lowerName = cardName.toLowerCase();

				// Track for 4-of validation
				if (!inCommander && !inSideboard) {
					const isBasicLand = ['plains', 'island', 'swamp', 'mountain', 'forest', 'wastes', 'snow-covered plains', 'snow-covered island', 'snow-covered swamp', 'snow-covered mountain', 'snow-covered forest'].includes(lowerName);

					if (!isBasicLand) {
						const currentQty = cardQuantities.get(lowerName) || 0;
						cardQuantities.set(lowerName, currentQty + 1);
					}
				}

				if (inSideboard) {
					sideboardCount += 1;
					uniqueSideCards.add(lowerName);
				} else if (inCommander) {
					mainDeckCount += 1;
					uniqueMainCards.add(lowerName);
				} else {
					mainDeckCount += 1;
					uniqueMainCards.add(lowerName);
				}
			}
		}

		// Check for 4-of violations (except in Commander format)
		if (selectedFormat !== 'Commander') {
			cardQuantities.forEach((qty, cardName) => {
				if (qty > 4) {
					parseErrors.push(`Too many copies of "${cardName}" (${qty}/4 max)`);
				}
			});
		}

		return {
			mainDeckCount,
			sideboardCount,
			totalCount: mainDeckCount + sideboardCount,
			uniqueCards: uniqueMainCards.size + uniqueSideCards.size,
			errors: parseErrors
		};
	}

	function validateDeck(): string[] {
		const validationErrors: string[] = [];

		// Check deck name
		if (!deckName.trim()) {
			validationErrors.push('❌ Deck name is required');
		}

		// Check deck list
		if (!deckList.trim()) {
			validationErrors.push('❌ Deck list is required');
		}

		// Format-specific validation
		if (selectedFormat === 'Commander') {
			// Commander: 1 commander + 99 other cards = 100 total
			// The commander should be marked with "Commander:" section
			const hasCommanderSection = deckList.toLowerCase().includes('commander:');

			if (stats.mainDeckCount !== 100) {
				if (hasCommanderSection) {
					validationErrors.push(`⚠️ Commander deck must be exactly 100 cards total (1 commander + 99 deck) (currently ${stats.mainDeckCount})`);
				} else {
					validationErrors.push(`⚠️ Commander deck must be exactly 100 cards. Add a "Commander:" section for your commander card.`);
				}
			}

			if (stats.sideboardCount > 0) {
				validationErrors.push('⚠️ Commander decks cannot have a sideboard');
			}
		} else if (selectedFormat === 'Standard' || selectedFormat === 'Modern' || selectedFormat === 'Pioneer' || selectedFormat === 'Legacy' || selectedFormat === 'Vintage' || selectedFormat === 'Pauper') {
			if (stats.mainDeckCount < 60) {
				validationErrors.push(`⚠️ Main deck must be at least 60 cards (currently ${stats.mainDeckCount})`);
			}
			if (stats.sideboardCount > 15) {
				validationErrors.push(`⚠️ Sideboard cannot exceed 15 cards (currently ${stats.sideboardCount})`);
			}
		} else if (selectedFormat === 'Historic') {
			if (stats.mainDeckCount < 60) {
				validationErrors.push(`⚠️ Main deck must be at least 60 cards (currently ${stats.mainDeckCount})`);
			}
			if (stats.sideboardCount > 15) {
				validationErrors.push(`⚠️ Sideboard cannot exceed 15 cards (currently ${stats.sideboardCount})`);
			}
		}

		// Add parsing errors
		if (stats.errors.length > 0) {
			validationErrors.push(...stats.errors.map(err => `🔴 ${err}`));
		}

		return validationErrors;
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
		errors = [];
		dispatch('close');
	}

	function handleClear() {
		deckList = '';
		errors = [];
	}

	function loadExample() {
		// Add format-specific example based on selected format (don't replace)
		let exampleText = '';

		if (selectedFormat === 'Commander') {
			exampleText = `Commander:
1 Atraxa, Praetors' Voice

# Ramp & Fixing
1 Sol Ring
1 Arcane Signet
1 Commander's Sphere
1 Cultivate
1 Kodama's Reach
1 Farseek
1 Three Visits
1 Nature's Lore
1 Rampant Growth
1 Birds of Paradise

# Card Draw
1 Rhystic Study
1 Mystic Remora
1 Consecrated Sphinx
1 Sylvan Library
1 Esper Sentinel

# Removal
1 Swords to Plowshares
1 Path to Exile
1 Cyclonic Rift
1 Assassin's Trophy
1 Beast Within
1 Generous Gift
1 Anguished Unmaking

# Tutors
1 Enlightened Tutor
1 Mystical Tutor
1 Worldly Tutor
1 Vampiric Tutor
1 Demonic Tutor

# Win Conditions
1 Approach of the Second Sun
1 Avenger of Zendikar
1 Craterhoof Behemoth
1 Triumph of the Hordes

# Lands (37 total with above = 100 cards)
1 Command Tower
1 Exotic Orchard
1 Breeding Pool
1 Temple Garden
1 Hallowed Fountain
1 Overgrown Tomb
1 Godless Shrine
1 Watery Grave
1 Misty Rainforest
1 Windswept Heath
1 Flooded Strand
4 Forest
4 Plains
4 Island
3 Swamp
10 basic lands and utility lands to reach 100 total`;
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

		// If deck list is empty, set name and add example
		// If deck list has content, just append the example
		if (!deckList.trim()) {
			if (selectedFormat === 'Commander') {
				deckName = 'Atraxa Superfriends';
			} else if (selectedFormat === 'Modern') {
				deckName = 'Burn';
			} else {
				deckName = 'Red Deck Wins';
			}
			deckList = exampleText;
		} else {
			// Append to existing content
			deckList += '\n\n# Example cards for ' + selectedFormat + ':\n' + exampleText;
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
				<select
					id="format"
					bind:value={selectedFormat}
					disabled={loading}
				>
					{#each formats as format}
						<option value={format}>{format}</option>
					{/each}
				</select>
			</div>

			<!-- Deck List Text Area -->
			<div class="form-group">
				<div class="label-row">
					<label for="deck-list">
						Deck List <span class="required">*</span>
					</label>
					<button
						type="button"
						class="example-link"
						on:click={loadExample}
						disabled={loading}
					>
						Load Example
					</button>
				</div>
				<textarea
					id="deck-list"
					bind:value={deckList}
					placeholder="4 Lightning Bolt&#10;20 Mountain&#10;&#10;Sideboard:&#10;2 Dragon's Claw"
					rows="12"
					disabled={loading}
				></textarea>
				<p class="hint">
					Format: <code>4 Card Name</code> or <code>4x Card Name</code>
					· Use "Sideboard:" to separate sideboard cards
				</p>
			</div>

			<!-- Real-time Stats -->
			{#if deckList.trim()}
				<div class="stats-box">
					<h3>Deck Statistics</h3>
					<div class="stats-grid">
						<div class="stat">
							<span class="stat-label">Main Deck:</span>
							<span class="stat-value">{stats.mainDeckCount} cards</span>
						</div>
						<div class="stat">
							<span class="stat-label">Sideboard:</span>
							<span class="stat-value">{stats.sideboardCount} cards</span>
						</div>
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
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
						<h3>{errors.length === 1 ? '1 Validation Error' : `${errors.length} Validation Errors`}</h3>
					</div>
					<div class="error-list">
						{#each errors as error}
							<div class="error-item">
								{error}
							</div>
						{/each}
					</div>
					<p class="error-hint">
						Fix these issues before saving your deck.
					</p>
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
					<button
						type="button"
						class="btn-secondary"
						on:click={handleClose}
						disabled={loading}
					>
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
