/**
 * Tests for decks.ts API functions
 * Tests the behavior that fetchUserDecks returns summary decks without card lists,
 * and getDeckDetails returns full decks with card lists.
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { fetchUserDecks, getDeckDetails } from '../decks';
import type { DeckListResponse, DeckGetResponse, DeckInfo, DeckCardLists } from '$lib/generated/mage/v1/table';

// Mock the gRPC client
vi.mock('$lib/grpc/client', () => {
	const mockClient = {
		ensureSessionId: vi.fn().mockResolvedValue('test-session-id'),
		call: vi.fn()
	};

	return {
		getMageClient: vi.fn(() => mockClient)
	};
});

describe('decks.ts API', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
		// Reset ensureSessionId to return a valid session by default
		const { getMageClient } = await import('$lib/grpc/client');
		const mockClient = getMageClient();
		mockClient.ensureSessionId.mockResolvedValue('test-session-id');
	});

	describe('fetchUserDecks', () => {
		it('should return decks with empty card arrays (summary only)', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockDeckInfo: DeckInfo = {
				id: 1,
				name: 'Test Commander Deck',
				format: 'Commander',
				description: '',
				mainDeckCount: 99,
				sideboardCount: 0,
				createdAt: 1000, // seconds
				updatedAt: 2000 // seconds
			};

			const mockResponse: DeckListResponse = {
				success: true,
				error: '',
				decks: [mockDeckInfo]
			};

			mockClient.call.mockResolvedValue(mockResponse);

			const decks = await fetchUserDecks();

			expect(decks).toHaveLength(1);
			expect(decks[0]).toEqual({
				id: '1',
				name: 'Test Commander Deck',
				format: 'Commander',
				cardCount: 99,
				createdAt: 1000000, // milliseconds
				updatedAt: 2000000, // milliseconds
				isValid: true,
				// These should be empty arrays - summary view doesn't include card details
				mainDeck: [],
				sideboard: [],
				commanders: []
			});

			// Verify the client was called correctly
			expect(mockClient.call).toHaveBeenCalledWith('DeckList', {
				sessionId: 'test-session-id',
				format: ''
			});
		});

		it('should filter decks by format when provided', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockResponse: DeckListResponse = {
				success: true,
				error: '',
				decks: []
			};

			mockClient.call.mockResolvedValue(mockResponse);

			await fetchUserDecks('Commander');

			expect(mockClient.call).toHaveBeenCalledWith('DeckList', {
				sessionId: 'test-session-id',
				format: 'Commander'
			});
		});

		it('should throw error when session is missing', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			mockClient.ensureSessionId.mockResolvedValue(null);

			await expect(fetchUserDecks()).rejects.toThrow('No active session');
		});

		it('should throw error when server returns error', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockResponse: DeckListResponse = {
				success: false,
				error: 'Session expired',
				decks: []
			};

			mockClient.call.mockResolvedValue(mockResponse);

			await expect(fetchUserDecks()).rejects.toThrow('Session expired');
		});
	});

	describe('getDeckDetails', () => {
		it('should return full deck with card lists', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockDeckInfo: DeckInfo = {
				id: 1,
				name: 'Test Commander Deck',
				format: 'Commander',
				description: '',
				mainDeckCount: 99,
				sideboardCount: 0,
				createdAt: 1000,
				updatedAt: 2000
			};

			const mockDeckCards: DeckCardLists = {
				mainDeck: [
					{
						name: 'Lightning Bolt',
						quantity: 4,
						manaCost: '{R}',
						cardType: 'Instant',
						types: ['INSTANT'],
						colors: ['R'],
						power: '',
						toughness: ''
					},
					{
						name: 'Mountain',
						quantity: 20,
						manaCost: '',
						cardType: 'Basic Land - Mountain',
						types: ['LAND'],
						colors: ['R'],
						power: '',
						toughness: ''
					}
				],
				sideboard: [],
				commanders: [
					{
						name: 'Atraxa, Praetors\' Voice',
						quantity: 1,
						manaCost: '{G}{W}{U}{B}',
						cardType: 'Legendary Creature - Phyrexian Horror',
						types: ['CREATURE'],
						colors: ['G', 'W', 'U', 'B'],
						power: '4',
						toughness: '4'
					}
				]
			};

			const mockResponse: DeckGetResponse = {
				success: true,
				error: '',
				info: mockDeckInfo,
				deck: mockDeckCards
			};

			mockClient.call.mockResolvedValue(mockResponse);

			const deck = await getDeckDetails('1');

			expect(deck).toEqual({
				id: '1',
				name: 'Test Commander Deck',
				format: 'Commander',
				cardCount: 99,
				createdAt: 1000000,
				updatedAt: 2000000,
				isValid: true,
				// These should contain the actual card data
				mainDeck: [
					{
						cardName: 'Lightning Bolt',
						quantity: 4,
						setCode: undefined,
						manaCost: '{R}',
						cardType: 'Instant',
						types: ['INSTANT'],
						colors: ['R'],
						power: '',
						toughness: ''
					},
					{
						cardName: 'Mountain',
						quantity: 20,
						setCode: undefined,
						manaCost: '',
						cardType: 'Basic Land - Mountain',
						types: ['LAND'],
						colors: ['R'],
						power: '',
						toughness: ''
					}
				],
				sideboard: [],
				commanders: [
					{
						cardName: 'Atraxa, Praetors\' Voice',
						quantity: 1,
						setCode: undefined,
						manaCost: '{G}{W}{U}{B}',
						cardType: 'Legendary Creature - Phyrexian Horror',
						types: ['CREATURE'],
						colors: ['G', 'W', 'U', 'B'],
						power: '4',
						toughness: '4'
					}
				]
			});

			// Verify the client was called correctly
			expect(mockClient.call).toHaveBeenCalledWith('DeckGet', {
				sessionId: 'test-session-id',
				deckId: 1
			});
		});

		it('should throw error when session is missing', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			mockClient.ensureSessionId.mockResolvedValue(null);

			await expect(getDeckDetails('1')).rejects.toThrow('No active session');
		});

		it('should throw error when server returns error', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockResponse: DeckGetResponse = {
				success: false,
				error: 'Deck not found',
				info: undefined,
				deck: undefined
			};

			mockClient.call.mockResolvedValue(mockResponse);

			await expect(getDeckDetails('1')).rejects.toThrow('Deck not found');
		});

		it('should throw error when deck info or cards are missing', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			const mockResponse: DeckGetResponse = {
				success: true,
				error: '',
				info: undefined,
				deck: undefined
			};

			mockClient.call.mockResolvedValue(mockResponse);

			await expect(getDeckDetails('1')).rejects.toThrow('Deck not found');
		});
	});

	describe('fetchUserDecks vs getDeckDetails behavior', () => {
		it('should demonstrate that fetchUserDecks returns empty card arrays while getDeckDetails returns full data', async () => {
			const { getMageClient } = await import('$lib/grpc/client');
			const mockClient = getMageClient();

			// Mock fetchUserDecks response (summary only)
			const mockDeckInfo: DeckInfo = {
				id: 1,
				name: 'My Commander Deck',
				format: 'Commander',
				description: '',
				mainDeckCount: 100,
				sideboardCount: 0,
				createdAt: 1000,
				updatedAt: 2000
			};

			const mockListResponse: DeckListResponse = {
				success: true,
				error: '',
				decks: [mockDeckInfo]
			};

			// Mock getDeckDetails response (full data)
			const mockDeckCards: DeckCardLists = {
				mainDeck: [
					{ name: 'Lightning Bolt', quantity: 4, manaCost: '{R}', cardType: 'Instant', types: [], colors: [], power: '', toughness: '' }
				],
				sideboard: [],
				commanders: [
					{ name: 'Atraxa', quantity: 1, manaCost: '{G}{W}{U}{B}', cardType: 'Creature', types: [], colors: [], power: '4', toughness: '4' }
				]
			};

			const mockGetResponse: DeckGetResponse = {
				success: true,
				error: '',
				info: mockDeckInfo,
				deck: mockDeckCards
			};

			// First call: fetchUserDecks
			mockClient.call.mockResolvedValueOnce(mockListResponse);
			const summaryDecks = await fetchUserDecks();

			// Second call: getDeckDetails
			mockClient.call.mockResolvedValueOnce(mockGetResponse);
			const fullDeck = await getDeckDetails('1');

			// Verify fetchUserDecks returns empty arrays
			expect(summaryDecks[0].mainDeck).toEqual([]);
			expect(summaryDecks[0].sideboard).toEqual([]);
			expect(summaryDecks[0].commanders).toEqual([]);
			expect(summaryDecks[0].cardCount).toBe(100); // But count is available

			// Verify getDeckDetails returns full card data
			expect(fullDeck.mainDeck).toHaveLength(1);
			expect(fullDeck.mainDeck[0].cardName).toBe('Lightning Bolt');
			expect(fullDeck.commanders).toHaveLength(1);
			expect(fullDeck.commanders[0].cardName).toBe('Atraxa');
			expect(fullDeck.cardCount).toBe(100); // Count is still available
		});
	});
});

