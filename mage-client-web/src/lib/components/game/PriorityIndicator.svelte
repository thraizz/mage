<script lang="ts">
	// Props
	let {
		hasPriority = false,
		activePlayerId = '',
		localPlayerId = '',
		playerName = 'Player',
		animated = true
	}: {
		hasPriority?: boolean;
		activePlayerId?: string;
		localPlayerId?: string;
		playerName?: string;
		animated?: boolean;
	} = $props();

	// Derived values
	const isYourTurn = $derived(activePlayerId === localPlayerId);
	const priorityText = $derived(
		hasPriority ? 'Your Priority' : isYourTurn ? `${playerName}'s Priority` : 'Waiting...'
	);
	const statusClass = $derived(hasPriority ? 'has-priority' : isYourTurn ? 'active' : 'waiting');
</script>

<div class="priority-indicator" class:animated class:has-priority={hasPriority}>
	<div class="priority-status {statusClass}">
		<div class="priority-icon">
			{#if hasPriority}
				⚡
			{:else if isYourTurn}
				⏳
			{:else}
				⏸️
			{/if}
		</div>
		<span class="priority-text">{priorityText}</span>
	</div>

	{#if hasPriority}
		<div class="priority-hint">
			<span class="hint-text">You may take an action</span>
		</div>
	{/if}
</div>

<style>
	.priority-indicator {
		background: linear-gradient(135deg, #1a1f2e 0%, #0f1419 100%);
		border: 2px solid #2a3441;
		border-radius: 8px;
		padding: 0.75rem 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
		transition: all 0.3s ease;
	}

	.priority-indicator.has-priority {
		border-color: #fbbf24;
		box-shadow: 0 4px 16px rgba(251, 191, 36, 0.4);
	}

	.priority-status {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.priority-icon {
		font-size: 1.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		background: #2a3441;
		border-radius: 50%;
		transition: all 0.3s ease;
	}

	.priority-status.has-priority .priority-icon {
		background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
	}

	.priority-status.active .priority-icon {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.priority-status.waiting .priority-icon {
		background: #374151;
		opacity: 0.6;
	}

	.priority-text {
		font-size: 1rem;
		font-weight: 600;
		color: #ffffff;
		transition: color 0.3s ease;
	}

	.priority-status.has-priority .priority-text {
		color: #fbbf24;
	}

	.priority-status.active .priority-text {
		color: #9ca3af;
	}

	.priority-status.waiting .priority-text {
		color: #6b7280;
		font-style: italic;
	}

	.priority-hint {
		padding: 0.5rem;
		background: rgba(251, 191, 36, 0.1);
		border-radius: 4px;
		border-left: 3px solid #fbbf24;
	}

	.hint-text {
		font-size: 0.875rem;
		color: #fbbf24;
		font-weight: 500;
	}

	/* Animated pulse effect */
	.priority-indicator.animated.has-priority {
		animation: priority-pulse 2s ease-in-out infinite;
	}

	@keyframes priority-pulse {
		0%,
		100% {
			box-shadow: 0 4px 16px rgba(251, 191, 36, 0.4);
		}
		50% {
			box-shadow: 0 4px 24px rgba(251, 191, 36, 0.6);
		}
	}

	.priority-indicator.animated.has-priority .priority-icon {
		animation: icon-pulse 1.5s ease-in-out infinite;
	}

	@keyframes icon-pulse {
		0%,
		100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.15);
		}
	}

	/* Responsive */
	@media (max-width: 768px) {
		.priority-indicator {
			padding: 0.5rem 0.75rem;
		}

		.priority-icon {
			font-size: 1.25rem;
			width: 28px;
			height: 28px;
		}

		.priority-text {
			font-size: 0.875rem;
		}

		.hint-text {
			font-size: 0.75rem;
		}
	}
</style>
