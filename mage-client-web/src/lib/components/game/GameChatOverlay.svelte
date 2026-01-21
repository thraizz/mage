<script lang="ts">
	/**
	 * Game Chat Overlay - Slide-out chat panel
	 */
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
		open = $bindable(false)
	}: {
		gameId: string;
		open?: boolean;
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
	const RATE_LIMIT_WINDOW_MS = 60000;
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
		if (rateLimitTimer) {
			clearInterval(rateLimitTimer);
		}

		rateLimitTimer = setInterval(() => {
			const seconds = getCooldownSeconds();
			rateLimitCooldownSeconds = seconds;

			if (seconds <= 0) {
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
			type: protoMsg.messageType === 0 ? 'user' : 'system',
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

		if (chatData.chatId !== chatId) {
			return;
		}

		if (chatData.message) {
			const localMessage = convertProtoToLocalMessage(chatData.message, chatData.chatId);
			messages.push(localMessage);
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

		if (!chatId) {
			error = 'Chat not connected';
			setTimeout(() => (error = null), 3000);
			return;
		}

		if (isRateLimited()) {
			const seconds = getCooldownSeconds();
			error = `Sending too fast, wait ${seconds} second${seconds !== 1 ? 's' : ''}`;
			setTimeout(() => (error = null), 3000);
			return;
		}

		sending = true;
		error = null;

		try {
			await sendChatMessage(chatId, content);
			recordMessageSent();
			messageInput = '';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to send message';
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
	 * Handle scroll event
	 */
	function handleScroll(): void {
		if (!messagesContainer) return;

		const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
		const threshold = 100;
		isAtBottom = scrollHeight - scrollTop - clientHeight < threshold;
	}

	/**
	 * Close the overlay
	 */
	function close(): void {
		open = false;
	}

	/**
	 * Handle backdrop click
	 */
	function handleBackdropClick(e: MouseEvent): void {
		if (e.target === e.currentTarget) {
			close();
		}
	}

	/**
	 * Handle escape key
	 */
	function handleKeydown(e: KeyboardEvent): void {
		if (e.key === 'Escape' && open) {
			close();
		}
	}

	/**
	 * Initialize chat connection
	 */
	async function initializeChat(): Promise<void> {
		// Don't initialize if gameId is empty or undefined
		if (!gameId) {
			console.log('[GameChat] Skipping initialization - no gameId');
			return;
		}

		loading = true;
		error = null;

		try {
			const gameChatId = await getGameChatId(gameId);
			chatId = gameChatId;

			unsubscribeWebSocket = websocketStore.on(CallbackMethod.CHATMESSAGE, handleChatMessage);

			await joinChat(gameChatId);

			addGameEvent('Game started. Good luck!');
		} catch (err) {
			console.error('[GameChat] Failed to initialize chat:', err);
			error = err instanceof Error ? err.message : 'Failed to connect to chat';
			addGameEvent('Unable to connect to chat');
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		initializeChat();
	});

	onDestroy(() => {
		if (rateLimitTimer) {
			clearInterval(rateLimitTimer);
		}

		if (unsubscribeWebSocket) {
			unsubscribeWebSocket();
			unsubscribeWebSocket = null;
		}

		if (chatId) {
			leaveChat(chatId).catch((err) => {
				console.error('[GameChat] Failed to leave chat:', err);
			});
		}
	});
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Backdrop -->
{#if open}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div class="overlay-backdrop" onclick={handleBackdropClick}></div>
{/if}

<!-- Slide-out Panel -->
<div class="chat-overlay" class:open>
	<div class="chat-header">
		<div class="header-left">
			<span class="header-icon">💬</span>
			<h3>Game Chat</h3>
			<span class="message-count">{messages.length}</span>
		</div>
		<button class="close-btn" onclick={close} title="Close (Esc)"> ✕ </button>
	</div>

	<div class="chat-messages" bind:this={messagesContainer} onscroll={handleScroll}>
		{#if loading}
			<div class="loading-state">
				<LoadingSpinner size="medium" />
				<p>Connecting to chat...</p>
			</div>
		{:else if messages.length === 0}
			<div class="empty-state">
				<span class="empty-icon">💬</span>
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
				⬇️
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
				title="Send (Enter)"
			>
				{#if sending}
					<LoadingSpinner size="small" color="white" />
				{:else}
					➤
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.overlay-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 89;
		animation: fade-in 0.2s ease;
	}

	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.chat-overlay {
		position: fixed;
		right: 0;
		top: 0;
		bottom: 0;
		width: 380px;
		max-width: 90vw;
		background: #141821;
		border-left: 2px solid #2a3441;
		display: flex;
		flex-direction: column;
		z-index: 90;
		transform: translateX(100%);
		transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: -4px 0 24px rgba(0, 0, 0, 0.5);
	}

	.chat-overlay.open {
		transform: translateX(0);
	}

	.chat-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.25rem;
		border-bottom: 2px solid #2a3441;
		background: #1a1f2e;
		flex-shrink: 0;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.header-icon {
		font-size: 1.25rem;
	}

	.chat-header h3 {
		margin: 0;
		font-size: 1.125rem;
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

	.close-btn {
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid #2a3441;
		border-radius: 6px;
		color: #9ca3af;
		cursor: pointer;
		transition: all 0.2s;
	}

	.close-btn:hover {
		background: #2a3441;
		border-color: #374151;
		color: #fff;
	}

	.chat-messages {
		flex: 1;
		overflow-y: auto;
		padding: 0.75rem;
		position: relative;
	}

	.loading-state,
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 200px;
		gap: 0.75rem;
		color: #6b7280;
	}

	.empty-icon {
		font-size: 2.5rem;
		opacity: 0.4;
	}

	.empty-state p {
		margin: 0;
		font-size: 0.875rem;
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
		background: #1a1f2e;
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
		border-top: 2px solid #2a3441;
		padding: 0.75rem;
		background: #1a1f2e;
		flex-shrink: 0;
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
		border-radius: 6px;
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
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 2.75rem;
		font-size: 1.125rem;
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
	@media (max-width: 480px) {
		.chat-overlay {
			width: 100%;
			max-width: 100%;
		}
	}
</style>
