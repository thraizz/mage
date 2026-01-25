<script lang="ts" module>
  export interface MenuAction {
    label?: string;
    icon?: string;
    divider?: boolean;
    submenu?: MenuAction[];
    onClick?: () => void;
    disabled?: boolean;
  }
</script>

<script lang="ts">
  interface Props {
    position: { x: number; y: number };
    deckCount: number;
    playerName: string;
    onClose: () => void;
    actions: MenuAction[];
  }

  let { position, deckCount, playerName, onClose, actions }: Props = $props();

  let openSubmenu: string | null = $state(null);
  let submenuPosition: { x: number; y: number } | null = $state(null);

  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.context-menu') && !target.closest('.submenu')) {
      onClose();
    }
  }

  function handleAction(action: MenuAction) {
    if (action.disabled) return;

    if (action.submenu) {
      // Toggle submenu
      if (openSubmenu === action.label) {
        openSubmenu = null;
        submenuPosition = null;
      } else {
        openSubmenu = action.label ?? null;
      }
    } else if (action.onClick) {
      action.onClick();
      onClose();
    }
  }

  function handleMouseEnter(action: MenuAction, event: MouseEvent) {
    if (!action.submenu) {
      openSubmenu = null;
      submenuPosition = null;
      return;
    }

    const target = event.currentTarget as HTMLElement;
    const rect = target.getBoundingClientRect();

    // Position submenu to the right of the menu item
    let x = rect.right;
    let y = rect.top;

    // Check if submenu would overflow right edge
    const menuWidth = 200; // Approximate submenu width
    if (x + menuWidth > window.innerWidth) {
      // Flip to left
      x = rect.left - menuWidth;
    }

    // Check if submenu would overflow bottom
    const estimatedHeight = (action.submenu?.length || 0) * 40; // Approximate item height
    if (y + estimatedHeight > window.innerHeight) {
      y = window.innerHeight - estimatedHeight - 10;
    }

    openSubmenu = action.label ?? null;
    submenuPosition = { x, y };
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      if (openSubmenu) {
        openSubmenu = null;
        submenuPosition = null;
      } else {
        onClose();
      }
    }
  }
</script>

<svelte:window onclick={handleClickOutside} onkeydown={handleKeydown} />

<div class="context-menu" style="left: {position.x}px; top: {position.y}px;" role="menu">
  <div class="menu-header">
    <span class="deck-info">
      {playerName}'s Deck
      <span class="deck-count">({deckCount} cards)</span>
    </span>
    <button class="close-btn" onclick={onClose}>×</button>
  </div>

  {#each actions as action}
    {#if action.divider}
      <div class="menu-divider"></div>
    {:else}
      <button
        class="menu-item {action.submenu ? 'has-submenu' : ''} {action.disabled ? 'disabled' : ''}"
        onclick={() => handleAction(action)}
        onmouseenter={(e) => handleMouseEnter(action, e)}
        disabled={action.disabled}
      >
        {#if action.icon}
          <span class="icon">{action.icon}</span>
        {/if}
        <span class="label">{action.label}</span>
        {#if action.submenu}
          <span class="submenu-arrow">▶</span>
        {/if}
      </button>
    {/if}
  {/each}
</div>

{#if openSubmenu && submenuPosition}
  {@const submenuActions = actions.find((a) => a.label === openSubmenu)?.submenu}
  {#if submenuActions}
    <div
      class="submenu"
      style="left: {submenuPosition.x}px; top: {submenuPosition.y}px;"
      role="menu"
    >
      {#each submenuActions as action}
        {#if action.divider}
          <div class="menu-divider"></div>
        {:else}
          <button
            class="menu-item {action.disabled ? 'disabled' : ''}"
            onclick={() => {
              if (action.onClick) {
                action.onClick();
                onClose();
              }
            }}
            disabled={action.disabled}
          >
            {#if action.icon}
              <span class="icon">{action.icon}</span>
            {/if}
            <span class="label">{action.label}</span>
          </button>
        {/if}
      {/each}
    </div>
  {/if}
{/if}

<style>
  .context-menu,
  .submenu {
    position: fixed;
    background: #1a1f2e;
    border: 2px solid #3a4451;
    border-radius: 8px;
    padding: 0.5rem 0;
    min-width: 200px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
    z-index: 2000;
    animation: slideIn 0.15s ease-out;
  }

  .submenu {
    z-index: 2001;
  }

  @keyframes slideIn {
    from {
      opacity: 0;
      transform: scale(0.95);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  .menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0.75rem;
    border-bottom: 1px solid #3a4451;
    margin-bottom: 0.25rem;
  }

  .deck-info {
    font-size: 0.875rem;
    font-weight: 600;
    color: #fff;
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .deck-count {
    font-size: 0.75rem;
    color: #9ca3af;
    font-weight: 400;
  }

  .close-btn {
    background: transparent;
    border: none;
    color: #9ca3af;
    font-size: 1.5rem;
    line-height: 1;
    cursor: pointer;
    padding: 0;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
  }

  .close-btn:hover {
    color: #fff;
  }

  .menu-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    background: transparent;
    border: none;
    color: #e5e7eb;
    font-size: 0.875rem;
    cursor: pointer;
    transition: background 0.15s;
    text-align: left;
  }

  .menu-item:hover:not(.disabled) {
    background: rgba(102, 126, 234, 0.15);
  }

  .menu-item.disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .menu-item .icon {
    font-size: 1rem;
    width: 1.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .menu-item .label {
    flex: 1;
  }

  .menu-item.has-submenu {
    position: relative;
  }

  .submenu-arrow {
    margin-left: auto;
    color: #9ca3af;
    font-size: 0.75rem;
  }

  .menu-divider {
    height: 1px;
    background: #3a4451;
    margin: 0.25rem 0;
  }
</style>
