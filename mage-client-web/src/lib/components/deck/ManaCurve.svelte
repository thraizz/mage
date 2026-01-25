<script lang="ts">
  interface Props {
    /** Array of card counts by CMC [0, 1, 2, 3, 4, 5, 6, 7+] */
    distribution: number[];
    size?: 'sm' | 'md' | 'lg';
  }

  let { distribution, size = 'md' }: Props = $props();

  const maxCount = $derived(Math.max(...distribution, 1));
  const labels = ['0', '1', '2', '3', '4', '5', '6', '7+'];

  function getHeight(count: number): number {
    return (count / maxCount) * 100;
  }
</script>

<div class="mana-curve mana-curve-{size}">
  <div class="bars">
    {#each distribution as count, i}
      <div class="bar-container">
        <div class="bar" style="height: {getHeight(count)}%">
          {#if count > 0}
            <span class="bar-count">{count}</span>
          {/if}
        </div>
        <span class="bar-label">{labels[i]}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .mana-curve {
    display: flex;
    flex-direction: column;
  }

  .bars {
    display: flex;
    align-items: flex-end;
    gap: var(--space-1);
    height: 100%;
  }

  .bar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex: 1;
  }

  .bar {
    width: 100%;
    background: linear-gradient(to top, var(--accent-gold-dim), var(--accent-gold));
    border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    min-height: 2px;
    position: relative;
    transition: height var(--transition-slow);
  }

  .bar-count {
    position: absolute;
    top: -1.25rem;
    left: 50%;
    transform: translateX(-50%);
    font-size: var(--text-xs);
    font-weight: var(--weight-semibold);
    color: var(--text-muted);
  }

  .bar-label {
    font-size: var(--text-xs);
    color: var(--text-dim);
    margin-top: var(--space-1);
  }

  /* Sizes */
  .mana-curve-sm .bars {
    height: 3rem;
  }

  .mana-curve-md .bars {
    height: 4rem;
  }

  .mana-curve-lg .bars {
    height: 6rem;
  }
</style>
