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
	on:mouseenter={() => (showTooltip = true)}
	on:mouseleave={() => (showTooltip = false)}
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
					<button class="reconnect-btn" on:click={handleReconnect}>Reconnect Now</button>
				{:else}
					<p>Not connected to server</p>
					{#if error}
						<p class="error">{error}</p>
					{/if}
					<button class="reconnect-btn" on:click={handleReconnect}>Reconnect</button>
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
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background-color: rgba(255, 255, 255, 0.1);
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: default;
		transition: background-color 0.2s;
	}

	.status-indicator:hover {
		background-color: rgba(255, 255, 255, 0.15);
	}

	.status-icon {
		font-size: 1rem;
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
		color: white;
		display: none;
	}

	/* Tooltip */
	.tooltip {
		position: absolute;
		top: calc(100% + 0.5rem);
		right: 0;
		background-color: #1f2937;
		color: white;
		padding: 0.75rem;
		border-radius: 0.5rem;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
		min-width: 200px;
		max-width: 300px;
		z-index: 1000;
		animation: fadeIn 0.2s ease-out;
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
		font-weight: 600;
		font-size: 0.875rem;
		margin-bottom: 0.5rem;
	}

	.tooltip-content {
		font-size: 0.75rem;
		color: #d1d5db;
		line-height: 1.4;
	}

	.tooltip-content p {
		margin: 0 0 0.25rem 0;
	}

	.tooltip-content p:last-child {
		margin-bottom: 0;
	}

	.latency {
		color: #10b981;
		font-weight: 500;
	}

	.attempt {
		color: #f59e0b;
		font-weight: 500;
	}

	.error {
		color: #ef4444;
		font-size: 0.7rem;
		font-style: italic;
	}

	.reconnect-btn {
		margin-top: 0.5rem;
		width: 100%;
		padding: 0.375rem 0.5rem;
		background-color: #3b82f6;
		color: white;
		border: none;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.reconnect-btn:hover {
		background-color: #2563eb;
	}

	/* Responsive */
	@media (min-width: 640px) {
		.status-text {
			display: inline;
		}
	}
</style>
