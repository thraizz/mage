<script lang="ts">
  import { getScryfallImageUrl } from '$lib/utils/scryfall';

  let {
    isDragging,
    cardName,
    position,
    isOverValidDrop,
    imageSize = 'normal'
  }: {
    isDragging: boolean;
    cardName: string | null;
    position: { x: number; y: number };
    isOverValidDrop: boolean;
    imageSize?: 'small' | 'normal';
  } = $props();

  const imageUrl = $derived(cardName ? getScryfallImageUrl(cardName, imageSize) : null);
</script>

{#if isDragging && cardName}
  <div class="drag-ghost" style="left: {position.x}px; top: {position.y}px;">
    <div class="drag-ghost-card" class:valid={isOverValidDrop}>
      {#if imageUrl}
        <img src={imageUrl} alt={cardName} class="drag-ghost-image" draggable="false" />
      {:else}
        <span class="drag-ghost-name">{cardName}</span>
      {/if}
    </div>
  </div>
{/if}

<style>
  .drag-ghost {
    position: fixed;
    pointer-events: none;
    z-index: 10000;
    transform: translate(-50%, -50%);
  }

  .drag-ghost-card {
    opacity: 0.8;
    transition: opacity 0.2s;
  }

  .drag-ghost-card.valid {
    opacity: 1;
  }

  .drag-ghost-image {
    width: 100px;
    height: 140px;
    object-fit: cover;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  }

  .drag-ghost-name {
    display: block;
    padding: 10px;
    background: rgba(0, 0, 0, 0.9);
    color: white;
    border-radius: 4px;
    font-size: 14px;
  }
</style>
