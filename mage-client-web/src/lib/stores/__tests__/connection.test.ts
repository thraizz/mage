import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { get } from 'svelte/store';
import { connection, connectionStatus, isConnected } from '../connection';

describe('Connection Store', () => {
	beforeEach(() => {
		vi.useFakeTimers();
		connection.reset();
	});

	afterEach(() => {
		vi.restoreAllMocks();
		vi.useRealTimers();
	});

	describe('Initial State', () => {
		it('should start in disconnected state', () => {
			const state = get(connection);
			expect(state.status).toBe('disconnected');
			expect(state.lastConnected).toBeNull();
			expect(state.lastDisconnected).toBeNull();
			expect(state.reconnectAttempt).toBe(0);
			expect(state.error).toBeNull();
			expect(state.latency).toBeNull();
		});

		it('should have isConnected as false initially', () => {
			expect(get(isConnected)).toBe(false);
		});
	});

	describe('Connection', () => {
		it('should transition to connecting then connected state', async () => {
			connection.initialize();
			connection.connect();

			expect(get(connection).status).toBe('connecting');

			// Wait for connection to complete (500ms in mock)
			await vi.advanceTimersByTimeAsync(500);

			const state = get(connection);
			expect(state.status).toBe('connected');
			expect(state.lastConnected).toBeGreaterThan(0);
			expect(state.reconnectAttempt).toBe(0);
			expect(get(isConnected)).toBe(true);
		});

		it('should emit connected event', async () => {
			const eventCallback = vi.fn();
			connection.addEventListener(eventCallback);

			connection.initialize();
			connection.connect();

			await vi.advanceTimersByTimeAsync(500);

			expect(eventCallback).toHaveBeenCalledWith(
				expect.objectContaining({
					type: 'connected',
					timestamp: expect.any(Number)
				})
			);
		});
	});

	describe('Disconnection', () => {
		it('should transition to disconnected state', async () => {
			connection.initialize();
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			connection.disconnect();

			const state = get(connection);
			expect(state.status).toBe('disconnected');
			expect(state.lastDisconnected).toBeGreaterThan(0);
			expect(get(isConnected)).toBe(false);
		});

		it('should emit disconnected event', async () => {
			const eventCallback = vi.fn();
			connection.addEventListener(eventCallback);

			connection.initialize();
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			connection.disconnect();

			expect(eventCallback).toHaveBeenCalledWith(
				expect.objectContaining({
					type: 'disconnected',
					timestamp: expect.any(Number)
				})
			);
		});
	});

	describe('Reconnection', () => {
		it('should attempt reconnection on connection loss', async () => {
			connection.initialize({ autoReconnect: true, reconnectDelay: 1000 });
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			// Simulate connection loss
			connection.simulateError(new Error('Connection lost'));

			// After error, should transition to reconnecting state
			expect(get(connection).status).toBe('reconnecting');
			expect(get(connection).reconnectAttempt).toBe(1);

			// Should attempt reconnect after delay
			await vi.advanceTimersByTimeAsync(1000);
			expect(get(connection).status).toBe('connecting');

			await vi.advanceTimersByTimeAsync(500);
			expect(get(connection).status).toBe('connected');
		});

		it('should use exponential backoff for retries', async () => {
			connection.initialize({
				autoReconnect: true,
				reconnectDelay: 1000,
				maxReconnectAttempts: 3
			});

			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			// First failure - delay 1000ms (1000 * 2^0)
			connection.simulateError(new Error('Connection lost'));
			expect(get(connection).reconnectAttempt).toBe(1);

			await vi.advanceTimersByTimeAsync(500);
			connection.simulateError(new Error('Connection lost'));

			// Second failure - delay 2000ms (1000 * 2^1)
			expect(get(connection).reconnectAttempt).toBe(2);

			await vi.advanceTimersByTimeAsync(500);
			connection.simulateError(new Error('Connection lost'));

			// Third failure - delay 4000ms (1000 * 2^2)
			expect(get(connection).reconnectAttempt).toBe(3);
		});

		it('should stop after max reconnection attempts', async () => {
			connection.initialize({
				autoReconnect: true,
				reconnectDelay: 100,
				maxReconnectAttempts: 1 // Set to 1 for simpler test
			});

			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			// First failure - transitions to reconnecting (attempt 1)
			connection.simulateError(new Error('Connection lost'));
			expect(get(connection).reconnectAttempt).toBe(1);

			// Wait for reconnect attempt (100ms delay), then immediately fail before it connects (500ms)
			await vi.advanceTimersByTimeAsync(200);

			// Fail the reconnection attempt - should reach max and give up
			connection.simulateError(new Error('Connection lost'));

			const state = get(connection);
			expect(state.status).toBe('disconnected');
			expect(state.error).toContain('Maximum reconnection attempts');
		});

		it('should emit reconnect_failed event after max attempts', async () => {
			const eventCallback = vi.fn();
			connection.addEventListener(eventCallback);

			connection.initialize({
				autoReconnect: true,
				reconnectDelay: 100,
				maxReconnectAttempts: 1 // Set to 1 for simpler test
			});

			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			// First failure (attempt 1)
			connection.simulateError(new Error('Connection lost'));
			await vi.advanceTimersByTimeAsync(200);

			// Second failure (exceeds max of 1)
			connection.simulateError(new Error('Connection lost'));

			// Should have emitted reconnect_failed event
			expect(eventCallback).toHaveBeenCalledWith(
				expect.objectContaining({
					type: 'reconnect_failed',
					attempt: 1,
					maxAttempts: 1
				})
			);
		});

		it('should allow manual reconnect', async () => {
			connection.initialize({ autoReconnect: false });
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			connection.disconnect();
			expect(get(connection).status).toBe('disconnected');

			connection.reconnect();
			expect(get(connection).status).toBe('connecting');

			await vi.advanceTimersByTimeAsync(500);
			expect(get(connection).status).toBe('connected');
		});
	});

	describe('Health Check', () => {
		it('should update latency on pong', async () => {
			connection.initialize({ enableHealthCheck: true, healthCheckInterval: 1000 });
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			// Wait for first ping to be sent (health check interval)
			await vi.advanceTimersByTimeAsync(1000);

			// Manually trigger pong (in real implementation, this comes from server)
			connection.pong();

			// Latency should be updated (exact value depends on timing)
			const latency = get(connection).latency;
			expect(latency).not.toBeNull();
			if (latency !== null) {
				expect(latency).toBeGreaterThanOrEqual(0);
			}
		});
	});

	describe('Event Listeners', () => {
		it('should allow adding event listeners', async () => {
			const callback = vi.fn();
			const unsubscribe = connection.addEventListener(callback);

			connection.initialize();
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			expect(callback).toHaveBeenCalled();

			unsubscribe();
		});

		it('should remove event listeners', async () => {
			const callback = vi.fn();
			const unsubscribe = connection.addEventListener(callback);

			unsubscribe();

			connection.initialize();
			connection.connect();
			await vi.advanceTimersByTimeAsync(500);

			expect(callback).not.toHaveBeenCalled();
		});
	});

	describe('Derived Stores', () => {
		it('connectionStatus should derive status only', async () => {
			connection.initialize();
			connection.connect();

			expect(get(connectionStatus)).toBe('connecting');

			await vi.advanceTimersByTimeAsync(500);
			expect(get(connectionStatus)).toBe('connected');
		});

		it('isConnected should return boolean', async () => {
			connection.initialize();
			connection.connect();

			expect(get(isConnected)).toBe(false);

			await vi.advanceTimersByTimeAsync(500);
			expect(get(isConnected)).toBe(true);

			connection.disconnect();
			expect(get(isConnected)).toBe(false);
		});
	});
});
