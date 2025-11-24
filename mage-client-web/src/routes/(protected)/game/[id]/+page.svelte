<script lang="ts">
	import { page } from '$app/stores';

	// Game view page placeholder
	// TODO: Implement full game functionality in T039-T076

	$: gameId = $page.params.id;
</script>

<svelte:head>
	<title>Game {gameId} - MAGE</title>
</svelte:head>

<div class="game-container">
	<div class="game-header">
		<div class="game-info">
			<h2>Game #{gameId}</h2>
			<p>Turn 1 - Main Phase</p>
		</div>
		<button class="btn-danger">Concede</button>
	</div>

	<div class="game-layout">
		<!-- Opponent Area -->
		<div class="opponent-area">
			<div class="opponent-info">
				<h3>Opponent</h3>
				<div class="stats">
					<span class="life">❤️ 20</span>
					<span class="library">📚 53</span>
					<span class="hand">🎴 7</span>
				</div>
			</div>
		</div>

		<!-- Battlefield -->
		<div class="battlefield">
			<div class="zone-label">Battlefield</div>
			<div class="cards-area">
				<p class="placeholder">Cards will appear here during gameplay...</p>
			</div>
		</div>

		<!-- Player Area -->
		<div class="player-area">
			<div class="player-info">
				<h3>You</h3>
				<div class="stats">
					<span class="life">❤️ 20</span>
					<span class="library">📚 53</span>
				</div>
			</div>
			<div class="hand">
				<p class="zone-label">Your Hand</p>
				<div class="hand-cards">
					<div class="card-placeholder">Card 1</div>
					<div class="card-placeholder">Card 2</div>
					<div class="card-placeholder">Card 3</div>
					<div class="card-placeholder">Card 4</div>
					<div class="card-placeholder">Card 5</div>
					<div class="card-placeholder">Card 6</div>
					<div class="card-placeholder">Card 7</div>
				</div>
			</div>
			<div class="actions">
				<button class="btn-primary">Pass Priority</button>
			</div>
		</div>
	</div>

	<!-- Game Chat -->
	<div class="game-chat">
		<h4>Game Log</h4>
		<div class="log-entries">
			<p class="log-entry">Game started</p>
			<p class="log-entry">You drew 7 cards</p>
			<p class="log-entry">Your turn begins</p>
		</div>
	</div>
</div>

<style>
	.game-container {
		min-height: 100vh;
		background: #1a1a1a;
		color: white;
		display: flex;
		flex-direction: column;
	}

	.game-header {
		background: #2a2a2a;
		padding: 1rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 2px solid #3a3a3a;
	}

	.game-info h2 {
		margin: 0;
		color: white;
	}

	.game-info p {
		margin: 0.25rem 0 0 0;
		color: #999;
	}

	.game-layout {
		flex: 1;
		display: grid;
		grid-template-rows: 150px 1fr 250px;
		padding: 1rem;
		gap: 1rem;
	}

	.opponent-area,
	.player-area {
		background: #2a2a2a;
		border-radius: 8px;
		padding: 1rem;
	}

	.opponent-info,
	.player-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.opponent-info h3,
	.player-info h3 {
		margin: 0;
	}

	.stats {
		display: flex;
		gap: 1rem;
		font-size: 1.25rem;
	}

	.life {
		color: #ef4444;
	}

	.library {
		color: #3b82f6;
	}

	.hand {
		color: #fbbf24;
	}

	.battlefield {
		background: #1a3a1a;
		border-radius: 8px;
		padding: 1rem;
		position: relative;
	}

	.zone-label {
		font-size: 0.875rem;
		color: #999;
		margin-bottom: 0.5rem;
	}

	.cards-area {
		min-height: 200px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.placeholder {
		color: #666;
		font-style: italic;
	}

	.hand {
		margin-top: 1rem;
	}

	.hand-cards {
		display: flex;
		gap: 0.5rem;
		overflow-x: auto;
		padding: 1rem 0;
	}

	.card-placeholder {
		min-width: 100px;
		height: 140px;
		background: #3a3a3a;
		border: 2px solid #4a4a4a;
		border-radius: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition:
			transform 0.2s,
			border-color 0.2s;
	}

	.card-placeholder:hover {
		transform: translateY(-10px);
		border-color: #667eea;
	}

	.actions {
		margin-top: 1rem;
		display: flex;
		justify-content: center;
		gap: 1rem;
	}

	.game-chat {
		position: fixed;
		right: 0;
		top: 0;
		bottom: 0;
		width: 300px;
		background: #2a2a2a;
		border-left: 2px solid #3a3a3a;
		padding: 1rem;
		display: flex;
		flex-direction: column;
	}

	.game-chat h4 {
		margin: 0 0 1rem 0;
	}

	.log-entries {
		flex: 1;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.log-entry {
		margin: 0;
		padding: 0.5rem;
		background: #1a1a1a;
		border-radius: 4px;
		font-size: 0.875rem;
		color: #ccc;
	}

	.btn-primary {
		padding: 0.75rem 1.5rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	.btn-danger {
		padding: 0.75rem 1.5rem;
		background: #ef4444;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-danger:hover {
		background: #dc2626;
	}

	@media (max-width: 1200px) {
		.game-chat {
			display: none;
		}
	}
</style>
