<script lang="ts">
  interface Tab {
    id: string;
    label: string;
    disabled?: boolean;
  }

  interface Props {
    tabs: Tab[];
    activeTab?: string;
    variant?: 'default' | 'pills';
    onchange?: (tabId: string) => void;
  }

  let {
    tabs,
    activeTab = $bindable(tabs[0]?.id || ''),
    variant = 'default',
    onchange
  }: Props = $props();

  function selectTab(tabId: string) {
    activeTab = tabId;
    onchange?.(tabId);
  }
</script>

<div class="tabs tabs-{variant}" role="tablist">
  {#each tabs as tab}
    <button
      type="button"
      role="tab"
      class="tab"
      class:active={activeTab === tab.id}
      disabled={tab.disabled}
      aria-selected={activeTab === tab.id}
      onclick={() => selectTab(tab.id)}
    >
      {tab.label}
    </button>
  {/each}
</div>

<style>
  .tabs {
    display: flex;
    gap: var(--space-1);
  }

  .tab {
    font-family: var(--font-body);
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-muted);
    background: transparent;
    border: none;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .tab:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .tab:focus-visible {
    outline: 2px solid var(--accent-gold);
    outline-offset: 2px;
  }

  /* Default variant - underline style */
  .tabs-default {
    border-bottom: 1px solid var(--border-subtle);
  }

  .tabs-default .tab {
    padding: var(--space-3) var(--space-4);
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
  }

  .tabs-default .tab:hover:not(:disabled) {
    color: var(--text-bright);
  }

  .tabs-default .tab.active {
    color: var(--accent-gold);
    border-bottom-color: var(--accent-gold);
  }

  /* Pills variant */
  .tabs-pills {
    background: var(--bg-obsidian);
    border-radius: var(--radius-md);
    padding: var(--space-1);
  }

  .tabs-pills .tab {
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
  }

  .tabs-pills .tab:hover:not(:disabled):not(.active) {
    background: var(--bg-iron);
    color: var(--text-bright);
  }

  .tabs-pills .tab.active {
    background: var(--accent-gold);
    color: var(--bg-void);
  }
</style>
