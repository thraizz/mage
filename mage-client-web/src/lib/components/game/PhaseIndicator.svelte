<script lang="ts">
  import type { GamePhase } from '$lib/types/game';

  // Props
  let {
    currentPhase = 'PRECOMBAT_MAIN',
    activePlayerId = '',
    localPlayerId = '',
    animated = true
  }: {
    currentPhase?: GamePhase;
    activePlayerId?: string;
    localPlayerId?: string;
    animated?: boolean;
  } = $props();

  // Phase configuration
  const phases = [
    { key: 'BEGINNING', label: 'Beginning', shortLabel: 'Begin', icon: '🌅' },
    { key: 'UNTAP', label: 'Untap', shortLabel: 'Untap', icon: '↻' },
    { key: 'UPKEEP', label: 'Upkeep', shortLabel: 'Upkeep', icon: '⏰' },
    { key: 'DRAW', label: 'Draw', shortLabel: 'Draw', icon: '🎴' },
    {
      key: 'PRECOMBAT_MAIN',
      label: 'Main Phase 1',
      shortLabel: 'Main 1',
      icon: '🎯',
      isMain: true
    },
    { key: 'COMBAT', label: 'Combat', shortLabel: 'Combat', icon: '⚔️', isCombat: true },
    {
      key: 'DECLARE_ATTACKERS',
      label: 'Declare Attackers',
      shortLabel: 'Attackers',
      icon: '🗡️',
      isCombat: true
    },
    {
      key: 'DECLARE_BLOCKERS',
      label: 'Declare Blockers',
      shortLabel: 'Blockers',
      icon: '🛡️',
      isCombat: true
    },
    {
      key: 'COMBAT_DAMAGE',
      label: 'Combat Damage',
      shortLabel: 'Damage',
      icon: '💥',
      isCombat: true
    },
    {
      key: 'END_OF_COMBAT',
      label: 'End of Combat',
      shortLabel: 'End Cmbt',
      icon: '🏁',
      isCombat: true
    },
    {
      key: 'POSTCOMBAT_MAIN',
      label: 'Main Phase 2',
      shortLabel: 'Main 2',
      icon: '🎯',
      isMain: true
    },
    { key: 'END', label: 'End Step', shortLabel: 'End', icon: '🌙' },
    { key: 'END_OF_TURN', label: 'End of Turn', shortLabel: 'EOT', icon: '🌃' },
    { key: 'CLEANUP', label: 'Cleanup', shortLabel: 'Cleanup', icon: '🧹' }
  ] as const;

  /**
   * Check if a phase is the current phase
   */
  function isCurrentPhase(phaseKey: string): boolean {
    return phaseKey === currentPhase;
  }

  /**
   * Get phase display class
   */
  function getPhaseClass(phase: (typeof phases)[number]): string {
    const classes = ['phase-item'];
    if (isCurrentPhase(phase.key)) classes.push('active');
    if ('isMain' in phase && phase.isMain) classes.push('main-phase');
    if ('isCombat' in phase && phase.isCombat) classes.push('combat-phase');
    return classes.join(' ');
  }

  // Derived values
  const isYourTurn = $derived(activePlayerId === localPlayerId);
  const currentPhaseInfo = $derived(phases.find((p) => p.key === currentPhase) || phases[0]);
</script>

<div class="phase-indicator" class:animated>
  <div class="phase-header">
    <span class="current-phase-label">
      {currentPhaseInfo.icon}
      {currentPhaseInfo.label}
    </span>
    <span class="turn-indicator" class:your-turn={isYourTurn}>
      {isYourTurn ? 'Your Turn' : "Opponent's Turn"}
    </span>
  </div>

  <div class="phase-track">
    {#each phases as phase}
      <div class={getPhaseClass(phase)} title={phase.label}>
        <div class="phase-dot"></div>
        <span class="phase-label">{phase.shortLabel}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .phase-indicator {
    background: #1a1f2e;
    border: 1px solid #2a3441;
    border-radius: 8px;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .phase-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
  }

  .current-phase-label {
    font-size: 1rem;
    font-weight: 600;
    color: #ffffff;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .turn-indicator {
    font-size: 0.875rem;
    padding: 0.375rem 0.75rem;
    border-radius: 4px;
    font-weight: 600;
    background: #374151;
    color: #9ca3af;
    transition: all 0.3s;
  }

  .turn-indicator.your-turn {
    background: #10b981;
    color: white;
  }

  .phase-track {
    display: flex;
    gap: 0.5rem;
    overflow-x: auto;
    padding-bottom: 0.5rem;
  }

  .phase-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.375rem;
    min-width: 60px;
    flex-shrink: 0;
    cursor: pointer;
    transition: all 0.2s;
  }

  .phase-item:hover {
    transform: translateY(-2px);
  }

  .phase-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #374151;
    transition: all 0.3s;
    position: relative;
  }

  .phase-item.active .phase-dot {
    width: 16px;
    height: 16px;
    background: #667eea;
    box-shadow: 0 0 12px rgba(102, 126, 234, 0.6);
  }

  .phase-item.main-phase .phase-dot {
    border: 2px solid #fbbf24;
  }

  .phase-item.main-phase.active .phase-dot {
    background: #fbbf24;
    box-shadow: 0 0 12px rgba(251, 191, 36, 0.6);
  }

  .phase-item.combat-phase .phase-dot {
    border: 2px solid #ef4444;
  }

  .phase-item.combat-phase.active .phase-dot {
    background: #ef4444;
    box-shadow: 0 0 12px rgba(239, 68, 68, 0.6);
  }

  .phase-label {
    font-size: 0.625rem;
    color: #6b7280;
    text-align: center;
    font-weight: 500;
    transition: color 0.3s;
  }

  .phase-item.active .phase-label {
    color: #ffffff;
    font-weight: 600;
  }

  /* Animated variant */
  .phase-indicator.animated .phase-item.active .phase-dot {
    animation: pulse-phase 1.5s ease-in-out infinite;
  }

  @keyframes pulse-phase {
    0%,
    100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.2);
      opacity: 0.8;
    }
  }

  /* Scrollbar */
  .phase-track::-webkit-scrollbar {
    height: 4px;
  }

  .phase-track::-webkit-scrollbar-track {
    background: #0d1117;
    border-radius: 2px;
  }

  .phase-track::-webkit-scrollbar-thumb {
    background: #3a4451;
    border-radius: 2px;
  }

  .phase-track::-webkit-scrollbar-thumb:hover {
    background: #4a5461;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .phase-indicator {
      padding: 0.75rem;
    }

    .current-phase-label {
      font-size: 0.875rem;
    }

    .turn-indicator {
      font-size: 0.75rem;
      padding: 0.25rem 0.5rem;
    }

    .phase-item {
      min-width: 48px;
    }

    .phase-label {
      font-size: 0.5rem;
    }
  }
</style>
