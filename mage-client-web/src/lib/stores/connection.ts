import { writable, derived } from 'svelte/store';
import type {
	ConnectionState,
	ConnectionOptions,
	ConnectionEvent,
	ConnectionEventCallback
} from '$lib/types/connection';

/**
 * Default connection options
 */
const DEFAULT_OPTIONS: Required<ConnectionOptions> = {
	autoReconnect: true,
	maxReconnectAttempts: 10,
	reconnectDelay: 1000,
	maxReconnectDelay: 30000,
	enableHealthCheck: true,
	healthCheckInterval: 30000,
	healthCheckTimeout: 5000
};

/**
 * Initial connection state
 */
const INITIAL_STATE: ConnectionState = {
	status: 'disconnected',
	lastConnected: null,
	lastDisconnected: null,
	reconnectAttempt: 0,
	error: null,
	latency: null
};

/**
 * Create a connection status store
 */
function createConnectionStore() {
	const { subscribe, set, update } = writable<ConnectionState>(INITIAL_STATE);

	// Store options
	let options = { ...DEFAULT_OPTIONS };

	// Event listeners
	const eventListeners: ConnectionEventCallback[] = [];

	// Timers
	let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	let healthCheckTimer: ReturnType<typeof setInterval> | null = null;
	let healthCheckTimeoutTimer: ReturnType<typeof setTimeout> | null = null;

	// Health check state
	let lastPingTime = 0;
	let waitingForPong = false;

	/**
	 * Emit a connection event
	 */
	function emitEvent(event: ConnectionEvent): void {
		eventListeners.forEach((callback) => callback(event));
	}

	/**
	 * Calculate reconnection delay with exponential backoff
	 */
	function getReconnectDelay(attempt: number): number {
		const delay = options.reconnectDelay * Math.pow(2, attempt);
		return Math.min(delay, options.maxReconnectDelay);
	}

	/**
	 * Clear all timers
	 */
	function clearTimers(): void {
		if (reconnectTimer) {
			clearTimeout(reconnectTimer);
			reconnectTimer = null;
		}
		if (healthCheckTimer) {
			clearInterval(healthCheckTimer);
			healthCheckTimer = null;
		}
		if (healthCheckTimeoutTimer) {
			clearTimeout(healthCheckTimeoutTimer);
			healthCheckTimeoutTimer = null;
		}
	}

	/**
	 * Start health check ping/pong mechanism
	 */
	function startHealthCheck(): void {
		if (!options.enableHealthCheck) return;

		// Clear existing timer
		if (healthCheckTimer) {
			clearInterval(healthCheckTimer);
		}

		// Send ping at regular intervals
		healthCheckTimer = setInterval(() => {
			sendPing();
		}, options.healthCheckInterval);
	}

	/**
	 * Stop health check
	 */
	function stopHealthCheck(): void {
		if (healthCheckTimer) {
			clearInterval(healthCheckTimer);
			healthCheckTimer = null;
		}
		if (healthCheckTimeoutTimer) {
			clearTimeout(healthCheckTimeoutTimer);
			healthCheckTimeoutTimer = null;
		}
		waitingForPong = false;
	}

	/**
	 * Send ping and wait for pong
	 */
	function sendPing(): void {
		if (waitingForPong) {
			// Previous ping timed out - connection may be dead
			handleConnectionLost(new Error('Health check timeout'));
			return;
		}

		lastPingTime = Date.now();
		waitingForPong = true;

		// Set timeout for pong response
		healthCheckTimeoutTimer = setTimeout(() => {
			if (waitingForPong) {
				handleConnectionLost(new Error('Health check timeout'));
			}
		}, options.healthCheckTimeout);

		// In a real implementation, this would send a ping to the server
		// For now, we'll simulate it with a mock
		if (import.meta.env.DEV) {
			console.log('[Connection] Ping sent');
		}
	}

	/**
	 * Handle pong response
	 */
	function handlePong(): void {
		if (!waitingForPong) return;

		const latency = Date.now() - lastPingTime;
		waitingForPong = false;

		if (healthCheckTimeoutTimer) {
			clearTimeout(healthCheckTimeoutTimer);
			healthCheckTimeoutTimer = null;
		}

		update((state) => ({
			...state,
			latency
		}));

		if (import.meta.env.DEV) {
			console.log(`[Connection] Pong received (${latency}ms)`);
		}
	}

	/**
	 * Handle connection lost
	 */
	function handleConnectionLost(error?: Error): void {
		stopHealthCheck();

		update((state) => ({
			...state,
			status: 'disconnected',
			lastDisconnected: Date.now(),
			error: error?.message || 'Connection lost'
		}));

		emitEvent({
			type: 'disconnected',
			timestamp: Date.now(),
			error
		});

		// Attempt reconnection if enabled
		if (options.autoReconnect) {
			scheduleReconnect();
		}
	}

	/**
	 * Schedule reconnection attempt
	 */
	function scheduleReconnect(): void {
		update((state) => {
			const attempt = state.reconnectAttempt;

			// Check if we've exceeded max attempts
			if (attempt >= options.maxReconnectAttempts) {
				emitEvent({
					type: 'reconnect_failed',
					timestamp: Date.now(),
					attempt,
					maxAttempts: options.maxReconnectAttempts
				});

				return {
					...state,
					status: 'disconnected',
					error: 'Maximum reconnection attempts exceeded'
				};
			}

			// Calculate delay and schedule reconnect
			const delay = getReconnectDelay(attempt);

			if (import.meta.env.DEV) {
				console.log(
					`[Connection] Scheduling reconnect attempt ${attempt + 1}/${options.maxReconnectAttempts} in ${delay}ms`
				);
			}

			reconnectTimer = setTimeout(() => {
				doConnect();
			}, delay);

			emitEvent({
				type: 'reconnecting',
				timestamp: Date.now(),
				attempt: attempt + 1,
				maxAttempts: options.maxReconnectAttempts
			});

			return {
				...state,
				status: 'reconnecting',
				reconnectAttempt: attempt + 1
			};
		});
	}

	/**
	 * Internal connect function
	 */
	function doConnect(): void {
		clearTimers();

		update((state) => ({
			...state,
			status: 'connecting',
			error: null
		}));

		// Simulate connection (in real implementation, this would establish gRPC/WebSocket connection)
		setTimeout(() => {
			update((state) => ({
				...state,
				status: 'connected',
				lastConnected: Date.now(),
				reconnectAttempt: 0,
				error: null
			}));

			emitEvent({
				type: 'connected',
				timestamp: Date.now()
			});

			// Start health check
			startHealthCheck();
		}, 500);
	}

	return {
		subscribe,

		/**
		 * Initialize connection with options
		 */
		initialize(opts: ConnectionOptions = {}): void {
			options = { ...DEFAULT_OPTIONS, ...opts };
		},

		/**
		 * Connect to server
		 */
		connect(): void {
			doConnect();
		},

		/**
		 * Disconnect from server
		 */
		disconnect(): void {
			clearTimers();
			stopHealthCheck();

			update((state) => ({
				...state,
				status: 'disconnected',
				lastDisconnected: Date.now(),
				reconnectAttempt: 0,
				error: null
			}));

			emitEvent({
				type: 'disconnected',
				timestamp: Date.now()
			});
		},

		/**
		 * Manually trigger reconnection
		 */
		reconnect(): void {
			clearTimers();
			stopHealthCheck();

			update((state) => ({
				...state,
				reconnectAttempt: 0,
				error: null
			}));

			doConnect();
		},

		/**
		 * Handle pong response from server
		 */
		pong: handlePong,

		/**
		 * Simulate connection error (for testing)
		 */
		simulateError(error: Error): void {
			handleConnectionLost(error);
		},

		/**
		 * Add event listener
		 */
		addEventListener(callback: ConnectionEventCallback): () => void {
			eventListeners.push(callback);

			// Return unsubscribe function
			return () => {
				const index = eventListeners.indexOf(callback);
				if (index > -1) {
					eventListeners.splice(index, 1);
				}
			};
		},

		/**
		 * Reset to initial state
		 */
		reset(): void {
			clearTimers();
			stopHealthCheck();
			set(INITIAL_STATE);
		}
	};
}

/**
 * Global connection store instance
 */
export const connection = createConnectionStore();

/**
 * Derived store for connection status only
 */
export const connectionStatus = derived(connection, ($connection) => $connection.status);

/**
 * Derived store for checking if connected
 */
export const isConnected = derived(connection, ($connection) => $connection.status === 'connected');

/**
 * Derived store for latency display
 */
export const connectionLatency = derived(connection, ($connection) => $connection.latency);
