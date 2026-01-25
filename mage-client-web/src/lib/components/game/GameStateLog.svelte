<script lang="ts">
  import { playtestGameStore } from '$lib/stores/playtest-game';
  import { toast } from '$lib/stores/toast';
  // Game components
  import Copy from '@lucide/svelte/icons/copy';
  import X from '@lucide/svelte/icons/x';

  let {
    open = $bindable(false),
    variant = 'slideout'
  }: {
    open?: boolean;
    variant?: 'slideout' | 'inline';
  } = $props();

  // Game log
  const gameLog = $derived($playtestGameStore.log || []);

  /**
   * Copy game log to clipboard
   */
  async function handleCopyLog(): Promise<void> {
    const logText = playtestGameStore.buildLogText($playtestGameStore);
    try {
      if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(logText);
        toast.success('Game log copied to clipboard!');
        return;
      }
      // Fallback for older browsers
      const textarea = document.createElement('textarea');
      textarea.value = logText;
      textarea.style.position = 'fixed';
      textarea.style.top = '0';
      textarea.style.left = '0';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.focus();
      textarea.select();
      const ok = document.execCommand('copy');
      document.body.removeChild(textarea);
      if (ok) {
        toast.success('Game log copied to clipboard!');
      } else {
        toast.error('Failed to copy log');
      }
    } catch (err) {
      console.error('Failed to copy log to clipboard:', err);
      toast.error('Failed to copy log');
    }
  }

  /**
   * Close the overlay
   */
  function close(): void {
    open = false;
  }

  /**
   * Handle click outside to close
   */
  function handleBackdropClick(e: MouseEvent): void {
    if (e.target === e.currentTarget) {
      close();
    }
  }

  /**
   * Handle escape key to close
   */
  function handleKeydown(e: KeyboardEvent): void {
    if (e.key === 'Escape' && open && variant === 'slideout') {
      close();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if variant === 'slideout'}
  <!-- Backdrop -->
  {#if open}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="log-overlay-backdrop" onclick={handleBackdropClick}></div>
  {/if}

  <!-- Slide-out Panel -->
  <div class="game-state-log-overlay" class:open>
    <div class="log-header">
      <div class="header-title">
        <span class="header-icon">📋</span>
        <h3>Game State Log</h3>
        <span class="entry-count">{gameLog.length}</span>
      </div>
      <div class="header-actions">
        <button
          class="copy-btn"
          onclick={handleCopyLog}
          title="Copy log to clipboard"
          aria-label="Copy log to clipboard"
        >
          <Copy size={16} aria-hidden="true" />
        </button>
        <button class="close-btn" onclick={close} title="Close (Esc)" aria-label="Close">
          <X size={20} aria-hidden="true" />
        </button>
      </div>
    </div>

    <div class="log-entries">
      {#if gameLog.length === 0}
        <div class="empty-state">
          <p>No events logged yet</p>
        </div>
      {:else}
        <div class="debug-log-entries">
          {#each gameLog.slice().reverse() as entry (entry.id)}
            <div class="debug-log-entry">
              <span class="debug-log-time">
                {new Date(entry.at).toLocaleTimeString([], {
                  hour: '2-digit',
                  minute: '2-digit',
                  second: '2-digit'
                })}
              </span>
              <span class="debug-log-turn">T{entry.turn}</span>
              <span class="debug-log-kind">{entry.kind}</span>
              <span class="debug-log-message">{entry.message}</span>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{:else}
  <!-- Inline variant (for debug overlay) -->
  <section class="debug-section">
    <div class="debug-section-header">
      <span>Game State Log ({gameLog.length} events)</span>
      <button
        class="debug-copy-btn"
        onclick={handleCopyLog}
        title="Copy log to clipboard"
        aria-label="Copy log to clipboard"
      >
        <Copy size={16} aria-hidden="true" />
        <span>Copy</span>
      </button>
    </div>
    <div class="debug-log-container">
      {#if gameLog.length === 0}
        <div class="debug-log-empty">No events logged yet</div>
      {:else}
        <div class="debug-log-entries">
          {#each gameLog.slice().reverse() as entry (entry.id)}
            <div class="debug-log-entry">
              <span class="debug-log-time">
                {new Date(entry.at).toLocaleTimeString([], {
                  hour: '2-digit',
                  minute: '2-digit',
                  second: '2-digit'
                })}
              </span>
              <span class="debug-log-turn">T{entry.turn}</span>
              <span class="debug-log-kind">{entry.kind}</span>
              <span class="debug-log-message">{entry.message}</span>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>
{/if}

<style>
  /* Backdrop */
  .log-overlay-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 89;
    animation: fade-in 0.2s ease;
  }

  @keyframes fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  /* Slide-out Panel */
  .game-state-log-overlay {
    position: fixed;
    right: 0;
    top: 0;
    bottom: 0;
    width: 500px;
    max-width: 90vw;
    background: #141821;
    border-left: 2px solid #2a3441;
    display: flex;
    flex-direction: column;
    z-index: 90;
    transform: translateX(100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: -4px 0 24px rgba(0, 0, 0, 0.5);
  }

  .game-state-log-overlay.open {
    transform: translateX(0);
  }

  /* Header */
  .log-header {
    padding: 1rem 1.25rem;
    background: #1a1f2e;
    border-bottom: 2px solid #2a3441;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-shrink: 0;
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .header-icon {
    font-size: 1.25rem;
  }

  .log-header h3 {
    margin: 0;
    font-size: 1.125rem;
    font-weight: 600;
    color: white;
  }

  .entry-count {
    font-size: 0.75rem;
    color: #9ca3af;
    font-weight: 600;
    padding: 0.25rem 0.5rem;
    background: #0f1419;
    border-radius: 4px;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .copy-btn,
  .close-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: 1px solid #2a3441;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    color: #9ca3af;
  }

  .copy-btn:hover,
  .close-btn:hover {
    background: #2a3441;
    border-color: #374151;
    color: #fff;
  }

  /* Entries List */
  .log-entries {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    background: #0f1419;
    padding: 0.5rem;
  }

  /* Scrollbar Styling */
  .log-entries::-webkit-scrollbar {
    width: 8px;
  }

  .log-entries::-webkit-scrollbar-track {
    background: #0f1419;
  }

  .log-entries::-webkit-scrollbar-thumb {
    background: #2a3441;
    border-radius: 4px;
  }

  .log-entries::-webkit-scrollbar-thumb:hover {
    background: #3a4451;
  }

  /* Empty State */
  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 2rem;
  }

  .empty-state p {
    color: #6b7280;
    font-size: 0.875rem;
    font-style: italic;
    text-align: center;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .game-state-log-overlay {
      width: 100%;
      max-width: 100%;
    }
  }
</style>
