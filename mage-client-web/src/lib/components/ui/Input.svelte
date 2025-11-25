<script lang="ts">
	interface Props {
		type?: 'text' | 'password' | 'email' | 'search' | 'number' | 'url';
		value?: string;
		placeholder?: string;
		label?: string;
		error?: string;
		hint?: string;
		disabled?: boolean;
		required?: boolean;
		readonly?: boolean;
		id?: string;
		name?: string;
		autocomplete?: AutoFill;
		oninput?: (e: Event) => void;
		onchange?: (e: Event) => void;
		onfocus?: (e: FocusEvent) => void;
		onblur?: (e: FocusEvent) => void;
	}

	let {
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		label = '',
		error = '',
		hint = '',
		disabled = false,
		required = false,
		readonly = false,
		id = crypto.randomUUID(),
		name = '',
		autocomplete,
		oninput,
		onchange,
		onfocus,
		onblur
	}: Props = $props();
</script>

<div class="input-group" class:has-error={error}>
	{#if label}
		<label for={id} class="input-label">
			{label}
			{#if required}<span class="required-mark">*</span>{/if}
		</label>
	{/if}

	<input
		{id}
		{type}
		{name}
		bind:value
		{placeholder}
		{disabled}
		{required}
		{readonly}
		{autocomplete}
		class="input"
		{oninput}
		{onchange}
		{onfocus}
		{onblur}
	/>

	{#if error}
		<span class="input-error">{error}</span>
	{:else if hint}
		<span class="input-hint">{hint}</span>
	{/if}
</div>

<style>
	.input-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.input-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	.required-mark {
		color: var(--status-error);
		margin-left: var(--space-1);
	}

	.input {
		width: 100%;
		padding: var(--space-2) var(--space-3);
		font-family: var(--font-body);
		font-size: var(--text-base);
		color: var(--text-bright);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		outline: none;
		transition: all var(--transition-fast);
	}

	.input::placeholder {
		color: var(--text-ghost);
	}

	.input:focus {
		border-color: var(--accent-gold);
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	.input:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		background: var(--bg-slate);
	}

	.input:read-only {
		background: var(--bg-slate);
		cursor: default;
	}

	.has-error .input {
		border-color: var(--status-error);
	}

	.has-error .input:focus {
		box-shadow: 0 0 0 3px var(--status-error-dim);
	}

	.input-error {
		font-size: var(--text-sm);
		color: var(--status-error);
	}

	.input-hint {
		font-size: var(--text-sm);
		color: var(--text-dim);
	}

	/* Search input specific styling */
	.input[type='search'] {
		padding-left: var(--space-4);
	}

	.input[type='search']::-webkit-search-cancel-button {
		-webkit-appearance: none;
		appearance: none;
		height: 1em;
		width: 1em;
		background: var(--text-dim);
		mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2'%3E%3Cline x1='18' y1='6' x2='6' y2='18'%3E%3C/line%3E%3Cline x1='6' y1='6' x2='18' y2='18'%3E%3C/line%3E%3C/svg%3E");
		cursor: pointer;
	}

	/* Number input - hide spinner buttons */
	.input[type='number']::-webkit-inner-spin-button,
	.input[type='number']::-webkit-outer-spin-button {
		-webkit-appearance: none;
		appearance: none;
		margin: 0;
	}

	.input[type='number'] {
		-moz-appearance: textfield;
	}
</style>
