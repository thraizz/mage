<script lang="ts">
  import ArrowDown from '@lucide/svelte/icons/arrow-down';
  import Clock from '@lucide/svelte/icons/clock';
  import Eye from '@lucide/svelte/icons/eye';
  import EyeOff from '@lucide/svelte/icons/eye-off';
  import FastForward from '@lucide/svelte/icons/fast-forward';
  import Keyboard from '@lucide/svelte/icons/keyboard';
  import Menu from '@lucide/svelte/icons/menu';
  import Plus from '@lucide/svelte/icons/plus';
  import RotateCcw from '@lucide/svelte/icons/rotate-ccw';
  import Search from '@lucide/svelte/icons/search';
  import Shuffle from '@lucide/svelte/icons/shuffle';

  interface Props {
    isMultiplayer: boolean;
    players: Array<{ playerId: string; name: string }>;
    activeControlSeat: string;
    availableSessions: number;
    turnNumber: number;
    activePlayerName: string;
    showAllHands: boolean;
    onBack: () => void;
    onSessionsClick: () => void;
    onSwitchPlayer: (playerId: string) => void;
    onToggleAllHands: () => void;
    onDrawCard: () => void;
    onUntapAll: () => void;
    onShuffleLibrary: () => void;
    onSearchLibrary: () => void;
    onCreateToken: () => void;
    onNextTurn: () => void;
    onShowKeyboardShortcuts: () => void;
    onShowDebug: () => void;
    onToggleMenu: () => void;
  }

  let {
    isMultiplayer,
    players,
    activeControlSeat,
    availableSessions,
    turnNumber,
    activePlayerName,
    showAllHands,
    onBack,
    onSessionsClick,
    onSwitchPlayer,
    onToggleAllHands,
    onDrawCard,
    onUntapAll,
    onShuffleLibrary,
    onSearchLibrary,
    onCreateToken,
    onNextTurn,
    onShowKeyboardShortcuts,
    onShowDebug,
    onToggleMenu
  }: Props = $props();
</script>

<div class="playtest-header">
  <div class="header-left">
    <button class="btn-back" onclick={onBack} title="Back to Lobby"> ← Back </button>
    {#if !isMultiplayer}
      <span class="mode-badge">Playtest</span>
      <button class="btn-sessions" onclick={onSessionsClick} title="Manage sessions">
        Sessions ({availableSessions})
      </button>
    {/if}
  </div>

  {#if !isMultiplayer}
    <div class="playtest-controls">
      <label for="playtest-controlling-select">Controlling:</label>
      <select
        id="playtest-controlling-select"
        class="player-select"
        value={activeControlSeat}
        onchange={(e) => onSwitchPlayer(e.currentTarget.value)}
      >
        {#each players as player}
          <option value={player.playerId}>{player.name}</option>
        {/each}
      </select>
      <button
        class="btn-secondary"
        onclick={onToggleAllHands}
        title="Toggle visibility of all hands"
      >
        {#if showAllHands}
          <EyeOff size={16} />
          Hide
        {:else}
          <Eye size={16} />
          Show
        {/if}
        Hands
      </button>
    </div>
  {/if}

  <div class="header-actions">
    <button class="btn-secondary" onclick={onDrawCard} title="Draw a card (C)">
      <ArrowDown size={16} />
      Draw
    </button>
    <button class="btn-secondary" onclick={onUntapAll} title="Untap all permanents (X)">
      <RotateCcw size={16} />
      Untap All
    </button>
    <button class="btn-secondary" onclick={onShuffleLibrary} title="Shuffle library (V)">
      <Shuffle size={16} />
      Shuffle
    </button>
    <button class="btn-secondary" onclick={onSearchLibrary} title="Search library (F)">
      <Search size={16} />
      Search
    </button>
    <button class="btn-secondary" onclick={onCreateToken} title="Create token (W)">
      <Plus size={16} />
      Token
    </button>
  </div>

  <div class="header-right">
    <div class="turn-indicator" title="Current turn">
      <Clock size={16} aria-hidden="true" />
      <span class="turn-text">
        Turn {turnNumber}{activePlayerName ? ` · ${activePlayerName}` : ''}
      </span>
    </div>
    <button class="btn-primary" onclick={onNextTurn} title="Next turn (E)">
      <FastForward size={16} />
      Next Turn
    </button>
    <button
      class="btn-icon"
      onclick={onShowKeyboardShortcuts}
      title="Keyboard shortcuts (?)"
      aria-label="Keyboard shortcuts"
    >
      <Keyboard size={20} aria-hidden="true" />
    </button>
    <button class="btn-icon" onclick={onShowDebug} title="Debug View"> 🔧 </button>
    <button
      class="btn-icon menu-toggle-btn-header"
      onclick={onToggleMenu}
      aria-label="Open menu"
      title="More options (M)"
    >
      <Menu size={20} />
    </button>
  </div>
</div>
