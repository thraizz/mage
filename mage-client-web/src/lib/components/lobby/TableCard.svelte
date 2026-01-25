<script lang="ts">
  import type { Table } from '$lib/types/table';
  import FormatBadge from '$lib/components/mtg/FormatBadge.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Lock from '@lucide/svelte/icons/lock';

  interface Props {
    table: Table;
    onclick?: (table: Table) => void;
  }

  let { table, onclick }: Props = $props();

  const isFull = $derived(table.players.length >= table.maxPlayers);
  const isOpen = $derived(table.status === 'waiting' && !isFull);

  const statusVariant = $derived.by((): 'warning' | 'success' | 'info' | 'muted' => {
    switch (table.status) {
      case 'waiting':
        return 'warning';
      case 'ready':
        return 'success';
      case 'playing':
        return 'info';
      default:
        return 'muted';
    }
  });

  const statusText = $derived.by(() => {
    switch (table.status) {
      case 'waiting':
        return 'Waiting';
      case 'ready':
        return 'Ready';
      case 'playing':
        return 'In Progress';
      case 'finished':
        return 'Finished';
      default:
        return 'Unknown';
    }
  });

  function handleClick() {
    onclick?.(table);
  }

  function handleKeypress(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleClick();
    }
  }
</script>

{#if onclick}
  <button
    type="button"
    class="table-card clickable"
    class:full={isFull}
    onclick={handleClick}
    onkeydown={handleKeypress}
  >
    <header class="card-header">
      <FormatBadge format={table.format} size="md" />
      <div class="header-icons">
        {#if table.hasPassword}
          <span class="lock-icon" title="Password protected">
            <Lock size={14} aria-hidden="true" />
          </span>
        {/if}
      </div>
    </header>

    <h3 class="table-name">{table.name}</h3>

    <div class="table-meta">
      <div class="meta-row">
        <span class="meta-label">Host</span>
        <span class="meta-value">{table.hostUsername}</span>
      </div>
    </div>

    <div class="player-slots">
      {#each Array(table.maxPlayers) as _, i}
        <div class="slot" class:filled={i < table.players.length}>
          {#if i < table.players.length}
            <span class="slot-filled"></span>
          {/if}
        </div>
      {/each}
      <span class="slot-count">{table.players.length}/{table.maxPlayers}</span>
    </div>

    <footer class="card-footer">
      <Badge variant={statusVariant}>{statusText}</Badge>
      {#if isOpen}
        <span class="join-hint">Click to join</span>
      {/if}
    </footer>
  </button>
{:else}
  <div class="table-card" class:full={isFull}>
    <header class="card-header">
      <FormatBadge format={table.format} size="md" />
      <div class="header-icons">
        {#if table.hasPassword}
          <span class="lock-icon" title="Password protected">
            <Lock size={14} aria-hidden="true" />
          </span>
        {/if}
      </div>
    </header>

    <h3 class="table-name">{table.name}</h3>

    <div class="table-meta">
      <div class="meta-row">
        <span class="meta-label">Host</span>
        <span class="meta-value">{table.hostUsername}</span>
      </div>
    </div>

    <div class="player-slots">
      {#each Array(table.maxPlayers) as _, i}
        <div class="slot" class:filled={i < table.players.length}>
          {#if i < table.players.length}
            <span class="slot-filled"></span>
          {/if}
        </div>
      {/each}
      <span class="slot-count">{table.players.length}/{table.maxPlayers}</span>
    </div>

    <footer class="card-footer">
      <Badge variant={statusVariant}>{statusText}</Badge>
    </footer>
  </div>
{/if}

<style>
  .table-card {
    display: block;
    width: 100%;
    text-align: left;
    font: inherit;
    color: inherit;
    background: var(--bg-obsidian);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-lg);
    padding: var(--space-4);
    transition: all var(--transition-base);
    position: relative;
    overflow: hidden;
  }

  .table-card.clickable {
    cursor: pointer;
  }

  .table-card.clickable:hover {
    border-color: var(--accent-gold-dim);
    box-shadow: var(--shadow-glow);
    transform: translateY(-2px);
  }

  .table-card.clickable:active {
    transform: translateY(0);
  }

  .table-card.full {
    opacity: 0.6;
  }

  .table-card:focus-visible {
    outline: 2px solid var(--accent-gold);
    outline-offset: 2px;
  }

  /* Header */
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: var(--space-3);
  }

  .header-icons {
    display: flex;
    gap: var(--space-2);
  }

  .lock-icon {
    color: var(--text-dim);
  }

  /* Table Name */
  .table-name {
    font-family: var(--font-display);
    font-size: var(--text-lg);
    font-weight: var(--weight-semibold);
    color: var(--text-bright);
    margin: 0 0 var(--space-3);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  /* Meta Info */
  .table-meta {
    margin-bottom: var(--space-3);
  }

  .meta-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: var(--text-sm);
  }

  .meta-label {
    color: var(--text-dim);
  }

  .meta-value {
    color: var(--text-muted);
    font-weight: var(--weight-medium);
  }

  /* Player Slots */
  .player-slots {
    display: flex;
    align-items: center;
    gap: var(--space-1);
    margin-bottom: var(--space-4);
  }

  .slot {
    width: 1.5rem;
    height: 0.375rem;
    background: var(--bg-iron);
    border-radius: var(--radius-full);
    overflow: hidden;
  }

  .slot.filled {
    background: var(--accent-gold-dim);
  }

  .slot-filled {
    display: block;
    width: 100%;
    height: 100%;
    background: var(--accent-gold);
  }

  .slot-count {
    margin-left: auto;
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-dim);
  }

  /* Footer */
  .card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: var(--space-3);
    border-top: 1px solid var(--border-subtle);
  }

  .join-hint {
    font-size: var(--text-xs);
    color: var(--text-ghost);
    opacity: 0;
    transition: opacity var(--transition-fast);
  }

  .table-card.clickable:hover .join-hint {
    opacity: 1;
  }
</style>
