<script lang="ts">
	import { page } from '$app/stores';
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
			<a href="/lobby" class="nav-link" class:active={$page.url.pathname === '/lobby'}> Lobby </a>
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
		background: var(--bg-obsidian);
		border-bottom: 1px solid var(--border-subtle);
		position: sticky;
		top: 0;
		z-index: var(--z-sticky);
	}

	.navbar-container {
		max-width: 1280px;
		margin: 0 auto;
		padding: 0 var(--space-4);
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: var(--navbar-height);
	}

	/* Brand */
	.navbar-brand {
		display: flex;
		align-items: center;
	}

	.brand-link {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		text-decoration: none;
		color: var(--accent-gold);
		font-family: var(--font-display);
		font-weight: var(--weight-bold);
		font-size: var(--text-xl);
		transition: color var(--transition-fast);
	}

	.brand-link:hover {
		color: var(--accent-gold-bright);
	}

	.brand-icon {
		font-size: 1.75rem;
	}

	.brand-text {
		display: none;
	}

	/* Navigation Links */
	.navbar-links {
		display: flex;
		align-items: center;
		gap: var(--space-1);
		flex: 1;
		justify-content: center;
	}

	.nav-link {
		color: var(--text-muted);
		text-decoration: none;
		font-weight: var(--weight-medium);
		font-size: var(--text-sm);
		padding: var(--space-2) var(--space-4);
		border-radius: var(--radius-md);
		transition: all var(--transition-fast);
		position: relative;
	}

	.nav-link:hover {
		color: var(--text-bright);
		background: var(--bg-iron);
	}

	.nav-link.active {
		color: var(--accent-gold);
		background: var(--bg-iron);
	}

	.nav-link.active::after {
		content: '';
		position: absolute;
		bottom: -1px;
		left: var(--space-4);
		right: var(--space-4);
		height: 2px;
		background: var(--accent-gold);
		border-radius: var(--radius-full);
	}

	/* Right Side */
	.navbar-right {
		display: flex;
		align-items: center;
		gap: var(--space-3);
	}

	/* Mobile Menu Button */
	.mobile-menu-button {
		display: none;
		flex-direction: column;
		justify-content: space-around;
		width: 28px;
		height: 28px;
		background: transparent;
		border: none;
		cursor: pointer;
		padding: 0;
	}

	.hamburger-line {
		width: 100%;
		height: 2px;
		background: var(--text-muted);
		border-radius: var(--radius-full);
		transition: all var(--transition-base);
	}

	.mobile-menu-button:hover .hamburger-line {
		background: var(--accent-gold);
	}

	/* Mobile Menu */
	.mobile-menu {
		display: none;
		background: var(--bg-slate);
		padding: var(--space-3);
		border-top: 1px solid var(--border-subtle);
	}

	.mobile-nav-link {
		display: block;
		color: var(--text-muted);
		text-decoration: none;
		font-weight: var(--weight-medium);
		font-size: var(--text-base);
		padding: var(--space-3) var(--space-4);
		border-radius: var(--radius-md);
		margin-bottom: var(--space-1);
		transition: all var(--transition-fast);
	}

	.mobile-nav-link:hover {
		color: var(--text-bright);
		background: var(--bg-iron);
	}

	.mobile-nav-link.active {
		color: var(--accent-gold);
		background: var(--bg-iron);
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
