<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let menuOpen = false;
	let menuButton: HTMLButtonElement;
	let menuDropdown: HTMLDivElement;

	function toggleMenu() {
		menuOpen = !menuOpen;
	}

	function closeMenu() {
		menuOpen = false;
	}

	function handleLogout() {
		auth.logout();
		goto('/login');
	}

	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (
			menuOpen &&
			menuButton &&
			menuDropdown &&
			!menuButton.contains(event.target as Node) &&
			!menuDropdown.contains(event.target as Node)
		) {
			closeMenu();
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

<div class="user-menu">
	<button
		bind:this={menuButton}
		class="user-button"
		on:click={toggleMenu}
		aria-label="User menu"
		aria-expanded={menuOpen}
	>
		<div class="user-avatar">
			{#if $auth.user}
				{$auth.user.username.charAt(0).toUpperCase()}
			{:else}
				?
			{/if}
		</div>
		<span class="user-name">
			{$auth.user?.username || 'Guest'}
		</span>
		<svg
			class="chevron"
			class:rotated={menuOpen}
			xmlns="http://www.w3.org/2000/svg"
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="6 9 12 15 18 9"></polyline>
		</svg>
	</button>

	{#if menuOpen}
		<div bind:this={menuDropdown} class="dropdown-menu" role="menu">
			<div class="dropdown-header">
				<div class="dropdown-username">{$auth.user?.username || 'Guest'}</div>
			</div>

			<div class="dropdown-divider"></div>

			<a href="/profile" class="dropdown-item" on:click={closeMenu} role="menuitem">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
					<circle cx="12" cy="7" r="4"></circle>
				</svg>
				Profile
			</a>

			<a href="/decks" class="dropdown-item" on:click={closeMenu} role="menuitem">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<rect x="2" y="7" width="20" height="15" rx="2" ry="2"></rect>
					<polyline points="17 2 12 7 7 2"></polyline>
				</svg>
				My Decks
			</a>

			<div class="dropdown-divider"></div>

			<button class="dropdown-item logout" on:click={handleLogout} role="menuitem">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="16"
					height="16"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
					<polyline points="16 17 21 12 16 7"></polyline>
					<line x1="21" y1="12" x2="9" y2="12"></line>
				</svg>
				Logout
			</button>
		</div>
	{/if}
</div>

<style>
	.user-menu {
		position: relative;
	}

	.user-button {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		background: var(--bg-iron);
		border: 1px solid var(--border-default);
		border-radius: var(--radius-md);
		color: var(--text-muted);
		cursor: pointer;
		transition: all var(--transition-fast);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
	}

	.user-button:hover {
		background: var(--bg-steel);
		border-color: var(--border-strong);
		color: var(--text-bright);
	}

	.user-avatar {
		width: 28px;
		height: 28px;
		border-radius: var(--radius-full);
		background: var(--accent-gold-dim);
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: var(--weight-bold);
		font-size: var(--text-sm);
		color: var(--bg-void);
	}

	.user-name {
		display: none;
	}

	.chevron {
		transition: transform var(--transition-fast);
		color: var(--text-dim);
	}

	.chevron.rotated {
		transform: rotate(180deg);
	}

	/* Dropdown Menu */
	.dropdown-menu {
		position: absolute;
		top: calc(100% + var(--space-2));
		right: 0;
		background: var(--bg-slate);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		box-shadow: var(--shadow-lg);
		min-width: 200px;
		z-index: var(--z-dropdown);
		animation: slideDown var(--transition-fast);
		overflow: hidden;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.dropdown-header {
		padding: var(--space-3) var(--space-4);
		border-bottom: 1px solid var(--border-subtle);
	}

	.dropdown-username {
		font-weight: var(--weight-semibold);
		color: var(--text-bright);
		font-size: var(--text-sm);
	}

	.dropdown-email {
		font-size: var(--text-xs);
		color: var(--text-dim);
		margin-top: var(--space-1);
	}

	.dropdown-divider {
		height: 1px;
		background: var(--border-subtle);
	}

	.dropdown-item {
		display: flex;
		align-items: center;
		gap: var(--space-3);
		width: 100%;
		padding: var(--space-3) var(--space-4);
		background: none;
		border: none;
		text-align: left;
		color: var(--text-muted);
		font-size: var(--text-sm);
		font-weight: var(--weight-medium);
		cursor: pointer;
		transition: all var(--transition-fast);
		text-decoration: none;
	}

	.dropdown-item:hover {
		background: var(--bg-iron);
		color: var(--text-bright);
	}

	.dropdown-item.logout {
		color: var(--status-error);
	}

	.dropdown-item.logout:hover {
		background: var(--status-error-dim);
	}

	.dropdown-item svg {
		flex-shrink: 0;
	}

	/* Responsive */
	@media (min-width: 640px) {
		.user-name {
			display: inline;
		}
	}
</style>
