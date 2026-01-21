import type { ChatMessage as ClientChatMessage, SendMessageRequest } from '$lib/types/chat';
import { getMageClient } from '$lib/grpc/client';
import type {
	ChatFindByRoomRequest,
	ChatFindByRoomResponse,
	ChatFindByTableRequest,
	ChatFindByTableResponse,
	ChatFindByGameRequest,
	ChatFindByGameResponse,
	ChatJoinRequest,
	ChatJoinResponse,
	ChatSendMessageRequest,
	ChatSendMessageResponse
} from '$lib/generated/mage/v1/chat';

/**
 * Fetch recent lobby chat messages
 *
 * Note: The server uses streaming/callbacks for real-time chat.
 * This function is a placeholder that shows how to find the chat ID.
 * Real chat messages should come through WebSocket callbacks.
 */
export async function fetchLobbyMessages(_limit: number = 50): Promise<ClientChatMessage[]> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	try {
		// Get main room ID
		const roomResponse = await client.getMainRoomId();
		if (!roomResponse.roomId) {
			throw new Error('Failed to get main room ID');
		}

		// Find chat ID for this room
		const chatResponse = await client.call<ChatFindByRoomRequest, ChatFindByRoomResponse>(
			'ChatFindByRoom',
			{
				sessionId,
				roomId: roomResponse.roomId
			}
		);

		// Note: The actual chat messages come through WebSocket callbacks
		// This is just a placeholder to join the chat room
		// Return empty array for now - messages will come via WebSocket
		console.log(`Chat ID for lobby: ${chatResponse.chatId}`);

		// In a real implementation, you would store the chatId and listen for
		// WebSocket callbacks of type CHATMESSAGE to populate messages
		return [];
	} catch (error) {
		console.error('Failed to fetch lobby messages:', error);
		// Return empty array instead of throwing - chat is not critical
		return [];
	}
}

/**
 * Send a message to lobby chat
 */
export async function sendLobbyMessage(request: SendMessageRequest): Promise<ClientChatMessage> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	// Find chat ID for this room
	const chatResponse = await client.call<ChatFindByRoomRequest, ChatFindByRoomResponse>(
		'ChatFindByRoom',
		{
			sessionId,
			roomId: roomResponse.roomId
		}
	);

	// Send the message
	const sendRequest: ChatSendMessageRequest = {
		sessionId,
		chatId: chatResponse.chatId,
		message: request.content
	};

	const response = await client.call<ChatSendMessageRequest, ChatSendMessageResponse>(
		'ChatSendMessage',
		sendRequest
	);

	if (!response.success) {
		throw new Error('Failed to send message');
	}

	// Return a client message representation
	// Note: The actual message will come back through WebSocket callback
	return {
		id: `msg-${Date.now()}`,
		type: 'user',
		username: 'You', // This would come from auth store in real implementation
		content: request.content,
		timestamp: Date.now()
	};
}

/**
 * Send a whisper message to a specific user
 *
 * Note: Whisper functionality may need to be implemented server-side
 * This is a placeholder implementation
 */
export async function sendWhisper(toUsername: string, content: string): Promise<ClientChatMessage> {
	// Note: The server might not have a separate whisper API
	// Whispers might be handled by prefixing the message or using a different chat type
	// This is a placeholder that sends a regular message
	// You may need to implement this differently based on server capabilities

	const message = `/whisper ${toUsername} ${content}`;
	return await sendLobbyMessage({ content: message });
}

/**
 * Join a chat room (lobby, table, tournament, etc.)
 * This should be called when entering a room to start receiving messages
 */
export async function joinChat(chatId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const joinRequest: ChatJoinRequest = {
		sessionId,
		chatId
	};

	const response = await client.call<ChatJoinRequest, ChatJoinResponse>('ChatJoin', joinRequest);

	if (!response.success) {
		throw new Error('Failed to join chat');
	}
}

/**
 * Leave a chat room
 */
export async function leaveChat(chatId: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const leaveRequest = {
		sessionId,
		chatId
	};

	await client.call('ChatLeave', leaveRequest);
}

/**
 * Get the chat ID for a table
 */
export async function getTableChatId(tableId: string): Promise<string> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Use ChatFindByTable to get the chat ID for this table
	const chatResponse = await client.call<ChatFindByTableRequest, ChatFindByTableResponse>(
		'ChatFindByTable',
		{
			sessionId,
			tableId
		}
	);

	if (!chatResponse.chatId) {
		throw new Error('Failed to get chat ID for table');
	}

	return chatResponse.chatId;
}

/**
 * Fetch recent table chat messages
 *
 * Note: The server uses streaming/callbacks for real-time chat.
 * This function is a placeholder that shows how to find the chat ID.
 * Real chat messages should come through WebSocket callbacks.
 */
export async function fetchTableMessages(
	_tableId: string,
	_limit: number = 50
): Promise<ClientChatMessage[]> {
	// Note: Real chat messages should come through WebSocket callbacks
	// This just returns empty array - the TableChat component should use WebSocket
	console.log(`fetchTableMessages is a placeholder - use WebSocket for live chat`);
	return [];
}

/**
 * Send a message to table chat
 */
export async function sendTableMessage(
	tableId: string,
	request: SendMessageRequest
): Promise<ClientChatMessage> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get the chat ID for this table
	const chatId = await getTableChatId(tableId);

	// Send the message
	const sendRequest: ChatSendMessageRequest = {
		sessionId,
		chatId: chatId,
		message: request.content
	};

	const response = await client.call<ChatSendMessageRequest, ChatSendMessageResponse>(
		'ChatSendMessage',
		sendRequest
	);

	if (!response.success) {
		throw new Error('Failed to send message');
	}

	// Return a client message representation
	// Note: The actual message will come back through WebSocket callback
	return {
		id: `msg-${Date.now()}`,
		type: 'user',
		username: 'You', // This would come from auth store in real implementation
		content: request.content,
		timestamp: Date.now()
	};
}

/**
 * Get the chat ID for a game
 */
export async function getGameChatId(gameId: string): Promise<string> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Use ChatFindByGame to get the chat ID for this game
	const chatResponse = await client.call<ChatFindByGameRequest, ChatFindByGameResponse>(
		'ChatFindByGame',
		{
			sessionId,
			gameId
		}
	);

	if (!chatResponse.chatId) {
		throw new Error('Failed to get chat ID for game');
	}

	return chatResponse.chatId;
}

/**
 * Send a message to game chat
 */
export async function sendGameMessage(gameId: string, content: string): Promise<ClientChatMessage> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	// Get the chat ID for this game
	const chatId = await getGameChatId(gameId);

	// Send the message
	const sendRequest: ChatSendMessageRequest = {
		sessionId,
		chatId: chatId,
		message: content
	};

	const response = await client.call<ChatSendMessageRequest, ChatSendMessageResponse>(
		'ChatSendMessage',
		sendRequest
	);

	if (!response.success) {
		throw new Error('Failed to send message');
	}

	// Return a client message representation
	// Note: The actual message will come back through WebSocket callback
	return {
		id: `msg-${Date.now()}`,
		type: 'user',
		username: 'You',
		content: content,
		timestamp: Date.now()
	};
}

/**
 * Send a chat message directly to a chat ID (for game/table/room chat)
 * This is a lower-level function used when you already have the chat ID
 */
export async function sendChatMessage(chatId: string, content: string): Promise<void> {
	const client = getMageClient();
	const sessionId = await client.ensureSessionId();

	if (!sessionId) {
		throw new Error('No active session - please login first');
	}

	const sendRequest: ChatSendMessageRequest = {
		sessionId,
		chatId,
		message: content
	};

	const response = await client.call<ChatSendMessageRequest, ChatSendMessageResponse>(
		'ChatSendMessage',
		sendRequest
	);

	if (!response.success) {
		throw new Error('Failed to send message');
	}
}
