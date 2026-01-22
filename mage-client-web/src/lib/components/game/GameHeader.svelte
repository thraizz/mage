<script lang="ts">
	/**
	 * GameHeader - Clean header for game page answering key questions:
	 * - What turn is it?
	 * - Whose turn is it?
	 * - Who has priority?
	 * - Where can I find what happened? (Log button)
	 */
	import type { GamePhase } from '$lib/types/game';

	// Props
	let {
		turn = 1,
		activePlayerName = 'Player',
		priorityPlayerName = '',
		localPlayerName = '',
		hasPriority = false,
		currentPhase = 'PRECOMBAT_MAIN' as GamePhase,
		onLogClick = () => {},
		onConcedeClick = () => {}
	}: {
		turn?: number;
		activePlayerName?: string;
		priorityPlayerName?: string;
		localPlayerName?: string;
		hasPriority?: boolean;
		currentPhase?: GamePhase;
		onLogClick?: () => void;
		onConcedeClick?: () => void;
	} = $props();

	// Phase configuration matching MTG turn structure
	const phases = [
		{ key: 'BEGINNING', label: 'Begin', type: 'normal' },
		{ key: 'UNTAP', label: 'Untap', type: 'normal' },
		{ key: 'UPKEEP', label: 'Upkeep', type: 'normal' },
		{ key: 'DRAW', label: 'Draw', type: 'normal' },
		{ key: 'PRECOMBAT_MAIN', label: 'Main 1', type: 'main' },
		{ key: 'COMBAT', label: 'Combat', type: 'combat' },
		{ key: 'DECLARE_ATTACKERS', label: 'Attackers', type: 'combat' },
		{ key: 'DECLARE_BLOCKERS', label: 'Blockers', type: 'combat' },
		{ key: 'COMBAT_DAMAGE', label: 'Damage', type: 'combat' },
		{ key: 'END_OF_COMBAT', label: 'End Cmbt', type: 'combat' },
		{ key: 'POSTCOMBAT_MAIN', label: 'Main 2', type: 'main' },
		{ key: 'END', label: 'End', type: 'normal' },
		{ key: 'END_OF_TURN', label: 'EOT', type: 'normal' },
		{ key: 'CLEANUP', label: 'Cleanup', type: 'normal' }
	] as const;

	// Derived: is it the local player's turn?
	const isYourTurn = $derived(
		localPlayerName && activePlayerName && localPlayerName === activePlayerName
	);

	// Derived: turn display text
	const turnDisplay = $derived(isYourTurn ? 'Your Turn' : `${activePlayerName}'s turn`);

	// Derived: who has priority display text
	const priorityDisplay = $derived(
		hasPriority
			? 'You have priority'
			: priorityPlayerName
				? `${priorityPlayerName} has priority`
				: 'Waiting...'
	);

	// Check if a phase is active
	function isActive(phaseKey: string): boolean {
		return phaseKey === currentPhase;
	}

	// Check if we're in a combat step
	const isInCombatStep = $derived(
		['COMBAT', 'DECLARE_ATTACKERS', 'DECLARE_BLOCKERS', 'COMBAT_DAMAGE', 'END_OF_COMBAT'].includes(
			currentPhase
		)
	);

	// Get specific combat step for styling
	const combatStepClass = $derived(
		currentPhase === 'DECLARE_ATTACKERS'
			? 'attackers'
			: currentPhase === 'DECLARE_BLOCKERS'
				? 'blockers'
				: currentPhase === 'COMBAT_DAMAGE'
					? 'damage'
					: ''
	);
</script>

<header
	class="game-header"
	class:has-priority={hasPriority}
	class:in-combat={isInCombatStep}
	class:attackers={combatStepClass === 'attackers'}
	class:blockers={combatStepClass === 'blockers'}
	class:damage={combatStepClass === 'damage'}
>
	<!-- Left section: Log button and turn info -->
	<div class="header-left">
		<button class="log-btn" onclick={onLogClick} title="View game log"> Log </button>

		<div class="turn-info">
			<span class="turn-number">Turn {turn}</span>
			<span class="turn-player" class:your-turn={isYourTurn}>{turnDisplay}</span>
		</div>

		<!-- Combat Phase Badge -->
		{#if isInCombatStep}
			<div
				class="combat-badge"
				class:attackers={combatStepClass === 'attackers'}
				class:blockers={combatStepClass === 'blockers'}
				class:damage={combatStepClass === 'damage'}
			>
				{#if combatStepClass === 'attackers'}
					⚔️ Declare Attackers
				{:else if combatStepClass === 'blockers'}
					🛡️ Declare Blockers
				{:else if combatStepClass === 'damage'}
					💥 Combat Damage
				{:else}
					⚔️ Combat
				{/if}
			</div>
		{/if}
	</div>

	<!-- Right section: Priority indicator -->
	<div class="header-right">
		<span class="priority-indicator" class:active={hasPriority}>
			{priorityDisplay}
		</span>
		<button class="concede-btn" onclick={onConcedeClick} title="Concede game"> 🏳️ Concede </button>
	</div>
</header>

<!-- Phase track below header -->
<div class="phase-track-container" class:has-priority={hasPriority}>
	<div class="phase-track">
		{#each phases as phase}
			<div
				class="phase-item"
				class:active={isActive(phase.key)}
				class:main={phase.type === 'main'}
				class:combat={phase.type === 'combat'}
			>
				<div class="phase-dot"></div>
				<span class="phase-label">{phase.label}</span>
			</div>
		{/each}
	</div>
</div>

<style>
	/* Header container */
	.game-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1.25rem;
		background: linear-gradient(180deg, #12161d 0%, #0d1117 100%);
		border-bottom: 1px solid #21262d;
		gap: 1rem;
	}

	.game-header.has-priority {
		border-bottom-color: rgba(251, 191, 36, 0.3);
	}

	/* Left section */
	.header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.log-btn {
		padding: 0.5rem 1rem;
		background: transparent;
		border: 1.5px solid #3d444d;
		border-radius: 6px;
		color: #e6edf3;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.log-btn:hover {
		background: rgba(255, 255, 255, 0.05);
		border-color: #8b949e;
	}

	.turn-info {
		display: flex;
		align-items: center;
		gap: 0.625rem;
	}

	.turn-number {
		color: #f0b429;
		font-weight: 700;
		font-size: 0.9375rem;
	}

	.turn-player {
		color: #e6edf3;
		font-size: 0.9375rem;
	}

	.turn-player.your-turn {
		color: #f0b429;
		font-weight: 600;
	}

	/* Combat Badge */
	.combat-badge {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.375rem 0.75rem;
		background: rgba(248, 81, 73, 0.15);
		border: 1px solid rgba(248, 81, 73, 0.4);
		border-radius: 9999px;
		color: #f85149;
		font-size: 0.8125rem;
		font-weight: 600;
		animation: combat-pulse 2s ease-in-out infinite;
	}

	.combat-badge.attackers {
		background: rgba(239, 68, 68, 0.15);
		border-color: rgba(239, 68, 68, 0.5);
		color: #ef4444;
	}

	.combat-badge.blockers {
		background: rgba(59, 130, 246, 0.15);
		border-color: rgba(59, 130, 246, 0.5);
		color: #3b82f6;
	}

	.combat-badge.damage {
		background: rgba(245, 158, 11, 0.15);
		border-color: rgba(245, 158, 11, 0.5);
		color: #f59e0b;
	}

	@keyframes combat-pulse {
		0%,
		100% {
			box-shadow: 0 0 0 0 rgba(248, 81, 73, 0.4);
		}
		50% {
			box-shadow: 0 0 0 4px rgba(248, 81, 73, 0);
		}
	}

	.combat-badge.attackers {
		animation-name: combat-pulse-red;
	}

	@keyframes combat-pulse-red {
		0%,
		100% {
			box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.4);
		}
		50% {
			box-shadow: 0 0 0 4px rgba(239, 68, 68, 0);
		}
	}

	.combat-badge.blockers {
		animation-name: combat-pulse-blue;
	}

	@keyframes combat-pulse-blue {
		0%,
		100% {
			box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.4);
		}
		50% {
			box-shadow: 0 0 0 4px rgba(59, 130, 246, 0);
		}
	}

	.combat-badge.damage {
		animation-name: combat-pulse-orange;
	}

	@keyframes combat-pulse-orange {
		0%,
		100% {
			box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.4);
		}
		50% {
			box-shadow: 0 0 0 4px rgba(245, 158, 11, 0);
		}
	}

	/* Header combat state styling */
	.game-header.in-combat {
		border-bottom-color: rgba(248, 81, 73, 0.4);
	}

	.game-header.in-combat.attackers {
		border-bottom-color: rgba(239, 68, 68, 0.5);
		background: linear-gradient(180deg, rgba(239, 68, 68, 0.08) 0%, #0d1117 100%);
	}

	.game-header.in-combat.blockers {
		border-bottom-color: rgba(59, 130, 246, 0.5);
		background: linear-gradient(180deg, rgba(59, 130, 246, 0.08) 0%, #0d1117 100%);
	}

	.game-header.in-combat.damage {
		border-bottom-color: rgba(245, 158, 11, 0.5);
		background: linear-gradient(180deg, rgba(245, 158, 11, 0.08) 0%, #0d1117 100%);
	}

	/* Right section */
	.header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.priority-indicator {
		color: #8b949e;
		font-size: 0.9375rem;
		font-weight: 500;
	}

	.priority-indicator.active {
		color: #f0b429;
		font-weight: 700;
		text-shadow: 0 0 12px rgba(240, 180, 41, 0.4);
	}

	.concede-btn {
		padding: 0.5rem 0.875rem;
		background: rgba(248, 81, 73, 0.1);
		border: 1px solid rgba(248, 81, 73, 0.3);
		border-radius: 6px;
		color: #f85149;
		font-size: 0.8125rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.concede-btn:hover {
		background: rgba(248, 81, 73, 0.2);
		border-color: rgba(248, 81, 73, 0.5);
	}

	/* Phase track container */
	.phase-track-container {
		padding: 0.75rem 1.25rem 1rem;
		background: #0d1117;
		border-bottom: 1px solid #21262d;
	}

	.phase-track-container.has-priority {
		border-bottom-color: rgba(251, 191, 36, 0.2);
	}

	/* Phase track */
	.phase-track {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		gap: 0.25rem;
		max-width: 900px;
		margin: 0 auto;
	}

	.phase-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.375rem;
		min-width: 52px;
		padding: 0 0.25rem;
	}

	.phase-dot {
		width: 12px;
		height: 12px;
		border-radius: 50%;
		background: #30363d;
		border: 2px solid transparent;
		transition: all 0.2s ease;
	}

	/* Main phase styling (yellow border) */
	.phase-item.main .phase-dot {
		background: transparent;
		border-color: #f0b429;
	}

	.phase-item.main.active .phase-dot {
		background: #f0b429;
		border-color: #f0b429;
		box-shadow: 0 0 10px rgba(240, 180, 41, 0.6);
	}

	/* Combat phase styling (red border) */
	.phase-item.combat .phase-dot {
		background: transparent;
		border-color: #f85149;
	}

	.phase-item.combat.active .phase-dot {
		background: #f85149;
		border-color: #f85149;
		box-shadow: 0 0 10px rgba(248, 81, 73, 0.6);
	}

	/* Normal phase active */
	.phase-item.active:not(.main):not(.combat) .phase-dot {
		background: #58a6ff;
		box-shadow: 0 0 10px rgba(88, 166, 255, 0.6);
	}

	.phase-label {
		font-size: 0.6875rem;
		color: #484f58;
		font-weight: 500;
		text-align: center;
		white-space: nowrap;
	}

	.phase-item.active .phase-label {
		color: #e6edf3;
		font-weight: 600;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.game-header {
			padding: 0.5rem 0.75rem;
		}

		.phase-track-container {
			padding: 0.5rem 0.5rem 0.75rem;
			overflow-x: auto;
		}

		.phase-track {
			justify-content: flex-start;
			gap: 0.125rem;
		}

		.phase-item {
			min-width: 44px;
		}

		.phase-label {
			font-size: 0.5625rem;
		}

		.concede-btn {
			padding: 0.375rem 0.5rem;
			font-size: 0.75rem;
		}
	}

	@media (max-width: 480px) {
		.priority-indicator {
			font-size: 0.8125rem;
		}

		.log-btn {
			padding: 0.375rem 0.75rem;
			font-size: 0.8125rem;
		}
	}
</style>
