<script lang="ts">
  interface Props {
    value?: string;
    placeholder?: string;
    label?: string;
    error?: string;
    hint?: string;
    disabled?: boolean;
    required?: boolean;
    readonly?: boolean;
    rows?: number;
    id?: string;
    name?: string;
    resize?: 'none' | 'vertical' | 'horizontal' | 'both';
    oninput?: (e: Event) => void;
    onchange?: (e: Event) => void;
  }

  let {
    value = $bindable(''),
    placeholder = '',
    label = '',
    error = '',
    hint = '',
    disabled = false,
    required = false,
    readonly = false,
    rows = 4,
    id = crypto.randomUUID(),
    name = '',
    resize = 'vertical',
    oninput,
    onchange
  }: Props = $props();
</script>

<div class="textarea-group" class:has-error={error}>
  {#if label}
    <label for={id} class="textarea-label">
      {label}
      {#if required}<span class="required-mark">*</span>{/if}
    </label>
  {/if}

  <textarea
    {id}
    {name}
    bind:value
    {placeholder}
    {disabled}
    {required}
    {readonly}
    {rows}
    class="textarea"
    style="resize: {resize}"
    {oninput}
    {onchange}
  ></textarea>

  {#if error}
    <span class="textarea-error">{error}</span>
  {:else if hint}
    <span class="textarea-hint">{hint}</span>
  {/if}
</div>

<style>
  .textarea-group {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .textarea-label {
    font-size: var(--text-sm);
    font-weight: var(--weight-medium);
    color: var(--text-muted);
  }

  .required-mark {
    color: var(--status-error);
    margin-left: var(--space-1);
  }

  .textarea {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    font-family: var(--font-body);
    font-size: var(--text-base);
    line-height: var(--leading-normal);
    color: var(--text-bright);
    background: var(--bg-iron);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-md);
    outline: none;
    transition:
      border-color var(--transition-fast),
      box-shadow var(--transition-fast);
  }

  .textarea::placeholder {
    color: var(--text-ghost);
  }

  .textarea:focus {
    border-color: var(--accent-gold);
    box-shadow: 0 0 0 3px var(--accent-gold-glow);
  }

  .textarea:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background: var(--bg-slate);
  }

  .textarea:read-only {
    background: var(--bg-slate);
    cursor: default;
  }

  .has-error .textarea {
    border-color: var(--status-error);
  }

  .has-error .textarea:focus {
    box-shadow: 0 0 0 3px var(--status-error-dim);
  }

  .textarea-error {
    font-size: var(--text-sm);
    color: var(--status-error);
  }

  .textarea-hint {
    font-size: var(--text-sm);
    color: var(--text-dim);
  }
</style>
