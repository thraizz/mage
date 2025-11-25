/**
 * WebSocket Store
 * Manages WebSocket connection and real-time event subscriptions
 */

import { writable, get } from 'svelte/store';
import type { CallbackMethod, ChatMessageData, StartGameData, GameUpdateData, GameInitData } from '$lib/generated/mage/v1/websocket';
import {
	callbackMethodFromJSON,
	ChatMessageData as ChatMessageDataCodec,
	StartGameData as StartGameDataCodec,
	GameUpdateData as GameUpdateDataCodec,
	GameInitData as GameInitDataCodec
} from '$lib/generated/mage/v1/websocket';
import { BinaryReader } from '@bufbuild/protobuf/wire';

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
 * Decode ChatMessageData from protobuf bytes
 */
function decodeChatMessageData(bytes: Uint8Array): ChatMessageData {
	const reader = new BinaryReader(bytes);
	return ChatMessageDataCodec.decode(reader);
}

/**
 * Decode StartGameData from protobuf bytes
 */
function decodeStartGameData(bytes: Uint8Array): StartGameData {
	const reader = new BinaryReader(bytes);
	return StartGameDataCodec.decode(reader);
}

/**
 * Decode GameUpdateData from protobuf bytes
 */
function decodeGameUpdateData(bytes: Uint8Array): GameUpdateData {
	const reader = new BinaryReader(bytes);
	return GameUpdateDataCodec.decode(reader);
}

/**
 * Decode GameInitData from protobuf bytes
 */
function decodeGameInitData(bytes: Uint8Array): GameInitData {
	const reader = new BinaryReader(bytes);
	return GameInitDataCodec.decode(reader);
}

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
			// Check if already connected or connecting
			if (ws) {
				if (ws.readyState === WebSocket.OPEN) {
					resolve();
					return;
				} else if (ws.readyState === WebSocket.CONNECTING) {
					// Already connecting, wait for it to finish
					const checkConnection = setInterval(() => {
						if (ws && ws.readyState === WebSocket.OPEN) {
							clearInterval(checkConnection);
							resolve();
						} else if (!ws || ws.readyState === WebSocket.CLOSED) {
							clearInterval(checkConnection);
							reject(new Error('WebSocket connection failed'));
						}
					}, 100);
					return;
				}
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
						console.log('[WebSocket] Raw message received:', {
							method: serverEvent.method,
							messageId: serverEvent.messageId,
							sessionId: serverEvent.sessionId,
							hasData: !!serverEvent.data
						});

						// Convert method string to enum if needed
						if (typeof serverEvent.method === 'string') {
							serverEvent.method = callbackMethodFromJSON(serverEvent.method);
						}

						// Unpack protobuf Any type if present
						let eventData = serverEvent.data;
						let decodedTypeUrl: string | null = null;
						if (eventData && typeof eventData === 'object') {
							// Check for protojson format (uses @type) or binary format (uses typeUrl + value)
							const dataObj = eventData as Record<string, unknown>;
							const typeUrl = (dataObj['@type'] as string) || (dataObj['typeUrl'] as string);
							
							if (typeUrl) {
								decodedTypeUrl = typeUrl;
								console.log('[WebSocket] Detected Any type:', typeUrl);
								
								// Check if this is JSON format (protojson) or binary format
								if ('@type' in dataObj) {
									// protojson format - data is already JSON decoded, just extract it
									// The data fields are siblings of @type, not nested under 'value'
									console.log('[WebSocket] Using protojson format (already decoded)');
									// For protojson, the wrapped message fields are at the same level as @type
									// We just pass the whole object through - handlers will access .game, .message, etc.
									// eventData is already correctly set
								} else if ('typeUrl' in dataObj && 'value' in dataObj) {
									// Binary format - decode from base64
									try {
										const base64Value = dataObj.value as string;
										const binaryString = atob(base64Value);
										const bytes = new Uint8Array(binaryString.length);
										for (let i = 0; i < binaryString.length; i++) {
											bytes[i] = binaryString.charCodeAt(i);
										}

										// Parse based on type_url
										if (typeUrl === 'type.googleapis.com/mage.v1.ChatMessageData') {
											eventData = decodeChatMessageData(bytes);
										} else if (typeUrl === 'type.googleapis.com/mage.v1.StartGameData') {
											eventData = decodeStartGameData(bytes);
										} else if (typeUrl === 'type.googleapis.com/mage.v1.GameUpdateData') {
											eventData = decodeGameUpdateData(bytes);
										} else if (typeUrl === 'type.googleapis.com/mage.v1.GameInitData') {
											eventData = decodeGameInitData(bytes);
										} else {
											console.warn('[WebSocket] Unknown typeUrl, unable to decode:', typeUrl);
										}
									} catch (err) {
										console.error('[WebSocket] Failed to decode binary Any type:', err);
									}
								}
							}
						}

						console.log('[WebSocket] Decoded event:', {
							method: serverEvent.method,
							typeUrl: decodedTypeUrl,
							data: eventData
						});

						// Call all registered handlers for this method
						const methodHandlers = handlers.get(serverEvent.method);
						const handlerCount = methodHandlers?.size || 0;
						console.log(`[WebSocket] Found ${handlerCount} handler(s) for method ${serverEvent.method}`);
						
						if (methodHandlers) {
							methodHandlers.forEach((handler) => {
								try {
									handler(eventData);
								} catch (err) {
									console.error(`[WebSocket] Handler error for method ${serverEvent.method}:`, err);
								}
							});
						} else {
							console.warn(`[WebSocket] No handlers registered for method: ${serverEvent.method}`);
						}
					} catch (err) {
						console.error('[WebSocket] Failed to parse message:', err, 'Raw data:', event.data);
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
		console.log(`[WebSocket] Handler registered for method ${method}, total handlers: ${handlers.get(method)!.size}`);

		// Return unsubscribe function
		return () => {
			const methodHandlers = handlers.get(method);
			if (methodHandlers) {
				methodHandlers.delete(handler);
				console.log(`[WebSocket] Handler unregistered for method ${method}, remaining: ${methodHandlers.size}`);
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
	return get(websocketStore).state;
}
