<script lang="ts">
  import type { Snippet } from 'svelte';
  import SearchIcon from '@lucide/svelte/icons/search';
  import TableIcon from '@lucide/svelte/icons/table';
  import LayersIcon from '@lucide/svelte/icons/layers';
  import UsersIcon from '@lucide/svelte/icons/users';
  import Gamepad2Icon from '@lucide/svelte/icons/gamepad-2';

  interface Props {
    title: string;
    description?: string;
    icon?: 'search' | 'table' | 'deck' | 'player' | 'game';
    children?: Snippet; // For action buttons
  }

  let { title, description = '', icon = 'search', children }: Props = $props();

  const iconComponents = {
    search: SearchIcon,
    table: TableIcon,
    deck: LayersIcon,
    player: UsersIcon,
    game: Gamepad2Icon
  } as const;

  const IconComponent = $derived(iconComponents[icon] ?? SearchIcon);
</script>

<div class="empty-state">
  <div class="empty-icon">
    <IconComponent size={64} aria-hidden="true" />
  </div>
  <h3 class="empty-title">{title}</h3>
  {#if description}
    <p class="empty-description">{description}</p>
  {/if}
  {#if children}
    <div class="empty-actions">
      {@render children()}
    </div>
  {/if}
</div>

<style>
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-12) var(--space-4);
    text-align: center;
  }

  .empty-icon {
    width: 4rem;
    height: 4rem;
    margin-bottom: var(--space-4);
    color: var(--text-ghost);
  }

  .empty-icon :global(svg) {
    width: 100%;
    height: 100%;
  }

  .empty-title {
    font-family: var(--font-display);
    font-size: var(--text-xl);
    font-weight: var(--weight-semibold);
    color: var(--text-muted);
    margin: 0 0 var(--space-2);
  }

  .empty-description {
    font-size: var(--text-base);
    color: var(--text-dim);
    margin: 0 0 var(--space-6);
    max-width: 24rem;
  }

  .empty-actions {
    display: flex;
    gap: var(--space-3);
  }
</style>
