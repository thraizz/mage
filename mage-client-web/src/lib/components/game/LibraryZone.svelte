<script lang="ts">
  import { toPossessiveName } from '$lib/utils/localization';
  /**
   * LibraryZone - Compact button to view library/deck contents
   * Styled like Graveyard and ExileZone buttons.
   * Clicking opens the LibrarySearch modal via the searchLibrary() API.
   */
  import Library from '@lucide/svelte/icons/library';

  // Props
  let {
    libraryCount = 0,
    playerName = 'Player',
    isOpponent = false,
    onSearch = () => {},
    onContextMenu = () => {}
  }: {
    libraryCount?: number;
    playerName?: string;
    isOpponent?: boolean;
    onSearch?: () => void;
    onContextMenu?: (event: MouseEvent) => void;
  } = $props();

  // Derived values
  const isEmpty = $derived(libraryCount === 0);

  /**
   * Handle click - trigger library search
   */
  function handleClick(): void {
    if (!isEmpty && !isOpponent) {
      onSearch();
    }
  }
</script>

<button
  class="library-compact"
  class:has-cards={!isEmpty}
  class:opponent={isOpponent}
  onclick={handleClick}
  oncontextmenu={(e) => {
    if (!isOpponent) {
      e.preventDefault();
      onContextMenu(e);
    }
  }}
  disabled={isOpponent}
  title="{toPossessiveName(playerName)} Library ({libraryCount} cards){isEmpty || isOpponent
    ? ''
    : ' - Click to view, Right-click for actions'}"
>
  <span class="library-icon">
    <Library size={14} />
  </span>
  <span class="library-label">Deck</span>
  <span class="card-count" class:zero={isEmpty}>{libraryCount}</span>
</button>

<style>
  /* Compact Library Button - styled like Graveyard */
  .library-compact {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: linear-gradient(135deg, rgba(26, 31, 46, 0.6) 0%, rgba(17, 24, 39, 0.6) 100%);
    border: 1px solid rgba(34, 197, 94, 0.3);
    border-radius: 8px;
    min-height: 36px;
    cursor: default;
    transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
    color: inherit;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  }

  .library-compact.has-cards:not(.opponent) {
    cursor: pointer;
    background: linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(26, 31, 46, 0.9) 100%);
    border-color: rgba(34, 197, 94, 0.5);
  }

  .library-compact.has-cards:not(.opponent):hover {
    background: linear-gradient(135deg, rgba(34, 197, 94, 0.25) 0%, rgba(42, 52, 65, 0.9) 100%);
    border-color: rgba(34, 197, 94, 0.7);
    box-shadow:
      0 2px 6px rgba(0, 0, 0, 0.3),
      0 1px 3px rgba(34, 197, 94, 0.2);
    transform: translateY(-1px);
  }

  .library-compact.opponent {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .library-compact:disabled {
    cursor: not-allowed;
  }

  .library-icon {
    font-size: 0.875rem;
    opacity: 0.7;
  }

  .library-compact.has-cards .library-icon {
    opacity: 1;
  }

  .library-label {
    font-size: 0.6875rem;
    color: #22c55e;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .card-count {
    font-size: 0.75rem;
    font-weight: 700;
    color: #22c55e;
    background: rgba(34, 197, 94, 0.2);
    padding: 0.125rem 0.375rem;
    border-radius: 4px;
    min-width: 1.25rem;
    text-align: center;
  }

  .card-count.zero {
    color: #4b5563;
    background: transparent;
  }
</style>
