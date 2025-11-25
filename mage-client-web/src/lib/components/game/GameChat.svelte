<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { onMount, onDestroy } from 'svelte';
	import LoadingSpinner from '../LoadingSpinner.svelte';
	import { websocketStore } from '$lib/stores/websocket';
	import { getGameChatId, joinChat, leaveChat, sendChatMessage } from '$lib/api/chat';
	import { CallbackMethod, type ChatMessageData } from '$lib/generated/mage/v1/websocket';
	import type { ChatMessage as ProtoChatMessage } from '$lib/generated/mage/v1/models';

	// Props
	let {
		gameId,
		collapsed = $bindable(false)
	}: {
		gameId: string;
		collapsed?: boolean;
	} = $props();

	// State
	let messages = $state<ChatMessage[]>([]);
	let messageInput = $state('');
	let loading = $state(false);
	let sending = $state(false);
	let error = $state<string | null>(null);
	let messagesContainer: HTMLDivElement | undefined;
	let isAtBottom = $state(true);
	let chatId = $state<string | null>(null);

	// Rate limiting state
	const RATE_LIMIT_MAX_MESSAGES = 10;
	const RATE_LIMIT_WINDOW_MS = 60000; // 60 seconds
	let messageTimestamps = $state<number[]>([]);
	let rateLimitCooldownSeconds = $state(0);
	let rateLimitTimer: ReturnType<typeof setInterval> | null = null;
	
	// WebSocket unsubscribe function
	let unsubscribeWebSocket: (() => void) | null = null;

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
	 * Add game event message
	 */
	export function addGameEvent(content: string): void {
		const message: ChatMessage = {
			id: `event-${Date.now()}-${Math.random()}`,
			type: 'system',
			username: 'Game',
			content,
			timestamp: Date.now()
		};

		messages.push(message);
		setTimeout(() => scrollToBottom(), 50);
	}

	/**
	 * Convert protobuf ChatMessage to our local ChatMessage type
	 */
	function convertProtoToLocalMessage(protoMsg: ProtoChatMessage, msgChatId: string): ChatMessage {
		return {
			id: `msg-${msgChatId}-${protoMsg.time?.getTime() || Date.now()}-${Math.random()}`,
			type: protoMsg.messageType === 0 ? 'user' : 'system', // 0 = USER in MessageType enum
			username: protoMsg.userName || 'Unknown',
			content: protoMsg.message,
			timestamp: protoMsg.time?.getTime() || Date.now()
		};
	}

	/**
	 * Handle incoming chat message from WebSocket
	 */
	function handleChatMessage(data: unknown): void {
		const chatData = data as ChatMessageData;
		
		// Only process messages for our chat room
		if (chatData.chatId !== chatId) {
			console.log('[GameChat] Ignoring message for different chat:', chatData.chatId, 'vs', chatId);
			return;
		}
		
		if (chatData.message) {
			const localMessage = convertProtoToLocalMessage(chatData.message, chatData.chatId);
			messages.push(localMessage);
			console.log('[GameChat] Received message:', localMessage);
			setTimeout(() => scrollToBottom(), 50);
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

		// Check if chat is connected
		if (!chatId) {
			error = 'Chat not connected';
			setTimeout(() => (error = null), 3000);
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
			// Send message via gRPC API
			await sendChatMessage(chatId, content);
			console.log('[GameChat] Sent message:', content);

			// Record message timestamp for rate limiting
			recordMessageSent();

			// Clear input (message will appear via WebSocket callback)
			messageInput = '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to send message';
			console.error('[GameChat] Failed to send message:', err);
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
	 * Toggle collapsed state
	 */
	function toggleCollapsed(): void {
		collapsed = !collapsed;
	}

	/**
	 * Initialize chat connection
	 */
	async function initializeChat(): Promise<void> {
		loading = true;
		error = null;
		
		try {
			// Get the chat ID for this game
			console.log('[GameChat] Getting chat ID for game:', gameId);
			const gameChatId = await getGameChatId(gameId);
			chatId = gameChatId;
			console.log('[GameChat] Got chat ID:', chatId);
			
			// Subscribe to WebSocket chat messages BEFORE joining
			unsubscribeWebSocket = websocketStore.on(CallbackMethod.CHATMESSAGE, handleChatMessage);
			console.log('[GameChat] Subscribed to CHATMESSAGE events');
			
			// Join the chat room
			await joinChat(gameChatId);
			console.log('[GameChat] Joined chat room:', gameChatId);
			
			// Add welcome message
			addGameEvent('Game started. Good luck!');
		} catch (err) {
			console.error('[GameChat] Failed to initialize chat:', err);
			error = err instanceof Error ? err.message : 'Failed to connect to chat';
			// Still show the component but in an error state
			addGameEvent('Unable to connect to chat');
		} finally {
			loading = false;
		}
	}

	/**
	 * Initialize on mount
	 */
	onMount(() => {
		console.log('[GameChat] Initialized for game:', gameId);
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
			unsubscribeWebSocket = null;
		}
		
		// Leave the chat room
		if (chatId) {
			leaveChat(chatId).catch((err) => {
				console.error('[GameChat] Failed to leave chat:', err);
			});
		}

		console.log('[GameChat] Cleaned up');
	});
</script>

<div class="game-chat" class:collapsed>
	<div class="chat-header">
		<div class="header-left">
			<h3>Game Chat</h3>
			{#if !collapsed}
				<span class="message-count">{messages.length}</span>
			{/if}
		</div>
		<button
			class="toggle-button"
			onclick={toggleCollapsed}
			title={collapsed ? 'Expand chat' : 'Collapse chat'}
		>
			{#if collapsed}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
				>
					<polyline points="18 15 12 9 6 15"></polyline>
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
				>
					<polyline points="6 9 12 15 18 9"></polyline>
				</svg>
			{/if}
		</button>
	</div>

	{#if !collapsed}
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
					>
						<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
					</svg>
					<p>No messages yet</p>
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
					title="Send message (Enter)"
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
						>
							<line x1="22" y1="2" x2="11" y2="13"></line>
							<polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
						</svg>
					{/if}
				</button>
			</form>
		</div>
	{/if}
</div>

<style>
	.game-chat {
		display: flex;
		flex-direction: column;
		height: 100%;
		background: #1a1f2e;
		border: 2px solid #2a3441;
		border-radius: 8px;
		overflow: hidden;
		transition: all 0.3s ease;
	}

	.game-chat.collapsed {
		height: auto;
	}

	.chat-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #2a3441;
		background: #141821;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.chat-header h3 {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: #ffffff;
	}

	.message-count {
		font-size: 0.75rem;
		color: #6b7280;
		background: #0d1117;
		padding: 0.125rem 0.5rem;
		border-radius: 1rem;
	}

	.toggle-button {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid #2a3441;
		border-radius: 4px;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
	}

	.toggle-button:hover {
		background: #2a3441;
		border-color: #374151;
		color: #ffffff;
	}

	.chat-messages {
		flex: 1;
		overflow-y: auto;
		padding: 0.75rem;
		position: relative;
		min-height: 200px;
		max-height: 400px;
	}

	.loading-state,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 0.75rem;
		color: #6b7280;
	}

	.empty-state svg {
		opacity: 0.3;
	}

	.empty-state p {
		margin: 0;
		font-size: 0.875rem;
		font-weight: 500;
		color: #6b7280;
		font-style: italic;
	}

	.messages-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.message {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		padding: 0.625rem;
		background: #0d1117;
		border-radius: 6px;
		border-left: 2px solid #2a3441;
		transition: all 0.2s;
	}

	.message:hover {
		background: #141821;
		border-left-color: #667eea;
	}

	.message.system {
		background: rgba(102, 126, 234, 0.1);
		border-left-color: #667eea;
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
		color: #ffffff;
		font-size: 0.8125rem;
	}

	.message.system .username {
		color: #667eea;
	}

	.timestamp {
		font-size: 0.6875rem;
		color: #6b7280;
	}

	.message-content {
		color: #9ca3af;
		font-size: 0.875rem;
		line-height: 1.5;
		word-wrap: break-word;
	}

	.scroll-to-bottom {
		position: absolute;
		bottom: 0.75rem;
		right: 0.75rem;
		width: 2.25rem;
		height: 2.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 50%;
		cursor: pointer;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
		transition: all 0.2s;
	}

	.scroll-to-bottom:hover {
		background: #5568d3;
		transform: scale(1.05);
	}

	.chat-input-container {
		border-top: 1px solid #2a3441;
		padding: 0.75rem;
		background: #141821;
	}

	.error-message,
	.rate-limit-warning {
		margin-bottom: 0.625rem;
		padding: 0.5rem 0.625rem;
		border-radius: 4px;
		font-size: 0.8125rem;
		font-weight: 500;
	}

	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: #ef4444;
		border: 1px solid rgba(239, 68, 68, 0.3);
	}

	.rate-limit-warning {
		background: rgba(251, 191, 36, 0.1);
		color: #fbbf24;
		border: 1px solid rgba(251, 191, 36, 0.3);
	}

	.chat-input-form {
		display: flex;
		gap: 0.5rem;
	}

	.chat-input {
		flex: 1;
		padding: 0.625rem 0.875rem;
		background: #0d1117;
		border: 1px solid #2a3441;
		border-radius: 4px;
		color: #ffffff;
		font-size: 0.875rem;
		transition: all 0.2s;
	}

	.chat-input::placeholder {
		color: #6b7280;
	}

	.chat-input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.chat-input:disabled {
		background: #1a1f2e;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.send-button {
		padding: 0.625rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 2.75rem;
	}

	.send-button:hover:not(:disabled) {
		background: #5568d3;
	}

	.send-button:disabled {
		background: #4b5563;
		cursor: not-allowed;
		opacity: 0.6;
	}

	/* Scrollbar */
	.chat-messages::-webkit-scrollbar {
		width: 6px;
	}

	.chat-messages::-webkit-scrollbar-track {
		background: #0d1117;
	}

	.chat-messages::-webkit-scrollbar-thumb {
		background: #3a4451;
		border-radius: 3px;
	}

	.chat-messages::-webkit-scrollbar-thumb:hover {
		background: #4a5461;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.chat-messages {
			min-height: 150px;
			max-height: 300px;
		}

		.chat-header {
			padding: 0.625rem 0.875rem;
		}

		.chat-input-container {
			padding: 0.625rem 0.875rem;
		}
	}
</style>
