<script lang="ts">
  import type { OnlinePlayer } from '$lib/types/player';
  import Panel from '$lib/components/ui/Panel.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Users from '@lucide/svelte/icons/users';
  import ChevronDown from '@lucide/svelte/icons/chevron-down';

  interface Props {
    players: OnlinePlayer[];
    currentUsername: string;
    isOpen?: boolean;
  }

  let { players, currentUsername, isOpen = $bindable(true) }: Props = $props();

  const playerCount = $derived(players.length);

  function toggleOpen() {
    isOpen = !isOpen;
  }
</script>

<div class="player-list-container">
  <button class="list-header" onclick={toggleOpen} aria-expanded={isOpen}>
    <div class="header-left">
      <Users class="header-icon" size={18} aria-hidden="true" />
      <span class="header-title">Online</span>
      <Badge variant="default" size="sm">{playerCount}</Badge>
    </div>
    <ChevronDown class={`toggle-icon ${isOpen ? 'open' : ''}`} size={16} aria-hidden="true" />
  </button>

  {#if isOpen}
    <div class="player-list" role="list">
      {#if players.length === 0}
        <div class="empty-state">
          <p>No players online</p>
        </div>
      {:else}
        {#each players as player (player.id)}
          <div class="player-item" role="listitem">
            <span class="status-dot"></span>
            <span class="player-name">
              {player.username}
              {#if player.username === currentUsername}
                <Badge variant="info" size="sm">You</Badge>
              {/if}
            </span>
          </div>
        {/each}
      {/if}
    </div>
  {/if}
</div>

<style>
  .player-list-container {
    background: var(--bg-obsidian);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  /* Header */
  .list-header {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--space-3) var(--space-4);
    background: transparent;
    border: none;
    cursor: pointer;
    transition: background var(--transition-fast);
  }

  .list-header:hover {
    background: var(--bg-slate);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  :global(svg.header-icon) {
    color: var(--accent-gold);
  }

  .header-title {
    font-size: var(--text-sm);
    font-weight: var(--weight-semibold);
    color: var(--text-bright);
  }

  :global(svg.toggle-icon) {
    color: var(--text-dim);
    transition: transform var(--transition-fast);
  }

  :global(svg.toggle-icon.open) {
    transform: rotate(180deg);
  }

  /* Player List */
  .player-list {
    max-height: 300px;
    overflow-y: auto;
    border-top: 1px solid var(--border-subtle);
  }

  .player-item {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-4);
    transition: background var(--transition-fast);
  }

  .player-item:hover {
    background: var(--bg-slate);
  }

  .player-item:not(:last-child) {
    border-bottom: 1px solid var(--border-subtle);
  }

  .status-dot {
    width: 0.5rem;
    height: 0.5rem;
    background: var(--status-success);
    border-radius: var(--radius-full);
    box-shadow: 0 0 6px var(--status-success);
    flex-shrink: 0;
  }

  .player-name {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    font-size: var(--text-sm);
    color: var(--text-muted);
  }

  /* Empty State */
  .empty-state {
    padding: var(--space-6) var(--space-4);
    text-align: center;
  }

  .empty-state p {
    color: var(--text-dim);
    font-size: var(--text-sm);
    margin: 0;
  }
</style>
