// Game state store with WebSocket integration

import { writable, derived } from 'svelte/store';
import type { GameState, Card, Player } from '../types';
import { MageWebSocket } from '../websocket';

// WebSocket connection
let ws: MageWebSocket | null = null;

// Initial empty state
const initialState: GameState = {
	game_id: '',
	current_player: '',
	active_player: '',
	priority_player: '',
	phase: '',
	step: '',
	turn: 0,
	players: [],
	battlefield: [],
	hand: [],
	graveyard: [],
	exile: [],
	stack: []
};

// Create the main game store
function createGameStore() {
	const { subscribe, set, update } = writable<GameState>(initialState);

	return {
		subscribe,
		
		// Connect to game server
		connect: async (url: string = 'ws://localhost:8080/ws') => {
			ws = new MageWebSocket(url);
			
			// Handle game state updates
			ws.on('game_state', (data: GameState) => {
				set(data);
			});
			
			// Handle incremental updates
			ws.on('card_moved', (data: any) => {
				update(state => {
					// Update card zone
					const card = state.battlefield.find(c => c.id === data.card_id);
					if (card) {
						card.zone = data.to_zone;
					}
					return state;
				});
			});
			
			await ws.connect();
		},
		
		// Join a game
		joinGame: (gameId: string, playerId: string) => {
			ws?.send({
				type: 'join_game',
				game_id: gameId,
				player_id: playerId
			});
		},
		
		// Create a new game
		createGame: (playerId: string, gameType: string = 'Duel') => {
			ws?.send({
				type: 'create_game',
				player_id: playerId,
				data: { game_type: gameType }
			});
		},
		
		// Play a card from hand
		playCard: (cardId: string) => {
			ws?.send({
				type: 'play_card',
				data: { card_id: cardId }
			});
		},
		
		// Declare attacker
		declareAttacker: (cardId: string, defenderId: string) => {
			ws?.send({
				type: 'declare_attacker',
				data: { card_id: cardId, defender_id: defenderId }
			});
		},
		
		// Declare blocker
		declareBlocker: (blockerId: string, attackerId: string) => {
			ws?.send({
				type: 'declare_blocker',
				data: { blocker_id: blockerId, attacker_id: attackerId }
			});
		},
		
		// Pass priority
		passPriority: () => {
			ws?.send({
				type: 'pass_priority'
			});
		},
		
		// Disconnect
		disconnect: () => {
			ws?.disconnect();
			ws = null;
			set(initialState);
		}
	};
}

export const game = createGameStore();

// Derived stores for convenience
export const myCards = derived(game, $game => 
	$game.battlefield.filter(c => c.controller === 'player1')
);

export const opponentCards = derived(game, $game => 
	$game.battlefield.filter(c => c.controller !== 'player1')
);

export const myPlayer = derived(game, $game => 
	$game.players.find(p => p.id === 'player1')
);

export const opponent = derived(game, $game => 
	$game.players.find(p => p.id !== 'player1')
);

export const isMyTurn = derived(game, $game => 
	$game.current_player === 'player1'
);
