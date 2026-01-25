<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    loading,
    error,
    isInitialized,
    onRetry,
    children
  }: {
    loading: boolean;
    error: string | null;
    isInitialized: boolean;
    onRetry: () => void;
    children?: Snippet;
  } = $props();
</script>

<div class="game-container">
  {#if loading}
    <div class="loading-overlay">
      <div class="spinner"></div>
      <p>Loading game...</p>
    </div>
  {:else if error}
    <div class="error-overlay">
      <div class="error-icon">⚠️</div>
      <h2>Error</h2>
      <p>{error}</p>
      <button class="btn-primary" onclick={onRetry}> Return to Lobby </button>
    </div>
  {:else if !isInitialized}
    <div class="loading-overlay">
      <p>Initializing game state...</p>
    </div>
  {:else}
    {@render children?.()}
  {/if}
</div>

<style>
  .game-container {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
    position: relative;
  }

  .loading-overlay,
  .error-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.9);
    z-index: 10000;
    color: white;
  }

  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid rgba(255, 255, 255, 0.1);
    border-top-color: #4a90e2;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 20px;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .error-icon {
    font-size: 4rem;
    margin-bottom: 20px;
  }

  .btn-primary {
    margin-top: 20px;
    padding: 10px 20px;
    background: #4a90e2;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }

  .btn-primary:hover {
    background: #357abd;
  }
</style>
