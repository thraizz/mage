<script lang="ts">
  import Card from './Card.svelte';
  import {
    visualStackStore,
    visualStackItems,
    visualStackIsEmpty,
    visualStackCount,
    type VisualStackItem
  } from '$lib/stores/visual-stack';
  import { getScryfallImageUrl } from '$lib/utils/scryfall';

  // Props
  let {
    isOpen = false,
    onResolve,
    onRemove,
    onClear
  }: {
    isOpen?: boolean;
    /** Called when resolving/removing a stack item - receives the stack item ID */
    onResolve?: (itemId: string) => void;
    onRemove?: (itemId: string) => void;
    onClear?: () => void;
  } = $props();

  // Internal state
  let draggedItemLocalId = $state<string | null>(null);
  let dragOverIndex = $state<number | null>(null);
  let editingNoteLocalId = $state<string | null>(null);
  let editingNoteText = $state('');

  // Derived from store
  const items = $derived($visualStackItems);
  const isEmpty = $derived($visualStackIsEmpty);
  const count = $derived($visualStackCount);

  // Items displayed in reverse order (top of stack = last item = displayed first)
  const displayItems = $derived([...items].reverse());

  /**
   * Handle resolve top button click - calls server to remove from stack
   */
  function handleResolveTop(): void {
    const topItem = displayItems[0];
    if (topItem && onResolve) {
      // Use localId which is the stack item ID from server
      onResolve(topItem.localId);
    }
  }

  /**
   * Handle remove specific item - calls server to remove from stack
   */
  function handleRemoveItem(localId: string): void {
    if (onRemove) {
      onRemove(localId);
    }
  }

  /**
   * Handle clear all - removes all items from server stack
   */
  function handleClearAll(): void {
    if (confirm('Clear all items from the stack?')) {
      if (onClear) {
        onClear();
      } else {
        // Fallback to removing each item individually
        for (const item of items) {
          if (onRemove) {
            onRemove(item.localId);
          }
        }
      }
    }
  }

  /**
   * Close the sidebar
   */
  function handleClose(): void {
    visualStackStore.setOpen(false);
  }

  // Internal drag-and-drop for reordering
  function handleDragStart(event: DragEvent, localId: string): void {
    if (!event.dataTransfer) return;
    event.dataTransfer.effectAllowed = 'move';
    event.dataTransfer.setData('text/plain', localId);
    draggedItemLocalId = localId;
  }

  function handleDragOver(event: DragEvent, displayIndex: number): void {
    event.preventDefault();
    if (!event.dataTransfer) return;
    event.dataTransfer.dropEffect = 'move';
    dragOverIndex = displayIndex;
  }

  function handleDragLeave(): void {
    dragOverIndex = null;
  }

  function handleDrop(event: DragEvent, displayIndex: number): void {
    event.preventDefault();
    if (!draggedItemLocalId) return;

    // Convert display indices to actual array indices (reverse mapping)
    const fromDisplayIndex = displayItems.findIndex((item) => item.localId === draggedItemLocalId);
    if (fromDisplayIndex === -1) return;

    // Convert from display order (reversed) to actual order
    const fromActualIndex = items.length - 1 - fromDisplayIndex;
    const toActualIndex = items.length - 1 - displayIndex;

    visualStackStore.reorderItems(fromActualIndex, toActualIndex);

    draggedItemLocalId = null;
    dragOverIndex = null;
  }

  function handleDragEnd(): void {
    draggedItemLocalId = null;
    dragOverIndex = null;
  }

  /**
   * Start editing a note
   */
  function startEditingNote(item: VisualStackItem): void {
    editingNoteLocalId = item.localId;
    editingNoteText = item.note || '';
  }

  /**
   * Save the note
   */
  function saveNote(): void {
    if (editingNoteLocalId) {
      visualStackStore.updateNote(editingNoteLocalId, editingNoteText);
      editingNoteLocalId = null;
      editingNoteText = '';
    }
  }

  /**
   * Cancel editing note
   */
  function cancelEditingNote(): void {
    editingNoteLocalId = null;
    editingNoteText = '';
  }

  /**
   * Get stack position label (1 = top/resolves next)
   */
  function getPositionLabel(displayIndex: number): number {
    return displayIndex + 1;
  }

  /**
   * Check if item is the top of stack
   */
  function isTopOfStack(displayIndex: number): boolean {
    return displayIndex === 0;
  }

  /**
   * Get zone display label
   */
  function getZoneLabel(zone: string): string {
    const labels: Record<string, string> = {
      hand: 'Hand',
      battlefield: 'Battlefield',
      graveyard: 'Graveyard',
      exile: 'Exile',
      library: 'Library',
      command: 'Command',
      stack: 'Stack'
    };
    return labels[zone] || zone;
  }
</script>

<aside class="visual-stack-sidebar" class:open={isOpen}>
  <div class="sidebar-header">
    <div class="header-title">
      <span class="stack-icon">📚</span>
      <h3>Stack</h3>
      {#if !isEmpty}
        <span class="item-count">{count}</span>
      {/if}
    </div>
    <div class="header-actions">
      {#if !isEmpty}
        <button class="resolve-btn" onclick={handleResolveTop} title="Resolve top item">
          ✓ Resolve
        </button>
        <button class="clear-btn" onclick={handleClearAll} title="Clear all"> 🗑️ </button>
      {/if}
      <button class="close-btn" onclick={handleClose} title="Close"> → </button>
    </div>
  </div>

  <div class="sidebar-content">
    {#if isEmpty}
      <div class="empty-state">
        <div class="empty-icon">📚</div>
        <p>Drag cards here</p>
        <p class="hint">to track the stack</p>
      </div>
    {:else}
      <div class="stack-list">
        {#each displayItems as item, displayIndex (item.localId)}
          {@const isTop = isTopOfStack(displayIndex)}
          {@const isDragging = draggedItemLocalId === item.localId}
          {@const isDragOver = dragOverIndex === displayIndex}
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            class="stack-item"
            class:is-top={isTop}
            class:is-dragging={isDragging}
            class:drag-over={isDragOver}
            draggable="true"
            ondragstart={(e) => handleDragStart(e, item.localId)}
            ondragover={(e) => handleDragOver(e, displayIndex)}
            ondragleave={handleDragLeave}
            ondrop={(e) => handleDrop(e, displayIndex)}
            ondragend={handleDragEnd}
          >
            <div class="item-position">
              <span class="position-number" class:top={isTop}>
                {getPositionLabel(displayIndex)}
              </span>
            </div>

            <div class="item-card">
              <Card
                cardId={item.cardId}
                cardName={item.cardName}
                imageUrl={item.imageUrl || getScryfallImageUrl(item.cardName)}
                size="small"
              />
            </div>

            <div class="item-info">
              <div class="card-name">{item.cardName}</div>
              <div class="source-zone">{getZoneLabel(item.sourceZone)}</div>

              {#if editingNoteLocalId === item.localId}
                <div class="note-editor">
                  <input
                    type="text"
                    bind:value={editingNoteText}
                    placeholder="Note..."
                    onkeydown={(e) => {
                      if (e.key === 'Enter') saveNote();
                      if (e.key === 'Escape') cancelEditingNote();
                    }}
                  />
                  <button class="note-save" onclick={saveNote}>✓</button>
                </div>
              {:else if item.note}
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div class="item-note" onclick={() => startEditingNote(item)}>
                  {item.note}
                </div>
              {:else}
                <button class="add-note-btn" onclick={() => startEditingNote(item)}>
                  + note
                </button>
              {/if}
            </div>

            <button
              class="remove-btn"
              onclick={() => handleRemoveItem(item.localId)}
              title="Remove"
            >
              ✕
            </button>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="sidebar-footer">
    {#if !isEmpty}
      <span class="resolves-label">↑ Resolves first</span>
    {:else}
      <span class="order-hint">Drag cards here or reorder</span>
    {/if}
  </div>
</aside>

<style>
  /* Sidebar container */
  .visual-stack-sidebar {
    width: 0;
    height: 100%;
    background: #1a1f2e;
    border-left: 2px solid #2a3441;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: width 0.2s ease-out;
    flex-shrink: 0;
  }

  .visual-stack-sidebar.open {
    width: 280px;
    border-left-color: #667eea;
  }

  /* Header */
  .sidebar-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.75rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-bottom: 1px solid #2a3441;
    min-height: 48px;
  }

  .header-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stack-icon {
    font-size: 1rem;
  }

  .sidebar-header h3 {
    margin: 0;
    font-size: 0.875rem;
    font-weight: 600;
    color: white;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .item-count {
    font-size: 0.75rem;
    color: white;
    background: rgba(0, 0, 0, 0.3);
    padding: 0.125rem 0.375rem;
    border-radius: 4px;
    font-weight: 600;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .resolve-btn {
    padding: 0.25rem 0.5rem;
    background: #10b981;
    border: none;
    border-radius: 4px;
    font-size: 0.625rem;
    font-weight: 600;
    color: white;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
  }

  .resolve-btn:hover {
    background: #059669;
  }

  .clear-btn,
  .close-btn {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.2);
    border: none;
    border-radius: 4px;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
    transition: all 0.2s;
  }

  .clear-btn:hover,
  .close-btn:hover {
    background: rgba(0, 0, 0, 0.4);
    color: white;
  }

  /* Content */
  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem;
  }

  /* Empty State */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem 1rem;
    gap: 0.5rem;
    height: 100%;
    min-height: 200px;
  }

  .empty-icon {
    font-size: 2rem;
    opacity: 0.3;
  }

  .empty-state p {
    color: #6b7280;
    font-size: 0.75rem;
    margin: 0;
    text-align: center;
  }

  .empty-state .hint {
    font-size: 0.625rem;
    font-style: italic;
  }

  /* Stack List */
  .stack-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  /* Stack Item */
  .stack-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    background: #141821;
    border: 1px solid #2a3441;
    border-radius: 6px;
    transition: all 0.15s;
    cursor: grab;
  }

  .stack-item:active {
    cursor: grabbing;
  }

  .stack-item.is-top {
    border-color: #fbbf24;
    background: rgba(251, 191, 36, 0.1);
  }

  .stack-item.is-dragging {
    opacity: 0.4;
  }

  .stack-item.drag-over {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.15);
  }

  .stack-item:hover {
    border-color: #4b5563;
  }

  /* Position indicator */
  .item-position {
    flex-shrink: 0;
  }

  .position-number {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #2a3441;
    border-radius: 50%;
    font-size: 0.625rem;
    font-weight: 700;
    color: #9ca3af;
  }

  .position-number.top {
    background: #fbbf24;
    color: #000;
  }

  /* Card thumbnail */
  .item-card {
    flex-shrink: 0;
    width: 80px;
  }

  /* Item info */
  .item-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .card-name {
    font-size: 0.75rem;
    font-weight: 600;
    color: white;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .source-zone {
    font-size: 0.625rem;
    color: #6b7280;
  }

  .item-note {
    font-size: 0.625rem;
    color: #9ca3af;
    background: rgba(102, 126, 234, 0.1);
    padding: 0.125rem 0.25rem;
    border-radius: 3px;
    cursor: pointer;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .item-note:hover {
    background: rgba(102, 126, 234, 0.2);
  }

  .add-note-btn {
    font-size: 0.5rem;
    color: #4b5563;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    text-align: left;
  }

  .add-note-btn:hover {
    color: #6b7280;
  }

  /* Note editor */
  .note-editor {
    display: flex;
    gap: 0.25rem;
  }

  .note-editor input {
    flex: 1;
    min-width: 0;
    padding: 0.125rem 0.25rem;
    background: #0d1117;
    border: 1px solid #2a3441;
    border-radius: 3px;
    font-size: 0.625rem;
    color: white;
  }

  .note-editor input:focus {
    outline: none;
    border-color: #667eea;
  }

  .note-save {
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #10b981;
    border: none;
    border-radius: 3px;
    font-size: 0.625rem;
    color: white;
    cursor: pointer;
  }

  /* Remove button */
  .remove-btn {
    flex-shrink: 0;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: 4px;
    font-size: 0.625rem;
    color: #4b5563;
    cursor: pointer;
    transition: all 0.15s;
  }

  .remove-btn:hover {
    background: #ef4444;
    color: white;
  }

  /* Footer */
  .sidebar-footer {
    padding: 0.5rem;
    background: #141821;
    border-top: 1px solid #2a3441;
    text-align: center;
  }

  .order-hint,
  .resolves-label {
    font-size: 0.625rem;
    color: #6b7280;
  }

  .resolves-label {
    color: #fbbf24;
  }

  /* Scrollbar */
  .sidebar-content::-webkit-scrollbar {
    width: 4px;
  }

  .sidebar-content::-webkit-scrollbar-track {
    background: #0d1117;
  }

  .sidebar-content::-webkit-scrollbar-thumb {
    background: #667eea;
    border-radius: 2px;
  }
</style>
