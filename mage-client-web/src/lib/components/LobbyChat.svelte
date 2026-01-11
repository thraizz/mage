<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { sendLobbyMessage, sendWhisper, joinChat } from '$lib/api/chat';
	import { websocketStore } from '$lib/stores/websocket';
	import { CallbackMethod } from '$lib/generated/mage/v1/websocket';
	import type { ChatMessageData } from '$lib/generated/mage/v1/websocket';
	import { getMageClient } from '$lib/grpc/client';
	import { onMount, onDestroy } from 'svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';
	import CircleAlert from '@lucide/svelte/icons/circle-alert';
	import MessageSquare from '@lucide/svelte/icons/message-square';
	import ArrowDown from '@lucide/svelte/icons/arrow-down';
	import Send from '@lucide/svelte/icons/send';
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
				<CircleAlert size={20} aria-hidden="true" />
				<span>{error}</span>
			</div>
		{:else if messages.length === 0}
			<div class="empty-state">
				<MessageSquare size={32} aria-hidden="true" />
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
				<ArrowDown size={18} aria-hidden="true" />
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
					<CircleAlert size={14} aria-hidden="true" />
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
				<Send size={18} aria-hidden="true" />
			{/if}
		</button>
	</form>
</div>

<style>
	.lobby-chat {
		display: flex;
		flex-direction: column;
		height: 100%;
		background-color: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	/* Header */
	.chat-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: var(--space-4);
		border-bottom: 1px solid var(--border-subtle);
		background-color: var(--bg-slate);
	}

	.chat-title {
		margin: 0;
		font-family: var(--font-display);
		font-size: var(--text-base);
		font-weight: var(--weight-semibold);
		color: var(--ci-scroll-parchment);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	.connection-status {
		display: inline-flex;
		align-items: center;
		gap: var(--space-1);
		padding: var(--space-1) var(--space-3);
		border-radius: var(--radius-full);
		font-size: var(--text-xs);
		font-weight: var(--weight-semibold);
	}

	.connection-status.connected {
		background: rgba(46, 204, 113, 0.2);
		color: var(--ci-forest-emerald);
	}

	.connection-status.connecting {
		background: rgba(245, 158, 11, 0.2);
		color: var(--status-warning);
	}

	.connection-status.disconnected {
		background: rgba(255, 77, 77, 0.2);
		color: var(--ci-mountain-ember);
	}

	.status-dot {
		width: 6px;
		height: 6px;
		border-radius: var(--radius-full);
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
		padding: var(--space-4);
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
		gap: var(--space-3);
		height: 100%;
		color: var(--ci-swamp-obsidian);
		text-align: center;
	}

	.error-state {
		color: var(--ci-mountain-ember);
	}

	.empty-state {
		padding: var(--space-8);
	}

	.empty-state :global(svg) {
		color: var(--border-default);
	}

	.empty-state p {
		margin: 0;
		font-weight: var(--weight-medium);
		color: var(--ci-scroll-parchment);
	}

	.empty-subtitle {
		font-size: var(--text-sm);
		color: var(--ci-swamp-obsidian);
		font-style: italic;
	}

	/* Message List */
	.message-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
	}

	.message {
		display: flex;
		align-items: baseline;
		gap: var(--space-2);
		flex-wrap: wrap;
		font-size: var(--text-sm);
		line-height: var(--leading-relaxed);
	}

	.message-time {
		color: var(--text-dim);
		font-size: var(--text-xs);
		flex-shrink: 0;
	}

	.message-username {
		font-weight: var(--weight-semibold);
		color: var(--ci-jace-cloak);
		flex-shrink: 0;
	}

	.message-content {
		color: var(--ci-scroll-parchment);
		word-break: break-word;
		flex: 1;
		min-width: 0;
	}

	/* System Messages */
	.message-system {
		color: var(--ci-swamp-obsidian);
		font-style: italic;
	}

	.message-system .message-content {
		color: var(--ci-swamp-obsidian);
	}

	/* Whisper Messages */
	.message-whisper {
		color: #a855f7;
	}

	.message-whisper-indicator {
		color: #a855f7;
		font-style: italic;
		font-size: var(--text-xs);
	}

	.message-whisper .message-username {
		color: #a855f7;
	}

	.message-whisper .message-content {
		color: #c084fc;
	}

	/* Scroll to Bottom Button */
	.scroll-to-bottom {
		position: absolute;
		bottom: var(--space-4);
		right: var(--space-4);
		width: 2.5rem;
		height: 2.5rem;
		border-radius: var(--radius-full);
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563EB 100%);
		color: var(--ci-scroll-parchment);
		border: none;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
		transition: all var(--transition-base);
	}

	.scroll-to-bottom:hover {
		background: linear-gradient(135deg, #2563EB 0%, #1D4ED8 100%);
		transform: translateY(-2px);
		box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
	}

	/* Input Container */
	.chat-input-container {
		display: flex;
		gap: var(--space-2);
		padding: var(--space-4);
		border-top: 1px solid var(--border-subtle);
		background-color: var(--bg-slate);
	}

	.message-input {
		flex: 1;
		padding: var(--space-3) var(--space-4);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		color: var(--ci-scroll-parchment);
		transition: all var(--transition-fast);
	}

	.message-input::placeholder {
		color: var(--text-ghost);
	}

	.message-input:focus {
		outline: none;
		border-color: var(--ci-jace-cloak);
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
		background: var(--bg-obsidian);
	}

	.message-input:disabled {
		background-color: var(--bg-iron);
		cursor: not-allowed;
		opacity: 0.5;
	}

	.send-button {
		padding: var(--space-3) var(--space-4);
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563EB 100%);
		color: var(--ci-scroll-parchment);
		border: none;
		border-radius: var(--radius-md);
		font-weight: var(--weight-semibold);
		cursor: pointer;
		transition: all var(--transition-base);
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 3rem;
		box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
	}

	.send-button:hover:not(:disabled) {
		background: linear-gradient(135deg, #2563EB 0%, #1D4ED8 100%);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
	}

	.send-button:disabled {
		background: var(--bg-steel);
		box-shadow: none;
		cursor: not-allowed;
	}

	/* Rate Limiting Styles */
	.input-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.message-input.rate-limited {
		border-color: var(--ci-mountain-ember);
		background-color: rgba(255, 77, 77, 0.1);
	}

	.message-input.rate-limited:focus {
		border-color: var(--ci-mountain-ember);
		box-shadow: 0 0 0 3px rgba(255, 77, 77, 0.2);
	}

	.rate-limit-warning {
		display: flex;
		align-items: center;
		gap: var(--space-1);
		font-size: var(--text-xs);
		color: var(--ci-mountain-ember);
		padding: 0 var(--space-1);
	}

	.rate-limit-warning :global(svg) {
		flex-shrink: 0;
	}

	.send-button.rate-limited {
		background: linear-gradient(135deg, var(--ci-mountain-ember) 0%, #DC2626 100%);
		cursor: not-allowed;
	}

	.send-button.rate-limited:hover {
		background: linear-gradient(135deg, var(--ci-mountain-ember) 0%, #DC2626 100%);
	}
</style>
