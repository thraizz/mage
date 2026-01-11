<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { getTableChatId, sendTableMessage, joinChat } from '$lib/api/chat';
	import { websocketStore } from '$lib/stores/websocket';
	import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
	import type { ChatMessageData } from '$lib/generated/mage/v1/websocket';
	import { onMount, onDestroy } from 'svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import MessageSquare from '@lucide/svelte/icons/message-square';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import Send from '@lucide/svelte/icons/send';
	import {
		convertProtoMessageToClientMessage,
		formatMessageTime,
		RateLimiter
	} from '$lib/utils/chat';

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
	let chatId = $state<string | null>(null);
	let connected = $state(false);

	// Rate limiting
	const rateLimiter = new RateLimiter(10, 60000);
	let rateLimitCooldownSeconds = $state(0);
	let rateLimitTimer: ReturnType<typeof setInterval> | null = null;

	// WebSocket unsubscribe function
	let unsubscribeWebSocket: (() => void) | null = null;

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
			const seconds = rateLimiter.getCooldownSeconds();
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
	 * Initialize chat connection
	 */
	async function initializeChat(): Promise<void> {
		loading = true;
		error = null;

		try {
			// Get the chat ID for this table
			chatId = await getTableChatId(tableId);
			console.log('[TableChat] Got chat ID:', chatId);

			// Join the chat room
			await joinChat(chatId);
			console.log('[TableChat] Joined chat');

			// Subscribe to WebSocket chat messages
			unsubscribeWebSocket = websocketStore.on(CallbackMethod.CHATMESSAGE, handleChatMessage);

			connected = true;
			console.log('[TableChat] Connected and listening for messages');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to connect to chat';
			console.error('Failed to initialize chat:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Handle incoming chat message from WebSocket
	 */
	function handleChatMessage(data: unknown): void {
		try {
			const messageData = data as ChatMessageData;

			// Only process messages for this chat
			if (messageData.chatId !== chatId) {
				return;
			}

			const protoMessage = messageData.message;
			if (!protoMessage) {
				return;
			}

			// Convert proto message to client message using utility
			const message = convertProtoMessageToClientMessage(protoMessage);

			// Add message to the list
			messages.push(message);

			// Scroll to bottom after new message
			setTimeout(() => scrollToBottom(), 50);
		} catch (err) {
			console.error('[TableChat] Error handling chat message:', err);
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
		if (rateLimiter.isLimited()) {
			const seconds = rateLimiter.getCooldownSeconds();
			error = `Sending too fast, wait ${seconds} second${seconds !== 1 ? 's' : ''}`;
			setTimeout(() => (error = null), 3000);
			return;
		}

		sending = true;
		error = null;

		try {
			await sendTableMessage(tableId, { content });

			// Record message timestamp for rate limiting
			rateLimiter.recordMessage();
			if (rateLimiter.isLimited()) {
				rateLimitCooldownSeconds = rateLimiter.getCooldownSeconds();
				startCooldownTimer();
			}

			// Clear input - message will be received through WebSocket
			messageInput = '';
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
	 * Initialize on mount
	 */
	onMount(() => {
		initializeChat();
	});

	/**
	 * Cleanup on unmount
	 */
	onDestroy(() => {
		// Clear rate limit timer
		if (rateLimitTimer) {
			clearInterval(rateLimitTimer);
		}

		// Unsubscribe from WebSocket events
		if (unsubscribeWebSocket) {
			unsubscribeWebSocket();
		}

		console.log('[TableChat] Cleaned up');
	});
</script>

<div class="table-chat">
	<div class="chat-header">
		<h3>Table Chat</h3>
		<div class="header-right">
			{#if connected}
				<span class="connection-status connected" title="Connected to chat">
					<span class="status-dot"></span>
					Live
				</span>
			{:else if loading}
				<span class="connection-status connecting">Connecting...</span>
			{:else if error}
				<span class="connection-status disconnected" title="Disconnected">
					<span class="status-dot"></span>
					Offline
				</span>
			{/if}
			<span class="message-count">{messages.length} messages</span>
		</div>
	</div>

	<div class="chat-messages" bind:this={messagesContainer} onscroll={handleScroll}>
		{#if loading}
			<div class="loading-state">
				<LoadingSpinner size="medium" />
				<p>Loading messages...</p>
			</div>
		{:else if messages.length === 0}
			<div class="empty-state">
				<MessageSquare size={48} aria-hidden="true" />
				<p>No messages yet</p>
				<span>Be the first to say something!</span>
			</div>
		{:else}
			<div class="messages-list">
				{#each messages as message (message.id)}
					<div class="message" class:system={message.type === 'system'}>
						<div class="message-header">
							<span class="username">{message.username}</span>
							<span class="timestamp">{formatMessageTime(message.timestamp)}</span>
						</div>
						<div class="message-content">{message.content}</div>
					</div>
				{/each}
			</div>
		{/if}

		{#if !isAtBottom}
			<button class="scroll-to-bottom" onclick={scrollToBottom} title="Scroll to bottom">
				<ChevronDown size={20} aria-hidden="true" />
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
					<Send size={20} aria-hidden="true" />
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
		min-height: fit-content;
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

	.header-right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.connection-status {
		display: inline-flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.25rem 0.625rem;
		border-radius: 1rem;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.connection-status.connected {
		background: #d1fae5;
		color: #065f46;
	}

	.connection-status.connecting {
		background: #fef3c7;
		color: #92400e;
	}

	.connection-status.disconnected {
		background: #fee2e2;
		color: #991b1b;
	}

	.status-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: currentColor;
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
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

	.empty-state :global(svg) {
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
