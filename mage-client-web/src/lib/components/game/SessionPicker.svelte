<script lang="ts">
  import type { PlaytestSessionMeta } from '$lib/stores/playtest-game';

  interface Props {
    show: boolean;
    sessions: PlaytestSessionMeta[];
    onRestore: (sessionId: string) => void;
    onDelete: (sessionId: string) => void;
    onClose: () => void;
    onNewPlaytest: () => void;
  }

  let { show, sessions, onRestore, onDelete, onClose, onNewPlaytest }: Props = $props();
</script>

{#if show}
  <div class="session-picker-overlay">
    <div class="session-picker-modal">
      <h2>Restore Playtest Session</h2>
      <p class="session-picker-hint">
        Select a recent playtest session to continue, or start a new one.
      </p>

      {#if sessions.length > 0}
        <div class="sessions-list">
          {#each sessions as session (session.id)}
            <div class="session-card">
              <div class="session-info">
                <div class="session-label">{session.label}</div>
                <div class="session-meta">
                  {session.playerCount} players · Turn {session.turn} ·
                  {new Date(session.savedAt).toLocaleDateString()}
                  {new Date(session.savedAt).toLocaleTimeString([], {
                    hour: '2-digit',
                    minute: '2-digit'
                  })}
                </div>
              </div>
              <div class="session-actions">
                <button class="btn-restore" onclick={() => onRestore(session.id)}> Restore </button>
                <button
                  class="btn-delete"
                  onclick={() => onDelete(session.id)}
                  title="Delete session"
                >
                  ✕
                </button>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <p class="no-sessions">No saved sessions found.</p>
      {/if}

      <div class="session-picker-actions">
        <button class="btn-back" onclick={onClose}> Back</button>
        <button class="btn-primary" onclick={onNewPlaytest}> Start New Playtest </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .session-picker-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.5);
    padding: 2rem;
    z-index: 100;
  }

  .session-picker-modal {
    background: rgba(26, 31, 46, 0.98);
    border: 1px solid #2a3441;
    border-radius: 12px;
    padding: 2rem;
    max-width: 700px;
    width: 100%;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .session-picker-modal h2 {
    margin: 0;
    color: #f8fafc;
    font-size: 1.5rem;
  }

  .session-picker-hint {
    margin: 0;
    color: #94a3b8;
    font-size: 0.875rem;
  }

  .sessions-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-height: 400px;
    overflow-y: auto;
    padding: 0.5rem;
    margin: -0.5rem;
  }

  .session-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: rgba(17, 24, 39, 0.8);
    border: 1px solid #2a3441;
    border-radius: 8px;
    transition: all 0.2s;
    gap: 1rem;
  }

  .session-card:hover {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.1);
  }

  .session-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .session-label {
    font-weight: 600;
    color: #f8fafc;
    font-size: 0.9375rem;
  }

  .session-meta {
    font-size: 0.8125rem;
    color: #94a3b8;
  }

  .session-actions {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  .btn-restore {
    padding: 0.5rem 1rem;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    transition: background 0.2s;
  }

  .btn-restore:hover {
    background: #5568d3;
  }

  .btn-delete {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 6px;
    color: #ef4444;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 1rem;
  }

  .btn-delete:hover {
    background: rgba(239, 68, 68, 0.2);
    border-color: rgba(239, 68, 68, 0.5);
  }

  .no-sessions {
    color: #94a3b8;
    text-align: center;
    padding: 2rem;
    font-style: italic;
  }

  .session-picker-actions {
    display: flex;
    gap: 0.75rem;
    justify-content: flex-end;
    padding-top: 1rem;
    border-top: 1px solid #2a3441;
  }
</style>
