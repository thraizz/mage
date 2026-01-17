<script lang="ts">
	import ManaSymbol from '$lib/components/mtg/ManaSymbol.svelte';

	interface Props {
		onCreateToken: (
			name: string,
			types: string,
			power: string,
			toughness: string,
			color: string
		) => void;
		onClose: () => void;
	}

	let { onCreateToken, onClose }: Props = $props();

	// Token state
	let tokenName = $state('');
	let tokenTypes = $state('Creature');
	let tokenSubtypes = $state('');
	let tokenPower = $state('1');
	let tokenToughness = $state('1');
	let tokenColor = $state('colorless');
	let tokenAbilities = $state('');

	// Common token presets for quick creation
	const commonTokens = [
		{ name: 'Soldier', types: 'Creature — Soldier', power: '1', toughness: '1', color: 'white' },
		{ name: 'Goblin', types: 'Creature — Goblin', power: '1', toughness: '1', color: 'red' },
		{ name: 'Zombie', types: 'Creature — Zombie', power: '2', toughness: '2', color: 'black' },
		{ name: 'Elf Warrior', types: 'Creature — Elf Warrior', power: '1', toughness: '1', color: 'green' },
		{ name: 'Bird', types: 'Creature — Bird', power: '1', toughness: '1', color: 'blue' },
		{ name: 'Treasure', types: 'Artifact — Treasure', power: '', toughness: '', color: 'colorless' },
		{ name: 'Food', types: 'Artifact — Food', power: '', toughness: '', color: 'colorless' },
		{ name: 'Clue', types: 'Artifact — Clue', power: '', toughness: '', color: 'colorless' }
	];

	// Card types
	const cardTypes = ['Creature', 'Artifact', 'Enchantment', 'Planeswalker'];

	// Colors
	const colors = [
		{ value: 'white', label: 'White', symbol: 'W' as const },
		{ value: 'blue', label: 'Blue', symbol: 'U' as const },
		{ value: 'black', label: 'Black', symbol: 'B' as const },
		{ value: 'red', label: 'Red', symbol: 'R' as const },
		{ value: 'green', label: 'Green', symbol: 'G' as const },
		{ value: 'colorless', label: 'Colorless', symbol: 'C' as const },
		{ value: 'multicolor', label: 'Multicolor', symbol: null }
	];

	function loadPreset(preset: typeof commonTokens[0]) {
		tokenName = preset.name;
		tokenTypes = preset.types.split(' — ')[0];
		tokenSubtypes = preset.types.includes(' — ') ? preset.types.split(' — ')[1] : '';
		tokenPower = preset.power;
		tokenToughness = preset.toughness;
		tokenColor = preset.color;
	}

	function handleCreate() {
		if (!tokenName.trim()) return;

		const fullTypes = tokenSubtypes.trim()
			? `${tokenTypes} — ${tokenSubtypes.trim()}`
			: tokenTypes;

		onCreateToken(
			tokenName.trim(),
			fullTypes,
			tokenPower || '0',
			tokenToughness || '0',
			tokenColor
		);

		// Reset form
		tokenName = '';
		tokenTypes = 'Creature';
		tokenSubtypes = '';
		tokenPower = '1';
		tokenToughness = '1';
		tokenColor = 'colorless';
		tokenAbilities = '';
	}

	const isCreature = $derived(tokenTypes === 'Creature');
</script>

<div class="token-overlay" role="dialog" aria-labelledby="token-dialog-title">
	<div class="token-dialog">
		<div class="dialog-header">
			<h2 id="token-dialog-title">Create Token</h2>
			<button class="close-button" onclick={onClose} aria-label="Close dialog">×</button>
		</div>

		<!-- Common Token Presets -->
		<div class="presets-section">
			<h3>Common Tokens</h3>
			<div class="preset-buttons">
				{#each commonTokens as preset}
					<button class="preset-btn" onclick={() => loadPreset(preset)}>
						{preset.name}
					</button>
				{/each}
			</div>
		</div>

		<!-- Custom Token Form -->
		<div class="form-section">
			<h3>Custom Token</h3>

			<div class="form-group">
				<label for="token-name">Token Name *</label>
				<input
					id="token-name"
					type="text"
					bind:value={tokenName}
					placeholder="e.g., Soldier, Zombie, Treasure"
					class="form-input"
				/>
			</div>

			<div class="form-row">
				<div class="form-group">
					<label for="token-type">Card Type *</label>
					<select id="token-type" bind:value={tokenTypes} class="form-select">
						{#each cardTypes as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>

				<div class="form-group">
					<label for="token-subtypes">Subtypes</label>
					<input
						id="token-subtypes"
						type="text"
						bind:value={tokenSubtypes}
						placeholder="e.g., Soldier, Elf Warrior"
						class="form-input"
					/>
				</div>
			</div>

			{#if isCreature}
				<div class="form-row">
					<div class="form-group">
						<label for="token-power">Power</label>
						<input
							id="token-power"
							type="text"
							bind:value={tokenPower}
							placeholder="1"
							class="form-input"
						/>
					</div>

					<div class="form-group">
						<label for="token-toughness">Toughness</label>
						<input
							id="token-toughness"
							type="text"
							bind:value={tokenToughness}
							placeholder="1"
							class="form-input"
						/>
					</div>
				</div>
			{/if}

			<div class="form-group">
				<label for="token-color">Color *</label>
				<div class="color-grid">
					{#each colors as color}
						<button
							class="color-btn {tokenColor === color.value ? 'active' : ''}"
							onclick={() => (tokenColor = color.value)}
						>
							<span class="color-symbol">
								{#if color.symbol}
									<ManaSymbol symbol={color.symbol} size="lg" />
								{:else}
									<!-- Multicolor - show all 5 colors -->
									<span class="multicolor-symbols">
										<ManaSymbol symbol="W" size="sm" />
										<ManaSymbol symbol="U" size="sm" />
										<ManaSymbol symbol="B" size="sm" />
										<ManaSymbol symbol="R" size="sm" />
										<ManaSymbol symbol="G" size="sm" />
									</span>
								{/if}
							</span>
							<span class="color-label">{color.label}</span>
						</button>
					{/each}
				</div>
			</div>
		</div>

		<div class="dialog-footer">
			<button class="btn-secondary" onclick={onClose}>Cancel</button>
			<button class="btn-primary" onclick={handleCreate} disabled={!tokenName.trim()}>
				Create Token
			</button>
		</div>
	</div>
</div>

<style>
	.token-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		backdrop-filter: blur(4px);
		animation: fadeIn 0.2s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.token-dialog {
		background: #1a1f2e;
		border: 2px solid #3a4451;
		border-radius: 12px;
		padding: 1.5rem;
		max-width: 600px;
		width: 90%;
		max-height: 85vh;
		overflow-y: auto;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
		animation: slideUp 0.3s ease-out;
	}

	@keyframes slideUp {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	.dialog-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid #3a4451;
	}

	.dialog-header h2 {
		margin: 0;
		color: #fff;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.close-button {
		background: transparent;
		border: none;
		color: #9ca3af;
		font-size: 2rem;
		line-height: 1;
		cursor: pointer;
		padding: 0;
		width: 2rem;
		height: 2rem;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: color 0.2s;
	}

	.close-button:hover {
		color: #fff;
	}

	.presets-section,
	.form-section {
		margin-bottom: 1.5rem;
	}

	.presets-section h3,
	.form-section h3 {
		font-size: 0.875rem;
		font-weight: 600;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-bottom: 0.75rem;
	}

	.preset-buttons {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 0.5rem;
	}

	.preset-btn {
		padding: 0.5rem;
		background: rgba(102, 126, 234, 0.1);
		border: 1px solid #667eea;
		border-radius: 4px;
		color: #667eea;
		font-size: 0.75rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.preset-btn:hover {
		background: rgba(102, 126, 234, 0.2);
		transform: translateY(-1px);
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	label {
		display: block;
		color: #9ca3af;
		font-size: 0.875rem;
		font-weight: 500;
		margin-bottom: 0.375rem;
	}

	.form-input,
	.form-select {
		width: 100%;
		padding: 0.5rem;
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid #3a4451;
		border-radius: 6px;
		color: #fff;
		font-size: 0.875rem;
		transition: border-color 0.2s;
	}

	.form-input:focus,
	.form-select:focus {
		outline: none;
		border-color: #667eea;
	}

	.form-select {
		cursor: pointer;
	}

	.color-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
		gap: 0.5rem;
	}

	.color-btn {
		padding: 0.625rem;
		background: rgba(255, 255, 255, 0.05);
		border: 2px solid transparent;
		border-radius: 6px;
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.25rem;
		transition: all 0.2s;
	}

	.color-btn.active {
		border-color: #667eea;
		background: rgba(102, 126, 234, 0.15);
	}

	.color-btn:hover:not(.active) {
		background: rgba(255, 255, 255, 0.08);
	}

	.color-symbol {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.multicolor-symbols {
		display: flex;
		gap: 0.125rem;
		align-items: center;
		justify-content: center;
		flex-wrap: wrap;
	}

	.color-label {
		font-size: 0.6875rem;
		font-weight: 500;
		color: #9ca3af;
	}

	.dialog-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		padding-top: 1rem;
		border-top: 1px solid #3a4451;
	}

	.btn-primary,
	.btn-secondary {
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		font-size: 0.875rem;
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
		transform: translateY(-1px);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: rgba(255, 255, 255, 0.1);
		color: #fff;
		border: 1px solid #3a4451;
	}

	.btn-secondary:hover {
		background: rgba(255, 255, 255, 0.15);
	}
</style>
