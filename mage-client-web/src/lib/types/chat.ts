/**
 * Chat type definitions for messaging
 */

export type MessageType = 'user' | 'system' | 'whisper';

export interface ChatMessage {
	id: string;
	type: MessageType;
	username: string;
	content: string;
	timestamp: number;
	// For whisper messages
	toUsername?: string;
	fromUsername?: string;
}

export interface ChatState {
	messages: ChatMessage[];
	isLoading: boolean;
	error: string | null;
}

export interface SendMessageRequest {
	content: string;
	roomId?: string;
}
