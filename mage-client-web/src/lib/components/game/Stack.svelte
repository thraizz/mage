<script lang="ts">
	import type { StackObject } from '$lib/types/game';

	// Props
	let {
		stackObjects = [],
		playerNames = new Map<string, string>(),
		// eslint-disable-next-line no-unused-vars
		onStackObjectClick = (stackId: string) => {}
	}: {
		stackObjects?: StackObject[];
		playerNames?: Map<string, string>;
		// eslint-disable-next-line no-unused-vars
		onStackObjectClick?: (stackId: string) => void;
	} = $props();

	/**
	 * Handle stack object click
	 */
	function handleClick(_stackId: string): void {
		onStackObjectClick(_stackId);
	}

	/**
	 * Get player name by ID
	 */
	function getPlayerName(playerId: string): string {
		return playerNames.get(playerId) || 'Unknown Player';
	}

	/**
	 * Get stack item type label
	 */
	function getTypeLabel(type: 'SPELL' | 'ABILITY'): string {
		return type === 'SPELL' ? 'Spell' : 'Ability';
	}

	/**
	 * Get stack item type color
	 */
	function getTypeColor(type: 'SPELL' | 'ABILITY'): string {
		return type === 'SPELL' ? '#667eea' : '#fbbf24';
	}

	// Derived values
	const isEmpty = $derived(stackObjects.length === 0);
	const stackCount = $derived(stackObjects.length);
</script>

<div class="stack-component">
	<div class="stack-header">
		<span class="stack-label">Stack</span>
		{#if !isEmpty}
			<span class="stack-count">({stackCount})</span>
		{/if}
	</div>

	{#if isEmpty}
		<div class="empty-state">
			<div class="empty-icon">📚</div>
			<p>No spells or abilities on the stack</p>
		</div>
	{:else}
		<div class="stack-items">
			{#each stackObjects as stackObj, index (stackObj.id)}
				{@const isTop = index === stackObjects.length - 1}
				<div
					class="stack-item"
					class:top={isTop}
					role="button"
					tabindex="0"
					onclick={() => handleClick(stackObj.id)}
					onkeydown={(e) => e.key === 'Enter' && handleClick(stackObj.id)}
				>
					<div class="stack-item-header">
						<span class="stack-type" style="background: {getTypeColor(stackObj.type)}">
							{getTypeLabel(stackObj.type)}
						</span>
						<span class="stack-position">
							{stackObjects.length - index}
						</span>
					</div>

					<div class="stack-item-content">
						<div class="stack-name">{stackObj.name}</div>
						<div class="stack-controller">
							{getPlayerName(stackObj.controllerId)}
						</div>
					</div>

					{#if stackObj.targets && stackObj.targets.length > 0}
						<div class="stack-targets">
							<span class="targets-label">Targets:</span>
							<span class="targets-count">{stackObj.targets.length}</span>
						</div>
					{/if}

					{#if isTop}
						<div class="resolving-indicator">
							<span>Will resolve next</span>
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<div class="stack-footer">
			<span class="resolution-order">Bottom ↓ resolves first ↓ Top</span>
		</div>
	{/if}
</div>

<style>
	.stack-component {
		background: #1a1f2e;
		border: 2px solid #667eea;
		border-radius: 8px;
		display: flex;
		flex-direction: column;
		max-height: 500px;
		box-shadow: 0 4px 8px rgba(102, 126, 234, 0.2);
	}

	.stack-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 6px 6px 0 0;
	}

	.stack-label {
		font-size: 0.875rem;
		color: white;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.stack-count {
		font-size: 0.875rem;
		color: rgba(255, 255, 255, 0.8);
	}

	/* Empty State */
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 2rem;
		gap: 1rem;
	}

	.empty-icon {
		font-size: 3rem;
		opacity: 0.3;
	}

	.empty-state p {
		color: #6b7280;
		font-size: 0.875rem;
		font-style: italic;
		margin: 0;
		text-align: center;
	}

	/* Stack Items */
	.stack-items {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		display: flex;
		flex-direction: column-reverse; /* Bottom to top */
		gap: 0.75rem;
	}

	.stack-item {
		background: #141821;
		border: 2px solid #2a3441;
		border-radius: 8px;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		cursor: pointer;
		transition:
			transform 0.2s,
			border-color 0.2s,
			box-shadow 0.2s;
	}

	.stack-item:hover {
		transform: translateX(4px);
		border-color: #667eea;
		box-shadow: 0 4px 8px rgba(102, 126, 234, 0.3);
	}

	.stack-item:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.3);
	}

	.stack-item.top {
		border-color: #fbbf24;
		box-shadow: 0 4px 12px rgba(251, 191, 36, 0.4);
	}

	.stack-item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.stack-type {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 600;
		color: white;
	}

	.stack-position {
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #2a3441;
		border-radius: 50%;
		font-size: 0.75rem;
		font-weight: 700;
		color: #9ca3af;
	}

	.stack-item.top .stack-position {
		background: #fbbf24;
		color: #000;
	}

	.stack-item-content {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.stack-name {
		font-size: 1rem;
		font-weight: 600;
		color: #ffffff;
	}

	.stack-controller {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.stack-targets {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem;
		background: #0d1117;
		border-radius: 4px;
		font-size: 0.875rem;
	}

	.targets-label {
		color: #9ca3af;
	}

	.targets-count {
		padding: 0.125rem 0.375rem;
		background: #667eea;
		border-radius: 4px;
		font-weight: 600;
		color: white;
		font-size: 0.75rem;
	}

	.resolving-indicator {
		padding: 0.5rem;
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
		border-radius: 4px;
		text-align: center;
	}

	.resolving-indicator span {
		font-size: 0.75rem;
		font-weight: 600;
		color: #000;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	/* Stack Footer */
	.stack-footer {
		padding: 0.75rem 1rem;
		background: #141821;
		border-top: 1px solid #2a3441;
		border-radius: 0 0 6px 6px;
		text-align: center;
	}

	.resolution-order {
		font-size: 0.75rem;
		color: #6b7280;
		font-style: italic;
	}

	/* Scrollbar */
	.stack-items::-webkit-scrollbar {
		width: 6px;
	}

	.stack-items::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.stack-items::-webkit-scrollbar-thumb {
		background: #667eea;
		border-radius: 3px;
	}

	.stack-items::-webkit-scrollbar-thumb:hover {
		background: #764ba2;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.stack-component {
			max-height: 400px;
		}

		.stack-items {
			padding: 0.75rem;
		}

		.stack-item {
			padding: 0.75rem;
		}

		.stack-name {
			font-size: 0.875rem;
		}

		.stack-controller {
			font-size: 0.75rem;
		}
	}
</style>
