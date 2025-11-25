<script lang="ts">
	/**
	 * Individual action log entry component
	 * Displays timestamped game actions with player colors and icons
	 */

	import type { ActionLogEntry } from '$lib/types/game';

	let {
		action,
		playerColor = '#667eea'
	}: {
		action: ActionLogEntry;
		playerColor?: string;
	} = $props();

	/**
	 * Get icon for action type
	 */
	function getActionIcon(type: string): string {
		const icons: Record<string, string> = {
			play: '🃏',
			cast: '✨',
			tap: '↻',
			untap: '↺',
			attack: '⚔️',
			block: '🛡️',
			damage: '💥',
			destroy: '💀',
			exile: '🌟',
			draw: '📥',
			discard: '🗑️',
			shuffle: '🔀',
			search: '🔍',
			counter: '🚫',
			trigger: '⚡',
			ability: '🎯',
			enchant: '💫',
			equip: '🔨',
			sacrifice: '⚰️',
			mill: '📚',
			scry: '🔮',
			surveil: '👁️',
			phase: '⏩',
			priority: '🎪',
			mana: '💎',
			life: '❤️'
		};

		return icons[type] || '•';
	}

	/**
	 * Format timestamp
	 */
	function formatTime(timestamp: number): string {
		const date = new Date(timestamp);
		const hours = date.getHours().toString().padStart(2, '0');
		const minutes = date.getMinutes().toString().padStart(2, '0');
		const seconds = date.getSeconds().toString().padStart(2, '0');
		return `${hours}:${minutes}:${seconds}`;
	}
</script>

<div class="action-log-item" class:is-system={action.type === 'system'}>
	<div class="action-time">{formatTime(action.timestamp)}</div>

	<div class="action-content">
		<div class="action-icon">{getActionIcon(action.actionType)}</div>

		<div class="action-details">
			{#if action.playerName && action.type !== 'system'}
				<span class="action-player" style="color: {playerColor}">
					{action.playerName}
				</span>
			{/if}

			<span class="action-text">{action.text}</span>

			{#if action.cardName}
				<span class="action-card">{action.cardName}</span>
			{/if}
		</div>
	</div>
</div>

<style>
	.action-log-item {
		display: flex;
		gap: 0.75rem;
		padding: 0.625rem 1rem;
		border-bottom: 1px solid #2a3441;
		transition: background 0.15s;
	}

	.action-log-item:hover {
		background: #1a1f2e;
	}

	.action-log-item:last-child {
		border-bottom: none;
	}

	.action-log-item.is-system {
		background: rgba(167, 139, 250, 0.05);
	}

	.action-time {
		font-size: 0.75rem;
		color: #6b7280;
		font-family: 'Courier New', monospace;
		flex-shrink: 0;
		width: 60px;
	}

	.action-content {
		display: flex;
		gap: 0.5rem;
		align-items: flex-start;
		flex: 1;
		min-width: 0;
	}

	.action-icon {
		font-size: 1rem;
		line-height: 1;
		flex-shrink: 0;
	}

	.action-details {
		display: flex;
		flex-wrap: wrap;
		gap: 0.375rem;
		align-items: baseline;
		font-size: 0.875rem;
		line-height: 1.4;
		min-width: 0;
	}

	.action-player {
		font-weight: 600;
		flex-shrink: 0;
	}

	.action-text {
		color: #d1d5db;
	}

	.action-card {
		font-weight: 600;
		color: #fbbf24;
		font-style: italic;
	}

	.is-system .action-text {
		color: #a78bfa;
		font-style: italic;
	}

	.is-system .action-player {
		display: none;
	}
</style>
