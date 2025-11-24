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
				{#if $auth.user?.email}
					<div class="dropdown-email">{$auth.user.email}</div>
				{/if}
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
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		background-color: rgba(255, 255, 255, 0.1);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 0.5rem;
		color: white;
		cursor: pointer;
		transition: all 0.2s;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.user-button:hover {
		background-color: rgba(255, 255, 255, 0.15);
	}

	.user-avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
		font-size: 1rem;
		color: white;
	}

	.user-name {
		display: none;
	}

	.chevron {
		transition: transform 0.2s;
		color: white;
	}

	.chevron.rotated {
		transform: rotate(180deg);
	}

	/* Dropdown Menu */
	.dropdown-menu {
		position: absolute;
		top: calc(100% + 0.5rem);
		right: 0;
		background-color: white;
		border-radius: 0.5rem;
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
		min-width: 200px;
		z-index: 1000;
		animation: slideDown 0.2s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.dropdown-header {
		padding: 0.75rem 1rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.dropdown-username {
		font-weight: 600;
		color: #1f2937;
		font-size: 0.875rem;
	}

	.dropdown-email {
		font-size: 0.75rem;
		color: #6b7280;
		margin-top: 0.25rem;
	}

	.dropdown-divider {
		height: 1px;
		background-color: #e5e7eb;
		margin: 0.25rem 0;
	}

	.dropdown-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		width: 100%;
		padding: 0.75rem 1rem;
		background: none;
		border: none;
		text-align: left;
		color: #374151;
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
		transition: background-color 0.2s;
		text-decoration: none;
	}

	.dropdown-item:hover {
		background-color: #f3f4f6;
	}

	.dropdown-item.logout {
		color: #dc2626;
	}

	.dropdown-item.logout:hover {
		background-color: #fee2e2;
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
