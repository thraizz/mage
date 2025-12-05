<script lang="ts">
	/**
	 * PriorityActionBar - Unified game control bar
	 * Contains: Status | Player Stats | Game Actions | Direct Controls | Settings
	 */
	import { modifyLife, setPlayerCounter, drawCards, shuffleLibrary, nextTurn, clearCombat } from '$lib/api/direct-actions';
	import { toast } from '$lib/stores/toast';

	// Props
	let {
		gameId = '',
		hasPriority = false,
		activePlayerId = '',
		localPlayerId = '',
		activePlayerName = 'Player',
		currentPhase = '',
		canPassPriority = true,
		isLoading = false,
		// Player stats
		playerLife = 20,
		playerPoison = 0,
		libraryCount = 0,
		// Callbacks
		onPassPriority = () => {},
		onPassUntilEOT = () => {},
		onCastSpell = () => {},
		onActivateAbility = () => {},
		onAdvancePhase = () => {},
		onCreateToken = () => {},
		// Auto-pass settings
		autoPassSettings = $bindable({
			noActions: false,
			opponentTurn: false
		})
	}: {
		gameId?: string;
		hasPriority?: boolean;
		activePlayerId?: string;
		localPlayerId?: string;
		activePlayerName?: string;
		currentPhase?: string;
		canPassPriority?: boolean;
		isLoading?: boolean;
		playerLife?: number;
		playerPoison?: number;
		libraryCount?: number;
		onPassPriority?: () => void;
		onPassUntilEOT?: () => void;
		onCastSpell?: () => void;
		onActivateAbility?: () => void;
		onAdvancePhase?: () => void;
		onCreateToken?: () => void;
		autoPassSettings?: {
			noActions: boolean;
			opponentTurn: boolean;
		};
	} = $props();

	// State
	let showSettings = $state(false);
	let showLifeMenu = $state(false);

	// Derived values
	const isYourTurn = $derived(activePlayerId === localPlayerId);
	const isInCombat = $derived(currentPhase === 'COMBAT');

	// Life/counter helpers
	async function handleLifeChange(delta: number): Promise<void> {
		if (!gameId) return;
		try {
			await modifyLife(gameId, localPlayerId, delta);
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handlePoisonChange(delta: number): Promise<void> {
		if (!gameId) return;
		const newValue = Math.max(0, playerPoison + delta);
		try {
			await setPlayerCounter(gameId, localPlayerId, 'poison', newValue);
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleDraw(count: number): Promise<void> {
		if (!gameId) return;
		try {
			await drawCards(gameId, count);
			showLifeMenu = false;
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleShuffle(): Promise<void> {
		if (!gameId) return;
		try {
			await shuffleLibrary(gameId);
			toast.success('Shuffled');
			showLifeMenu = false;
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	// Game action handlers
	function handlePassPriority(): void {
		if (!hasPriority || !canPassPriority || isLoading) return;
		onPassPriority();
	}

	function handlePassUntilEOT(): void {
		if (!hasPriority || isLoading) return;
		onPassUntilEOT();
	}

	function handleCastSpell(): void {
		if (!hasPriority || isLoading) return;
		onCastSpell();
	}

	function handleActivateAbility(): void {
		if (!hasPriority || isLoading) return;
		onActivateAbility();
	}

	function handleAdvancePhase(): void {
		if (!hasPriority || isLoading) return;
		onAdvancePhase();
	}

	async function handleNextTurn(): Promise<void> {
		if (!gameId) return;
		try {
			await nextTurn(gameId);
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	async function handleClearCombat(): Promise<void> {
		if (!gameId) return;
		try {
			await clearCombat(gameId);
			toast.success('Combat cleared');
		} catch (error) {
			toast.error(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}

	// Keyboard shortcuts
	function handleKeydown(event: KeyboardEvent): void {
		if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
			return;
		}

		if (event.key === ' ' || event.key === 'Spacebar') {
			event.preventDefault();
			handlePassPriority();
		} else if (event.key === 'c' || event.key === 'C') {
			event.preventDefault();
			handleCastSpell();
		} else if (event.key === 'a' || event.key === 'A') {
			event.preventDefault();
			handleActivateAbility();
		} else if (event.key === 'n' || event.key === 'N') {
			event.preventDefault();
			handleAdvancePhase();
		} else if (event.key === 'F6') {
			event.preventDefault();
			handlePassUntilEOT();
		} else if (event.key === 't' || event.key === 'T') {
			event.preventDefault();
			onCreateToken();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="action-bar" class:has-priority={hasPriority}>
	<!-- Priority Status -->
	<div class="status-section">
		{#if hasPriority}
			<div class="status-badge active">
				<span class="status-icon">⚡</span>
				<span class="status-text">Priority</span>
			</div>
		{:else if isYourTurn}
			<div class="status-badge waiting">
				<span class="status-icon">⏳</span>
				<span class="status-text">Wait</span>
			</div>
		{:else}
			<div class="status-badge opponent">
				<span class="status-icon">⏸</span>
				<span class="status-text">{activePlayerName}</span>
			</div>
		{/if}
	</div>

	<div class="divider"></div>

	<!-- Player Stats (Life, Poison, Library) -->
	<div class="player-stats-section">
		<div class="life-group">
			<button class="stat-btn life-btn" onclick={() => handleLifeChange(-1)} title="Lose 1 life">−</button>
			<button 
				class="stat-display life" 
				onclick={() => showLifeMenu = !showLifeMenu}
				title="Click for more options"
			>
				<span class="stat-icon">❤️</span>
				<span class="stat-value">{playerLife}</span>
			</button>
			<button class="stat-btn life-btn" onclick={() => handleLifeChange(1)} title="Gain 1 life">+</button>
		</div>

		{#if playerPoison > 0}
			<div class="stat-display poison" title="Poison counters">
				<span class="stat-icon">☠️</span>
				<span class="stat-value">{playerPoison}</span>
			</div>
		{/if}

		<div class="stat-display library" title="Cards in library">
			<span class="stat-icon">📚</span>
			<span class="stat-value">{libraryCount}</span>
		</div>

		<!-- Life/Library Quick Menu -->
		{#if showLifeMenu}
			<div class="quick-menu">
				<div class="menu-section">
					<span class="menu-label">Life</span>
					<div class="menu-row">
						<button onclick={() => handleLifeChange(-5)}>−5</button>
						<button onclick={() => handleLifeChange(-1)}>−1</button>
						<button onclick={() => handleLifeChange(1)}>+1</button>
						<button onclick={() => handleLifeChange(5)}>+5</button>
					</div>
				</div>
				<div class="menu-section">
					<span class="menu-label">Poison</span>
					<div class="menu-row">
						<button onclick={() => handlePoisonChange(-1)}>−1</button>
						<span class="menu-value">{playerPoison}</span>
						<button onclick={() => handlePoisonChange(1)}>+1</button>
					</div>
				</div>
				<div class="menu-section">
					<span class="menu-label">Library</span>
					<div class="menu-row">
						<button onclick={() => handleDraw(1)}>Draw 1</button>
						<button onclick={() => handleDraw(7)}>Draw 7</button>
						<button onclick={handleShuffle}>Shuffle</button>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<div class="divider"></div>

	<!-- Game Actions -->
	<div class="actions-group">
		<button
			class="action-btn pass"
			class:pulse={hasPriority}
			disabled={!hasPriority || !canPassPriority || isLoading}
			onclick={handlePassPriority}
			title="Pass priority (Spacebar)"
		>
			{#if isLoading}
				<span class="spinner"></span>
			{:else}
				<span class="btn-icon">→</span>
			{/if}
			<span class="btn-text">Pass</span>
			<kbd>Space</kbd>
		</button>

		<button
			class="action-btn"
			disabled={!hasPriority || isLoading}
			onclick={handleCastSpell}
			title="Cast spell (C)"
		>
			<span class="btn-icon">🎴</span>
			<span class="btn-text">Cast</span>
			<kbd>C</kbd>
		</button>

		<button
			class="action-btn"
			disabled={!hasPriority || isLoading}
			onclick={handleActivateAbility}
			title="Activate ability (A)"
		>
			<span class="btn-icon">⚡</span>
			<span class="btn-text">Activate</span>
			<kbd>A</kbd>
		</button>

		<button
			class="action-btn phase"
			disabled={!hasPriority || isLoading}
			onclick={handleAdvancePhase}
			title="Next phase (N)"
		>
			<span class="btn-icon">⏩</span>
			<span class="btn-text">Phase</span>
			<kbd>N</kbd>
		</button>
	</div>

	<div class="divider"></div>

	<!-- Direct Controls -->
	<div class="actions-group">
		{#if isInCombat}
			<button
				class="action-btn danger"
				onclick={handleClearCombat}
				title="Clear combat"
			>
				<span class="btn-icon">🛑</span>
			</button>
		{/if}

		<button
			class="action-btn"
			onclick={handleNextTurn}
			title="End turn"
		>
			<span class="btn-icon">⏭</span>
			<span class="btn-text">End Turn</span>
		</button>

		<button
			class="action-btn accent"
			onclick={onCreateToken}
			title="Create token (T)"
		>
			<span class="btn-icon">✨</span>
			<span class="btn-text">Token</span>
			<kbd>T</kbd>
		</button>

		<button
			class="action-btn muted"
			disabled={!hasPriority || isLoading}
			onclick={handlePassUntilEOT}
			title="Pass turn (F6)"
		>
			<span class="btn-icon">⏭️</span>
			<kbd>F6</kbd>
		</button>
	</div>

	<!-- Settings -->
	<button
		class="settings-btn"
		class:active={showSettings}
		onclick={() => showSettings = !showSettings}
		title="Settings"
	>
		⚙️
	</button>

	{#if showSettings}
		<div class="settings-dropdown">
			<div class="settings-title">Auto-Pass</div>
			<label class="setting-row">
				<input type="checkbox" bind:checked={autoPassSettings.noActions} />
				<span>When no actions</span>
			</label>
			<label class="setting-row">
				<input type="checkbox" bind:checked={autoPassSettings.opponentTurn} />
				<span>Opponent's turn</span>
			</label>
		</div>
	{/if}
</div>

<style>
	.action-bar {
		position: fixed;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background: rgba(18, 20, 26, 0.98);
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-bottom: none;
		border-radius: 12px 12px 0 0;
		backdrop-filter: blur(12px);
		box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.5);
		z-index: 100;
	}

	.action-bar.has-priority {
		border-color: rgba(201, 162, 39, 0.5);
		box-shadow: 
			0 -4px 24px rgba(0, 0, 0, 0.5),
			0 0 30px rgba(201, 162, 39, 0.1);
	}

	/* Status */
	.status-section {
		min-width: 90px;
	}

	.status-badge {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.625rem;
		border-radius: 6px;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.status-badge.active {
		background: rgba(201, 162, 39, 0.15);
		color: #c9a227;
		animation: status-pulse 2s ease-in-out infinite;
	}

	.status-badge.waiting {
		background: rgba(59, 130, 246, 0.15);
		color: #3b82f6;
	}

	.status-badge.opponent {
		background: rgba(113, 113, 122, 0.15);
		color: #a1a1aa;
	}

	@keyframes status-pulse {
		0%, 100% { box-shadow: 0 0 0 0 rgba(201, 162, 39, 0.3); }
		50% { box-shadow: 0 0 0 6px rgba(201, 162, 39, 0); }
	}

	.status-icon {
		font-size: 0.875rem;
	}

	.status-text {
		white-space: nowrap;
		max-width: 80px;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	/* Divider */
	.divider {
		width: 1px;
		height: 28px;
		background: rgba(63, 63, 70, 0.5);
	}

	/* Player Stats */
	.player-stats-section {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		position: relative;
	}

	.life-group {
		display: flex;
		align-items: center;
		gap: 0.125rem;
	}

	.stat-btn {
		width: 24px;
		height: 24px;
		border: 1px solid rgba(63, 63, 70, 0.4);
		border-radius: 4px;
		background: rgba(36, 40, 51, 0.6);
		color: #a1a1aa;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.stat-btn:hover {
		background: rgba(63, 63, 70, 0.6);
		color: #f4f4f5;
	}

	.stat-btn.life-btn:first-child:hover {
		background: rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}

	.stat-btn.life-btn:last-child:hover {
		background: rgba(34, 197, 94, 0.3);
		color: #22c55e;
	}

	.stat-display {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.8125rem;
		font-weight: 600;
		background: transparent;
		border: none;
		cursor: pointer;
		transition: background 0.15s ease;
	}

	.stat-display:hover {
		background: rgba(63, 63, 70, 0.3);
	}

	.stat-display.life {
		color: #f4f4f5;
	}

	.stat-display.poison {
		color: #a855f7;
	}

	.stat-display.library {
		color: #a1a1aa;
	}

	.stat-icon {
		font-size: 0.75rem;
	}

	.stat-value {
		font-family: 'JetBrains Mono', monospace;
		min-width: 20px;
		text-align: center;
	}

	/* Quick Menu */
	.quick-menu {
		position: absolute;
		bottom: calc(100% + 8px);
		left: 0;
		min-width: 200px;
		background: rgba(18, 20, 26, 0.98);
		border: 1px solid rgba(63, 63, 70, 0.6);
		border-radius: 8px;
		padding: 0.75rem;
		box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.5);
		display: flex;
		flex-direction: column;
		gap: 0.625rem;
	}

	.menu-section {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.menu-label {
		font-size: 0.625rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #71717a;
	}

	.menu-row {
		display: flex;
		gap: 0.25rem;
		align-items: center;
	}

	.menu-row button {
		padding: 0.375rem 0.5rem;
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 4px;
		background: rgba(36, 40, 51, 0.8);
		color: #a1a1aa;
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.menu-row button:hover {
		background: rgba(201, 162, 39, 0.2);
		border-color: rgba(201, 162, 39, 0.4);
		color: #f4f4f5;
	}

	.menu-value {
		min-width: 24px;
		text-align: center;
		font-weight: 600;
		color: #f4f4f5;
		font-family: 'JetBrains Mono', monospace;
	}

	/* Actions Group */
	.actions-group {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	/* Action Button */
	.action-btn {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		padding: 0.4rem 0.625rem;
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 6px;
		background: rgba(36, 40, 51, 0.8);
		color: #a1a1aa;
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.15s ease;
		white-space: nowrap;
	}

	.action-btn:hover:not(:disabled) {
		background: rgba(63, 63, 70, 0.6);
		color: #f4f4f5;
		transform: translateY(-1px);
	}

	.action-btn:disabled {
		opacity: 0.35;
		cursor: not-allowed;
	}

	.action-btn.pass {
		background: rgba(34, 197, 94, 0.2);
		border-color: rgba(34, 197, 94, 0.4);
		color: #22c55e;
	}

	.action-btn.pass:hover:not(:disabled) {
		background: rgba(34, 197, 94, 0.3);
	}

	.action-btn.pass.pulse:not(:disabled) {
		animation: pass-pulse 2s ease-in-out infinite;
	}

	@keyframes pass-pulse {
		0%, 100% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.4); }
		50% { box-shadow: 0 0 0 6px rgba(34, 197, 94, 0); }
	}

	.action-btn.phase {
		background: rgba(59, 130, 246, 0.2);
		border-color: rgba(59, 130, 246, 0.4);
		color: #3b82f6;
	}

	.action-btn.phase:hover:not(:disabled) {
		background: rgba(59, 130, 246, 0.3);
	}

	.action-btn.accent {
		background: rgba(201, 162, 39, 0.15);
		border-color: rgba(201, 162, 39, 0.3);
		color: #c9a227;
	}

	.action-btn.accent:hover:not(:disabled) {
		background: rgba(201, 162, 39, 0.25);
		color: #e4b82a;
	}

	.action-btn.danger {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}

	.action-btn.danger:hover:not(:disabled) {
		background: rgba(239, 68, 68, 0.25);
	}

	.action-btn.muted {
		background: rgba(63, 63, 70, 0.3);
		color: #71717a;
	}

	.action-btn.muted:hover:not(:disabled) {
		background: rgba(63, 63, 70, 0.5);
		color: #a1a1aa;
	}

	.btn-icon {
		font-size: 0.875rem;
	}

	.btn-text {
		font-weight: 600;
	}

	kbd {
		padding: 0.0625rem 0.25rem;
		background: rgba(0, 0, 0, 0.3);
		border-radius: 3px;
		font-size: 0.5625rem;
		font-weight: 700;
		color: #71717a;
		font-family: 'JetBrains Mono', monospace;
	}

	.action-btn.pass kbd,
	.action-btn.phase kbd,
	.action-btn.accent kbd {
		color: rgba(255, 255, 255, 0.6);
	}

	/* Spinner */
	.spinner {
		width: 12px;
		height: 12px;
		border: 2px solid rgba(255, 255, 255, 0.2);
		border-top-color: currentColor;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Settings */
	.settings-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(36, 40, 51, 0.8);
		border: 1px solid rgba(63, 63, 70, 0.5);
		border-radius: 6px;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.settings-btn:hover,
	.settings-btn.active {
		background: rgba(63, 63, 70, 0.5);
	}

	.settings-dropdown {
		position: absolute;
		bottom: calc(100% + 8px);
		right: 0.75rem;
		width: 180px;
		padding: 0.625rem;
		background: rgba(18, 20, 26, 0.98);
		border: 1px solid rgba(63, 63, 70, 0.6);
		border-radius: 8px;
		box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.5);
	}

	.settings-title {
		font-size: 0.625rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: #71717a;
		margin-bottom: 0.375rem;
		padding-bottom: 0.25rem;
		border-bottom: 1px solid rgba(63, 63, 70, 0.5);
	}

	.setting-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.25rem 0;
		font-size: 0.75rem;
		color: #a1a1aa;
		cursor: pointer;
	}

	.setting-row:hover {
		color: #f4f4f5;
	}

	.setting-row input[type="checkbox"] {
		width: 12px;
		height: 12px;
		accent-color: #c9a227;
		cursor: pointer;
	}

	/* Responsive */
	@media (max-width: 1024px) {
		.btn-text {
			display: none;
		}

		.status-text {
			display: none;
		}

		.status-section {
			min-width: auto;
		}
	}

	@media (max-width: 768px) {
		kbd {
			display: none;
		}

		.action-bar {
			gap: 0.375rem;
			padding: 0.375rem 0.5rem;
		}

		.divider {
			display: none;
		}
	}

	@media (max-width: 600px) {
		.action-bar {
			left: 0;
			right: 0;
			transform: none;
			border-radius: 0;
		}

		.stat-display.library {
			display: none;
		}
	}
</style>
