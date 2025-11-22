<script lang="ts">
	import { page } from '$app/stores';

	// Table lobby page placeholder
	// TODO: Implement full table lobby functionality in T026-T030

	$: tableId = $page.params.id;
</script>

<svelte:head>
	<title>Table {tableId} - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>Table #{tableId}</h1>
		<p>Pre-game lobby - waiting for players</p>
	</header>

	<div class="table-content">
		<div class="main-area">
			<div class="table-info">
				<h2>Standard Game</h2>
				<p>Host: Player1</p>
			</div>

			<div class="players-grid">
				<div class="player-slot occupied">
					<div class="player-info">
						<h3>Player1</h3>
						<p>Host</p>
						<span class="ready-badge">Ready</span>
					</div>
				</div>

				<div class="player-slot occupied">
					<div class="player-info">
						<h3>You</h3>
						<p>Player</p>
						<button class="btn-secondary">Ready Up</button>
					</div>
				</div>

				<div class="player-slot empty">
					<p>Waiting for player...</p>
				</div>

				<div class="player-slot empty">
					<p>Waiting for player...</p>
				</div>
			</div>

			<div class="actions">
				<button class="btn-danger">Leave Table</button>
			</div>
		</div>

		<aside class="chat-sidebar">
			<h3>Table Chat</h3>
			<div class="chat-messages">
				<div class="chat-message system">
					<p>You joined the table</p>
				</div>
				<div class="chat-message">
					<strong>Player1:</strong> Welcome!
				</div>
			</div>
			<div class="chat-input">
				<input type="text" placeholder="Type a message..." />
				<button class="btn-primary">Send</button>
			</div>
		</aside>
	</div>
</div>

<style>
	.container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		margin-bottom: 2rem;
	}

	h1 {
		margin: 0;
		font-size: 2.5rem;
		color: #333;
	}

	header p {
		margin: 0.5rem 0 0 0;
		color: #666;
	}

	.table-content {
		display: grid;
		grid-template-columns: 1fr 350px;
		gap: 2rem;
	}

	.main-area {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.table-info {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.table-info h2 {
		margin: 0 0 0.5rem 0;
		color: #333;
	}

	.table-info p {
		margin: 0;
		color: #666;
	}

	.players-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
	}

	.player-slot {
		background: white;
		border: 2px solid #ddd;
		border-radius: 8px;
		padding: 2rem;
		text-align: center;
		min-height: 150px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.player-slot.occupied {
		border-color: #667eea;
	}

	.player-slot.empty {
		background: #f9fafb;
		color: #999;
	}

	.player-info h3 {
		margin: 0 0 0.5rem 0;
		color: #333;
	}

	.player-info p {
		margin: 0 0 1rem 0;
		color: #666;
		font-size: 0.875rem;
	}

	.ready-badge {
		display: inline-block;
		padding: 0.25rem 0.75rem;
		background: #10b981;
		color: white;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.actions {
		display: flex;
		justify-content: center;
		gap: 1rem;
	}

	.chat-sidebar {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		height: fit-content;
		max-height: 600px;
	}

	.chat-sidebar h3 {
		margin: 0 0 1rem 0;
		color: #333;
	}

	.chat-messages {
		flex: 1;
		overflow-y: auto;
		margin-bottom: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.chat-message {
		padding: 0.5rem;
		border-radius: 4px;
		font-size: 0.875rem;
	}

	.chat-message.system {
		color: #999;
		font-style: italic;
	}

	.chat-input {
		display: flex;
		gap: 0.5rem;
	}

	.chat-input input {
		flex: 1;
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	.btn-primary {
		padding: 0.5rem 1rem;
		background: #667eea;
		color: white;
		border: none;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background: #5568d3;
	}

	.btn-secondary {
		padding: 0.5rem 1rem;
		background: white;
		color: #667eea;
		border: 1px solid #667eea;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover {
		background: #667eea;
		color: white;
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

	@media (max-width: 768px) {
		.table-content {
			grid-template-columns: 1fr;
		}

		.players-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
