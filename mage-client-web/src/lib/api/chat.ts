import type { ChatMessage, SendMessageRequest } from '$lib/types/chat';

/**
 * Mock chat messages for development
 */
const MOCK_MESSAGES: ChatMessage[] = [
	{
		id: 'msg1',
		type: 'system',
		username: 'System',
		content: 'Welcome to the Mage lobby!',
		timestamp: Date.now() - 600000
	},
	{
		id: 'msg2',
		type: 'user',
		username: 'alice',
		content: 'Anyone up for Commander?',
		timestamp: Date.now() - 540000
	},
	{
		id: 'msg3',
		type: 'user',
		username: 'bob',
		content: 'I am! Let me grab my deck',
		timestamp: Date.now() - 480000
	},
	{
		id: 'msg4',
		type: 'system',
		username: 'System',
		content: 'alice created table "Commander Night"',
		timestamp: Date.now() - 420000
	},
	{
		id: 'msg5',
		type: 'user',
		username: 'charlie',
		content: 'What format are you playing?',
		timestamp: Date.now() - 360000
	},
	{
		id: 'msg6',
		type: 'user',
		username: 'alice',
		content: 'Commander! Join my table',
		timestamp: Date.now() - 300000
	},
	{
		id: 'msg7',
		type: 'user',
		username: 'dave',
		content: 'Looking for Standard 1v1',
		timestamp: Date.now() - 240000
	},
	{
		id: 'msg8',
		type: 'user',
		username: 'eve',
		content: "I'll play Standard with you",
		timestamp: Date.now() - 180000
	},
	{
		id: 'msg9',
		type: 'system',
		username: 'System',
		content: 'dave created table "Standard Ranked"',
		timestamp: Date.now() - 120000
	},
	{
		id: 'msg10',
		type: 'user',
		username: 'frank',
		content: 'Anyone have a spare Pauper deck?',
		timestamp: Date.now() - 60000
	}
];

// Track messages for this session
let sessionMessages = [...MOCK_MESSAGES];

/**
 * Fetch recent lobby chat messages
 */
export async function fetchLobbyMessages(limit: number = 50): Promise<ChatMessage[]> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 300));

	// In production, this would be:
	// const response = await grpcCall(chatService.getLobbyMessages, { limit }, 'ChatService.getLobbyMessages');
	// return response.messages;

	return sessionMessages.slice(-limit);
}

/**
 * Send a message to lobby chat
 */
export async function sendLobbyMessage(request: SendMessageRequest): Promise<ChatMessage> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 200));

	// In production, this would be:
	// const response = await grpcCall(chatService.sendMessage, request, 'ChatService.sendMessage');
	// return response.message;

	const newMessage: ChatMessage = {
		id: `msg-${Date.now()}`,
		type: 'user',
		username: 'currentuser', // This would come from auth store
		content: request.content,
		timestamp: Date.now()
	};

	// Add to session messages
	sessionMessages.push(newMessage);

	return newMessage;
}

/**
 * Send a whisper message to a specific user
 */
export async function sendWhisper(
	toUsername: string,
	content: string
): Promise<ChatMessage> {
	// Simulate network delay
	await new Promise((resolve) => setTimeout(resolve, 200));

	// In production, this would be:
	// const response = await grpcCall(chatService.sendWhisper, { toUsername, content }, 'ChatService.sendWhisper');
	// return response.message;

	const newMessage: ChatMessage = {
		id: `whisper-${Date.now()}`,
		type: 'whisper',
		username: 'currentuser',
		content,
		timestamp: Date.now(),
		toUsername,
		fromUsername: 'currentuser'
	};

	// Add to session messages
	sessionMessages.push(newMessage);

	return newMessage;
}
