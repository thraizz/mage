<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title?: string;
    padding?: 'none' | 'sm' | 'md' | 'lg';
    variant?: 'default' | 'elevated' | 'bordered';
    children?: Snippet;
    header?: Snippet;
    footer?: Snippet;
  }

  let {
    title = '',
    padding = 'md',
    variant = 'default',
    children,
    header,
    footer
  }: Props = $props();
</script>

<div class="panel panel-{variant}">
  {#if title || header}
    <div class="panel-header panel-pad-{padding}">
      {#if header}
        {@render header()}
      {:else}
        <h3 class="panel-title">{title}</h3>
      {/if}
    </div>
  {/if}

  <div class="panel-content panel-pad-{padding}">
    {#if children}
      {@render children()}
    {/if}
  </div>

  {#if footer}
    <div class="panel-footer panel-pad-{padding}">
      {@render footer()}
    </div>
  {/if}
</div>

<style>
  .panel {
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  /* Variants */
  .panel-default {
    background: var(--bg-obsidian);
  }

  .panel-elevated {
    background: var(--bg-slate);
    box-shadow: var(--shadow-md);
  }

  .panel-bordered {
    background: var(--bg-obsidian);
    border: 1px solid var(--border-subtle);
  }

  /* Header */
  .panel-header {
    border-bottom: 1px solid var(--border-subtle);
  }

  .panel-title {
    font-family: var(--font-display);
    font-size: var(--text-lg);
    font-weight: var(--weight-semibold);
    color: var(--text-bright);
    margin: 0;
  }

  /* Footer */
  .panel-footer {
    border-top: 1px solid var(--border-subtle);
    background: var(--bg-slate);
  }

  /* Padding variants */
  .panel-pad-none {
    padding: 0;
  }

  .panel-pad-sm {
    padding: var(--space-3);
  }

  .panel-pad-md {
    padding: var(--space-4);
  }

  .panel-pad-lg {
    padding: var(--space-6);
  }
</style>
