/**
 * Chat utility functions
 * Shared logic for handling chat messages across different chat components
 */

import type { ChatMessage as ProtoMessage } from '$lib/generated/mage/v1/models';
import type { ChatMessage } from '$lib/types/chat';

/**
 * Extract timestamp from protobuf message
 * Handles multiple formats: Date, number, or protobuf Timestamp
 */
export function extractTimestamp(protoTime: Date | number | undefined | null): number {
  if (!protoTime) {
    return Date.now();
  }

  if (protoTime instanceof Date) {
    return protoTime.getTime();
  }

  if (typeof protoTime === 'number') {
    return protoTime;
  }

  // Handle protobuf Timestamp object with getTime method
  if (typeof (protoTime as any).getTime === 'function') {
    return (protoTime as any).getTime();
  }

  // Fallback
  return Date.now();
}

/**
 * Convert protobuf ChatMessage to client ChatMessage
 */
export function convertProtoMessageToClientMessage(protoMessage: ProtoMessage): ChatMessage {
  const timestamp = extractTimestamp(protoMessage.time);

  return {
    id: `msg-${Date.now()}-${Math.random()}`,
    type: protoMessage.userName.toLowerCase() === 'system' ? 'system' : 'user',
    username: protoMessage.userName || 'Unknown',
    content: protoMessage.message,
    timestamp
  };
}

/**
 * Check if a message is a system message
 */
export function isSystemMessage(protoMessage: ProtoMessage): boolean {
  return protoMessage.userName.toLowerCase() === 'system';
}

/**
 * Format timestamp as HH:MM
 */
export function formatMessageTime(timestamp: number): string {
  const date = new Date(timestamp);
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  });
}

/**
 * Parse whisper command (/w username message)
 */
export function parseWhisperCommand(content: string): {
  isWhisper: boolean;
  username?: string;
  message?: string;
} {
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
 * Validate whisper command
 */
export function validateWhisperCommand(
  whisperInfo: { isWhisper: boolean; username?: string; message?: string },
  currentUsername?: string
): string | null {
  if (!whisperInfo.isWhisper) {
    return null;
  }

  if (!whisperInfo.username || !whisperInfo.message) {
    return 'Invalid whisper format. Use: /w username message';
  }

  // Cannot whisper to self
  if (currentUsername && whisperInfo.username.toLowerCase() === currentUsername.toLowerCase()) {
    return 'You cannot whisper to yourself';
  }

  return null;
}

/**
 * Rate limiting tracker
 */
export class RateLimiter {
  private messageTimestamps: number[] = [];
  private readonly maxMessages: number;
  private readonly windowMs: number;

  constructor(maxMessages: number = 10, windowMs: number = 60000) {
    this.maxMessages = maxMessages;
    this.windowMs = windowMs;
  }

  /**
   * Check if rate limited
   */
  isLimited(): boolean {
    const now = Date.now();
    // Remove timestamps older than the window
    this.messageTimestamps = this.messageTimestamps.filter((ts) => now - ts < this.windowMs);
    return this.messageTimestamps.length >= this.maxMessages;
  }

  /**
   * Record a message sent
   */
  recordMessage(): void {
    this.messageTimestamps.push(Date.now());
  }

  /**
   * Get cooldown time remaining in seconds
   */
  getCooldownSeconds(): number {
    if (this.messageTimestamps.length < this.maxMessages) {
      return 0;
    }

    const now = Date.now();
    const oldestTimestamp = this.messageTimestamps[0];
    const timeElapsed = now - oldestTimestamp;
    const timeRemaining = this.windowMs - timeElapsed;

    return Math.ceil(timeRemaining / 1000);
  }

  /**
   * Clear all recorded timestamps
   */
  reset(): void {
    this.messageTimestamps = [];
  }
}
