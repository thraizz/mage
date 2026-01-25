<script lang="ts">
  import type { Player } from '$lib/types/gamestore';
  import Clock from '@lucide/svelte/icons/clock';
  import Eye from '@lucide/svelte/icons/eye';
  import EyeOff from '@lucide/svelte/icons/eye-off';
  import Keyboard from '@lucide/svelte/icons/keyboard';
  import X from '@lucide/svelte/icons/x';

  interface GameMenuProps {
    // State
    isOpen: boolean;

    // Data
    isMultiplayer: boolean;
    players?: Player[];
    activeControlSeat?: string;
    turnNumber?: number;
    activePlayerName?: string;
    availableSessions?: number;

    // UI State (for playtest controls)
    showAllHands?: boolean;

    // Event handlers
    onClose: () => void;
    onBackToLobby: () => void;
    onShowKeyboardShortcuts: () => void;

    // Playtest-specific handlers (optional)
    onSwitchPlayer?: (playerId: string) => void;
    onToggleAllHands?: () => void;
    onNextTurn?: () => void;
    onShowDebug?: () => void;
    onSessionsClick?: () => void;
  }

  let {
    isOpen,
    isMultiplayer,
    players,
    activeControlSeat,
    turnNumber,
    activePlayerName,
    availableSessions,
    showAllHands,
    onClose,
    onBackToLobby,
    onShowKeyboardShortcuts,
    onSwitchPlayer,
    onToggleAllHands,
    onNextTurn,
    onShowDebug,
    onSessionsClick
  }: GameMenuProps = $props();
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div
    class="menu-backdrop"
    role="button"
    tabindex="0"
    onclick={onClose}
    onkeydown={(e) => e.key === 'Escape' && onClose()}
  ></div>

  <!-- Menu Panel -->
  <div class="menu-overlay open">
    <div class="menu-header">
      <h2>Menu</h2>
      <button class="menu-close-btn" onclick={onClose} aria-label="Close menu">
        <X size={24} />
      </button>
    </div>

    <div class="menu-content">
      <!-- Playtest-specific: Controls Section -->
      {#if !isMultiplayer && players && activeControlSeat && onSwitchPlayer && onToggleAllHands}
        <div class="menu-section">
          <h3 class="menu-section-title">Controls</h3>
          <div class="menu-section-content">
            <label>
              <span class="menu-label">Controlling:</span>
              <select
                class="control-select"
                value={activeControlSeat}
                onchange={(e) => onSwitchPlayer(e.currentTarget.value)}
              >
                {#each players as player}
                  <option value={player.playerId}>{player.name}</option>
                {/each}
              </select>
            </label>

            <button class="menu-btn" onclick={onToggleAllHands}>
              {#if showAllHands}
                <EyeOff size={16} />
                Hide
              {:else}
                <Eye size={16} />
                Show
              {/if}
              All Hands
            </button>
          </div>
        </div>
      {/if}

      <!-- Playtest-specific: Turn Info Section -->
      {#if !isMultiplayer && turnNumber !== undefined && onNextTurn}
        <div class="menu-section">
          <h3 class="menu-section-title">Turn Info</h3>
          <div class="menu-section-content">
            <div class="turn-info">
              <Clock size={18} />
              <span>Turn {turnNumber}</span>
              {#if activePlayerName}
                <span class="active-player">· {activePlayerName}</span>
              {/if}
            </div>
            <button class="menu-btn primary" onclick={onNextTurn}>Next Turn</button>
          </div>
        </div>
      {/if}

      <!-- Common: Utilities Section -->
      <div class="menu-section">
        <h3 class="menu-section-title">Utilities</h3>
        <div class="menu-section-content">
          <button
            class="menu-btn"
            onclick={() => {
              onShowKeyboardShortcuts();
              onClose();
            }}
          >
            <Keyboard size={18} />
            Keyboard Shortcuts
          </button>
          {#if !isMultiplayer && onShowDebug}
            <button
              class="menu-btn"
              onclick={() => {
                onShowDebug();
                onClose();
              }}
            >
              🔧 Debug View
            </button>
          {/if}
        </div>
      </div>

      <!-- Common: Navigation Section -->
      <div class="menu-section">
        <h3 class="menu-section-title">Navigation</h3>
        <div class="menu-section-content">
          <button class="menu-btn" onclick={onBackToLobby}> ← Back to Lobby </button>
          {#if !isMultiplayer && availableSessions && availableSessions > 0 && onSessionsClick}
            <button
              class="menu-btn"
              onclick={() => {
                onSessionsClick();
                onClose();
              }}
            >
              <Clock size={18} />
              Sessions
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .menu-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 999;
  }

  .menu-overlay {
    position: fixed;
    top: 0;
    right: -400px;
    width: 400px;
    height: 100vh;
    background: #1a1a2e;
    box-shadow: -2px 0 10px rgba(0, 0, 0, 0.3);
    z-index: 1000;
    transition: right 0.3s ease;
    overflow-y: auto;
  }

  .menu-overlay.open {
    right: 0;
  }

  .menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .menu-header h2 {
    margin: 0;
    color: white;
    font-size: 24px;
  }

  .menu-close-btn {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    padding: 5px;
  }

  .menu-content {
    padding: 20px;
  }

  .menu-section {
    margin-bottom: 30px;
  }

  .menu-section-title {
    color: #4a90e2;
    font-size: 14px;
    text-transform: uppercase;
    margin-bottom: 10px;
    letter-spacing: 1px;
  }

  .menu-section-content {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .menu-btn {
    padding: 10px 15px;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 4px;
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: background 0.2s;
  }

  .menu-btn:hover {
    background: rgba(255, 255, 255, 0.2);
  }

  .menu-btn.primary {
    background: #4a90e2;
    border-color: #4a90e2;
  }

  .menu-btn.primary:hover {
    background: #357abd;
  }

  .menu-label {
    color: #aaa;
    font-size: 14px;
    margin-bottom: 5px;
    display: block;
  }

  .control-select {
    width: 100%;
    padding: 8px;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 4px;
    color: white;
    font-size: 14px;
  }

  .turn-info {
    display: flex;
    align-items: center;
    gap: 8px;
    color: white;
    padding: 10px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }

  .active-player {
    color: #4a90e2;
    font-weight: bold;
  }

  label {
    display: block;
  }
</style>
