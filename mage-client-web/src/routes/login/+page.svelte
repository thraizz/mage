<script lang="ts">
	// Login page placeholder
	// TODO: Implement full login functionality in T007
</script>

<svelte:head>
	<title>Login - MAGE</title>
</svelte:head>

<div class="container">
	<div class="card">
		<h1>Login</h1>
		<p>Welcome to MAGE - Magic: The Gathering Online</p>

		<form class="login-form">
			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					name="username"
					placeholder="Enter your username"
					required
				/>
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					name="password"
					placeholder="Enter your password"
					required
				/>
			</div>

			<button type="submit" class="btn-primary">Login</button>
		</form>

		<div class="links">
			<a href="/register">Don't have an account? Register</a>
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		padding: 2rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.card {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 400px;
	}

	h1 {
		margin: 0 0 0.5rem 0;
		font-size: 2rem;
		color: #333;
	}

	p {
		margin: 0 0 2rem 0;
		color: #666;
	}

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	label {
		font-weight: 500;
		color: #333;
	}

	input {
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
	}

	input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
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

	.links {
		margin-top: 1rem;
		text-align: center;
	}

	.links a {
		color: #667eea;
		text-decoration: none;
	}

	.links a:hover {
		text-decoration: underline;
	}
</style>
