<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Button from '$lib/components/ui/Button.svelte';

	interface ChatMessage {
		id: string;
		username: string;
		content: string;
		timestamp: Date;
		isSystem?: boolean;
	}

	interface Props {
		title?: string;
		messages: ChatMessage[];
		currentUsername: string;
		onsend: (message: string) => void;
		placeholder?: string;
		maxHeight?: string;
	}

	let {
		title = 'Chat',
		messages,
		currentUsername,
		onsend,
		placeholder = 'Type a message...',
		maxHeight = '400px'
	}: Props = $props();

	let inputValue = $state('');
	let messagesContainer: HTMLDivElement;
	let shouldAutoScroll = $state(true);

	function handleSubmit(e: Event) {
		e.preventDefault();
		const trimmed = inputValue.trim();
		if (trimmed) {
			onsend(trimmed);
			inputValue = '';
		}
	}

	function handleScroll() {
		if (messagesContainer) {
			const { scrollTop, scrollHeight, clientHeight } = messagesContainer;
			shouldAutoScroll = scrollHeight - scrollTop - clientHeight < 50;
		}
	}

	function scrollToBottom() {
		if (messagesContainer && shouldAutoScroll) {
			messagesContainer.scrollTop = messagesContainer.scrollHeight;
		}
	}

	function formatTime(date: Date): string {
		return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}

	$effect(() => {
		// Scroll to bottom when messages change
		messages;
		scrollToBottom();
	});

	onMount(() => {
		scrollToBottom();
	});
</script>

<div class="chat-panel">
	{#if title}
		<header class="chat-header">
			<svg class="chat-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
			</svg>
			<span class="chat-title">{title}</span>
		</header>
	{/if}

	<div
		class="messages-container"
		style="max-height: {maxHeight}"
		bind:this={messagesContainer}
		onscroll={handleScroll}
	>
		{#if messages.length === 0}
			<div class="empty-messages">
				<p>No messages yet</p>
			</div>
		{:else}
			{#each messages as message (message.id)}
				<div
					class="message"
					class:own={message.username === currentUsername}
					class:system={message.isSystem}
				>
					{#if message.isSystem}
						<span class="message-system">{message.content}</span>
					{:else}
						<div class="message-header">
							<span class="message-author" class:self={message.username === currentUsername}>
								{message.username}
							</span>
							<span class="message-time">{formatTime(message.timestamp)}</span>
						</div>
						<p class="message-content">{message.content}</p>
					{/if}
				</div>
			{/each}
		{/if}
	</div>

	<form class="chat-input-form" onsubmit={handleSubmit}>
		<input
			type="text"
			class="chat-input"
			{placeholder}
			bind:value={inputValue}
			maxlength="500"
		/>
		<Button type="submit" variant="primary" size="sm" disabled={!inputValue.trim()}>
			Send
		</Button>
	</form>
</div>

<style>
	.chat-panel {
		display: flex;
		flex-direction: column;
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		overflow: hidden;
		height: 100%;
	}

	/* Header */
	.chat-header {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-4);
		border-bottom: 1px solid var(--border-subtle);
	}

	.chat-icon {
		color: var(--accent-gold);
	}

	.chat-title {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
	}

	/* Messages */
	.messages-container {
		flex: 1;
		overflow-y: auto;
		padding: var(--space-3);
	}

	.empty-messages {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100%;
		min-height: 100px;
	}

	.empty-messages p {
		color: var(--text-dim);
		font-size: var(--text-sm);
		margin: 0;
	}

	.message {
		margin-bottom: var(--space-3);
	}

	.message:last-child {
		margin-bottom: 0;
	}

	.message.system {
		text-align: center;
	}

	.message-system {
		font-size: var(--text-xs);
		color: var(--text-dim);
		font-style: italic;
	}

	.message-header {
		display: flex;
		align-items: baseline;
		gap: var(--space-2);
		margin-bottom: var(--space-1);
	}

	.message-author {
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		color: var(--text-muted);
	}

	.message-author.self {
		color: var(--accent-gold);
	}

	.message-time {
		font-size: var(--text-xs);
		color: var(--text-ghost);
	}

	.message-content {
		font-size: var(--text-sm);
		color: var(--text-bright);
		margin: 0;
		word-wrap: break-word;
		line-height: var(--leading-normal);
	}

	/* Input Form */
	.chat-input-form {
		display: flex;
		gap: var(--space-2);
		padding: var(--space-3);
		border-top: 1px solid var(--border-subtle);
		background: var(--bg-slate);
	}

	.chat-input {
		flex: 1;
		padding: var(--space-2) var(--space-3);
		font-family: var(--font-body);
		font-size: var(--text-sm);
		color: var(--text-bright);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		outline: none;
		transition: all var(--transition-fast);
	}

	.chat-input::placeholder {
		color: var(--text-ghost);
	}

	.chat-input:focus {
		border-color: var(--accent-gold);
		box-shadow: 0 0 0 2px var(--accent-gold-glow);
	}
</style>
