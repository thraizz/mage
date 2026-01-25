<script lang="ts">
  import { fade } from 'svelte/transition';
  import type { LoadingSize } from '$lib/types/loading';

  // Props
  export let size: LoadingSize = 'medium';
  export let label: string | undefined = undefined;
  export let overlay = false;
  export let color = 'var(--accent-gold)';

  function getSizeClass(): string {
    switch (size) {
      case 'small':
        return 'spinner-small';
      case 'large':
        return 'spinner-large';
      case 'medium':
      default:
        return 'spinner-medium';
    }
  }
</script>

{#if overlay}
  <!-- Overlay Mode -->
  <div class="spinner-overlay" transition:fade={{ duration: 200 }} role="status" aria-live="polite">
    <div class="spinner-wrapper">
      <div class="spinner {getSizeClass()}" style="border-top-color: {color}"></div>
      {#if label}
        <div class="spinner-label">{label}</div>
      {/if}
    </div>
  </div>
{:else}
  <!-- Inline Mode -->
  <div class="spinner-inline" role="status" aria-live="polite">
    <div class="spinner {getSizeClass()}" style="border-top-color: {color}"></div>
    {#if label}
      <div class="spinner-label">{label}</div>
    {/if}
  </div>
{/if}

<style>
  /* Overlay Mode */
  .spinner-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--modal-backdrop);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: var(--z-modal);
  }

  .spinner-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-4);
  }

  /* Inline Mode */
  .spinner-inline {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-2);
  }

  /* Spinner Animation */
  .spinner {
    border-radius: var(--radius-full);
    border-style: solid;
    border-color: var(--border-subtle);
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  /* Size Variants */
  .spinner-small {
    width: 1rem;
    height: 1rem;
    border-width: 2px;
  }

  .spinner-medium {
    width: 2rem;
    height: 2rem;
    border-width: 3px;
  }

  .spinner-large {
    width: 3rem;
    height: 3rem;
    border-width: 4px;
  }

  /* Label */
  .spinner-label {
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-bright);
    text-align: center;
  }

  .spinner-inline .spinner-label {
    color: var(--text-muted);
  }

  /* Inline Mode spinner border */
  .spinner-inline .spinner {
    border-color: var(--bg-steel);
  }
</style>
