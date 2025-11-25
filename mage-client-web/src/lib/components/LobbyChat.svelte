<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { sendLobbyMessage, sendWhisper, joinChat } from '$lib/api/chat';
	import { websocketStore } from '$lib/stores/websocket';
	import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
	import type { ChatMessageData } from '$lib/generated/mage/v1/websocket';
	import { getMageClient } from '$lib/grpc/client';
	import { onMount, onDestroy } from 'svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import {
		convertProtoMessageToClientMessage,
		formatMessageTime,
		parseWhisperCommand,
		validateWhisperCommand,
		RateLimiter
	} from '$lib/utils/chat';

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

	// WebSocket unsubscribe function
	let unsubscribeWebSocket: (() => void) | null = null;

	// Rate limiting
	const rateLimiter = new RateLimiter(10, 60000);
	let rateLimitCooldownSeconds = $state(0);
	let rateLimitTimer: ReturnType<typeof setInterval> | null = null;

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
			const client = getMageClient();

			// Get main room ID
			const roomResponse = await client.getMainRoomId();
			if (!roomResponse.roomId) {
				throw new Error('Failed to get main room ID');
			}

			// Find chat ID for the lobby (main room)
			const sessionId = await client.ensureSessionId();
			if (!sessionId) {
				throw new Error('No active session - please login first');
			}

			const chatResponse = await client.call<
				{ sessionId: string; roomId: string },
				{ chatId?: string }
			>('ChatFindByRoom', {
				sessionId,
				roomId: roomResponse.roomId
			});

			if (!chatResponse.chatId) {
				throw new Error('Failed to get chat ID for lobby');
			}

			chatId = chatResponse.chatId;
			console.log('[LobbyChat] Got chat ID:', chatId);

			// Join the chat room
			await joinChat(chatResponse.chatId);
			console.log('[LobbyChat] Joined chat');

			// Subscribe to WebSocket chat messages
			unsubscribeWebSocket = websocketStore.on(CallbackMethod.CHATMESSAGE, handleChatMessage);

			connected = true;
			console.log('[LobbyChat] Connected and listening for messages');
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
		console.log('[LobbyChat] Received WebSocket message:', data);

		try {
			const messageData = data as ChatMessageData;
			console.log('[LobbyChat] Parsed as ChatMessageData:', messageData);
			console.log('[LobbyChat] Current chatId:', chatId, 'Message chatId:', messageData.chatId);

			// Only process messages for this chat
			if (messageData.chatId !== chatId) {
				console.log('[LobbyChat] Ignoring message for different chat');
				return;
			}

			const protoMessage = messageData.message;
			if (!protoMessage) {
				console.log('[LobbyChat] No message in data');
				return;
			}

			console.log('[LobbyChat] Proto message:', protoMessage);

			// Convert proto message to client message using utility
			const message = convertProtoMessageToClientMessage(protoMessage);

			console.log('[LobbyChat] Adding message to list:', message);

			// Add message to the list
			messages.push(message);
			messages = [...messages]; // Force reactivity

			// Scroll to bottom after new message
			setTimeout(() => scrollToBottom(), 50);
		} catch (err) {
			console.error('[LobbyChat] Error handling chat message:', err);
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

		// Check if it's a whisper command
		const whisperInfo = parseWhisperCommand(content);

		// Validate whisper command
		const whisperError = validateWhisperCommand(whisperInfo);
		if (whisperError) {
			error = whisperError;
			setTimeout(() => (error = null), 3000);
			return;
		}

		sending = true;
		error = null;

		console.log('[LobbyChat] Sending message:', content);

		try {
			if (whisperInfo.isWhisper && whisperInfo.username && whisperInfo.message) {
				// Send whisper
				console.log('[LobbyChat] Sending whisper to:', whisperInfo.username);
				await sendWhisper(whisperInfo.username, whisperInfo.message);
			} else {
				// Send regular message
				console.log('[LobbyChat] Sending regular message');
				await sendLobbyMessage({ content });
			}

			console.log('[LobbyChat] Message sent successfully');

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
			console.error('[LobbyChat] Failed to send message:', err);
			setTimeout(() => (error = null), 5000);
		} finally {
			sending = false;
		}
	}

	/**
	 * Scroll to bottom of messages
	 */
	function scrollToBottom(): void {
		if (messagesContainer) {
			messagesContainer.scrollTop = messagesContainer.scrollHeight;
		}
	}

	/**
	 * Check if user is at bottom of messages
	 */
	function handleScroll(): void {
		if (messagesContainer) {
			const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
			isAtBottom = scrollHeight - scrollTop - clientHeight < 10;
		}
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
		// Unsubscribe from WebSocket events
		if (unsubscribeWebSocket) {
			unsubscribeWebSocket();
		}

		console.log('[LobbyChat] Cleaned up');
	});
</script>

<div class="lobby-chat">
	<!-- Header -->
	<div class="chat-header">
		<h3 class="chat-title">Lobby Chat</h3>
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
		</div>
	</div>

	<!-- Messages Container -->
	<div class="messages-container" bind:this={messagesContainer} onscroll={handleScroll}>
		{#if loading}
			<div class="loading-state">
				<LoadingSpinner size="small" />
				<span>Loading messages...</span>
			</div>
		{:else if error}
			<div class="error-state">
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
					<circle cx="12" cy="12" r="10"></circle>
					<line x1="12" y1="8" x2="12" y2="12"></line>
					<line x1="12" y1="16" x2="12.01" y2="16"></line>
				</svg>
				<span>{error}</span>
			</div>
		{:else if messages.length === 0}
			<div class="empty-state">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="32"
					height="32"
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
				<span class="empty-subtitle">Be the first to say something!</span>
			</div>
		{:else}
			<!-- Message List -->
			<div class="message-list">
				{#each messages as message (message.id)}
					<div class="message message-{message.type}">
						<span class="message-time">{formatMessageTime(message.timestamp)}</span>
						{#if message.type === 'system'}
							<span class="message-content">{message.content}</span>
						{:else if message.type === 'whisper'}
							<span class="message-whisper-indicator">(whisper)</span>
							{#if message.toUsername}
								<span class="message-username">To {message.toUsername}:</span>
							{:else}
								<span class="message-username">From {message.fromUsername}:</span>
							{/if}
							<span class="message-content">{message.content}</span>
						{:else}
							<span class="message-username">{message.username}:</span>
							<span class="message-content">{message.content}</span>
						{/if}
					</div>
				{/each}
			</div>
		{/if}

		<!-- Scroll to Bottom Button -->
		{#if !isAtBottom && messages.length > 0}
			<button class="scroll-to-bottom" onclick={scrollToBottom} title="Scroll to bottom">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="18"
					height="18"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M12 5v14M19 12l-7 7-7-7"></path>
				</svg>
			</button>
		{/if}
	</div>

	<!-- Input Container -->
	<form class="chat-input-container" onsubmit={handleSendMessage}>
		<div class="input-wrapper">
			<input
				type="text"
				class="message-input"
				class:rate-limited={rateLimitCooldownSeconds > 0}
				placeholder="Type a message... (use /w username message to whisper)"
				bind:value={messageInput}
				disabled={sending || rateLimitCooldownSeconds > 0}
			/>
			{#if rateLimitCooldownSeconds > 0}
				<div class="rate-limit-warning">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="14"
						height="14"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="12" y1="8" x2="12" y2="12"></line>
						<line x1="12" y1="16" x2="12.01" y2="16"></line>
					</svg>
					<span
						>Sending too fast, wait {rateLimitCooldownSeconds} second{rateLimitCooldownSeconds !== 1
							? 's'
							: ''}</span
					>
				</div>
			{/if}
		</div>
		<button
			type="submit"
			class="send-button"
			class:rate-limited={rateLimitCooldownSeconds > 0}
			disabled={!messageInput.trim() || sending || rateLimitCooldownSeconds > 0}
		>
			{#if sending}
				<LoadingSpinner size="small" color="white" />
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="18"
					height="18"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z"></path>
				</svg>
			{/if}
		</button>
	</form>
</div>

<style>
	.lobby-chat {
		display: flex;
		flex-direction: column;
		height: 100%;
		background-color: white;
		border-radius: 0.5rem;
		overflow: hidden;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	/* Header */
	.chat-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-bottom: 1px solid #e5e7eb;
		background-color: #f9fafb;
	}

	.chat-title {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
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

	/* Messages Container */
	.messages-container {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
		position: relative;
	}

	/* Loading/Error/Empty States */
	.loading-state,
	.error-state,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		height: 100%;
		color: #6b7280;
		text-align: center;
	}

	.error-state {
		color: #dc2626;
	}

	.empty-state {
		padding: 2rem;
	}

	.empty-state svg {
		color: #d1d5db;
	}

	.empty-state p {
		margin: 0;
		font-weight: 500;
		color: #374151;
	}

	.empty-subtitle {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	/* Message List */
	.message-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.message {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		flex-wrap: wrap;
		font-size: 0.875rem;
		line-height: 1.5;
	}

	.message-time {
		color: #9ca3af;
		font-size: 0.75rem;
		flex-shrink: 0;
	}

	.message-username {
		font-weight: 600;
		color: #667eea;
		flex-shrink: 0;
	}

	.message-content {
		color: #374151;
		word-break: break-word;
		flex: 1;
		min-width: 0;
	}

	/* System Messages */
	.message-system {
		color: #6b7280;
		font-style: italic;
	}

	.message-system .message-content {
		color: #6b7280;
	}

	/* Whisper Messages */
	.message-whisper {
		color: #a855f7;
	}

	.message-whisper-indicator {
		color: #a855f7;
		font-style: italic;
		font-size: 0.75rem;
	}

	.message-whisper .message-username {
		color: #a855f7;
	}

	.message-whisper .message-content {
		color: #9333ea;
	}

	/* Scroll to Bottom Button */
	.scroll-to-bottom {
		position: absolute;
		bottom: 1rem;
		right: 1rem;
		width: 2.5rem;
		height: 2.5rem;
		border-radius: 50%;
		background-color: #667eea;
		color: white;
		border: none;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		transition: all 0.2s;
	}

	.scroll-to-bottom:hover {
		background-color: #5568d3;
		transform: translateY(-2px);
		box-shadow: 0 6px 8px -1px rgba(0, 0, 0, 0.15);
	}

	/* Input Container */
	.chat-input-container {
		display: flex;
		gap: 0.5rem;
		padding: 1rem;
		border-top: 1px solid #e5e7eb;
		background-color: white;
	}

	.message-input {
		flex: 1;
		padding: 0.625rem 0.875rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		transition: all 0.2s;
	}

	.message-input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.message-input:disabled {
		background-color: #f3f4f6;
		cursor: not-allowed;
		opacity: 0.6;
	}

	.send-button {
		padding: 0.625rem 1rem;
		background-color: #667eea;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 3rem;
	}

	.send-button:hover:not(:disabled) {
		background-color: #5568d3;
	}

	.send-button:disabled {
		background-color: #9ca3af;
		cursor: not-allowed;
	}

	/* Rate Limiting Styles */
	.input-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}

	.message-input.rate-limited {
		border-color: #ef4444;
		background-color: #fef2f2;
	}

	.message-input.rate-limited:focus {
		border-color: #ef4444;
		box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
	}

	.rate-limit-warning {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		font-size: 0.75rem;
		color: #dc2626;
		padding: 0 0.25rem;
	}

	.rate-limit-warning svg {
		flex-shrink: 0;
	}

	.send-button.rate-limited {
		background-color: #ef4444;
		cursor: not-allowed;
	}

	.send-button.rate-limited:hover {
		background-color: #ef4444;
	}

	/* Scrollbar Styling */
	.messages-container::-webkit-scrollbar {
		width: 8px;
	}

	.messages-container::-webkit-scrollbar-track {
		background-color: #f3f4f6;
		border-radius: 4px;
	}

	.messages-container::-webkit-scrollbar-thumb {
		background-color: #d1d5db;
		border-radius: 4px;
	}

	.messages-container::-webkit-scrollbar-thumb:hover {
		background-color: #9ca3af;
	}
</style>
