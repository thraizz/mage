<script lang="ts">
  import type { CounterView } from '$lib/generated/mage/v1/models';

  interface Props {
    cardName: string;
    cardId: string;
    currentCounters: CounterView[];
    onAddCounter: (counterName: string, amount: number) => void;
    onRemoveCounter: (counterName: string, amount: number) => void;
    onSetCounter: (counterName: string, amount: number) => void;
    onClose: () => void;
  }

  let {
    cardName,
    cardId,
    currentCounters,
    onAddCounter,
    onRemoveCounter,
    onSetCounter,
    onClose
  }: Props = $props();

  // State for new counter
  let newCounterName = $state('');
  let newCounterAmount = $state(1);

  // Common counter types for quick selection
  const commonCounters = [
    '+1/+1',
    '-1/-1',
    'loyalty',
    'charge',
    'time',
    'poison',
    'energy',
    'age',
    'quest',
    'fade'
  ];

  function handleAddCounter() {
    if (!newCounterName.trim()) return;
    onAddCounter(newCounterName.trim(), newCounterAmount);
    newCounterName = '';
    newCounterAmount = 1;
  }

  function handleQuickAdd(counterName: string) {
    onAddCounter(counterName, 1);
  }

  function handleRemove(counterName: string) {
    onRemoveCounter(counterName, 1);
  }

  function handleSet(counterName: string, amount: number) {
    onSetCounter(counterName, amount);
  }
</script>

<div class="counter-overlay" role="dialog" aria-labelledby="counter-dialog-title">
  <div class="counter-dialog">
    <div class="dialog-header">
      <h2 id="counter-dialog-title">Counters - {cardName}</h2>
      <button class="close-button" onclick={onClose} aria-label="Close dialog">×</button>
    </div>

    <!-- Current Counters -->
    <div class="current-counters">
      <h3>Current Counters</h3>
      {#if currentCounters.length === 0}
        <p class="no-counters">No counters on this card</p>
      {:else}
        <div class="counter-list">
          {#each currentCounters as counter (counter.name)}
            <div class="counter-item">
              <div class="counter-info">
                <span class="counter-name">{counter.name}</span>
                <span class="counter-count">×{counter.count}</span>
              </div>
              <div class="counter-actions">
                <button class="btn-small btn-secondary" onclick={() => handleRemove(counter.name)}>
                  −1
                </button>
                <button
                  class="btn-small btn-secondary"
                  onclick={() => onAddCounter(counter.name, 1)}
                >
                  +1
                </button>
                <button class="btn-small btn-danger" onclick={() => handleSet(counter.name, 0)}>
                  Remove
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Quick Add Common Counters -->
    <div class="quick-add">
      <h3>Quick Add</h3>
      <div class="quick-buttons">
        {#each commonCounters as counterName}
          <button class="btn-quick" onclick={() => handleQuickAdd(counterName)}>
            {counterName}
          </button>
        {/each}
      </div>
    </div>

    <!-- Custom Counter -->
    <div class="custom-counter">
      <h3>Add Custom Counter</h3>
      <div class="custom-form">
        <input
          type="text"
          bind:value={newCounterName}
          placeholder="Counter name"
          class="counter-input"
          onkeydown={(e) => e.key === 'Enter' && handleAddCounter()}
        />
        <input type="number" bind:value={newCounterAmount} min="1" class="amount-input" />
        <button class="btn-primary" onclick={handleAddCounter}>Add</button>
      </div>
    </div>

    <div class="dialog-footer">
      <button class="btn-secondary" onclick={onClose}>Close</button>
    </div>
  </div>
</div>

<style>
  .counter-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
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

  .counter-dialog {
    background: #1a1f2e;
    border: 2px solid #3a4451;
    border-radius: 12px;
    padding: 1.5rem;
    max-width: 500px;
    width: 90%;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
    animation: slideUp 0.3s ease-out;
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
    color: #fff;
    font-size: 1.25rem;
    font-weight: 600;
  }

  .close-button {
    background: transparent;
    border: none;
    color: #9ca3af;
    font-size: 2rem;
    line-height: 1;
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

  .current-counters,
  .quick-add,
  .custom-counter {
    margin-bottom: 1.5rem;
  }

  .current-counters h3,
  .quick-add h3,
  .custom-counter h3 {
    font-size: 0.875rem;
    font-weight: 600;
    color: #9ca3af;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.75rem;
  }

  .no-counters {
    color: #6b7280;
    font-style: italic;
    padding: 1rem;
    text-align: center;
    background: rgba(255, 255, 255, 0.02);
    border-radius: 6px;
  }

  .counter-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .counter-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid #3a4451;
    border-radius: 6px;
    transition: background 0.2s;
  }

  .counter-item:hover {
    background: rgba(255, 255, 255, 0.08);
  }

  .counter-info {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .counter-name {
    color: #fff;
    font-weight: 600;
    font-size: 0.875rem;
  }

  .counter-count {
    color: #10b981;
    font-weight: 700;
    font-size: 1rem;
  }

  .counter-actions {
    display: flex;
    gap: 0.5rem;
  }

  .quick-buttons {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 0.5rem;
  }

  .btn-quick {
    padding: 0.5rem;
    background: rgba(102, 126, 234, 0.1);
    border: 1px solid #667eea;
    border-radius: 4px;
    color: #667eea;
    font-size: 0.75rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-quick:hover {
    background: rgba(102, 126, 234, 0.2);
    border-color: #7c93ee;
    transform: translateY(-1px);
  }

  .custom-form {
    display: flex;
    gap: 0.5rem;
  }

  .counter-input {
    flex: 1;
    padding: 0.5rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid #3a4451;
    border-radius: 6px;
    color: #fff;
    font-size: 0.875rem;
    transition: border-color 0.2s;
  }

  .counter-input:focus {
    outline: none;
    border-color: #667eea;
  }

  .amount-input {
    width: 70px;
    padding: 0.5rem;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid #3a4451;
    border-radius: 6px;
    color: #fff;
    font-size: 0.875rem;
    transition: border-color 0.2s;
  }

  .amount-input:focus {
    outline: none;
    border-color: #667eea;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    padding-top: 1rem;
    border-top: 1px solid #3a4451;
  }

  .btn-primary,
  .btn-secondary,
  .btn-small,
  .btn-danger {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 600;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover {
    background: #5568d3;
    transform: translateY(-1px);
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
    border: 1px solid #3a4451;
  }

  .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.15);
  }

  .btn-small {
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
  }

  .btn-danger {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid #ef4444;
  }

  .btn-danger:hover {
    background: rgba(239, 68, 68, 0.2);
  }
</style>
