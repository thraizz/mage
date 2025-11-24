<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { fetchTableMessages, sendTableMessage } from '$lib/api/chat';
	import { onMount, onDestroy } from 'svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';

	// Props
	let { tableId }: { tableId: string } = $props();

	// State
	let messages = $state<ChatMessage[]>([]);
	let messageInput = $state('');
	let loading = $state(false);
	let sending = $state(false);
	let error = $state<string | null>(null);
	let messagesContainer: HTMLDivElement | undefined;
	let isAtBottom = $state(true);

	// Rate limiting state
	const RATE_LIMIT_MAX_MESSAGES = 10;
	const RATE_LIMIT_WINDOW_MS = 60000; // 60 seconds
	let messageTimestamps = $state<number[]>([]);
	let rateLimitCooldownSeconds = $state(0);
	let rateLimitTimer: ReturnType<typeof setInterval> | null = null;

	/**
	 * Format timestamp as HH:MM
	 */
	function formatTime(timestamp: number): string {
		const date = new Date(timestamp);
		return date.toLocaleTimeString('en-US', {
			hour: '2-digit',
			minute: '2-digit',
			hour12: false
		});
	}

	/**
	 * Check if user is rate limited
	 */
	function isRateLimited(): boolean {
		const now = Date.now();
		// Remove timestamps older than the window
		messageTimestamps = messageTimestamps.filter((ts) => now - ts < RATE_LIMIT_WINDOW_MS);
		return messageTimestamps.length >= RATE_LIMIT_MAX_MESSAGES;
	}

	/**
	 * Calculate cooldown time remaining
	 */
	function getCooldownSeconds(): number {
		if (messageTimestamps.length < RATE_LIMIT_MAX_MESSAGES) {
			return 0;
		}
		const now = Date.now();
		const oldestTimestamp = messageTimestamps[0];
		const timeElapsed = now - oldestTimestamp;
		const timeRemaining = RATE_LIMIT_WINDOW_MS - timeElapsed;
		return Math.ceil(timeRemaining / 1000);
	}

	/**
	 * Start cooldown timer
	 */
	function startCooldownTimer(): void {
		// Clear existing timer
		if (rateLimitTimer) {
			clearInterval(rateLimitTimer);
		}

		// Update cooldown display every second
		rateLimitTimer = setInterval(() => {
			const seconds = getCooldownSeconds();
			rateLimitCooldownSeconds = seconds;

			if (seconds <= 0) {
				// Cooldown expired
				if (rateLimitTimer) {
					clearInterval(rateLimitTimer);
					rateLimitTimer = null;
				}
			}
		}, 1000);
	}

	/**
	 * Record message timestamp
	 */
	function recordMessageSent(): void {
		const now = Date.now();
		messageTimestamps.push(now);

		// Check if now rate limited
		if (isRateLimited()) {
			rateLimitCooldownSeconds = getCooldownSeconds();
			startCooldownTimer();
		}
	}

	/**
	 * Load initial messages
	 */
	async function loadMessages(): Promise<void> {
		loading = true;
		error = null;

		try {
			messages = await fetchTableMessages(tableId, 50);
			// Scroll to bottom after loading
			setTimeout(() => scrollToBottom(), 100);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load messages';
			console.error('Failed to load table messages:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Send a message
	 */
	async function handleSendMessage(e: Event): Promise<void> {
		e.preventDefault();

		const content = messageInput.trim();
		if (!content || sending) {
			return;
		}

		// Check rate limiting
		if (isRateLimited()) {
			const seconds = getCooldownSeconds();
			error = `Sending too fast, wait ${seconds} second${seconds !== 1 ? 's' : ''}`;
			setTimeout(() => (error = null), 3000);
			return;
		}

		sending = true;
		error = null;

		try {
			const message = await sendTableMessage(tableId, { content });

			// Record message timestamp for rate limiting
			recordMessageSent();

			messages.push(message);
			messageInput = '';

			// Scroll to bottom after sending
			setTimeout(() => scrollToBottom(), 50);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to send message';
			console.error('Failed to send message:', err);
			setTimeout(() => (error = null), 5000);
		} finally {
			sending = false;
		}
	}

	/**
	 * Scroll to bottom of messages
	 */
	function scrollToBottom(): void {
		if (messagesContainer && isAtBottom) {
			messagesContainer.scrollTop = messagesContainer.scrollHeight;
		}
	}

	/**
	 * Handle scroll event to detect if user is at bottom
	 */
	function handleScroll(): void {
		if (!messagesContainer) return;

		const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
		const threshold = 100;
		isAtBottom = scrollHeight - scrollTop - clientHeight < threshold;
	}

	/**
	 * Focus input on mount
	 */
	onMount(() => {
		loadMessages();
	});

	/**
	 * Cleanup on unmount
	 */
	onDestroy(() => {
		if (rateLimitTimer) {
			clearInterval(rateLimitTimer);
		}
	});
</script>

<div class="table-chat">
	<div class="chat-header">
		<h3>Table Chat</h3>
		<span class="message-count">{messages.length} messages</span>
	</div>

	<div class="chat-messages" bind:this={messagesContainer} onscroll={handleScroll}>
		{#if loading}
			<div class="loading-state">
				<LoadingSpinner size="medium" />
				<p>Loading messages...</p>
			</div>
		{:else if messages.length === 0}
			<div class="empty-state">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="48"
					height="48"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
				</svg>
				<p>No messages yet</p>
				<span>Be the first to say something!</span>
			</div>
		{:else}
			<div class="messages-list">
				{#each messages as message (message.id)}
					<div class="message" class:system={message.type === 'system'}>
						<div class="message-header">
							<span class="username">{message.username}</span>
							<span class="timestamp">{formatTime(message.timestamp)}</span>
						</div>
						<div class="message-content">{message.content}</div>
					</div>
				{/each}
			</div>
		{/if}

		{#if !isAtBottom}
			<button class="scroll-to-bottom" onclick={scrollToBottom} title="Scroll to bottom">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<polyline points="6 9 12 15 18 9"></polyline>
				</svg>
			</button>
		{/if}
	</div>

	<div class="chat-input-container">
		{#if error}
			<div class="error-message">{error}</div>
		{/if}

		{#if rateLimitCooldownSeconds > 0}
			<div class="rate-limit-warning">
				Too many messages! Wait {rateLimitCooldownSeconds}s
			</div>
		{/if}

		<form onsubmit={handleSendMessage} class="chat-input-form">
			<input
				type="text"
				bind:value={messageInput}
				placeholder="Type a message..."
				class="chat-input"
				disabled={sending || rateLimitCooldownSeconds > 0}
				maxlength="500"
			/>
			<button
				type="submit"
				class="send-button"
				disabled={sending || !messageInput.trim() || rateLimitCooldownSeconds > 0}
				title="Send message"
			>
				{#if sending}
					<LoadingSpinner size="small" color="white" />
				{:else}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<line x1="22" y1="2" x2="11" y2="13"></line>
						<polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
					</svg>
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.table-chat {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 0.75rem;
		overflow: hidden;
	}

	.chat-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.25rem;
		border-bottom: 1px solid #e5e7eb;
		background: #f9fafb;
	}

	.chat-header h3 {
		margin: 0;
		font-size: 1.125rem;
		font-weight: 700;
		color: #111827;
	}

	.message-count {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.chat-messages {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		position: relative;
		min-height: 300px;
	}

	.loading-state,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 0.75rem;
		color: #9ca3af;
	}

	.empty-state svg {
		opacity: 0.5;
	}

	.empty-state p {
		margin: 0;
		font-size: 1.125rem;
		font-weight: 600;
		color: #6b7280;
	}

	.empty-state span {
		font-size: 0.875rem;
	}

	.messages-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.message {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding: 0.75rem;
		background: #f9fafb;
		border-radius: 0.5rem;
		transition: background 0.2s;
	}

	.message:hover {
		background: #f3f4f6;
	}

	.message.system {
		background: #eff6ff;
		font-style: italic;
	}

	.message-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
	}

	.username {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
	}

	.timestamp {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.message-content {
		color: #111827;
		font-size: 0.9375rem;
		line-height: 1.5;
		word-wrap: break-word;
	}

	.scroll-to-bottom {
		position: absolute;
		bottom: 1rem;
		right: 1rem;
		width: 2.5rem;
		height: 2.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 50%;
		cursor: pointer;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		transition: all 0.2s;
	}

	.scroll-to-bottom:hover {
		background: #5568d3;
		transform: scale(1.05);
	}

	.chat-input-container {
		border-top: 1px solid #e5e7eb;
		padding: 1rem;
		background: white;
	}

	.error-message,
	.rate-limit-warning {
		margin-bottom: 0.75rem;
		padding: 0.625rem;
		border-radius: 0.375rem;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.error-message {
		background: #fef2f2;
		color: #dc2626;
		border: 1px solid #fecaca;
	}

	.rate-limit-warning {
		background: #fef3c7;
		color: #92400e;
		border: 1px solid #fde68a;
	}

	.chat-input-form {
		display: flex;
		gap: 0.75rem;
	}

	.chat-input {
		flex: 1;
		padding: 0.75rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.9375rem;
		transition: all 0.2s;
	}

	.chat-input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.chat-input:disabled {
		background: #f3f4f6;
		cursor: not-allowed;
	}

	.send-button {
		padding: 0.75rem 1.25rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 0.5rem;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 3rem;
	}

	.send-button:hover:not(:disabled) {
		background: #5568d3;
	}

	.send-button:disabled {
		background: #9ca3af;
		cursor: not-allowed;
		opacity: 0.6;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.table-chat {
			border-radius: 0.5rem;
		}

		.chat-messages {
			min-height: 200px;
		}

		.chat-header {
			padding: 0.875rem 1rem;
		}

		.chat-input-container {
			padding: 0.875rem 1rem;
		}
	}
</style>
