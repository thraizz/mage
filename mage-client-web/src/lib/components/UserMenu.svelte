<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import User from '@lucide/svelte/icons/user';
	import Layers from '@lucide/svelte/icons/layers';
	import LogOut from '@lucide/svelte/icons/log-out';

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
		onclick={toggleMenu}
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
		<ChevronDown class={`chevron ${menuOpen ? 'rotated' : ''}`} size={16} aria-hidden="true" />
	</button>

	{#if menuOpen}
		<div bind:this={menuDropdown} class="dropdown-menu" role="menu">
			<div class="dropdown-header">
				<div class="dropdown-username">{$auth.user?.username || 'Guest'}</div>
			</div>

			<div class="dropdown-divider"></div>

			<a href="/profile" class="dropdown-item" onclick={closeMenu} role="menuitem">
				<User size={16} aria-hidden="true" />
				Profile
			</a>

			<a href="/decks" class="dropdown-item" onclick={closeMenu} role="menuitem">
				<Layers size={16} aria-hidden="true" />
				My Decks
			</a>

			<div class="dropdown-divider"></div>

			<button class="dropdown-item logout" onclick={handleLogout} role="menuitem">
				<LogOut size={16} aria-hidden="true" />
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

	:global(svg.chevron) {
		transition: transform var(--transition-fast);
		color: var(--text-dim);
	}

	:global(svg.chevron.rotated) {
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

	.dropdown-item :global(svg) {
		flex-shrink: 0;
	}

	/* Responsive */
	@media (min-width: 640px) {
		.user-name {
			display: inline;
		}
	}
</style>
