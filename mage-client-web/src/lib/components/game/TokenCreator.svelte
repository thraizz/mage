<script lang="ts">
	import { createToken } from '$lib/api/direct-actions';
	import { toast } from '$lib/stores/toast';

	interface Props {
		gameId: string;
		onClose: () => void;
	}

	let { gameId, onClose }: Props = $props();

	// Token properties
	let name = $state('');
	let types = $state('Creature');
	let power = $state('1');
	let toughness = $state('1');
	let color = $state('colorless');
	let customAbilities = $state('');

	// Common token presets
	const presets = [
		{ name: 'Soldier', types: 'Creature — Soldier', power: '1', toughness: '1', color: 'white' },
		{ name: 'Spirit', types: 'Creature — Spirit', power: '1', toughness: '1', color: 'white' },
		{ name: 'Zombie', types: 'Creature — Zombie', power: '2', toughness: '2', color: 'black' },
		{ name: 'Goblin', types: 'Creature — Goblin', power: '1', toughness: '1', color: 'red' },
		{ name: 'Saproling', types: 'Creature — Saproling', power: '1', toughness: '1', color: 'green' },
		{ name: 'Beast', types: 'Creature — Beast', power: '3', toughness: '3', color: 'green' },
		{ name: 'Angel', types: 'Creature — Angel', power: '4', toughness: '4', color: 'white' },
		{ name: 'Dragon', types: 'Creature — Dragon', power: '5', toughness: '5', color: 'red' },
		{ name: 'Elemental', types: 'Creature — Elemental', power: '*', toughness: '*', color: 'green' },
		{ name: 'Treasure', types: 'Artifact — Treasure', power: '', toughness: '', color: 'colorless' },
		{ name: 'Clue', types: 'Artifact — Clue', power: '', toughness: '', color: 'colorless' },
		{ name: 'Food', types: 'Artifact — Food', power: '', toughness: '', color: 'colorless' },
		{ name: 'Copy', types: 'Copy of any permanent', power: '*', toughness: '*', color: 'colorless' }
	];

	const colorOptions = [
		{ value: 'colorless', label: 'Colorless', class: 'colorless' },
		{ value: 'white', label: 'White', class: 'white' },
		{ value: 'blue', label: 'Blue', class: 'blue' },
		{ value: 'black', label: 'Black', class: 'black' },
		{ value: 'red', label: 'Red', class: 'red' },
		{ value: 'green', label: 'Green', class: 'green' },
		{ value: 'multicolor', label: 'Multicolor', class: 'multicolor' }
	];

	function applyPreset(preset: (typeof presets)[0]) {
		name = preset.name;
		types = preset.types;
		power = preset.power;
		toughness = preset.toughness;
		color = preset.color;
	}

	async function handleCreate() {
		if (!name.trim()) {
			toast.error('Please enter a token name');
			return;
		}

		try {
			const abilities = customAbilities
				.split('\n')
				.map((a) => a.trim())
				.filter((a) => a.length > 0);

		await createToken(gameId, name, types, power, toughness, color, abilities, 1);

		toast.success(`Created ${name} token`);
			onClose();
		} catch (error) {
			toast.error(`Failed to create token: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		} else if (event.key === 'Enter' && event.ctrlKey) {
			handleCreate();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="token-creator-overlay" role="dialog" aria-modal="true">
	<div class="token-creator">
		<div class="creator-header">
			<h2>Create Token</h2>
			<button class="close-btn" onclick={onClose} aria-label="Close">×</button>
		</div>

		<div class="creator-body">
			<!-- Presets Section -->
			<div class="presets-section">
				<h3>Quick Presets</h3>
				<div class="presets-grid">
					{#each presets as preset}
						<button class="preset-btn {preset.color}" onclick={() => applyPreset(preset)}>
							{preset.name}
						</button>
					{/each}
				</div>
			</div>

			<!-- Custom Token Form -->
			<div class="form-section">
				<h3>Token Details</h3>

				<div class="form-row">
					<label for="token-name">Name</label>
					<input
						id="token-name"
						type="text"
						bind:value={name}
						placeholder="Token name"
						class="form-input"
					/>
				</div>

				<div class="form-row">
					<label for="token-types">Types</label>
					<input
						id="token-types"
						type="text"
						bind:value={types}
						placeholder="e.g., Creature — Zombie"
						class="form-input"
					/>
				</div>

				<div class="form-row two-col">
					<div>
						<label for="token-power">Power</label>
						<input
							id="token-power"
							type="text"
							bind:value={power}
							placeholder="*"
							class="form-input"
						/>
					</div>
					<div>
						<label for="token-toughness">Toughness</label>
						<input
							id="token-toughness"
							type="text"
							bind:value={toughness}
							placeholder="*"
							class="form-input"
						/>
					</div>
				</div>

				<div class="form-row">
					<label>Color</label>
					<div class="color-options">
						{#each colorOptions as option}
							<button
								class="color-btn {option.class}"
								class:selected={color === option.value}
								onclick={() => (color = option.value)}
								title={option.label}
							>
								{option.label.charAt(0)}
							</button>
						{/each}
					</div>
				</div>

				<div class="form-row">
					<label for="token-abilities">Abilities (one per line)</label>
					<textarea
						id="token-abilities"
						bind:value={customAbilities}
						placeholder="Flying&#10;Haste&#10;When this creature dies, draw a card."
						class="form-textarea"
						rows="3"
					></textarea>
				</div>
			</div>

			<!-- Token Preview -->
			<div class="preview-section">
				<h3>Preview</h3>
				<div class="token-preview {color}">
					<div class="preview-name">{name || 'Token'}</div>
					<div class="preview-types">{types || 'Creature'}</div>
					{#if power || toughness}
						<div class="preview-pt">{power || '*'}/{toughness || '*'}</div>
					{/if}
				</div>
			</div>
		</div>

		<div class="creator-footer">
			<button class="cancel-btn" onclick={onClose}>Cancel</button>
			<button class="create-btn" onclick={handleCreate}>Create Token</button>
		</div>
	</div>
</div>

<style>
	.token-creator-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.token-creator {
		background: var(--surface-1, #1a1a2e);
		border: 1px solid var(--border-color, #444);
		border-radius: 12px;
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	.creator-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 16px 20px;
		background: var(--surface-2, #252540);
		border-bottom: 1px solid var(--border-color, #333);
	}

	.creator-header h2 {
		margin: 0;
		font-size: 18px;
		font-weight: 600;
		color: var(--text-color, #fff);
	}

	.close-btn {
		background: none;
		border: none;
		color: var(--text-muted, #888);
		font-size: 24px;
		cursor: pointer;
		padding: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 6px;
	}

	.close-btn:hover {
		background: var(--surface-3, #353550);
		color: var(--text-color, #fff);
	}

	.creator-body {
		padding: 20px;
		overflow-y: auto;
		display: grid;
		gap: 20px;
	}

	h3 {
		margin: 0 0 12px;
		font-size: 14px;
		font-weight: 600;
		color: var(--text-color, #fff);
	}

	.presets-section {
		padding-bottom: 16px;
		border-bottom: 1px solid var(--border-color, #333);
	}

	.presets-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
		gap: 8px;
	}

	.preset-btn {
		padding: 8px 12px;
		border: 1px solid var(--border-color, #333);
		border-radius: 6px;
		background: var(--surface-2, #252540);
		color: var(--text-color, #fff);
		font-size: 12px;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.preset-btn:hover {
		border-color: var(--accent-color, #00d4ff);
	}

	.preset-btn.white {
		border-left: 3px solid #f8f8dc;
	}
	.preset-btn.blue {
		border-left: 3px solid #0077ff;
	}
	.preset-btn.black {
		border-left: 3px solid #555;
	}
	.preset-btn.red {
		border-left: 3px solid #ff4444;
	}
	.preset-btn.green {
		border-left: 3px solid #00cc66;
	}
	.preset-btn.colorless {
		border-left: 3px solid #888;
	}

	.form-section {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.form-row {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.form-row.two-col {
		flex-direction: row;
		gap: 12px;
	}

	.form-row.two-col > div {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	label {
		font-size: 12px;
		color: var(--text-muted, #888);
		font-weight: 500;
	}

	.form-input,
	.form-textarea {
		padding: 10px 12px;
		border: 1px solid var(--border-color, #333);
		border-radius: 6px;
		background: var(--surface-2, #252540);
		color: var(--text-color, #fff);
		font-size: 14px;
	}

	.form-input:focus,
	.form-textarea:focus {
		outline: none;
		border-color: var(--accent-color, #00d4ff);
	}

	.form-textarea {
		resize: vertical;
		min-height: 80px;
		font-family: inherit;
	}

	.color-options {
		display: flex;
		gap: 6px;
	}

	.color-btn {
		width: 36px;
		height: 36px;
		border: 2px solid transparent;
		border-radius: 6px;
		cursor: pointer;
		font-weight: 700;
		font-size: 14px;
		transition: all 0.15s ease;
	}

	.color-btn.selected {
		border-color: var(--accent-color, #00d4ff);
		transform: scale(1.1);
	}

	.color-btn.colorless {
		background: linear-gradient(135deg, #888, #666);
		color: #fff;
	}
	.color-btn.white {
		background: linear-gradient(135deg, #f8f8dc, #ffffd0);
		color: #333;
	}
	.color-btn.blue {
		background: linear-gradient(135deg, #0077ff, #0044cc);
		color: #fff;
	}
	.color-btn.black {
		background: linear-gradient(135deg, #444, #222);
		color: #fff;
	}
	.color-btn.red {
		background: linear-gradient(135deg, #ff4444, #cc2222);
		color: #fff;
	}
	.color-btn.green {
		background: linear-gradient(135deg, #00cc66, #008844);
		color: #fff;
	}
	.color-btn.multicolor {
		background: linear-gradient(135deg, #ffd700, #ff8c00, #ff69b4, #8a2be2);
		color: #fff;
	}

	.preview-section {
		padding-top: 16px;
		border-top: 1px solid var(--border-color, #333);
	}

	.token-preview {
		background: var(--surface-2, #252540);
		border: 2px solid var(--border-color, #333);
		border-radius: 12px;
		padding: 16px;
		text-align: center;
		min-height: 100px;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 4px;
	}

	.token-preview.white {
		border-color: #f8f8dc;
		background: linear-gradient(135deg, rgba(248, 248, 220, 0.2), transparent);
	}
	.token-preview.blue {
		border-color: #0077ff;
		background: linear-gradient(135deg, rgba(0, 119, 255, 0.2), transparent);
	}
	.token-preview.black {
		border-color: #555;
		background: linear-gradient(135deg, rgba(85, 85, 85, 0.2), transparent);
	}
	.token-preview.red {
		border-color: #ff4444;
		background: linear-gradient(135deg, rgba(255, 68, 68, 0.2), transparent);
	}
	.token-preview.green {
		border-color: #00cc66;
		background: linear-gradient(135deg, rgba(0, 204, 102, 0.2), transparent);
	}
	.token-preview.multicolor {
		border-color: #ffd700;
		background: linear-gradient(135deg, rgba(255, 215, 0, 0.2), transparent);
	}

	.preview-name {
		font-size: 18px;
		font-weight: 700;
		color: var(--text-color, #fff);
	}

	.preview-types {
		font-size: 12px;
		color: var(--text-muted, #888);
	}

	.preview-pt {
		font-size: 16px;
		font-weight: 700;
		color: var(--text-color, #fff);
		margin-top: 8px;
	}

	.creator-footer {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
		padding: 16px 20px;
		background: var(--surface-2, #252540);
		border-top: 1px solid var(--border-color, #333);
	}

	.cancel-btn,
	.create-btn {
		padding: 10px 20px;
		border: none;
		border-radius: 6px;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.cancel-btn {
		background: var(--surface-3, #353550);
		color: var(--text-color, #fff);
	}

	.cancel-btn:hover {
		background: var(--surface-4, #454560);
	}

	.create-btn {
		background: var(--accent-color, #00d4ff);
		color: #000;
	}

	.create-btn:hover {
		background: var(--accent-hover, #33ddff);
	}
</style>

