<script lang="ts">
  import type { Table } from '$lib/types/table';
  import FormatBadge from '$lib/components/mtg/FormatBadge.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Panel from '$lib/components/ui/Panel.svelte';
  import Lock from '@lucide/svelte/icons/lock';

  interface Props {
    table: Table;
    onleave: () => void;
  }

  let { table, onleave }: Props = $props();

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
        return 'Waiting for Players';
      case 'ready':
        return 'Ready to Start';
      case 'playing':
        return 'Game in Progress';
      case 'finished':
        return 'Finished';
      default:
        return table.status;
    }
  });
</script>

<Panel variant="bordered" padding="md">
  <div class="table-info">
    <div class="info-header">
      <div class="info-title">
        <h1>{table.name || `Table #${table.id.slice(0, 8)}`}</h1>
        <div class="info-badges">
          <FormatBadge format={table.format} size="lg" />
          {#if table.hasPassword}
            <Badge variant="warning" size="sm">
              <Lock size={12} aria-hidden="true" />
              Protected
            </Badge>
          {/if}
        </div>
      </div>
      <Button variant="danger" size="sm" onclick={onleave}>Leave Table</Button>
    </div>

    <div class="info-details">
      <div class="detail-row">
        <span class="detail-label">Host</span>
        <span class="detail-value">{table.hostUsername}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Players</span>
        <span class="detail-value">{table.players.length} / {table.maxPlayers}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Status</span>
        <Badge variant={statusVariant}>{statusText}</Badge>
      </div>
    </div>
  </div>
</Panel>

<style>
  .table-info {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .info-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-4);
  }

  .info-title h1 {
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    font-weight: var(--weight-bold);
    color: var(--text-bright);
    margin: 0 0 var(--space-2);
  }

  .info-badges {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .info-details {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-6);
  }

  .detail-row {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .detail-label {
    font-size: var(--text-xs);
    font-weight: var(--weight-medium);
    color: var(--text-dim);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .detail-value {
    font-size: var(--text-base);
    font-weight: var(--weight-medium);
    color: var(--text-bright);
  }

  @media (max-width: 640px) {
    .info-header {
      flex-direction: column;
    }

    .info-details {
      flex-direction: column;
      gap: var(--space-3);
    }
  }
</style>
