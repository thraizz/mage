<script lang="ts">
	/**
	 * Action Log Overlay Component
	 * Displays a scrollable log of game actions as a slide-out overlay
	 */

	import { onMount } from 'svelte';
	import type { ActionLogEntry } from '$lib/types/game';
	import ActionLogItem from './ActionLogItem.svelte';

	let {
		open = $bindable(false),
		maxEntries = 500
	}: {
		open?: boolean;
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
	 * Close the overlay
	 */
	function close(): void {
		open = false;
	}

	/**
	 * Get player color
	 */
	function getPlayerColor(playerId?: string): string {
		if (!playerId) return '#667eea';
		return playerColors.get(playerId) || '#667eea';
	}

	/**
	 * Handle click outside to close
	 */
	function handleBackdropClick(e: MouseEvent): void {
		if (e.target === e.currentTarget) {
			close();
		}
	}

	/**
	 * Handle escape key to close
	 */
	function handleKeydown(e: KeyboardEvent): void {
		if (e.key === 'Escape' && open) {
			close();
		}
	}

	onMount(() => {
		// Add some initial system messages
		addAction('system', 'Game started', { type: 'system' });
		addAction('phase', 'Beginning Phase', { type: 'system' });
	});
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Backdrop -->
{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="overlay-backdrop" onclick={handleBackdropClick}></div>
{/if}

<!-- Slide-out Panel -->
<div class="action-log-overlay" class:open>
	<!-- Header -->
	<div class="log-header">
		<div class="header-title">
			<span class="header-icon">📋</span>
			<h3>Action Log</h3>
			<span class="entry-count">{entries.length}</span>
		</div>
		<div class="header-actions">
			<button
				class="clear-btn"
				onclick={clearLog}
				title="Clear log"
				disabled={entries.length === 0}
			>
				🗑️
			</button>
			<button class="close-btn" onclick={close} title="Close (Esc)">
				✕
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
	{#if isUserScrolling && open}
		<button class="scroll-to-bottom" onclick={scrollToBottom} title="Scroll to bottom">
			⬇️ New Actions
		</button>
	{/if}
</div>

<style>
	.overlay-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 89;
		animation: fade-in 0.2s ease;
	}

	@keyframes fade-in {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	.action-log-overlay {
		position: fixed;
		left: 0;
		top: 0;
		bottom: 0;
		width: 380px;
		max-width: 90vw;
		background: #141821;
		border-right: 2px solid #2a3441;
		display: flex;
		flex-direction: column;
		z-index: 90;
		transform: translateX(-100%);
		transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 4px 0 24px rgba(0, 0, 0, 0.5);
	}

	.action-log-overlay.open {
		transform: translateX(0);
	}

	/* Header */
	.log-header {
		padding: 1rem 1.25rem;
		background: #1a1f2e;
		border-bottom: 2px solid #2a3441;
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-shrink: 0;
	}

	.header-title {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.header-icon {
		font-size: 1.25rem;
	}

	.log-header h3 {
		margin: 0;
		font-size: 1.125rem;
		font-weight: 600;
		color: white;
	}

	.entry-count {
		font-size: 0.75rem;
		color: #9ca3af;
		font-weight: 600;
		padding: 0.25rem 0.5rem;
		background: #0f1419;
		border-radius: 4px;
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.clear-btn,
	.close-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid #2a3441;
		border-radius: 6px;
		font-size: 1rem;
		cursor: pointer;
		transition: all 0.2s;
		color: #9ca3af;
	}

	.clear-btn:hover:not(:disabled),
	.close-btn:hover {
		background: #2a3441;
		border-color: #374151;
		color: #fff;
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
		padding: 0.5rem;
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
	@media (max-width: 480px) {
		.action-log-overlay {
			width: 100%;
			max-width: 100%;
		}
	}
</style>


