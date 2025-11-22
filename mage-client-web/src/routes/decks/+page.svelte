<script lang="ts">
	// Deck management page placeholder
	// TODO: Implement full deck management functionality in T031-T035
</script>

<svelte:head>
	<title>My Decks - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>My Decks</h1>
		<button class="btn-primary">Upload New Deck</button>
	</header>

	<div class="decks-grid">
		<div class="deck-card">
			<div class="deck-header">
				<h3>Mono-Red Aggro</h3>
				<span class="format-badge">Standard</span>
			</div>
			<div class="deck-stats">
				<p>60 cards</p>
				<p>Last modified: 2 days ago</p>
			</div>
			<div class="deck-colors">
				<span class="mana-symbol red">🔴</span>
			</div>
			<div class="deck-actions">
				<button class="btn-secondary">View</button>
				<button class="btn-secondary">Edit</button>
				<button class="btn-danger">Delete</button>
			</div>
		</div>

		<div class="deck-card">
			<div class="deck-header">
				<h3>Azorius Control</h3>
				<span class="format-badge">Standard</span>
			</div>
			<div class="deck-stats">
				<p>60 cards</p>
				<p>Last modified: 1 week ago</p>
			</div>
			<div class="deck-colors">
				<span class="mana-symbol white">⚪</span>
				<span class="mana-symbol blue">🔵</span>
			</div>
			<div class="deck-actions">
				<button class="btn-secondary">View</button>
				<button class="btn-secondary">Edit</button>
				<button class="btn-danger">Delete</button>
			</div>
		</div>

		<div class="deck-card">
			<div class="deck-header">
				<h3>Elves Tribal</h3>
				<span class="format-badge">Commander</span>
			</div>
			<div class="deck-stats">
				<p>100 cards</p>
				<p>Last modified: 2 weeks ago</p>
			</div>
			<div class="deck-colors">
				<span class="mana-symbol green">🟢</span>
			</div>
			<div class="deck-actions">
				<button class="btn-secondary">View</button>
				<button class="btn-secondary">Edit</button>
				<button class="btn-danger">Delete</button>
			</div>
		</div>

		<div class="empty-state">
			<p>📦 Upload more decks to see them here</p>
			<button class="btn-primary">Upload Deck</button>
		</div>
	</div>
</div>

<style>
	.container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h1 {
		margin: 0;
		font-size: 2.5rem;
		color: #333;
	}

	.decks-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.deck-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		transition: box-shadow 0.2s;
	}

	.deck-card:hover {
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.deck-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 1rem;
	}

	.deck-header h3 {
		margin: 0;
		font-size: 1.25rem;
		color: #333;
		flex: 1;
	}

	.format-badge {
		padding: 0.25rem 0.75rem;
		background: #667eea;
		color: white;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 500;
		white-space: nowrap;
	}

	.deck-stats {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.deck-stats p {
		margin: 0;
		font-size: 0.875rem;
		color: #666;
	}

	.deck-colors {
		display: flex;
		gap: 0.5rem;
		font-size: 1.5rem;
	}

	.deck-actions {
		display: flex;
		gap: 0.5rem;
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
		flex: 1;
	}

	.btn-secondary:hover {
		background: #667eea;
		color: white;
	}

	.btn-danger {
		padding: 0.5rem 1rem;
		background: white;
		color: #ef4444;
		border: 1px solid #ef4444;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-danger:hover {
		background: #ef4444;
		color: white;
	}

	.empty-state {
		grid-column: 1 / -1;
		text-align: center;
		padding: 3rem;
		background: #f9fafb;
		border: 2px dashed #ddd;
		border-radius: 8px;
	}

	.empty-state p {
		margin: 0 0 1rem 0;
		font-size: 1.125rem;
		color: #666;
	}
</style>
