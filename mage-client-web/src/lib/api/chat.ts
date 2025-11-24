import type { ChatMessage as ClientChatMessage, SendMessageRequest } from '$lib/types/chat';
import { getMageClient } from '$lib/grpc/client';
import type { ChatMessage as ProtoChatMessage } from '$lib/generated/mage/v1/models';
import type {
	ChatFindByRoomRequest,
	ChatFindByRoomResponse,
	ChatJoinRequest,
	ChatJoinResponse,
	ChatSendMessageRequest,
	ChatSendMessageResponse
} from '$lib/generated/mage/v1/chat';

/**
 * Convert proto ChatMessage to our client ChatMessage type
 */
function convertProtoChatMessage(msg: ProtoChatMessage, id: string): ClientChatMessage {
	// Determine message type based on message content or color
	let type: 'system' | 'user' | 'whisper' = 'user';
	if (msg.userName.toLowerCase() === 'system' || msg.userName === '') {
		type = 'system';
	}
	// Note: Whispers would need additional logic from server to identify

	return {
		id,
		type,
		username: msg.userName || 'System',
		content: msg.message,
		timestamp: msg.time?.getTime() || Date.now()
	};
}

/**
 * Fetch recent lobby chat messages
 *
 * Note: The server uses streaming/callbacks for real-time chat.
 * This function is a placeholder that shows how to find the chat ID.
 * Real chat messages should come through WebSocket callbacks.
 */
export async function fetchLobbyMessages(limit: number = 50): Promise<ClientChatMessage[]> {
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
		chatId,
		userName: 'Player' // This should come from auth store
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
 * Fetch recent table chat messages
 *
 * Note: The server uses streaming/callbacks for real-time chat.
 * This function is a placeholder that shows how to find the chat ID.
 * Real chat messages should come through WebSocket callbacks.
 */
export async function fetchTableMessages(
	tableId: string,
	limit: number = 50
): Promise<ClientChatMessage[]> {
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

		// Note: In a real implementation, you would need to get the chat ID for this table
		// For now, return empty array - messages will come via WebSocket
		console.log(`Fetching chat for table: ${tableId}`);

		// In a real implementation, you would:
		// 1. Get the table's chat ID from the server
		// 2. Store the chatId and listen for WebSocket callbacks of type CHATMESSAGE
		return [];
	} catch (error) {
		console.error('Failed to fetch table messages:', error);
		return [];
	}
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

	// Get main room ID
	const roomResponse = await client.getMainRoomId();
	if (!roomResponse.roomId) {
		throw new Error('Failed to get main room ID');
	}

	// Note: In a real implementation, you would need to get the chat ID for this table
	// For now, we'll construct a table-specific chat ID (server may handle differently)
	const chatId = `table-${tableId}`;

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
