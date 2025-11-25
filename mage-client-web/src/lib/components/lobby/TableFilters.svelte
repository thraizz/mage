<script lang="ts">
	import Input from '$lib/components/ui/Input.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Checkbox from '$lib/components/ui/Checkbox.svelte';
	import Button from '$lib/components/ui/Button.svelte';

	interface Props {
		search: string;
		format: string;
		openOnly: boolean;
		formats: string[];
		onchange?: () => void;
	}

	let {
		search = $bindable(''),
		format = $bindable('All'),
		openOnly = $bindable(false),
		formats,
		onchange
	}: Props = $props();

	const formatOptions = $derived([
		{ value: 'All', label: 'All Formats' },
		...formats.map((f) => ({ value: f, label: f }))
	]);

	const hasFilters = $derived(search.trim() !== '' || format !== 'All' || openOnly);

	function clearFilters() {
		search = '';
		format = 'All';
		openOnly = false;
		onchange?.();
	}

	function handleSearchInput() {
		onchange?.();
	}

	function handleFormatChange() {
		onchange?.();
	}

	function handleOpenOnlyChange() {
		onchange?.();
	}
</script>

<div class="filters">
	<div class="filter-search">
		<Input
			type="search"
			bind:value={search}
			placeholder="Search tables..."
			oninput={handleSearchInput}
		/>
	</div>

	<div class="filter-format">
		<Select bind:value={format} options={formatOptions} onchange={handleFormatChange} />
	</div>

	<div class="filter-toggle">
		<Checkbox bind:checked={openOnly} label="Open only" onchange={handleOpenOnlyChange} />
	</div>

	{#if hasFilters}
		<Button variant="ghost" size="sm" onclick={clearFilters}>Clear</Button>
	{/if}
</div>

<style>
	.filters {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		flex-wrap: wrap;
	}

	.filter-search {
		flex: 1;
		min-width: 200px;
		max-width: 320px;
	}

	.filter-format {
		min-width: 160px;
	}

	.filter-toggle {
		padding: var(--space-2) var(--space-3);
		background: var(--bg-iron);
		border-radius: var(--radius-md);
	}

	@media (max-width: 768px) {
		.filters {
			flex-direction: column;
			align-items: stretch;
		}

		.filter-search {
			max-width: none;
		}

		.filter-format {
			width: 100%;
		}
	}
</style>
