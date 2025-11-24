import { describe, it, expect, vi, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { grpcCall, grpcCallWithToast, createGrpcMetadata } from '../service-factory';
import { GrpcStatusCode } from '$lib/types/grpc';
import { auth } from '$lib/stores/auth';
import { toast } from '$lib/stores/toast';

// Mock stores
vi.mock('$lib/stores/auth', () => ({
	auth: {
		subscribe: vi.fn(),
		logout: vi.fn()
	}
}));

vi.mock('$lib/stores/toast', () => ({
	toast: {
		error: vi.fn(),
		success: vi.fn(),
		warning: vi.fn(),
		info: vi.fn()
	}
}));

// Mock svelte/store
vi.mock('svelte/store', () => ({
	get: vi.fn()
}));

describe('gRPC Service Factory', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	describe('createGrpcMetadata', () => {
		it('should create metadata without auth when not authenticated', () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const metadata = createGrpcMetadata();
			expect(metadata).toEqual({});
		});

		it('should include auth token when authenticated', () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: true,
				token: 'test-token-123',
				user: { id: '1', username: 'testuser', email: 'test@example.com' }
			});

			const metadata = createGrpcMetadata();
			expect(metadata).toEqual({
				authorization: 'Bearer test-token-123'
			});
		});

		it('should merge custom metadata', () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: true,
				token: 'test-token-123',
				user: { id: '1', username: 'testuser', email: 'test@example.com' }
			});

			const metadata = createGrpcMetadata({ 'x-custom': 'value' });
			expect(metadata).toEqual({
				authorization: 'Bearer test-token-123',
				'x-custom': 'value'
			});
		});
	});

	describe('grpcCall', () => {
		it('should successfully complete a call', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const mockFn = vi.fn().mockResolvedValue({ success: true });
			const request = { data: 'test' };

			const result = await grpcCall(mockFn, request, 'TestService.testMethod');

			expect(result).toEqual({ success: true });
			expect(mockFn).toHaveBeenCalledWith(request, {});
		});

		it('should timeout after specified duration', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const mockFn = vi.fn().mockImplementation(
				() =>
					new Promise((resolve) => {
						setTimeout(() => resolve({ success: true }), 2000);
					})
			);

			const request = { data: 'test' };

			await expect(
				grpcCall(mockFn, request, 'TestService.testMethod', { timeout: 100 })
			).rejects.toMatchObject({
				code: GrpcStatusCode.DEADLINE_EXCEEDED
			});
		});

		it('should logout on authentication error', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: true,
				token: 'invalid-token',
				user: { id: '1', username: 'testuser', email: 'test@example.com' }
			});

			const mockFn = vi.fn().mockRejectedValue({
				code: GrpcStatusCode.UNAUTHENTICATED,
				message: 'Invalid token'
			});

			const request = { data: 'test' };

			await expect(grpcCall(mockFn, request, 'TestService.testMethod')).rejects.toMatchObject({
				code: GrpcStatusCode.UNAUTHENTICATED
			});

			expect(auth.logout).toHaveBeenCalled();
			expect(toast.error).toHaveBeenCalledWith('Session expired. Please log in again.');
		});

		it('should retry on retryable errors', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			let callCount = 0;
			const mockFn = vi.fn().mockImplementation(() => {
				callCount++;
				if (callCount < 3) {
					return Promise.reject({
						code: GrpcStatusCode.UNAVAILABLE,
						message: 'Service unavailable'
					});
				}
				return Promise.resolve({ success: true });
			});

			const request = { data: 'test' };

			const result = await grpcCall(mockFn, request, 'TestService.testMethod', {
				retry: true,
				maxRetries: 3
			});

			expect(result).toEqual({ success: true });
			expect(mockFn).toHaveBeenCalledTimes(3);
		});

		it('should not retry on non-retryable errors', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const mockFn = vi.fn().mockRejectedValue({
				code: GrpcStatusCode.INVALID_ARGUMENT,
				message: 'Invalid request'
			});

			const request = { data: 'test' };

			await expect(
				grpcCall(mockFn, request, 'TestService.testMethod', {
					retry: true,
					maxRetries: 3
				})
			).rejects.toMatchObject({
				code: GrpcStatusCode.INVALID_ARGUMENT
			});

			expect(mockFn).toHaveBeenCalledTimes(1);
		});
	});

	describe('grpcCallWithToast', () => {
		it('should return result on success', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const mockFn = vi.fn().mockResolvedValue({ success: true });
			const request = { data: 'test' };

			const result = await grpcCallWithToast(mockFn, request, 'TestService.testMethod');

			expect(result).toEqual({ success: true });
			expect(toast.error).not.toHaveBeenCalled();
		});

		it('should show error toast and return null on failure', async () => {
			vi.mocked(get).mockReturnValue({
				isAuthenticated: false,
				token: null,
				user: null
			});

			const mockFn = vi.fn().mockRejectedValue({
				code: GrpcStatusCode.NOT_FOUND,
				message: 'Resource not found'
			});

			const request = { data: 'test' };

			const result = await grpcCallWithToast(mockFn, request, 'TestService.testMethod');

			expect(result).toBeNull();
			expect(toast.error).toHaveBeenCalledWith(expect.stringContaining('Not Found'));
		});
	});
});
