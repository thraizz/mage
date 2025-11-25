<script lang="ts">
	/**
	 * Action Log Component
	 * Displays a scrollable log of game actions with timestamps and player colors
	 */

	import { onMount } from 'svelte';
	import type { ActionLogEntry } from '$lib/types/game';
	import ActionLogItem from './ActionLogItem.svelte';

	let {
		collapsed = $bindable(false),
		maxEntries = 500
	}: {
		collapsed?: boolean;
		maxEntries?: number;
	} = $props();

	// Action log entries
	let entries = $state<ActionLogEntry[]>([]);

	// Scroll container reference
	let scrollContainer: HTMLDivElement;
	let isUserScrolling = $state(false);
	let scrollTimeout: ReturnType<typeof setTimeout> | undefined;

	// Player color mapping (will be set from game state)
	let playerColors = $state<Map<string, string>>(
		new Map([
			['player-1', '#667eea'],
			['player-2', '#ef4444']
		])
	);

	/**
	 * Add action to log
	 */
	export function addAction(
		actionType: ActionLogEntry['actionType'],
		text: string,
		options: {
			playerName?: string;
			playerId?: string;
			cardName?: string;
			cardId?: string;
			type?: 'player' | 'system';
		} = {}
	): void {
		const entry: ActionLogEntry = {
			id: `action-${Date.now()}-${Math.random()}`,
			timestamp: Date.now(),
			type: options.type || 'player',
			playerName: options.playerName,
			playerId: options.playerId,
			actionType,
			text,
			cardName: options.cardName,
			cardId: options.cardId
		};

		entries.push(entry);

		// Limit entries to maxEntries
		if (entries.length > maxEntries) {
			entries = entries.slice(-maxEntries);
		}

		// Auto-scroll to bottom if user isn't scrolling
		if (!isUserScrolling) {
			setTimeout(() => scrollToBottom(), 50);
		}
	}

	/**
	 * Set player colors from game state
	 */
	export function setPlayerColors(colors: Map<string, string>): void {
		playerColors = colors;
	}

	/**
	 * Clear all entries
	 */
	export function clearLog(): void {
		entries = [];
	}

	/**
	 * Scroll to bottom
	 */
	function scrollToBottom(): void {
		if (scrollContainer) {
			scrollContainer.scrollTop = scrollContainer.scrollHeight;
		}
	}

	/**
	 * Handle scroll event
	 */
	function handleScroll(): void {
		if (!scrollContainer) return;

		// Check if user is at bottom
		const isAtBottom =
			scrollContainer.scrollHeight - scrollContainer.scrollTop - scrollContainer.clientHeight < 50;

		// Set user scrolling flag
		isUserScrolling = !isAtBottom;

		// Reset timeout
		if (scrollTimeout) {
			clearTimeout(scrollTimeout);
		}

		// Auto-reset user scrolling after 3 seconds
		scrollTimeout = setTimeout(() => {
			isUserScrolling = false;
		}, 3000);
	}

	/**
	 * Toggle collapsed state
	 */
	function toggleCollapsed(): void {
		collapsed = !collapsed;
	}

	/**
	 * Get player color
	 */
	function getPlayerColor(playerId?: string): string {
		if (!playerId) return '#667eea';
		return playerColors.get(playerId) || '#667eea';
	}

	onMount(() => {
		// Add some initial system messages
		addAction('system', 'Game started', { type: 'system' });
		addAction('phase', 'Beginning Phase', { type: 'system' });
	});
</script>

<div class="action-log-container" class:collapsed>
	<!-- Collapse Toggle Button -->
	<button
		class="toggle-btn"
		onclick={toggleCollapsed}
		aria-label={collapsed ? 'Expand' : 'Collapse'}
	>
		<span class="toggle-icon">{collapsed ? '▶' : '◀'}</span>
	</button>

	<!-- Log Content -->
	{#if !collapsed}
		<div class="action-log-content">
			<!-- Header -->
			<div class="log-header">
				<div class="header-title">
					<span class="header-icon">📋</span>
					<h3>Action Log</h3>
				</div>
				<div class="header-actions">
					<span class="entry-count">{entries.length}</span>
					<button
						class="clear-btn"
						onclick={clearLog}
						title="Clear log"
						disabled={entries.length === 0}
					>
						🗑️
					</button>
				</div>
			</div>

			<!-- Entries List -->
			<div class="log-entries" bind:this={scrollContainer} onscroll={handleScroll}>
				{#if entries.length === 0}
					<div class="empty-state">
						<p>No actions yet</p>
					</div>
				{:else}
					{#each entries as entry (entry.id)}
						<ActionLogItem action={entry} playerColor={getPlayerColor(entry.playerId)} />
					{/each}
				{/if}
			</div>

			<!-- Scroll to Bottom Button -->
			{#if isUserScrolling}
				<button class="scroll-to-bottom" onclick={scrollToBottom} title="Scroll to bottom">
					⬇️ New Actions
				</button>
			{/if}
		</div>
	{/if}
</div>

<style>
	.action-log-container {
		position: fixed;
		left: 0;
		top: 0;
		bottom: 0;
		width: 320px;
		background: #141821;
		border-right: 2px solid #2a3441;
		display: flex;
		flex-direction: column;
		z-index: 20;
		transition: width 0.3s ease;
	}

	.action-log-container.collapsed {
		width: 48px;
	}

	/* Toggle Button */
	.toggle-btn {
		position: absolute;
		right: -16px;
		top: 50%;
		transform: translateY(-50%);
		width: 32px;
		height: 48px;
		background: #1a1f2e;
		border: 2px solid #2a3441;
		border-radius: 0 8px 8px 0;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s;
		z-index: 1;
	}

	.toggle-btn:hover {
		background: #2a3441;
		border-color: #3a4451;
	}

	.toggle-icon {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	/* Content */
	.action-log-content {
		display: flex;
		flex-direction: column;
		height: 100%;
		overflow: hidden;
	}

	/* Header */
	.log-header {
		padding: 1rem;
		background: #1a1f2e;
		border-bottom: 2px solid #2a3441;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
	}

	.header-title {
		display: flex;
		align-items: center;
		gap: 0.625rem;
	}

	.header-icon {
		font-size: 1.25rem;
	}

	.log-header h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: white;
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.entry-count {
		font-size: 0.875rem;
		color: #9ca3af;
		font-weight: 600;
		padding: 0.25rem 0.5rem;
		background: #0f1419;
		border-radius: 4px;
	}

	.clear-btn {
		background: transparent;
		border: none;
		font-size: 1rem;
		cursor: pointer;
		padding: 0.25rem;
		opacity: 0.6;
		transition: opacity 0.2s;
	}

	.clear-btn:hover:not(:disabled) {
		opacity: 1;
	}

	.clear-btn:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	/* Entries List */
	.log-entries {
		flex: 1;
		overflow-y: auto;
		overflow-x: hidden;
		background: #0f1419;
	}

	/* Scrollbar Styling */
	.log-entries::-webkit-scrollbar {
		width: 8px;
	}

	.log-entries::-webkit-scrollbar-track {
		background: #0f1419;
	}

	.log-entries::-webkit-scrollbar-thumb {
		background: #2a3441;
		border-radius: 4px;
	}

	.log-entries::-webkit-scrollbar-thumb:hover {
		background: #3a4451;
	}

	/* Empty State */
	.empty-state {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		padding: 2rem;
	}

	.empty-state p {
		color: #6b7280;
		font-size: 0.875rem;
		font-style: italic;
		text-align: center;
	}

	/* Scroll to Bottom Button */
	.scroll-to-bottom {
		position: absolute;
		bottom: 1rem;
		left: 50%;
		transform: translateX(-50%);
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 20px;
		font-size: 0.875rem;
		font-weight: 600;
		cursor: pointer;
		box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
		transition: all 0.2s;
		z-index: 10;
	}

	.scroll-to-bottom:hover {
		background: #5568d3;
		transform: translateX(-50%) translateY(-2px);
		box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
	}

	/* Responsive */
	@media (max-width: 1400px) {
		.action-log-container:not(.collapsed) {
			width: 280px;
		}
	}

	@media (max-width: 1024px) {
		.action-log-container {
			width: 280px;
		}

		.action-log-container.collapsed {
			width: 0;
			border-right: none;
		}

		.toggle-btn {
			right: 0;
			border-radius: 0 8px 8px 0;
		}

		.action-log-container.collapsed .toggle-btn {
			right: -32px;
		}
	}

	@media (max-width: 768px) {
		.action-log-container:not(.collapsed) {
			width: 100%;
			max-width: 280px;
		}
	}
</style>
