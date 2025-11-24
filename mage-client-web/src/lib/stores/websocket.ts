/**
 * WebSocket Store
 * Manages WebSocket connection and real-time event subscriptions
 */

import { writable, get } from 'svelte/store';
import type { CallbackMethod } from '$lib/generated/mage/v1/websocket';
import { callbackMethodFromJSON } from '$lib/generated/mage/v1/websocket';

/**
 * WebSocket connection state
 */
export type WebSocketState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

/**
 * WebSocket message structure matching server format
 */
export interface ServerEvent {
	method: CallbackMethod;
	messageId: string;
	sessionId: string;
	data?: unknown;
}

/**
 * Event handler function
 */
type EventHandler = (data: unknown) => void;

/**
 * WebSocket store state
 */
interface WebSocketStore {
	state: WebSocketState;
	error: string | null;
	lastConnected: number | null;
	reconnectAttempts: number;
}

const WEBSOCKET_URL = import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:17179/ws';

const MAX_RECONNECT_ATTEMPTS = 10;
const BASE_RECONNECT_DELAY = 1000;
const MAX_RECONNECT_DELAY = 30000;

/**
 * WebSocket store
 */
function createWebSocketStore() {
	const { subscribe, update } = writable<WebSocketStore>({
		state: 'disconnected',
		error: null,
		lastConnected: null,
		reconnectAttempts: 0
	});

	let ws: WebSocket | null = null;
	let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	let sessionId: string | null = null;

	// Event handlers map: method -> handler
	const handlers = new Map<CallbackMethod, Set<EventHandler>>();

	/**
	 * Calculate reconnect delay with exponential backoff
	 */
	function getReconnectDelay(attempts: number): number {
		const delay = Math.min(BASE_RECONNECT_DELAY * Math.pow(2, attempts), MAX_RECONNECT_DELAY);
		// Add jitter to prevent thundering herd
		return delay + Math.random() * 1000;
	}

	/**
	 * Connect to WebSocket server
	 */
	function connect(newSessionId: string): Promise<void> {
		return new Promise((resolve, reject) => {
			if (ws && ws.readyState === WebSocket.OPEN) {
				resolve();
				return;
			}

			sessionId = newSessionId;

			update((s) => ({
				...s,
				state: s.reconnectAttempts > 0 ? 'reconnecting' : 'connecting',
				error: null
			}));

			try {
				// Include sessionId in URL for server authentication
				const url = `${WEBSOCKET_URL}?sessionId=${encodeURIComponent(sessionId)}`;
				ws = new WebSocket(url);

				ws.onopen = () => {
					console.log('[WebSocket] Connected successfully');
					update((s) => ({
						state: 'connected',
						error: null,
						lastConnected: Date.now(),
						reconnectAttempts: 0
					}));
					resolve();
				};

				ws.onmessage = (event) => {
					try {
						const serverEvent = JSON.parse(event.data) as ServerEvent;

						// Convert method string to enum if needed
						if (typeof serverEvent.method === 'string') {
							serverEvent.method = callbackMethodFromJSON(serverEvent.method);
						}

						// Call all registered handlers for this method
						const methodHandlers = handlers.get(serverEvent.method);
						if (methodHandlers) {
							methodHandlers.forEach((handler) => {
								try {
									handler(serverEvent.data);
								} catch (err) {
									console.error(`[WebSocket] Handler error for method ${serverEvent.method}:`, err);
								}
							});
						}
					} catch (err) {
						console.error('[WebSocket] Failed to parse message:', err);
					}
				};

				ws.onerror = (event) => {
					console.error('[WebSocket] Error:', event);
					update((s) => ({
						...s,
						state: 'disconnected',
						error: 'WebSocket connection error'
					}));
					reject(new Error('WebSocket connection failed'));
				};

				ws.onclose = (event) => {
					console.log('[WebSocket] Disconnected:', event.code, event.reason);

					const currentState = get({ subscribe });
					if (currentState.reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
						// Attempt reconnection
						const delay = getReconnectDelay(currentState.reconnectAttempts);
						console.log(
							`[WebSocket] Reconnecting in ${delay}ms (attempt ${currentState.reconnectAttempts + 1}/${MAX_RECONNECT_ATTEMPTS})`
						);

						update((s) => ({
							...s,
							state: 'reconnecting',
							reconnectAttempts: s.reconnectAttempts + 1
						}));

						reconnectTimer = setTimeout(() => {
							if (sessionId) {
								connect(sessionId).catch(console.error);
							}
						}, delay);
					} else {
						// Max attempts reached
						update((s) => ({
							...s,
							state: 'disconnected',
							error: 'Failed to reconnect after maximum attempts'
						}));
					}
				};
			} catch (err) {
				update((s) => ({
					...s,
					state: 'disconnected',
					error: err instanceof Error ? err.message : 'Failed to create WebSocket'
				}));
				reject(err);
			}
		});
	}

	/**
	 * Disconnect from WebSocket server
	 */
	function disconnect() {
		if (reconnectTimer) {
			clearTimeout(reconnectTimer);
			reconnectTimer = null;
		}

		if (ws) {
			ws.close();
			ws = null;
		}

		update((s) => ({
			...s,
			state: 'disconnected',
			error: null,
			reconnectAttempts: 0
		}));

		sessionId = null;
	}

	/**
	 * Subscribe to a specific event method
	 */
	function on(method: CallbackMethod, handler: EventHandler): () => void {
		if (!handlers.has(method)) {
			handlers.set(method, new Set());
		}
		handlers.get(method)!.add(handler);

		// Return unsubscribe function
		return () => {
			const methodHandlers = handlers.get(method);
			if (methodHandlers) {
				methodHandlers.delete(handler);
				if (methodHandlers.size === 0) {
					handlers.delete(method);
				}
			}
		};
	}

	/**
	 * Send a message to the server
	 */
	function send(data: unknown): void {
		if (ws && ws.readyState === WebSocket.OPEN) {
			ws.send(JSON.stringify(data));
		} else {
			console.error('[WebSocket] Cannot send message: not connected');
		}
	}

	/**
	 * Check if connected
	 */
	function isConnected(): boolean {
		return ws !== null && ws.readyState === WebSocket.OPEN;
	}

	/**
	 * Force reconnection
	 */
	function reconnect(): void {
		if (sessionId) {
			disconnect();
			update((s) => ({ ...s, reconnectAttempts: 0 }));
			connect(sessionId).catch(console.error);
		}
	}

	return {
		subscribe,
		connect,
		disconnect,
		on,
		send,
		isConnected,
		reconnect
	};
}

/**
 * Global WebSocket store instance
 */
export const websocketStore = createWebSocketStore();

/**
 * Helper function to get current connection state
 */
export function getWebSocketState(): WebSocketState {
	let state: WebSocketState = 'disconnected';
	websocketStore.subscribe((s) => {
		state = s.state;
	})();
	return state;
}
