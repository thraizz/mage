<script lang="ts">
  import type { CardView } from '$lib/generated/mage/v1/models';
  import { X } from '@lucide/svelte';
  import Card from './Card.svelte';

  interface Props {
    cards: CardView[];
    onClose: () => void;
  }

  let { cards, onClose }: Props = $props();
  let dialogRef: HTMLDialogElement | undefined = $state();

  $effect(() => {
    // Native dialog handles focus trapping and Esc key automatically
    dialogRef?.showModal();
    document.body.style.overflow = 'hidden';

    return () => {
      document.body.style.overflow = '';
    };
  });

  function handleCancel(e: Event) {
    e.preventDefault();
    onClose();
  }

  function handleClickOutside(e: MouseEvent) {
    // In a <dialog>, a click on the element itself is a click on the backdrop
    if (e.target === dialogRef) {
      onClose();
    }
  }
</script>

<dialog
  bind:this={dialogRef}
  oncancel={handleCancel}
  onclick={handleClickOutside}
  class="reveal-modal"
  aria-labelledby="reveal-dialog-title"
>
  <div class="dialog-content">
    <div class="dialog-header">
      <h2 id="reveal-dialog-title">Revealed Cards ({cards.length})</h2>
      <button type="button" class="close-button" onclick={onClose} aria-label="Close dialog">
        <X size={20} aria-hidden="true" />
        <span class="sr-only">Close dialog</span>
      </button>
    </div>

    <div class="cards-section" role="region" aria-live="polite">
      {#if cards.length === 0}
        <p class="no-cards">No cards to reveal</p>
      {:else}
        <div class="card-grid">
          {#each cards as card (card.id)}
            <div class="card-wrapper">
              <Card
                cardId={card.id}
                cardName={card.name}
                manaCost={card.manaCost}
                cardType={card.type}
                power={card.power}
                toughness={card.toughness}
                isSelected={false}
                size="normal"
              />
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <div class="dialog-footer">
      <button type="button" class="btn-primary" onclick={onClose}>Close</button>
    </div>
  </div>
</dialog>

<style>
  /* 1. THE BACKDROP (Replaces .overlay) */
  .reveal-modal::backdrop {
    background: rgba(0, 0, 0, 0.85);
    animation: fadeIn 0.2s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  /* 2. THE DIALOG CONTAINER (Replaces .dialog positioning) */
  .reveal-modal {
    background: #1a1f2e;
    border: 2px solid #3a4451;
    border-radius: 12px;
    padding: 0; /* Padding moved to .dialog-content */
    max-width: 900px;
    width: 90%;
    max-height: 85vh;
    color: #fff;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
    overflow: hidden; /* Let the inner content handle scroll */

    /* Animating the dialog itself */
    animation: slideUp 0.3s ease-out;
  }

  /* Native dialogs are centered by default, but this ensures it */
  .reveal-modal[open] {
    display: flex;
    flex-direction: column;
  }

  @keyframes slideUp {
    from {
      transform: translateY(20px);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }

  /* 3. INTERNAL STYLES (Kept mostly the same) */
  .dialog-content {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    max-height: inherit;
    overflow-y: auto;
  }

  .dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid #3a4451;
  }

  .dialog-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
  }

  .close-button {
    background: transparent;
    border: none;
    color: #9ca3af;
    cursor: pointer;
    padding: 0;
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
  }

  .close-button:hover {
    color: #fff;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    border: 0;
  }

  .cards-section {
    margin-bottom: 1.5rem;
  }

  .card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    padding: 0.5rem;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding-top: 1rem;
    border-top: 1px solid #3a4451;
  }

  .btn-primary {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    border: none;
    background: #667eea;
    color: white;
    transition: all 0.2s;
  }

  .btn-primary:hover {
    background: #5568d3;
    transform: translateY(-1px);
  }

  .no-cards {
    color: #6b7280;
    font-style: italic;
    padding: 2rem;
    text-align: center;
    background: rgba(255, 255, 255, 0.02);
    border-radius: 6px;
  }
</style>
