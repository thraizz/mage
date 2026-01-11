<script lang="ts">
	import Check from '@lucide/svelte/icons/check';

	interface Props {
		checked?: boolean;
		label?: string;
		disabled?: boolean;
		id?: string;
		name?: string;
		onchange?: (e: Event) => void;
	}

	let {
		checked = $bindable(false),
		label = '',
		disabled = false,
		id = crypto.randomUUID(),
		name = '',
		onchange
	}: Props = $props();
</script>

<label class="checkbox-wrapper" class:disabled>
	<input
		type="checkbox"
		{id}
		{name}
		bind:checked
		{disabled}
		class="checkbox-input"
		{onchange}
	/>
	<span class="checkbox-box">
		<Check class="checkbox-icon" size={12} aria-hidden="true" />
	</span>
	{#if label}
		<span class="checkbox-label">{label}</span>
	{/if}
</label>

<style>
	.checkbox-wrapper {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		cursor: pointer;
		user-select: none;
	}

	.checkbox-wrapper.disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.checkbox-input {
		position: absolute;
		opacity: 0;
		width: 0;
		height: 0;
	}

	.checkbox-box {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 1.25rem;
		height: 1.25rem;
		background: var(--bg-iron);
		border: 2px solid var(--border-default);
		border-radius: var(--radius-sm);
		transition: all var(--transition-fast);
	}

	:global(svg.checkbox-icon) {
		width: 0.75rem;
		height: 0.75rem;
		opacity: 0;
		transform: scale(0.5);
		transition: all var(--transition-fast);
		color: var(--bg-void);
	}

	.checkbox-input:checked + .checkbox-box {
		background: var(--accent-gold);
		border-color: var(--accent-gold);
	}

	.checkbox-input:checked + .checkbox-box :global(svg.checkbox-icon) {
		opacity: 1;
		transform: scale(1);
	}

	.checkbox-input:focus-visible + .checkbox-box {
		outline: 2px solid var(--accent-gold);
		outline-offset: 2px;
	}

	.checkbox-wrapper:hover:not(.disabled) .checkbox-box {
		border-color: var(--accent-gold-dim);
	}

	.checkbox-label {
		font-size: var(--text-base);
		color: var(--text-bright);
	}
</style>
