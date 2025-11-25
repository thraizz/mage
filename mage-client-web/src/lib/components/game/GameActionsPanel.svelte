<script lang="ts">
	// Props
	let {
		hasPriority = false,
		canPassPriority = true,
		isLoading = false,
		onPassPriority = () => {},
		onCastSpell = () => {},
		onActivateAbility = () => {}
	}: {
		hasPriority?: boolean;
		canPassPriority?: boolean;
		isLoading?: boolean;
		onPassPriority?: () => void;
		onCastSpell?: () => void;
		onActivateAbility?: () => void;
	} = $props();

	/**
	 * Handle pass priority button click
	 */
	function handlePassPriority(): void {
		if (!hasPriority || !canPassPriority || isLoading) return;
		onPassPriority();
	}

	/**
	 * Handle cast spell button click
	 */
	function handleCastSpell(): void {
		if (!hasPriority || isLoading) return;
		onCastSpell();
	}

	/**
	 * Handle activate ability button click
	 */
	function handleActivateAbility(): void {
		if (!hasPriority || isLoading) return;
		onActivateAbility();
	}

	/**
	 * Handle keyboard shortcuts
	 */
	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === ' ' || event.key === 'Spacebar') {
			event.preventDefault();
			handlePassPriority();
		} else if (event.key === 'c' || event.key === 'C') {
			event.preventDefault();
			handleCastSpell();
		} else if (event.key === 'a' || event.key === 'A') {
			event.preventDefault();
			handleActivateAbility();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="game-actions-panel" class:has-priority={hasPriority}>
	<div class="actions-header">
		<span class="actions-label">Actions</span>
		{#if hasPriority}
			<span class="priority-badge">Priority</span>
		{/if}
	</div>

	<div class="actions-buttons">
		<button
			class="action-btn pass-priority"
			disabled={!hasPriority || !canPassPriority || isLoading}
			onclick={handlePassPriority}
			aria-label="Pass priority (Spacebar)"
			title="Pass priority (Spacebar)"
		>
			{#if isLoading}
				<span class="loading-spinner"></span>
			{:else}
				<span class="btn-icon">→</span>
			{/if}
			<span class="btn-text">Pass Priority</span>
			<span class="btn-shortcut">Space</span>
		</button>

		<button
			class="action-btn cast-spell"
			disabled={!hasPriority || isLoading}
			onclick={handleCastSpell}
			aria-label="Cast spell (C)"
			title="Cast spell from hand (C)"
		>
			<span class="btn-icon">🎴</span>
			<span class="btn-text">Cast Spell</span>
			<span class="btn-shortcut">C</span>
		</button>

		<button
			class="action-btn activate-ability"
			disabled={!hasPriority || isLoading}
			onclick={handleActivateAbility}
			aria-label="Activate ability (A)"
			title="Activate permanent ability (A)"
		>
			<span class="btn-icon">⚡</span>
			<span class="btn-text">Activate</span>
			<span class="btn-shortcut">A</span>
		</button>
	</div>

	{#if !hasPriority}
		<div class="waiting-message">
			<span class="waiting-icon">⏳</span>
			<span class="waiting-text">Waiting for priority...</span>
		</div>
	{/if}
</div>

<style>
	.game-actions-panel {
		background: linear-gradient(135deg, #1a1f2e 0%, #0f1419 100%);
		border: 2px solid #2a3441;
		border-radius: 8px;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
		transition: all 0.3s ease;
	}

	.game-actions-panel.has-priority {
		border-color: #667eea;
		box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
	}

	.actions-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
	}

	.actions-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.priority-badge {
		padding: 0.25rem 0.5rem;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 700;
		color: #000;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		animation: badge-pulse 2s ease-in-out infinite;
	}

	@keyframes badge-pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.7;
		}
	}

	.actions-buttons {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.action-btn {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		background: #2a3441;
		border: 2px solid #374151;
		border-radius: 6px;
		font-size: 0.875rem;
		font-weight: 600;
		color: #ffffff;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		overflow: hidden;
	}

	.action-btn:hover:not(:disabled) {
		background: #374151;
		border-color: #4b5563;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
	}

	.action-btn:active:not(:disabled) {
		transform: translateY(0);
	}

	.action-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.action-btn.pass-priority:not(:disabled) {
		background: linear-gradient(135deg, #10b981 0%, #059669 100%);
		border-color: #10b981;
	}

	.action-btn.pass-priority:hover:not(:disabled) {
		background: linear-gradient(135deg, #059669 0%, #047857 100%);
		border-color: #059669;
	}

	.action-btn.cast-spell:not(:disabled):hover {
		border-color: #667eea;
	}

	.action-btn.activate-ability:not(:disabled):hover {
		border-color: #fbbf24;
	}

	.btn-icon {
		font-size: 1.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.btn-text {
		flex: 1;
		text-align: left;
	}

	.btn-shortcut {
		padding: 0.25rem 0.5rem;
		background: rgba(0, 0, 0, 0.3);
		border-radius: 4px;
		font-size: 0.625rem;
		font-weight: 700;
		color: #9ca3af;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.loading-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: #ffffff;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.waiting-message {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.75rem;
		background: rgba(107, 114, 128, 0.1);
		border-radius: 4px;
		border: 1px dashed #374151;
	}

	.waiting-icon {
		font-size: 1.25rem;
		opacity: 0.6;
	}

	.waiting-text {
		font-size: 0.875rem;
		color: #6b7280;
		font-style: italic;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.game-actions-panel {
			padding: 0.75rem;
		}

		.action-btn {
			padding: 0.625rem 0.875rem;
			font-size: 0.8125rem;
		}

		.btn-icon {
			font-size: 1.125rem;
		}

		.btn-shortcut {
			display: none;
		}
	}

	/* Compact mode for smaller screens */
	@media (max-width: 480px) {
		.actions-buttons {
			flex-direction: row;
			flex-wrap: wrap;
		}

		.action-btn {
			flex: 1 1 calc(50% - 0.25rem);
			min-width: 0;
		}

		.action-btn.pass-priority {
			flex-basis: 100%;
		}

		.btn-text {
			font-size: 0.75rem;
		}
	}
</style>
