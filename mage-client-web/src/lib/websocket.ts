// WebSocket client for Mage server

import type { WSMessage } from './types';

export class MageWebSocket {
	private ws: WebSocket | null = null;
	private url: string;
	private reconnectAttempts = 0;
	private maxReconnectAttempts = 5;
	private reconnectDelay = 1000;
	private messageHandlers: Map<string, (data: any) => void> = new Map();
	private onConnectCallbacks: (() => void)[] = [];
	private onDisconnectCallbacks: (() => void)[] = [];

	constructor(url: string) {
		this.url = url;
	}

	connect(): Promise<void> {
		return new Promise((resolve, reject) => {
			try {
				this.ws = new WebSocket(this.url);

				this.ws.onopen = () => {
					console.log('WebSocket connected');
					this.reconnectAttempts = 0;
					this.onConnectCallbacks.forEach(cb => cb());
					resolve();
				};

				this.ws.onmessage = (event) => {
					try {
						const message: WSMessage = JSON.parse(event.data);
						const handler = this.messageHandlers.get(message.type);
						if (handler) {
							handler(message.data);
						}
					} catch (error) {
						console.error('Failed to parse message:', error);
					}
				};

				this.ws.onerror = (error) => {
					console.error('WebSocket error:', error);
					reject(error);
				};

				this.ws.onclose = () => {
					console.log('WebSocket disconnected');
					this.onDisconnectCallbacks.forEach(cb => cb());
					this.attemptReconnect();
				};
			} catch (error) {
				reject(error);
			}
		});
	}

	private attemptReconnect() {
		if (this.reconnectAttempts < this.maxReconnectAttempts) {
			this.reconnectAttempts++;
			console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
			setTimeout(() => {
				this.connect().catch(console.error);
			}, this.reconnectDelay * this.reconnectAttempts);
		}
	}

	send(message: WSMessage) {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(message));
		} else {
			console.error('WebSocket not connected');
		}
	}

	on(type: string, handler: (data: any) => void) {
		this.messageHandlers.set(type, handler);
	}

	onConnect(callback: () => void) {
		this.onConnectCallbacks.push(callback);
	}

	onDisconnect(callback: () => void) {
		this.onDisconnectCallbacks.push(callback);
	}

	disconnect() {
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
	}

	isConnected(): boolean {
		return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
	}
}
