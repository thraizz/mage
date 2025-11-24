<script lang="ts">
	import type { ChatMessage } from '$lib/types/chat';
	import { fetchLobbyMessages, sendLobbyMessage, sendWhisper } from '$lib/api/chat';
	import { onMount } from 'svelte';
	import LoadingSpinner from './LoadingSpinner.svelte';

	// State
	let messages = $state<ChatMessage[]>([]);
	let messageInput = $state('');
	let loading = $state(false);
	let sending = $state(false);
	let error = $state<string | null>(null);
	let messagesContainer: HTMLDivElement | undefined;
	let isAtBottom = $state(true);

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
	 * Load initial messages
	 */
	async function loadMessages(): Promise<void> {
		loading = true;
		error = null;

		try {
			messages = await fetchLobbyMessages(50);
			// Scroll to bottom after loading
			setTimeout(() => scrollToBottom(), 100);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load messages';
			console.error('Failed to load messages:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Parse whisper command (/w username message)
	 */
	function parseWhisperCommand(content: string): { isWhisper: boolean; username?: string; message?: string } {
		const whisperRegex = /^\/w\s+(\S+)\s+(.+)$/;
		const match = content.match(whisperRegex);

		if (match) {
			return {
				isWhisper: true,
				username: match[1],
				message: match[2]
			};
		}

		return { isWhisper: false };
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

		// Check if it's a whisper command
		const whisperInfo = parseWhisperCommand(content);

		// Validation for whisper
		if (whisperInfo.isWhisper) {
			if (!whisperInfo.username || !whisperInfo.message) {
				error = 'Invalid whisper format. Use: /w username message';
				setTimeout(() => (error = null), 3000);
				return;
			}

			// Cannot whisper to self
			if (whisperInfo.username.toLowerCase() === 'currentuser') {
				error = 'You cannot whisper to yourself';
				setTimeout(() => (error = null), 3000);
				return;
			}
		}

		sending = true;
		error = null;

		try {
			let message: ChatMessage;

			if (whisperInfo.isWhisper && whisperInfo.username && whisperInfo.message) {
				// Send whisper
				message = await sendWhisper(whisperInfo.username, whisperInfo.message);
			} else {
				// Send regular message
				message = await sendLobbyMessage({ content });
			}

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
	 * Load messages on mount
	 */
	onMount(() => {
		loadMessages();
	});
</script>

<div class="lobby-chat">
	<!-- Header -->
	<div class="chat-header">
		<h3 class="chat-title">Lobby Chat</h3>
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
						<span class="message-time">{formatTime(message.timestamp)}</span>
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
		<input
			type="text"
			class="message-input"
			placeholder="Type a message... (use /w username message to whisper)"
			bind:value={messageInput}
			disabled={sending}
		/>
		<button type="submit" class="send-button" disabled={!messageInput.trim() || sending}>
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
