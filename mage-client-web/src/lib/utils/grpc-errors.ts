import { GrpcStatusCode, type GrpcError, type UserError } from '$lib/types/grpc';

/**
 * Error message mappings for gRPC status codes
 */
const ERROR_MESSAGES: Record<
  GrpcStatusCode,
  { title: string; message: string; canRetry: boolean }
> = {
  [GrpcStatusCode.OK]: {
    title: 'Success',
    message: 'Operation completed successfully',
    canRetry: false
  },
  [GrpcStatusCode.CANCELLED]: {
    title: 'Request Cancelled',
    message: 'The request was cancelled',
    canRetry: true
  },
  [GrpcStatusCode.UNKNOWN]: {
    title: 'Unknown Error',
    message: 'An unknown error occurred. Please try again.',
    canRetry: true
  },
  [GrpcStatusCode.INVALID_ARGUMENT]: {
    title: 'Invalid Request',
    message: 'The request contains invalid data',
    canRetry: false
  },
  [GrpcStatusCode.DEADLINE_EXCEEDED]: {
    title: 'Request Timeout',
    message: 'The request took too long to complete. Please try again.',
    canRetry: true
  },
  [GrpcStatusCode.NOT_FOUND]: {
    title: 'Not Found',
    message: 'The requested resource was not found',
    canRetry: false
  },
  [GrpcStatusCode.ALREADY_EXISTS]: {
    title: 'Already Exists',
    message: 'The resource already exists',
    canRetry: false
  },
  [GrpcStatusCode.PERMISSION_DENIED]: {
    title: 'Permission Denied',
    message: 'You do not have permission to perform this action',
    canRetry: false
  },
  [GrpcStatusCode.RESOURCE_EXHAUSTED]: {
    title: 'Rate Limited',
    message: 'Too many requests. Please wait and try again.',
    canRetry: true
  },
  [GrpcStatusCode.FAILED_PRECONDITION]: {
    title: 'Precondition Failed',
    message: 'The operation cannot be performed in the current state',
    canRetry: false
  },
  [GrpcStatusCode.ABORTED]: {
    title: 'Request Aborted',
    message: 'The request was aborted. Please try again.',
    canRetry: true
  },
  [GrpcStatusCode.OUT_OF_RANGE]: {
    title: 'Out of Range',
    message: 'The request is outside the valid range',
    canRetry: false
  },
  [GrpcStatusCode.UNIMPLEMENTED]: {
    title: 'Not Implemented',
    message: 'This feature is not yet implemented',
    canRetry: false
  },
  [GrpcStatusCode.INTERNAL]: {
    title: 'Server Error',
    message: 'An internal server error occurred. Please try again later.',
    canRetry: true
  },
  [GrpcStatusCode.UNAVAILABLE]: {
    title: 'Service Unavailable',
    message: 'The service is temporarily unavailable. Please try again.',
    canRetry: true
  },
  [GrpcStatusCode.DATA_LOSS]: {
    title: 'Data Loss',
    message: 'Unrecoverable data loss or corruption',
    canRetry: false
  },
  [GrpcStatusCode.UNAUTHENTICATED]: {
    title: 'Authentication Required',
    message: 'You must be logged in to perform this action',
    canRetry: false
  }
};

/**
 * Convert a raw error to a GrpcError
 */
export function toGrpcError(error: unknown): GrpcError {
  // Already a GrpcError
  if (isGrpcError(error)) {
    return error;
  }

  // Standard Error object
  if (error instanceof Error) {
    // Check if it has a code property (common in grpc-web errors)
    const errWithCode = error as Error & { code?: number; metadata?: Record<string, string> };
    return {
      code: errWithCode.code ?? GrpcStatusCode.UNKNOWN,
      message: error.message,
      metadata: errWithCode.metadata
    };
  }

  // Network/fetch errors
  if (typeof error === 'object' && error !== null) {
    const obj = error as Record<string, unknown>;
    return {
      code: GrpcStatusCode.UNAVAILABLE,
      message: obj.message?.toString() ?? 'Network error occurred',
      details: JSON.stringify(obj)
    };
  }

  // Unknown error type
  return {
    code: GrpcStatusCode.UNKNOWN,
    message: String(error)
  };
}

/**
 * Type guard for GrpcError
 */
export function isGrpcError(error: unknown): error is GrpcError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'code' in error &&
    'message' in error &&
    typeof (error as GrpcError).code === 'number' &&
    typeof (error as GrpcError).message === 'string'
  );
}

/**
 * Convert a GrpcError to a user-friendly error message
 */
export function toUserError(error: GrpcError | unknown): UserError {
  const grpcError = toGrpcError(error);
  const errorInfo = ERROR_MESSAGES[grpcError.code] ?? ERROR_MESSAGES[GrpcStatusCode.UNKNOWN];

  return {
    title: errorInfo.title,
    message: grpcError.message || errorInfo.message,
    code: grpcError.code,
    canRetry: errorInfo.canRetry
  };
}

/**
 * Check if an error is retryable
 */
export function isRetryableError(error: GrpcError | unknown): boolean {
  const grpcError = toGrpcError(error);
  const errorInfo = ERROR_MESSAGES[grpcError.code];
  return errorInfo?.canRetry ?? false;
}

/**
 * Check if error is authentication related
 * Includes status codes and "session not found" errors
 */
export function isAuthError(error: GrpcError | unknown): boolean {
  const grpcError = toGrpcError(error);
  const message = grpcError.message?.toLowerCase() || '';

  // Check for authentication status codes
  if (
    grpcError.code === GrpcStatusCode.UNAUTHENTICATED ||
    grpcError.code === GrpcStatusCode.PERMISSION_DENIED
  ) {
    return true;
  }

  // Check for "session not found" in error messages (even if status is NOT_FOUND)
  if (
    message.includes('session not found') ||
    message.includes('invalid or expired session') ||
    message.includes('missing session') ||
    message.includes('session expired')
  ) {
    return true;
  }

  return false;
}

/**
 * Log gRPC error to console (in dev mode only)
 */
export function logGrpcError(error: GrpcError | unknown, context?: string): void {
  if (import.meta.env.DEV) {
    const grpcError = toGrpcError(error);
    const prefix = context ? `[gRPC Error - ${context}]` : '[gRPC Error]';

    console.error(prefix, {
      code: grpcError.code,
      codeName: GrpcStatusCode[grpcError.code],
      message: grpcError.message,
      details: grpcError.details,
      metadata: grpcError.metadata
    });
  }
}
