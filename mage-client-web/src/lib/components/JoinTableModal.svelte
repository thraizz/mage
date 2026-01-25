<script lang="ts">
  import type { Table } from '$lib/types/table';
  import type { Deck } from '$lib/types/deck';
  import { joinTable } from '$lib/api/table';
  import { fetchUserDecks, getDeckDetails } from '$lib/api/decks';
  import { toast } from '$lib/stores/toast';
  import Modal from './Modal.svelte';
  import LoadingSpinner from './LoadingSpinner.svelte';
  import { structuredCardsToText, type CardEntry } from '$lib/utils/deck-parser';

  // Props
  let {
    open = $bindable(false),
    table,

    onSuccess
  }: {
    open: boolean;
    table: Table | null;

    onSuccess: (tableId: string) => void;
  } = $props();

  // State
  let decks = $state<Deck[]>([]);
  let selectedDeckId = $state<string>('');
  let password = $state<string>('');
  let loading = $state(false);
  let loadingDecks = $state(false);
  let error = $state<string | null>(null);

  // Load decks when modal opens
  $effect(() => {
    if (open && table) {
      loadDecks();
    }
  });

  /**
   * Load user's decks for this format
   */
  async function loadDecks(): Promise<void> {
    if (!table) return;

    loadingDecks = true;
    error = null;
    selectedDeckId = '';

    try {
      const allDecks = await fetchUserDecks();
      // Filter decks by format
      decks = allDecks.filter((d: Deck) => d.format === table.format);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load decks';
      console.error('Failed to load decks:', err);
    } finally {
      loadingDecks = false;
    }
  }

  /**
   * Handle join
   */
  async function handleJoin(): Promise<void> {
    if (!table) return;

    if (!selectedDeckId) {
      error = 'Please select a deck';
      return;
    }

    if (table.hasPassword && !password.trim()) {
      error = 'Password is required';
      return;
    }

    loading = true;
    error = null;

    try {
      // Find selected deck
      const deck = decks.find((d) => d.id === selectedDeckId);
      if (!deck) {
        throw new Error('Selected deck not found');
      }

      console.log('[JoinTable] Selected deck summary:', {
        id: deck.id,
        name: deck.name,
        format: deck.format
      });

      // Fetch full deck details (including card lists)
      // fetchUserDecks only returns summary info without card details
      const fullDeck = await getDeckDetails(deck.id);

      console.log('[JoinTable] Full deck loaded:', {
        id: fullDeck.id,
        name: fullDeck.name,
        format: fullDeck.format,
        mainDeckCount: fullDeck.mainDeck.length,
        sideboardCount: fullDeck.sideboard.length,
        commanderCount: fullDeck.commanders.length
      });

      // Convert deck to text format
      const deckCards: CardEntry[] = [
        ...fullDeck.commanders.map((c) => ({
          name: c.cardName,
          quantity: c.quantity,
          section: 'commander' as const
        })),
        ...fullDeck.mainDeck.map((c) => ({
          name: c.cardName,
          quantity: c.quantity,
          section: 'main' as const
        })),
        ...fullDeck.sideboard.map((c) => ({
          name: c.cardName,
          quantity: c.quantity,
          section: 'sideboard' as const
        }))
      ];
      const deckList = structuredCardsToText(deckCards);

      console.log('[JoinTable] Deck text being sent:', {
        textLength: deckList.length,
        preview: deckList.substring(0, 500),
        fullText: deckList
      });

      // Join table with deck
      await joinTable(table.id, deckList, table.hasPassword ? password : undefined);

      console.log('[JoinTable] Join successful for table:', table.id);
      toast.success(`Joined table: ${table.name}`);
      onSuccess(table.id);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to join table';
      console.error('Failed to join table:', err);
      toast.error(error);
    } finally {
      loading = false;
    }
  }

  /**
   * Handle close
   */
  function handleClose(): void {
    if (loading) return;
    open = false;
    password = '';
    error = null;
  }
</script>

<Modal bind:open size="medium" closeOnBackdrop={!loading}>
  <div class="join-modal">
    <h2 class="modal-title">Join Table</h2>

    {#if table}
      <div class="table-info">
        <div class="info-row">
          <span class="info-label">Table:</span>
          <span class="info-value">{table.name}</span>
        </div>
        <div class="info-row">
          <span class="info-label">Format:</span>
          <span class="info-value format-badge">{table.format}</span>
        </div>
        <div class="info-row">
          <span class="info-label">Host:</span>
          <span class="info-value">{table.hostUsername}</span>
        </div>
        <div class="info-row">
          <span class="info-label">Players:</span>
          <span class="info-value">{table.players.length}/{table.maxPlayers}</span>
        </div>
      </div>

      {#if loadingDecks}
        <div class="loading-container">
          <LoadingSpinner size="medium" />
          <p class="loading-text">Loading your decks...</p>
        </div>
      {:else if decks.length === 0}
        <div class="no-decks">
          <p class="no-decks-message">
            You don't have any {table.format} decks. Please create one in the Deck Manager first.
          </p>
          <a href="/decks" class="create-deck-link">Go to Deck Manager</a>
        </div>
      {:else}
        <div class="form-group">
          <label for="deck-select" class="form-label">Select Deck</label>
          <select
            id="deck-select"
            class="form-select"
            bind:value={selectedDeckId}
            disabled={loading}
          >
            <option value="">-- Choose a deck --</option>
            {#each decks as deck}
              <option value={deck.id}>{deck.name} ({deck.cardCount} cards)</option>
            {/each}
          </select>
        </div>

        {#if table.hasPassword}
          <div class="form-group">
            <label for="password-input" class="form-label">Password</label>
            <input
              id="password-input"
              type="password"
              class="form-input"
              placeholder="Enter table password"
              bind:value={password}
              disabled={loading}
            />
          </div>
        {/if}

        {#if error}
          <div class="error-message">{error}</div>
        {/if}

        <div class="modal-actions">
          <button class="btn btn-secondary" onclick={handleClose} disabled={loading}>
            Cancel
          </button>
          <button
            class="btn btn-primary"
            onclick={handleJoin}
            disabled={loading || !selectedDeckId}
          >
            {#if loading}
              <LoadingSpinner size="small" />
              <span>Joining...</span>
            {:else}
              Join Table
            {/if}
          </button>
        </div>
      {/if}
    {/if}
  </div>
</Modal>

<style>
  .join-modal {
    padding: var(--space-6);
  }

  .modal-title {
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    font-weight: var(--weight-bold);
    color: var(--ci-scroll-parchment);
    margin: 0 0 var(--space-6) 0;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .table-info {
    background-color: var(--bg-slate);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    padding: var(--space-4);
    margin-bottom: var(--space-6);
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    padding: var(--space-2) 0;
  }

  .info-row:not(:last-child) {
    border-bottom: 1px solid var(--border-subtle);
  }

  .info-label {
    font-weight: var(--weight-semibold);
    color: var(--ci-swamp-obsidian);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-size: var(--text-xs);
  }

  .info-value {
    color: var(--ci-scroll-parchment);
    font-weight: var(--weight-medium);
  }

  .format-badge {
    background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
    color: var(--ci-scroll-parchment);
    padding: var(--space-1) var(--space-2);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    font-weight: var(--weight-semibold);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-8) 0;
  }

  .loading-text {
    color: var(--ci-swamp-obsidian);
    font-style: italic;
    margin: 0;
  }

  .no-decks {
    text-align: center;
    padding: var(--space-8) 0;
  }

  .no-decks-message {
    color: var(--ci-swamp-obsidian);
    margin: 0 0 var(--space-4) 0;
  }

  .create-deck-link {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
    color: var(--ci-scroll-parchment);
    border-radius: var(--radius-md);
    text-decoration: none;
    font-weight: var(--weight-semibold);
    font-size: var(--text-sm);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    transition: all var(--transition-base);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  }

  .create-deck-link:hover {
    background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
    box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
    transform: translateY(-1px);
  }

  .form-group {
    margin-bottom: var(--space-5);
  }

  .form-label {
    display: block;
    font-weight: var(--weight-semibold);
    color: var(--ci-scroll-parchment);
    margin-bottom: var(--space-2);
    font-size: var(--text-sm);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .form-select,
  .form-input {
    width: 100%;
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-md);
    font-size: var(--text-sm);
    transition: all var(--transition-fast);
    background-color: var(--bg-iron);
    color: var(--ci-scroll-parchment);
  }

  .form-select option {
    background: var(--bg-slate);
    color: var(--ci-scroll-parchment);
  }

  .form-input::placeholder {
    color: var(--text-ghost);
  }

  .form-select:focus,
  .form-input:focus {
    outline: none;
    border-color: var(--ci-jace-cloak);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
    background-color: var(--bg-obsidian);
  }

  .form-select:disabled,
  .form-input:disabled {
    background-color: var(--bg-slate);
    cursor: not-allowed;
    opacity: 0.5;
  }

  .error-message {
    background-color: rgba(255, 77, 77, 0.1);
    border: 1px solid var(--ci-mountain-ember);
    color: var(--ci-mountain-ember);
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-md);
    margin-bottom: var(--space-5);
    font-size: var(--text-sm);
  }

  .modal-actions {
    display: flex;
    gap: var(--space-3);
    justify-content: flex-end;
    margin-top: var(--space-6);
  }

  .btn {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-5);
    border-radius: var(--radius-md);
    font-weight: var(--weight-semibold);
    font-size: var(--text-sm);
    cursor: pointer;
    transition: all var(--transition-base);
    border: none;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-secondary {
    background-color: var(--bg-iron);
    color: var(--ci-scroll-parchment);
    border: 1px solid var(--border-default);
  }

  .btn-secondary:hover:not(:disabled) {
    background-color: var(--bg-steel);
    border-color: var(--border-strong);
  }

  .btn-primary {
    background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
    color: var(--ci-scroll-parchment);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
  }

  .btn-primary:hover:not(:disabled) {
    background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
    box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
    transform: translateY(-1px);
  }
</style>
