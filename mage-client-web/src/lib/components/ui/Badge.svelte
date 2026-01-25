<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'muted';
    size?: 'sm' | 'md';
    children?: Snippet;
  }

  let { variant = 'default', size = 'md', children }: Props = $props();
</script>

<span class="badge badge-{variant} badge-{size}">
  {#if children}
    {@render children()}
  {/if}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
    font-family: var(--font-body);
    font-weight: var(--weight-medium);
    border-radius: var(--radius-full);
    white-space: nowrap;
    vertical-align: middle;
  }

  /* Sizes */
  .badge-sm {
    padding: 0.125rem var(--space-2);
    font-size: var(--text-xs);
  }

  .badge-md {
    padding: var(--space-1) var(--space-3);
    font-size: var(--text-sm);
  }

  /* Variants */
  .badge-default {
    background: var(--bg-iron);
    color: var(--text-muted);
  }

  .badge-success {
    background: var(--status-success-dim);
    color: var(--status-success);
  }

  .badge-warning {
    background: var(--status-warning-dim);
    color: var(--status-warning);
  }

  .badge-error {
    background: var(--status-error-dim);
    color: var(--status-error);
  }

  .badge-info {
    background: var(--status-info-dim);
    color: var(--status-info);
  }

  .badge-muted {
    background: var(--bg-slate);
    color: var(--text-dim);
  }
</style>
