/**
 * Connection status states
 */
export type ConnectionStatus = 'connected' | 'connecting' | 'disconnected' | 'reconnecting';

/**
 * Connection event types
 */
export type ConnectionEventType =
  | 'connected'
  | 'disconnected'
  | 'reconnecting'
  | 'reconnect_failed'
  | 'error';

/**
 * Connection event
 */
export interface ConnectionEvent {
  type: ConnectionEventType;
  timestamp: number;
  error?: Error;
  attempt?: number;
  maxAttempts?: number;
}

/**
 * Connection state interface
 */
export interface ConnectionState {
  status: ConnectionStatus;
  lastConnected: number | null;
  lastDisconnected: number | null;
  reconnectAttempt: number;
  error: string | null;
  latency: number | null; // Ping latency in milliseconds
}

/**
 * Connection options
 */
export interface ConnectionOptions {
  /**
   * Enable automatic reconnection
   * @default true
   */
  autoReconnect?: boolean;

  /**
   * Maximum reconnection attempts before giving up
   * @default 10
   */
  maxReconnectAttempts?: number;

  /**
   * Initial delay before first reconnection attempt (ms)
   * @default 1000
   */
  reconnectDelay?: number;

  /**
   * Maximum delay between reconnection attempts (ms)
   * @default 30000
   */
  maxReconnectDelay?: number;

  /**
   * Enable periodic ping/pong health checks
   * @default true
   */
  enableHealthCheck?: boolean;

  /**
   * Interval between health check pings (ms)
   * @default 30000
   */
  healthCheckInterval?: number;

  /**
   * Timeout for health check response (ms)
   * @default 5000
   */
  healthCheckTimeout?: number;
}

/**
 * Connection event callback
 */
export type ConnectionEventCallback = (event: ConnectionEvent) => void;
