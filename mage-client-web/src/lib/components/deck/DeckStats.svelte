<script lang="ts">
  import type { Deck } from '$lib/types/deck';
  import FormatBadge from '$lib/components/mtg/FormatBadge.svelte';
  import ManaCurve from './ManaCurve.svelte';
  import Panel from '$lib/components/ui/Panel.svelte';

  interface Props {
    deck: Deck;
  }

  let { deck }: Props = $props();

  // Calculate mana curve distribution
  // This is a placeholder - actual implementation would parse card CMCs
  const manaCurve = $derived([0, 8, 12, 10, 8, 4, 2, 2]);
</script>

<Panel title="Deck Statistics" variant="bordered">
  <div class="stats-content">
    <div class="stat-row">
      <span class="stat-label">Format</span>
      <FormatBadge format={deck.format} />
    </div>

    <div class="stat-row">
      <span class="stat-label">Total Cards</span>
      <span class="stat-value">{deck.cardCount}</span>
    </div>

    <div class="stat-row">
      <span class="stat-label">Main Deck</span>
      <span class="stat-value">{deck.mainDeck.length} cards</span>
    </div>

    {#if deck.sideboard && deck.sideboard.length > 0}
      <div class="stat-row">
        <span class="stat-label">Sideboard</span>
        <span class="stat-value">{deck.sideboard.length} cards</span>
      </div>
    {/if}

    {#if deck.commanders && deck.commanders.length > 0}
      <div class="stat-row">
        <span class="stat-label">Commander</span>
        <span class="stat-value">{deck.commanders.length}</span>
      </div>
    {/if}

    <div class="curve-section">
      <span class="stat-label">Mana Curve</span>
      <ManaCurve distribution={manaCurve} size="md" />
    </div>
  </div>
</Panel>

<style>
  .stats-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .stat-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-label {
    font-size: var(--text-sm);
    color: var(--text-dim);
  }

  .stat-value {
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-bright);
  }

  .curve-section {
    padding-top: var(--space-3);
    border-top: 1px solid var(--border-subtle);
  }

  .curve-section .stat-label {
    display: block;
    margin-bottom: var(--space-3);
  }
</style>
