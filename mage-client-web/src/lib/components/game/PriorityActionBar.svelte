<script lang="ts">
	/**
	 * PriorityActionBar - Unified priority indicator + action buttons
	 * Docked at the bottom of the screen for clear, accessible game actions
	 */

	// Props
	let {
		hasPriority = false,
		activePlayerId = '',
		localPlayerId = '',
		activePlayerName = 'Player',
		canPassPriority = true,
		isLoading = false,
		onPassPriority = () => {},
		onPassUntilEOT = () => {},
		onCastSpell = () => {},
		onActivateAbility = () => {},
		onAdvancePhase = () => {},
		// Auto-pass settings
		autoPassSettings = $bindable({
			noActions: false,
			opponentTurn: false
		})
	}: {
		hasPriority?: boolean;
		activePlayerId?: string;
		localPlayerId?: string;
		activePlayerName?: string;
		canPassPriority?: boolean;
		isLoading?: boolean;
		onPassPriority?: () => void;
		onPassUntilEOT?: () => void;
		onCastSpell?: () => void;
		onActivateAbility?: () => void;
		onAdvancePhase?: () => void;
		autoPassSettings?: {
			noActions: boolean;
			opponentTurn: boolean;
		};
	} = $props();

	// State
	let showSettings = $state(false);

	// Derived values
	const isYourTurn = $derived(activePlayerId === localPlayerId);

	/**
	 * Handle pass priority
	 */
	function handlePassPriority(): void {
		if (!hasPriority || !canPassPriority || isLoading) return;
		onPassPriority();
	}

	/**
	 * Handle F6 - pass until end of turn
	 */
	function handlePassUntilEOT(): void {
		if (!hasPriority || isLoading) return;
		onPassUntilEOT();
	}

	/**
	 * Handle cast spell
	 */
	function handleCastSpell(): void {
		console.log('[PriorityActionBar.handleCastSpell] Called', { hasPriority, isLoading });
		if (!hasPriority) {
			console.log('[PriorityActionBar.handleCastSpell] No priority, returning');
			return;
		}
		if (isLoading) {
			console.log('[PriorityActionBar.handleCastSpell] Loading, returning');
			return;
		}
		console.log('[PriorityActionBar.handleCastSpell] Calling onCastSpell callback');
		onCastSpell();
	}

	/**
	 * Handle activate ability
	 */
	function handleActivateAbility(): void {
		if (!hasPriority || isLoading) return;
		onActivateAbility();
	}

	/**
	 * Handle advance phase
	 */
	function handleAdvancePhase(): void {
		if (!hasPriority || isLoading) return;
		onAdvancePhase();
	}

	/**
	 * Toggle settings panel
	 */
	function toggleSettings(): void {
		showSettings = !showSettings;
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent): void {
		// Don't trigger if user is typing in an input
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
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="priority-action-bar" class:has-priority={hasPriority}>
	<!-- Priority Status -->
	<div class="priority-status">
		{#if hasPriority}
			<div class="status-indicator active">
				<span class="status-icon">⚡</span>
				<span class="status-text">Your Priority</span>
			</div>
		{:else if isYourTurn}
			<div class="status-indicator your-turn">
				<span class="status-icon">⏳</span>
				<span class="status-text">Waiting...</span>
			</div>
		{:else}
			<div class="status-indicator waiting">
				<span class="status-icon">⏸️</span>
				<span class="status-text">{activePlayerName}'s turn</span>
			</div>
		{/if}
	</div>

	<!-- Action Buttons -->
	<div class="action-buttons">
		<button
			class="action-btn primary"
			class:pulse={hasPriority}
			disabled={!hasPriority || !canPassPriority || isLoading}
			onclick={handlePassPriority}
			title="Pass priority (Spacebar)"
		>
			{#if isLoading}
				<span class="loading-spinner"></span>
			{:else}
				<span class="btn-icon">→</span>
			{/if}
			<span class="btn-label">Pass Priority</span>
			<kbd class="shortcut">Space</kbd>
		</button>

		<button
			class="action-btn secondary"
			disabled={!hasPriority || isLoading}
			onclick={handleCastSpell}
			title="Cast spell from hand (C)"
		>
			<span class="btn-icon">🎴</span>
			<span class="btn-label">Cast</span>
			<kbd class="shortcut">C</kbd>
		</button>

		<button
			class="action-btn secondary"
			disabled={!hasPriority || isLoading}
			onclick={handleActivateAbility}
			title="Activate ability (A)"
		>
			<span class="btn-icon">⚡</span>
			<span class="btn-label">Activate</span>
			<kbd class="shortcut">A</kbd>
		</button>

		<button
			class="action-btn next-phase"
			disabled={!hasPriority || isLoading}
			onclick={handleAdvancePhase}
			title="Advance to next phase (N)"
		>
			<span class="btn-icon">⏩</span>
			<span class="btn-label">Next Phase</span>
			<kbd class="shortcut">N</kbd>
		</button>

		<div class="divider"></div>

		<button
			class="action-btn f6"
			disabled={!hasPriority || isLoading}
			onclick={handlePassUntilEOT}
			title="Pass until your next upkeep (F6)"
		>
			<span class="btn-icon">⏭️</span>
			<span class="btn-label">Pass Turn</span>
			<kbd class="shortcut">F6</kbd>
		</button>
	</div>

	<!-- Auto-pass Settings Toggle -->
	<div class="settings-section">
		<button
			class="settings-toggle"
			class:active={showSettings}
			onclick={toggleSettings}
			title="Auto-pass settings"
		>
			<span class="settings-icon">⚙️</span>
		</button>

		{#if showSettings}
			<div class="settings-panel">
				<div class="settings-title">Auto-Pass Settings</div>
				<label class="setting-option">
					<input
						type="checkbox"
						bind:checked={autoPassSettings.noActions}
					/>
					<span>Auto-pass when no actions available</span>
				</label>
				<label class="setting-option">
					<input
						type="checkbox"
						bind:checked={autoPassSettings.opponentTurn}
					/>
					<span>Auto-pass on opponent's turn</span>
				</label>
			</div>
		{/if}
	</div>
</div>

<style>
	.priority-action-bar {
		position: fixed;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);
		display: flex;
		align-items: center;
		gap: 1.5rem;
		padding: 0.75rem 1.5rem;
		background: linear-gradient(to top, rgba(15, 20, 25, 0.98), rgba(26, 31, 46, 0.95));
		border: 1px solid #2a3441;
		border-bottom: none;
		border-radius: 12px 12px 0 0;
		backdrop-filter: blur(12px);
		box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.4);
		z-index: 100;
		transition: all 0.3s ease;
	}

	.priority-action-bar.has-priority {
		border-color: rgba(251, 191, 36, 0.5);
		box-shadow: 
			0 -4px 24px rgba(0, 0, 0, 0.4),
			0 0 40px rgba(251, 191, 36, 0.15);
	}

	/* Priority Status */
	.priority-status {
		min-width: 160px;
	}

	.status-indicator {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 1rem;
		border-radius: 8px;
		background: #1a1f2e;
		border: 1px solid #2a3441;
	}

	.status-indicator.active {
		background: linear-gradient(135deg, rgba(251, 191, 36, 0.15), rgba(245, 158, 11, 0.1));
		border-color: rgba(251, 191, 36, 0.4);
		animation: status-glow 2s ease-in-out infinite;
	}

	.status-indicator.your-turn {
		background: linear-gradient(135deg, rgba(102, 126, 234, 0.15), rgba(118, 75, 162, 0.1));
		border-color: rgba(102, 126, 234, 0.4);
	}

	.status-indicator.waiting {
		opacity: 0.7;
	}

	@keyframes status-glow {
		0%, 100% { box-shadow: 0 0 8px rgba(251, 191, 36, 0.2); }
		50% { box-shadow: 0 0 16px rgba(251, 191, 36, 0.4); }
	}

	.status-icon {
		font-size: 1.25rem;
	}

	.status-text {
		font-size: 0.875rem;
		font-weight: 600;
		color: #fff;
		white-space: nowrap;
	}

	.status-indicator.active .status-text {
		color: #fbbf24;
	}

	.status-indicator.waiting .status-text {
		color: #9ca3af;
	}

	/* Action Buttons */
	.action-buttons {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.action-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1rem;
		background: #2a3441;
		border: 1px solid #374151;
		border-radius: 8px;
		color: #fff;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		white-space: nowrap;
	}

	.action-btn:hover:not(:disabled) {
		background: #374151;
		border-color: #4b5563;
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.action-btn:active:not(:disabled) {
		transform: translateY(0);
	}

	.action-btn:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.action-btn.primary {
		background: linear-gradient(135deg, #10b981, #059669);
		border-color: #10b981;
		padding: 0.75rem 1.25rem;
	}

	.action-btn.primary:hover:not(:disabled) {
		background: linear-gradient(135deg, #059669, #047857);
		border-color: #059669;
	}

	.action-btn.primary.pulse:not(:disabled) {
		animation: pulse-btn 2s ease-in-out infinite;
	}

	@keyframes pulse-btn {
		0%, 100% { box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.4); }
		50% { box-shadow: 0 0 0 8px rgba(16, 185, 129, 0); }
	}

	.action-btn.f6 {
		background: #374151;
		border-color: #4b5563;
	}

	.action-btn.f6:hover:not(:disabled) {
		background: #4b5563;
		border-color: #6b7280;
	}

	.action-btn.next-phase {
		background: linear-gradient(135deg, #3b82f6, #2563eb);
		border-color: #3b82f6;
	}

	.action-btn.next-phase:hover:not(:disabled) {
		background: linear-gradient(135deg, #2563eb, #1d4ed8);
		border-color: #2563eb;
	}

	.btn-icon {
		font-size: 1.125rem;
	}

	.btn-label {
		font-weight: 600;
	}

	.shortcut {
		padding: 0.125rem 0.375rem;
		background: rgba(0, 0, 0, 0.3);
		border-radius: 4px;
		font-size: 0.625rem;
		font-weight: 700;
		color: #9ca3af;
		font-family: inherit;
		border: none;
	}

	.action-btn.primary .shortcut {
		background: rgba(0, 0, 0, 0.2);
		color: rgba(255, 255, 255, 0.8);
	}

	.divider {
		width: 1px;
		height: 32px;
		background: #374151;
		margin: 0 0.5rem;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: #fff;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Settings Section */
	.settings-section {
		position: relative;
	}

	.settings-toggle {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.settings-toggle:hover,
	.settings-toggle.active {
		background: #2a3441;
		border-color: #374151;
	}

	.settings-icon {
		font-size: 1.25rem;
		transition: transform 0.3s;
	}

	.settings-toggle.active .settings-icon {
		transform: rotate(90deg);
	}

	.settings-panel {
		position: absolute;
		bottom: calc(100% + 0.75rem);
		right: 0;
		width: 280px;
		padding: 1rem;
		background: #1a1f2e;
		border: 1px solid #2a3441;
		border-radius: 8px;
		box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.4);
	}

	.settings-title {
		font-size: 0.875rem;
		font-weight: 700;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		margin-bottom: 0.75rem;
		padding-bottom: 0.5rem;
		border-bottom: 1px solid #2a3441;
	}

	.setting-option {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem 0;
		cursor: pointer;
		font-size: 0.875rem;
		color: #d1d5db;
		transition: color 0.2s;
	}

	.setting-option:hover {
		color: #fff;
	}

	.setting-option input[type="checkbox"] {
		width: 16px;
		height: 16px;
		accent-color: #10b981;
		cursor: pointer;
	}

	/* Responsive */
	@media (max-width: 900px) {
		.priority-action-bar {
			left: 0;
			right: 0;
			transform: none;
			border-radius: 0;
			padding: 0.5rem 1rem;
			gap: 0.75rem;
		}

		.btn-label {
			display: none;
		}

		.action-btn {
			padding: 0.625rem;
		}

		.action-btn.primary {
			padding: 0.75rem;
		}

		.priority-status {
			min-width: auto;
		}

		.status-text {
			display: none;
		}
	}

	@media (max-width: 600px) {
		.shortcut {
			display: none;
		}

		.divider {
			display: none;
		}

		.action-buttons {
			gap: 0.25rem;
		}
	}
</style>

