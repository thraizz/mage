<script lang="ts">
	interface Option {
		value: string;
		label: string;
		disabled?: boolean;
	}

	interface Props {
		value?: string;
		options: Option[];
		placeholder?: string;
		label?: string;
		error?: string;
		disabled?: boolean;
		required?: boolean;
		id?: string;
		name?: string;
		onchange?: (e: Event) => void;
	}

	let {
		value = $bindable(''),
		options,
		placeholder = 'Select an option',
		label = '',
		error = '',
		disabled = false,
		required = false,
		id = crypto.randomUUID(),
		name = '',
		onchange
	}: Props = $props();
</script>

<div class="select-group" class:has-error={error}>
	{#if label}
		<label for={id} class="select-label">
			{label}
			{#if required}<span class="required-mark">*</span>{/if}
		</label>
	{/if}

	<div class="select-wrapper">
		<select {id} {name} bind:value {disabled} {required} class="select" {onchange}>
			{#if placeholder}
				<option value="" disabled>{placeholder}</option>
			{/if}
			{#each options as option}
				<option value={option.value} disabled={option.disabled}>
					{option.label}
				</option>
			{/each}
		</select>
		<span class="select-arrow">
			<svg
				width="16"
				height="16"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
			>
				<polyline points="6 9 12 15 18 9"></polyline>
			</svg>
		</span>
	</div>

	{#if error}
		<span class="select-error">{error}</span>
	{/if}
</div>

<style>
	.select-group {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.select-label {
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		color: var(--text-muted);
	}

	.required-mark {
		color: var(--status-error);
		margin-left: var(--space-1);
	}

	.select-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.select {
		width: 100%;
		padding: var(--space-2) var(--space-8) var(--space-2) var(--space-3);
		font-family: var(--font-body);
		font-size: var(--text-base);
		color: var(--text-bright);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		outline: none;
		cursor: pointer;
		appearance: none;
		transition: all var(--transition-fast);
	}

	.select:focus {
		border-color: var(--accent-gold);
		box-shadow: 0 0 0 3px var(--accent-gold-glow);
	}

	.select:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		background: var(--bg-slate);
	}

	.select option {
		background: var(--bg-slate);
		color: var(--text-bright);
		padding: var(--space-2);
	}

	.select option:disabled {
		color: var(--text-ghost);
	}

	.select-arrow {
		position: absolute;
		right: var(--space-3);
		pointer-events: none;
		color: var(--text-dim);
		display: flex;
		align-items: center;
	}

	.has-error .select {
		border-color: var(--status-error);
	}

	.has-error .select:focus {
		box-shadow: 0 0 0 3px var(--status-error-dim);
	}

	.select-error {
		font-size: var(--text-sm);
		color: var(--status-error);
	}
</style>
