<script lang="ts">
	// User profile page placeholder
	// TODO: Implement full profile functionality in T036-T038
</script>

<svelte:head>
	<title>Profile - MAGE</title>
</svelte:head>

<div class="container">
	<header>
		<h1>Player Profile</h1>
	</header>

	<div class="profile-content">
		<div class="profile-card">
			<div class="avatar">
				<div class="avatar-placeholder">👤</div>
			</div>
			<h2>YourUsername</h2>
			<p>player@example.com</p>
			<p class="joined">Joined: November 2024</p>
		</div>

		<div class="stats-grid">
			<div class="stat-card">
				<h3>Games Played</h3>
				<p class="stat-value">42</p>
			</div>

			<div class="stat-card">
				<h3>Wins</h3>
				<p class="stat-value">28</p>
			</div>

			<div class="stat-card">
				<h3>Losses</h3>
				<p class="stat-value">14</p>
			</div>

			<div class="stat-card">
				<h3>Win Rate</h3>
				<p class="stat-value">66.7%</p>
			</div>
		</div>

		<div class="section">
			<h3>Recent Matches</h3>
			<div class="matches-list">
				<div class="match-item win">
					<div class="match-info">
						<strong>vs Player1</strong>
						<span>Standard</span>
					</div>
					<span class="result">Win</span>
				</div>

				<div class="match-item loss">
					<div class="match-info">
						<strong>vs Player2</strong>
						<span>Commander</span>
					</div>
					<span class="result">Loss</span>
				</div>

				<div class="match-item win">
					<div class="match-info">
						<strong>vs Player3</strong>
						<span>Modern</span>
					</div>
					<span class="result">Win</span>
				</div>
			</div>
		</div>

		<div class="section">
			<h3>Settings</h3>
			<button class="btn-secondary">Change Password</button>
			<button class="btn-secondary">Edit Profile</button>
		</div>
	</div>
</div>

<style>
	.container {
		max-width: 1200px;
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

	.profile-content {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.profile-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 2rem;
		text-align: center;
	}

	.avatar {
		margin: 0 auto 1rem;
		width: 100px;
		height: 100px;
	}

	.avatar-placeholder {
		width: 100%;
		height: 100%;
		background: #667eea;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 3rem;
	}

	.profile-card h2 {
		margin: 0 0 0.5rem 0;
		color: #333;
	}

	.profile-card p {
		margin: 0.25rem 0;
		color: #666;
	}

	.joined {
		font-size: 0.875rem;
		color: #999;
	}

	.stats-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 1rem;
	}

	.stat-card {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
		text-align: center;
	}

	.stat-card h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1rem;
		color: #666;
		font-weight: 500;
	}

	.stat-value {
		margin: 0;
		font-size: 2rem;
		font-weight: 700;
		color: #667eea;
	}

	.section {
		background: white;
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1.5rem;
	}

	.section h3 {
		margin: 0 0 1rem 0;
		color: #333;
	}

	.matches-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.match-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		border-radius: 4px;
		background: #f9fafb;
	}

	.match-item.win {
		border-left: 4px solid #10b981;
	}

	.match-item.loss {
		border-left: 4px solid #ef4444;
	}

	.match-info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.match-info strong {
		color: #333;
	}

	.match-info span {
		font-size: 0.875rem;
		color: #666;
	}

	.result {
		font-weight: 600;
	}

	.match-item.win .result {
		color: #10b981;
	}

	.match-item.loss .result {
		color: #ef4444;
	}

	.btn-secondary {
		padding: 0.75rem 1.5rem;
		background: white;
		color: #667eea;
		border: 1px solid #667eea;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
		margin-right: 1rem;
	}

	.btn-secondary:hover {
		background: #667eea;
		color: white;
	}
</style>
