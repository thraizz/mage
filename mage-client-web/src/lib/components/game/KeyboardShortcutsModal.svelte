<script lang="ts">
	import Modal from '$lib/components/Modal.svelte';

	type KeyboardShortcutsMode = 'game' | 'playtest';

	export let open = false;
	export let mode: KeyboardShortcutsMode = 'game';

	function handleClose(): void {
		open = false;
	}

	type Shortcut = { keys: string; action: string; note?: string };

	const gameGlobal: Shortcut[] = [
		{ keys: '?', action: 'Show/hide keyboard shortcuts' },
		{ keys: 'X', action: 'Untap all permanents' },
		{ keys: 'C', action: 'Draw a card' },
		{ keys: 'V', action: 'Shuffle your library' },
		{ keys: 'M', action: 'Mulligan (mulligan phase only)' },
		{ keys: 'B', action: 'Open game chat' },
		{ keys: 'N', action: 'Advance to next phase/step' },
		{ keys: 'E', action: 'End turn (pass until next turn)' },
		{ keys: 'W', action: 'Open token creator' },
		{ keys: 'F', action: 'Search your library (when available)' }
	];

	const gamePriority: Shortcut[] = [
		{ keys: 'Space', action: 'Pass priority' },
		{ keys: 'F6', action: 'Pass until next turn' },
		{ keys: 'C', action: 'Cast spell / play land (selected card)' },
		{ keys: 'N', action: 'Advance to next phase/step' },
		{ keys: 'T', action: 'Open token creator' }
	];

	const gameHovered: Shortcut[] = [
		{ keys: 'J', action: 'Flip face down/up' },
		{ keys: 'L', action: 'Transform (alt/default face)' },
		{ keys: 'D', action: 'Move to graveyard' },
		{ keys: 'S', action: 'Move to exile' },
		{ keys: 'R', action: 'Move to hand' },
		{ keys: 'T', action: 'Put on top of library' },
		{ keys: '.', action: 'Put on bottom of library' },
		{ keys: 'U', action: 'Add +1/+1 counter' }
	];

	const playtestGlobal: Shortcut[] = [
		{ keys: '?', action: 'Show/hide keyboard shortcuts' },
		{ keys: 'X', action: 'Untap all permanents' },
		{ keys: 'C', action: 'Draw a card' },
		{ keys: 'V', action: 'Shuffle your library' },
		{ keys: 'E', action: 'Next turn' },
		{ keys: 'W', action: 'Open token creator' },
		{ keys: 'F', action: 'Search your deck' }
	];

	const playtestHovered: Shortcut[] = [
		{ keys: 'D', action: 'Move to graveyard' },
		{ keys: 'S', action: 'Move to exile' },
		{ keys: 'R', action: 'Move to hand' },
		{ keys: 'T', action: 'Put on top of library' }
	];
</script>

<Modal bind:open title="Keyboard shortcuts" size="medium" onClose={handleClose}>
	<div class="shortcuts">
		<p class="hint">Shortcuts are ignored while you’re typing in an input/textarea.</p>

		{#if mode === 'game'}
			<section class="section">
				<h3>Global</h3>
				<div class="grid">
					{#each gameGlobal as s}
						<div class="row">
							<kbd>{s.keys}</kbd>
							<span class="action">{s.action}</span>
						</div>
					{/each}
				</div>
			</section>

			<section class="section">
				<h3>Priority bar</h3>
				<div class="grid">
					{#each gamePriority as s}
						<div class="row">
							<kbd>{s.keys}</kbd>
							<span class="action">{s.action}</span>
						</div>
					{/each}
				</div>
			</section>

			<section class="section">
				<h3>Hovered card</h3>
				<div class="grid">
					{#each gameHovered as s}
						<div class="row">
							<kbd>{s.keys}</kbd>
							<span class="action">{s.action}</span>
						</div>
					{/each}
				</div>
			</section>
		{:else}
			<section class="section">
				<h3>Global</h3>
				<div class="grid">
					{#each playtestGlobal as s}
						<div class="row">
							<kbd>{s.keys}</kbd>
							<span class="action">{s.action}</span>
						</div>
					{/each}
				</div>
			</section>

			<section class="section">
				<h3>Hovered card</h3>
				<div class="grid">
					{#each playtestHovered as s}
						<div class="row">
							<kbd>{s.keys}</kbd>
							<span class="action">{s.action}</span>
						</div>
					{/each}
				</div>
			</section>
		{/if}
	</div>
</Modal>

<style>
	.shortcuts {
		display: flex;
		flex-direction: column;
		gap: var(--space-6);
	}

	.hint {
		margin: 0;
		color: var(--text-muted);
		font-size: var(--text-sm);
	}

	.section h3 {
		margin: 0 0 var(--space-3) 0;
		color: var(--text-bright);
		font-weight: var(--weight-semibold);
		font-size: var(--text-base);
	}

	.grid {
		display: grid;
		grid-template-columns: 1fr;
		gap: var(--space-2);
	}

	.row {
		display: grid;
		grid-template-columns: 90px 1fr;
		gap: var(--space-3);
		align-items: center;
		padding: var(--space-2) var(--space-3);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		background: var(--bg-iron);
	}

	kbd {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 2px 8px;
		border-radius: var(--radius-md);
		border: 1px solid var(--border-subtle);
		background: var(--bg-obsidian);
		color: var(--text-bright);
		font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
			monospace;
		font-size: var(--text-xs);
		line-height: 1.4;
		white-space: nowrap;
	}

	.action {
		color: var(--text-muted);
		font-size: var(--text-sm);
	}

	@media (max-width: 520px) {
		.row {
			grid-template-columns: 76px 1fr;
		}
	}
</style>
