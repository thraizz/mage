<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { connection, connectionLatency } from '$lib/stores/connection';

	let showTooltip = false;
	let unsubscribe: (() => void) | null = null;

	$: connectionState = $connection.status;
	$: reconnectAttempts = $connection.reconnectAttempt;
	$: latency = $connectionLatency;
	$: error = $connection.error;

	onMount(() => {
		// Initialize connection with default options
		connection.initialize({
			autoReconnect: true,
			maxReconnectAttempts: 10,
			reconnectDelay: 1000,
			maxReconnectDelay: 30000,
			enableHealthCheck: true,
			healthCheckInterval: 30000,
			healthCheckTimeout: 5000
		});

		// Listen to connection events
		unsubscribe = connection.addEventListener((event) => {
			if (import.meta.env.DEV) {
				console.log('[ConnectionStatus] Event:', event);
			}
		});

		// Auto-connect on mount
		connection.connect();
	});

	onDestroy(() => {
		if (unsubscribe) {
			unsubscribe();
		}
	});

	function getStatusColor(): string {
		switch (connectionState) {
			case 'connected':
				return '#10b981'; // Green
			case 'connecting':
			case 'reconnecting':
				return '#f59e0b'; // Orange
			case 'disconnected':
				return '#ef4444'; // Red
			default:
				return '#6b7280'; // Gray
		}
	}

	function getStatusText(): string {
		switch (connectionState) {
			case 'connected':
				return 'Online';
			case 'connecting':
				return 'Connecting...';
			case 'reconnecting':
				return reconnectAttempts > 0 ? `Reconnecting... (${reconnectAttempts})` : 'Reconnecting...';
			case 'disconnected':
				return 'Offline';
			default:
				return 'Unknown';
		}
	}

	function getStatusIcon(): string {
		switch (connectionState) {
			case 'connected':
				return '●';
			case 'connecting':
			case 'reconnecting':
				return '◐';
			case 'disconnected':
				return '●';
			default:
				return '○';
		}
	}

	function handleReconnect(): void {
		connection.reconnect();
	}
</script>

<div
	class="connection-status"
	onmouseenter={() => (showTooltip = true)}
	onmouseleave={() => (showTooltip = false)}
	role="status"
	aria-live="polite"
>
	<div class="status-indicator" style="color: {getStatusColor()}">
		<span class="status-icon" class:pulse={connectionState !== 'connected'}>
			{getStatusIcon()}
		</span>
		<span class="status-text">{getStatusText()}</span>
	</div>

	{#if showTooltip}
		<div class="tooltip">
			<div class="tooltip-title">Connection Status</div>
			<div class="tooltip-content">
				{#if connectionState === 'connected'}
					<p>You are connected to the server</p>
					{#if latency !== null}
						<p class="latency">Latency: {latency}ms</p>
					{/if}
				{:else if connectionState === 'connecting'}
					<p>Establishing connection to server...</p>
				{:else if connectionState === 'reconnecting'}
					<p>Lost connection. Attempting to reconnect</p>
					<p class="attempt">Attempt {reconnectAttempts} of 10</p>
					<button class="reconnect-btn" onclick={handleReconnect}>Reconnect Now</button>
				{:else}
					<p>Not connected to server</p>
					{#if error}
						<p class="error">{error}</p>
					{/if}
					<button class="reconnect-btn" onclick={handleReconnect}>Reconnect</button>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.connection-status {
		position: relative;
		display: flex;
		align-items: center;
	}

	.status-indicator {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		cursor: default;
		transition: all var(--transition-fast);
	}

	.status-indicator:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
	}

	.status-icon {
		font-size: var(--text-base);
		line-height: 1;
	}

	.status-icon.pulse {
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}

	.status-text {
		color: var(--text-muted);
		display: none;
	}

	/* Tooltip */
	.tooltip {
		position: absolute;
		top: calc(100% + var(--space-2));
		right: 0;
		background: var(--bg-slate);
		border: 1px solid var(--border-subtle);
		color: var(--text-bright);
		padding: var(--space-3);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-lg);
		min-width: 200px;
		max-width: 300px;
		z-index: var(--z-tooltip);
		animation: fadeIn var(--transition-fast);
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-5px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.tooltip-title {
		font-weight: var(--weight-semibold);
		font-size: var(--text-sm);
		color: var(--text-bright);
		margin-bottom: var(--space-2);
	}

	.tooltip-content {
		font-size: var(--text-xs);
		color: var(--text-muted);
		line-height: var(--leading-normal);
	}

	.tooltip-content p {
		margin: 0 0 var(--space-1) 0;
	}

	.tooltip-content p:last-child {
		margin-bottom: 0;
	}

	.latency {
		color: var(--status-success);
		font-weight: var(--weight-medium);
	}

	.attempt {
		color: var(--status-warning);
		font-weight: var(--weight-medium);
	}

	.error {
		color: var(--status-error);
		font-size: var(--text-xs);
		font-style: italic;
	}

	.reconnect-btn {
		margin-top: var(--space-2);
		width: 100%;
		padding: var(--space-2) var(--space-3);
		background: var(--accent-gold);
		color: var(--bg-void);
		border: none;
		border-radius: var(--radius-md);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: background var(--transition-fast);
	}

	.reconnect-btn:hover {
		background: var(--accent-gold-bright);
	}

	/* Responsive */
	@media (min-width: 640px) {
		.status-text {
			display: inline;
		}
	}
</style>
