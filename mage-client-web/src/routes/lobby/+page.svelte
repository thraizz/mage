<script lang="ts">
	// Lobby page placeholder
	// TODO: Implement full lobby functionality in T017-T020
</script>

<svelte:head>
	<title>Lobby - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>Game Lobby</h1>
		<p>Find or create a game table</p>
	</header>

	<div class="lobby-content">
		<div class="tables-section">
			<div class="section-header">
				<h2>Active Tables</h2>
				<button class="btn-primary">Create Table</button>
			</div>

			<div class="table-list">
				<div class="table-card">
					<div class="table-info">
						<h3>Standard Game</h3>
						<p>Host: Player1</p>
						<p>Players: 2/4</p>
					</div>
					<button class="btn-secondary">Join</button>
				</div>

				<div class="table-card">
					<div class="table-info">
						<h3>Commander</h3>
						<p>Host: Player2</p>
						<p>Players: 3/4</p>
					</div>
					<button class="btn-secondary">Join</button>
				</div>

				<div class="empty-state">
					<p>More tables will appear here...</p>
				</div>
			</div>
		</div>

		<aside class="sidebar">
			<h3>Online Players (5)</h3>
			<ul class="player-list">
				<li><span class="status-dot"></span> You</li>
				<li><span class="status-dot"></span> Player1</li>
				<li><span class="status-dot"></span> Player2</li>
				<li><span class="status-dot"></span> Player3</li>
				<li><span class="status-dot"></span> Player4</li>
			</ul>
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

	.lobby-content {
		display: grid;
		grid-template-columns: 1fr 300px;
		gap: 2rem;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	h2 {
		margin: 0;
		font-size: 1.5rem;
		color: #333;
	}

	.table-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.table-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		transition: box-shadow 0.2s;
	}

	.table-card:hover {
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	.table-info h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1.25rem;
		color: #333;
	}

	.table-info p {
		margin: 0.25rem 0;
		color: #666;
	}

	.sidebar {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		height: fit-content;
	}

	.sidebar h3 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
		color: #333;
	}

	.player-list {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.player-list li {
		padding: 0.5rem 0;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #333;
	}

	.status-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: #10b981;
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

	.empty-state {
		text-align: center;
		padding: 3rem;
		color: #999;
	}

	@media (max-width: 768px) {
		.lobby-content {
			grid-template-columns: 1fr;
		}
	}
</style>
