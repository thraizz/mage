import { get } from 'svelte/store';
import { auth } from '$lib/stores/auth';
import type { GrpcClientOptions, GrpcMetadata } from '$lib/types/grpc';
import { GrpcStatusCode } from '$lib/types/grpc';
import {
  toGrpcError,
  toUserError,
  isRetryableError,
  isAuthError,
  logGrpcError
} from '$lib/utils/grpc-errors';
import { toast } from '$lib/stores/toast';
import { handleSessionError } from '$lib/utils/session-error-handler';

/**
 * Default client options
 */
const DEFAULT_OPTIONS: Required<GrpcClientOptions> = {
  timeout: 30000, // 30 seconds
  retry: false,
  maxRetries: 3,
  metadata: {}
};

/**
 * Sleep utility for retry delays
 */
function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Calculate exponential backoff delay
 */
function getRetryDelay(attempt: number): number {
  return Math.min(1000 * Math.pow(2, attempt), 10000); // Max 10 seconds
}

/**
 * Create gRPC metadata with authentication token
 */
export function createGrpcMetadata(customMetadata?: Record<string, string>): GrpcMetadata {
  const metadata: GrpcMetadata = {
    ...customMetadata
  };

  // Add authentication token if available
  const authState = get(auth);
  if (authState.isAuthenticated && authState.token) {
    metadata.authorization = `Bearer ${authState.token}`;
  }

  return metadata;
}

/**
 * Wrapper for gRPC calls with error handling, timeout, and retry logic
 */
export async function grpcCall<TRequest, TResponse>(
  callFn: (request: TRequest, metadata: GrpcMetadata) => Promise<TResponse>,
  request: TRequest,
  context: string,
  options: GrpcClientOptions = {}
): Promise<TResponse> {
  const opts = { ...DEFAULT_OPTIONS, ...options };
  let lastError: unknown;

  // Retry loop
  for (let attempt = 0; attempt <= (opts.retry ? opts.maxRetries : 0); attempt++) {
    try {
      // Create metadata with auth token
      const metadata = createGrpcMetadata(opts.metadata);

      // Create timeout promise
      const timeoutPromise = new Promise<never>((_, reject) => {
        setTimeout(() => {
          reject({
            code: GrpcStatusCode.DEADLINE_EXCEEDED,
            message: `Request timeout after ${opts.timeout}ms`
          });
        }, opts.timeout);
      });

      // Race between actual call and timeout
      const response = await Promise.race([callFn(request, metadata), timeoutPromise]);

      // Success - return response
      return response;
    } catch (error) {
      lastError = error;
      const grpcError = toGrpcError(error);

      // Log error in dev mode
      logGrpcError(grpcError, context);

      // Handle authentication errors - logout user and redirect to login
      if (isAuthError(grpcError)) {
        handleSessionError('Session expired. Please log in again.');
        throw grpcError;
      }

      // Check if we should retry
      if (opts.retry && attempt < opts.maxRetries && isRetryableError(grpcError)) {
        const delay = getRetryDelay(attempt);
        if (import.meta.env.DEV) {
          console.log(`[gRPC Retry] Attempt ${attempt + 1}/${opts.maxRetries}, waiting ${delay}ms`);
        }
        await sleep(delay);
        continue;
      }

      // No more retries - throw error
      throw grpcError;
    }
  }

  // All retries exhausted
  throw lastError;
}

/**
 * Wrapper for gRPC calls with user-friendly error toast
 */
export async function grpcCallWithToast<TRequest, TResponse>(
  callFn: (request: TRequest, metadata: GrpcMetadata) => Promise<TResponse>,
  request: TRequest,
  context: string,
  options: GrpcClientOptions = {}
): Promise<TResponse | null> {
  try {
    return await grpcCall(callFn, request, context, options);
  } catch (error) {
    const userError = toUserError(error);
    toast.error(`${userError.title}: ${userError.message}`);
    return null;
  }
}

/**
 * Create a service client wrapper with automatic auth injection and error handling
 */
export function createServiceClient<TService extends object>(
  serviceImpl: TService,
  defaultOptions: GrpcClientOptions = {}
): TService {
  // Create proxy that wraps all methods
  return new Proxy(serviceImpl, {
    get(target, prop) {
      const original = target[prop as keyof TService];

      // Only wrap functions
      if (typeof original !== 'function') {
        return original;
      }

      // Return wrapped function
      return function (request: unknown, callOptions?: GrpcClientOptions) {
        const opts = { ...defaultOptions, ...callOptions };
        const context = `${String(prop)}`;

        return grpcCall(
          (req, metadata) => {
            // Call original method with metadata
            // Note: This assumes the generated client accepts metadata as second parameter
            // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
            return (original as Function).call(target, req, metadata);
          },
          request,
          context,
          opts
        );
      };
    }
  }) as TService;
}
