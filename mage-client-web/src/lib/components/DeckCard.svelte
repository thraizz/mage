<script lang="ts">
  import type { Deck } from '$lib/types/deck';
  import DeckCard from '$lib/components/deck/DeckCard.svelte';

  export let deck: Deck;
  export let onclick: ((e: MouseEvent) => void) | undefined = undefined;
</script>

<DeckCard {deck} {onclick} />
