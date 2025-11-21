<script lang="ts">
	import { onMount } from 'svelte';
	import { game, myCards, opponentCards, myPlayer, opponent, isMyTurn } from '$lib/stores/game';
	import Battlefield from '$lib/components/Battlefield.svelte';
	import PlayerInfo from '$lib/components/PlayerInfo.svelte';
	import type { Card } from '$lib/types';

	let connected = $state(false);
	let gameId = $state('test-game-1');
	let playerId = $state('player1');
	let selectedCard: Card | null = $state(null);

	onMount(async () => {
		try {
			await game.connect('ws://localhost:8080/ws');
			connected = true;
		} catch (error) {
			console.error('Failed to connect:', error);
		}
	});

	function handleCreateGame() {
		game.createGame(playerId, 'Duel');
	}

	function handleJoinGame() {
		game.joinGame(gameId, playerId);
	}

	function handleCardClick(card: Card) {
		selectedCard = card;
		console.log('Selected card:', card);
	}

	function handleDeclareAttacker(card: Card) {
		// Find opponent player ID
		const opp = $opponent;
		if (opp) {
			game.declareAttacker(card.id, opp.id);
		}
	}

	function handlePassPriority() {
		game.passPriority();
	}
</script>

<svelte:head>
	<title>Mage - Web Client</title>
</svelte:head>

<div class="game-container">
	<header>
		<h1>⚔️ Mage - Magic: The Gathering</h1>
		<div class="connection-status" class:connected>
			{connected ? '🟢 Connected' : '🔴 Disconnected'}
		</div>
	</header>

	{#if !connected}
		<div class="connecting">
			<p>Connecting to game server...</p>
			<p class="hint">Make sure the Go server is running on localhost:8080</p>
		</div>
	{:else if !$game.game_id}
		<div class="lobby">
			<h2>Game Lobby</h2>
			
			<div class="lobby-section">
				<input type="text" bind:value={playerId} placeholder="Your player ID" />
			</div>

			<div class="lobby-section">
				<button onclick={handleCreateGame} class="primary">Create New Game</button>
			</div>

			<div class="lobby-section">
				<input type="text" bind:value={gameId} placeholder="Game ID" />
				<button onclick={handleJoinGame}>Join Game</button>
			</div>
		</div>
	{:else}
		<div class="game-board">
			<!-- Opponent Info -->
			<div class="opponent-section">
				<PlayerInfo player={$opponent} isActive={!$isMyTurn} />
			</div>

			<!-- Opponent Battlefield -->
			<Battlefield 
				cards={$opponentCards} 
				title="Opponent's Battlefield"
				onCardClick={handleCardClick}
			/>

			<!-- Game Info -->
			<div class="game-info">
				<div class="turn-info">
					<span>Turn {$game.turn}</span>
					<span>{$game.phase} - {$game.step}</span>
					<span class:active={$isMyTurn}>
						{$isMyTurn ? '🎯 Your Turn' : '⏳ Opponent\'s Turn'}
					</span>
				</div>
				
				{#if selectedCard}
					<div class="selected-card">
						<strong>Selected:</strong> {selectedCard.name}
						{#if $isMyTurn && !selectedCard.attacking}
							<button onclick={() => selectedCard && handleDeclareAttacker(selectedCard)}>
								⚔️ Attack
							</button>
						{/if}
					</div>
				{/if}

				<div class="actions">
					<button onclick={handlePassPriority} disabled={!$isMyTurn}>
						Pass Priority
					</button>
				</div>
			</div>

			<!-- My Battlefield -->
			<Battlefield 
				cards={$myCards} 
				title="Your Battlefield"
				onCardClick={handleCardClick}
			/>

			<!-- My Info -->
			<div class="player-section">
				<PlayerInfo player={$myPlayer} isActive={$isMyTurn} />
			</div>
		</div>
	{/if}
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: system-ui, -apple-system, sans-serif;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		min-height: 100vh;
	}

	.game-container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 20px;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: white;
		padding: 16px 24px;
		border-radius: 8px;
		margin-bottom: 20px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	h1 {
		margin: 0;
		font-size: 28px;
		color: #333;
	}

	.connection-status {
		font-weight: bold;
		color: #e74c3c;
	}

	.connection-status.connected {
		color: #27ae60;
	}

	.connecting {
		background: white;
		padding: 64px;
		border-radius: 8px;
		text-align: center;
	}

	.hint {
		color: #999;
		font-size: 14px;
		margin-top: 8px;
	}

	.lobby {
		background: white;
		padding: 48px;
		border-radius: 8px;
		max-width: 500px;
		margin: 0 auto;
	}

	.lobby h2 {
		margin-top: 0;
	}

	.lobby-section {
		margin-bottom: 24px;
	}

	input {
		width: 100%;
		padding: 12px;
		border: 2px solid #ddd;
		border-radius: 4px;
		font-size: 16px;
		box-sizing: border-box;
	}

	button {
		padding: 12px 24px;
		border: none;
		border-radius: 4px;
		font-size: 16px;
		cursor: pointer;
		background: #667eea;
		color: white;
		font-weight: bold;
		transition: all 0.2s ease;
	}

	button:hover:not(:disabled) {
		background: #5568d3;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	button.primary {
		width: 100%;
		background: #27ae60;
	}

	button.primary:hover {
		background: #229954;
	}

	.game-board {
		display: grid;
		gap: 20px;
	}

	.opponent-section,
	.player-section {
		background: white;
		border-radius: 8px;
		padding: 16px;
	}

	.game-info {
		background: white;
		border-radius: 8px;
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.turn-info {
		display: flex;
		justify-content: space-around;
		align-items: center;
		font-size: 18px;
		font-weight: bold;
	}

	.turn-info .active {
		color: #27ae60;
	}

	.selected-card {
		padding: 12px;
		background: #f0f0f0;
		border-radius: 4px;
		display: flex;
		gap: 12px;
		align-items: center;
	}

	.actions {
		display: flex;
		gap: 12px;
		justify-content: center;
	}
</style>
