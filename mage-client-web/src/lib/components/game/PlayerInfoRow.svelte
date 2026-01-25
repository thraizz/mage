<script lang="ts">
  import type { CardView } from '$lib/generated/mage/v1/models';
  import { toPossessiveName } from '$lib/utils/localization';
  import Heart from '@lucide/svelte/icons/heart';
  import Skull from '@lucide/svelte/icons/skull';
  import type { Action } from 'svelte/action';
  import ExileZone from './ExileZone.svelte';
  import Graveyard from './Graveyard.svelte';
  import LibraryZone from './LibraryZone.svelte';
  import ManaPoolComponent from './ManaPool.svelte';

  interface Player {
    name: string;
    life: number;
    poison: number;
    libraryCount: number;
  }

  interface ManaPoolData {
    white: number;
    blue: number;
    black: number;
    red: number;
    green: number;
    colorless: number;
  }

  interface Props {
    player: Player;
    graveyard: CardView[];
    exile: CardView[];
    mana: ManaPoolData;
    showLifeMenu: boolean;
    onLifeChange: (delta: number) => void;
    onPoisonChange: (delta: number) => void;
    onToggleLifeMenu: () => void;
    onSearchLibrary: () => void;
    onDeckContextMenu: (e: MouseEvent) => void;
    libraryDropZoneRef?: (el: HTMLElement | null) => void;
    graveyardDropZoneRef?: (el: HTMLElement | null) => void;
    exileDropZoneRef?: (el: HTMLElement | null) => void;
    // Playtest-specific props
    isPlaytest?: boolean;
    players?: Array<{ playerId: string; name: string }>;
    activeControlSeat?: string;
    onSwitchPlayer?: (playerId: string) => void;
  }

  let {
    player,
    graveyard,
    exile,
    mana,
    showLifeMenu,
    onLifeChange,
    onPoisonChange,
    onToggleLifeMenu,
    onSearchLibrary,
    onDeckContextMenu,
    libraryDropZoneRef,
    graveyardDropZoneRef,
    exileDropZoneRef,
    isPlaytest = false,
    players = [],
    activeControlSeat = '',
    onSwitchPlayer
  }: Props = $props();

  let lifeMenuEl: HTMLDivElement | null = $state(null);

  // Create actions for drop zones
  const libraryDropZone: Action<HTMLElement> = (node) => {
    if (libraryDropZoneRef) libraryDropZoneRef(node);
    return {
      destroy() {
        if (libraryDropZoneRef) libraryDropZoneRef(null);
      }
    };
  };

  const graveyardDropZone: Action<HTMLElement> = (node) => {
    if (graveyardDropZoneRef) graveyardDropZoneRef(node);
    return {
      destroy() {
        if (graveyardDropZoneRef) graveyardDropZoneRef(null);
      }
    };
  };

  const exileDropZone: Action<HTMLElement> = (node) => {
    if (exileDropZoneRef) exileDropZoneRef(node);
    return {
      destroy() {
        if (exileDropZoneRef) exileDropZoneRef(null);
      }
    };
  };
</script>

<div class="player-info-row">
  <div class="player-identity">
    {#if isPlaytest && players.length > 0 && onSwitchPlayer}
      <div class="controlling-dropdown">
        <select
          id="player-controlling-select"
          class="player-select"
          value={activeControlSeat}
          onchange={(e) => onSwitchPlayer?.(e.currentTarget.value)}
        >
          {#each players as p}
            <option value={p.playerId}>{p.name}</option>
          {/each}
        </select>
      </div>
    {:else if player}
      <span class="player-name">{player.name}</span>
    {:else}
      <span class="player-name">Waiting for player...</span>
    {/if}
  </div>

  <div class="player-stats-inline">
    <div class="life-group">
      <button class="stat-btn minus" onclick={() => onLifeChange(-1)}>−</button>
      <button
        title={`${toPossessiveName(player.name)} Life total`}
        class="stat-display life"
        onclick={onToggleLifeMenu}
      >
        <span class="stat-icon"><Heart size={14} /></span>
        <span class="stat-value">{player.life}</span>
      </button>
      <button class="stat-btn plus" onclick={() => onLifeChange(1)}>+</button>
    </div>

    {#if player.poison > 0}
      <div class="stat-display poison">
        <span class="stat-icon"><Skull size={14} /></span>
        <span class="stat-value">{player.poison}</span>
      </div>
    {/if}

    <div class="library-drop-zone" use:libraryDropZone>
      <LibraryZone
        libraryCount={player.libraryCount}
        playerName="You"
        onSearch={onSearchLibrary}
        onContextMenu={onDeckContextMenu}
      />
    </div>

    {#if showLifeMenu}
      <div bind:this={lifeMenuEl} class="quick-menu">
        <div class="menu-section">
          <span class="menu-label">Life</span>
          <div class="menu-row">
            <button onclick={() => onLifeChange(-5)}>−5</button>
            <button onclick={() => onLifeChange(-1)}>−1</button>
            <button onclick={() => onLifeChange(1)}>+1</button>
            <button onclick={() => onLifeChange(5)}>+5</button>
          </div>
        </div>
        <div class="menu-section">
          <span class="menu-label">Poison</span>
          <div class="menu-row">
            <button onclick={() => onPoisonChange(-1)}>−1</button>
            <span class="menu-value">{player.poison}</span>
            <button onclick={() => onPoisonChange(1)}>+1</button>
          </div>
        </div>
        <button class="menu-close" onclick={onToggleLifeMenu}>✕</button>
      </div>
    {/if}
  </div>

  <div class="player-zones">
    <div class="graveyard-drop-zone" use:graveyardDropZone>
      <Graveyard
        cards={graveyard.map((c) => ({
          id: c.id,
          name: c.name,
          manaCost: c.manaCost,
          cardType: c.type,
          power: c.power,
          toughness: c.toughness,
          imageUrl: '',
          isTapped: false,
          isSelected: false
        }))}
        playerName="You"
        isOpponent={false}
        canDrag={true}
        onCardClick={() => {}}
      />
    </div>
    <div class="exile-drop-zone" use:exileDropZone>
      <ExileZone
        cards={exile.map((c) => ({
          id: c.id,
          name: c.name,
          manaCost: c.manaCost,
          cardType: c.type,
          power: c.power,
          toughness: c.toughness,
          imageUrl: '',
          isTapped: false,
          isSelected: false
        }))}
        playerName="You"
        isOpponent={false}
        canDrag={true}
        onCardClick={() => {}}
        compact={true}
      />
    </div>
    <ManaPoolComponent {mana} showEmpty={false} size="small" onManaClick={() => {}} />
  </div>
</div>

<style>
  .controlling-dropdown {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .player-select {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
    border-radius: 4px;
    border: 1px solid #374151;
    background: #1e293b;
    color: #e2e8f0;
    cursor: pointer;
  }

  .player-select:hover {
    border-color: #4b5563;
  }

  .player-select:focus {
    outline: none;
    border-color: #667eea;
  }
</style>
