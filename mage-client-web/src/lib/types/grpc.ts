/**
 * gRPC error codes mapping
 * Based on grpc-web status codes
 */
export enum GrpcStatusCode {
	OK = 0,
	CANCELLED = 1,
	UNKNOWN = 2,
	INVALID_ARGUMENT = 3,
	DEADLINE_EXCEEDED = 4,
	NOT_FOUND = 5,
	ALREADY_EXISTS = 6,
	PERMISSION_DENIED = 7,
	RESOURCE_EXHAUSTED = 8,
	FAILED_PRECONDITION = 9,
	ABORTED = 10,
	OUT_OF_RANGE = 11,
	UNIMPLEMENTED = 12,
	INTERNAL = 13,
	UNAVAILABLE = 14,
	DATA_LOSS = 15,
	UNAUTHENTICATED = 16
}

/**
 * gRPC error interface
 */
export interface GrpcError {
	code: GrpcStatusCode;
	message: string;
	details?: string;
	metadata?: Record<string, string>;
}

/**
 * User-friendly error message
 */
export interface UserError {
	title: string;
	message: string;
	code: GrpcStatusCode;
	canRetry: boolean;
}

/**
 * gRPC client options
 */
export interface GrpcClientOptions {
	/**
	 * Request timeout in milliseconds
	 * @default 30000
	 */
	timeout?: number;

	/**
	 * Whether to automatically retry on certain errors
	 * @default false
	 */
	retry?: boolean;

	/**
	 * Maximum number of retry attempts
	 * @default 3
	 */
	maxRetries?: number;

	/**
	 * Custom metadata to include in all requests
	 */
	metadata?: Record<string, string>;
}

/**
 * gRPC metadata interface
 */
export interface GrpcMetadata {
	authorization?: string;
	[key: string]: string | undefined;
}
