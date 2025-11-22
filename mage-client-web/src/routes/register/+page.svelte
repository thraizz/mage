<script lang="ts">
	// Registration page placeholder
	// TODO: Implement full registration functionality in T008
</script>

<svelte:head>
	<title>Register - MAGE</title>
</svelte:head>

<div class="container">
	<div class="card">
		<h1>Create Account</h1>
		<p>Join MAGE and start playing Magic: The Gathering online</p>

		<form class="register-form">
			<div class="form-group">
				<label for="username">Username</label>
				<input
					type="text"
					id="username"
					name="username"
					placeholder="Choose a username"
					required
				/>
			</div>

			<div class="form-group">
				<label for="email">Email</label>
				<input type="email" id="email" name="email" placeholder="Enter your email" required />
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<input
					type="password"
					id="password"
					name="password"
					placeholder="Create a password"
					required
				/>
			</div>

			<div class="form-group">
				<label for="confirm-password">Confirm Password</label>
				<input
					type="password"
					id="confirm-password"
					name="confirm-password"
					placeholder="Confirm your password"
					required
				/>
			</div>

			<button type="submit" class="btn-primary">Create Account</button>
		</form>

		<div class="links">
			<a href="/login">Already have an account? Login</a>
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

	.register-form {
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
