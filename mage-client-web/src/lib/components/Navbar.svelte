<script lang="ts">
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth';
	import UserMenu from './UserMenu.svelte';
	import ConnectionStatus from './ConnectionStatus.svelte';

	let mobileMenuOpen = false;

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	// Close mobile menu when route changes
	$: if ($page.url.pathname) {
		closeMobileMenu();
	}
</script>

<nav class="navbar">
	<div class="navbar-container">
		<!-- Logo/Brand -->
		<div class="navbar-brand">
			<a href="/lobby" class="brand-link">
				<span class="brand-icon">🎴</span>
				<span class="brand-text">Mage</span>
			</a>
		</div>

		<!-- Desktop Navigation Links -->
		<div class="navbar-links desktop-only">
			<a href="/lobby" class="nav-link" class:active={$page.url.pathname === '/lobby'}>
				Lobby
			</a>
			<a href="/decks" class="nav-link" class:active={$page.url.pathname === '/decks'}>
				My Decks
			</a>
			<a href="/profile" class="nav-link" class:active={$page.url.pathname === '/profile'}>
				Profile
			</a>
		</div>

		<!-- Right Side: Connection Status + User Menu -->
		<div class="navbar-right">
			<ConnectionStatus />
			<UserMenu />

			<!-- Mobile Menu Button -->
			<button
				class="mobile-menu-button mobile-only"
				on:click={toggleMobileMenu}
				aria-label="Toggle menu"
			>
				<span class="hamburger-line"></span>
				<span class="hamburger-line"></span>
				<span class="hamburger-line"></span>
			</button>
		</div>
	</div>

	<!-- Mobile Menu -->
	{#if mobileMenuOpen}
		<div class="mobile-menu">
			<a
				href="/lobby"
				class="mobile-nav-link"
				class:active={$page.url.pathname === '/lobby'}
				on:click={closeMobileMenu}
			>
				Lobby
			</a>
			<a
				href="/decks"
				class="mobile-nav-link"
				class:active={$page.url.pathname === '/decks'}
				on:click={closeMobileMenu}
			>
				My Decks
			</a>
			<a
				href="/profile"
				class="mobile-nav-link"
				class:active={$page.url.pathname === '/profile'}
				on:click={closeMobileMenu}
			>
				Profile
			</a>
		</div>
	{/if}
</nav>

<style>
	.navbar {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		position: sticky;
		top: 0;
		z-index: 100;
	}

	.navbar-container {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 1rem;
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 64px;
	}

	/* Brand */
	.navbar-brand {
		display: flex;
		align-items: center;
	}

	.brand-link {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		text-decoration: none;
		color: white;
		font-weight: 700;
		font-size: 1.5rem;
		transition: opacity 0.2s;
	}

	.brand-link:hover {
		opacity: 0.9;
	}

	.brand-icon {
		font-size: 2rem;
	}

	.brand-text {
		display: none;
	}

	/* Navigation Links */
	.navbar-links {
		display: flex;
		align-items: center;
		gap: 2rem;
		flex: 1;
		justify-content: center;
	}

	.nav-link {
		color: white;
		text-decoration: none;
		font-weight: 500;
		font-size: 1rem;
		padding: 0.5rem 1rem;
		border-radius: 0.375rem;
		transition: background-color 0.2s;
		position: relative;
	}

	.nav-link:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.nav-link.active {
		background-color: rgba(255, 255, 255, 0.2);
	}

	.nav-link.active::after {
		content: '';
		position: absolute;
		bottom: -8px;
		left: 0;
		right: 0;
		height: 3px;
		background-color: white;
		border-radius: 2px;
	}

	/* Right Side */
	.navbar-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	/* Mobile Menu Button */
	.mobile-menu-button {
		display: none;
		flex-direction: column;
		justify-content: space-around;
		width: 30px;
		height: 30px;
		background: transparent;
		border: none;
		cursor: pointer;
		padding: 0;
	}

	.hamburger-line {
		width: 100%;
		height: 3px;
		background-color: white;
		border-radius: 2px;
		transition: all 0.3s;
	}

	/* Mobile Menu */
	.mobile-menu {
		display: none;
		background-color: rgba(102, 126, 234, 0.98);
		padding: 1rem;
		border-top: 1px solid rgba(255, 255, 255, 0.1);
	}

	.mobile-nav-link {
		display: block;
		color: white;
		text-decoration: none;
		font-weight: 500;
		font-size: 1rem;
		padding: 0.75rem 1rem;
		border-radius: 0.375rem;
		margin-bottom: 0.5rem;
		transition: background-color 0.2s;
	}

	.mobile-nav-link:hover {
		background-color: rgba(255, 255, 255, 0.1);
	}

	.mobile-nav-link.active {
		background-color: rgba(255, 255, 255, 0.2);
	}

	/* Responsive */
	.desktop-only {
		display: flex;
	}

	.mobile-only {
		display: none;
	}

	@media (min-width: 640px) {
		.brand-text {
			display: inline;
		}
	}

	@media (max-width: 768px) {
		.desktop-only {
			display: none;
		}

		.mobile-only {
			display: flex;
		}

		.mobile-menu {
			display: block;
		}

		.navbar-links {
			display: none;
		}
	}
</style>
