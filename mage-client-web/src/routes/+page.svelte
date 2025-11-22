<script lang="ts">
	// Home/Landing page with navigation
</script>

<svelte:head>
	<title>MAGE - Magic: The Gathering Online</title>
</svelte:head>

<div class="container">
	<div class="hero">
		<h1>🎮 MAGE</h1>
		<p>Magic: The Gathering Online - Play, Compete, Collect</p>
	</div>

	<div class="nav-grid">
		<a href="/login" class="nav-card">
			<div class="icon">🔐</div>
			<h2>Login</h2>
			<p>Sign in to your account</p>
		</a>

		<a href="/register" class="nav-card">
			<div class="icon">📝</div>
			<h2>Register</h2>
			<p>Create a new account</p>
		</a>

		<a href="/lobby" class="nav-card">
			<div class="icon">🎯</div>
			<h2>Lobby</h2>
			<p>Find and join games</p>
		</a>

		<a href="/decks" class="nav-card">
			<div class="icon">🎴</div>
			<h2>My Decks</h2>
			<p>Manage your deck collection</p>
		</a>

		<a href="/profile" class="nav-card">
			<div class="icon">👤</div>
			<h2>Profile</h2>
			<p>View your stats and settings</p>
		</a>

		<a href="/table/demo" class="nav-card">
			<div class="icon">🪑</div>
			<h2>Table (Demo)</h2>
			<p>Pre-game table lobby</p>
		</a>

		<a href="/game/demo" class="nav-card">
			<div class="icon">⚔️</div>
			<h2>Game (Demo)</h2>
			<p>Active game view</p>
		</a>
	</div>

	<div class="footer">
		<p>All routes are placeholder pages - full functionality coming soon!</p>
	</div>
</div>

<style>
	.container {
		min-height: 100vh;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 2rem;
	}

	.hero {
		text-align: center;
		margin-bottom: 3rem;
		color: white;
	}

	h1 {
		font-size: 4rem;
		margin: 0 0 1rem 0;
		text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	.hero p {
		font-size: 1.5rem;
		margin: 0;
		opacity: 0.9;
	}

	.nav-grid {
		max-width: 1200px;
		margin: 0 auto;
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1.5rem;
	}

	.nav-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		text-align: center;
		text-decoration: none;
		color: #333;
		transition: all 0.3s;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.nav-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
	}

	.icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.nav-card h2 {
		margin: 0 0 0.5rem 0;
		font-size: 1.5rem;
		color: #667eea;
	}

	.nav-card p {
		margin: 0;
		color: #666;
		font-size: 0.875rem;
	}

	.footer {
		text-align: center;
		margin-top: 3rem;
		color: white;
		opacity: 0.8;
	}

	.footer p {
		margin: 0;
		font-size: 0.875rem;
	}

	@media (max-width: 768px) {
		h1 {
			font-size: 2.5rem;
		}

		.hero p {
			font-size: 1.125rem;
		}

		.nav-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
