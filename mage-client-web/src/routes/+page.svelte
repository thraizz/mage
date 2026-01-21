<script lang="ts">
	export const ssr = false;

	import { toast } from '$lib/stores/toast';
	import { confirm } from '$lib/stores/confirm';
	import Modal from '$lib/components/Modal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

	// Available background images (shared with auth pages)
	const backgroundImages = ['Boros.jpg', 'Golgari.jpg', 'Gruul.jpg', 'Izzet.jpg', 'Rakdos.jpg'];
	let selectedBackground = backgroundImages[Math.floor(Math.random() * backgroundImages.length)];

	// Modal state
	let showBasicModal = false;
	let showSmallModal = false;
	let showLargeModal = false;
	let showNoBackdropModal = false;
	let showNoCloseModal = false;

	// Confirm dialog state (for component-based usage)
	let showComponentConfirm = false;

	// Loading spinner state
	let showOverlaySpinner = false;
	let showInlineSpinners = false;

	// Test toast notifications
	function testSuccess() {
		toast.success('Successfully connected to server!');
	}

	function testError() {
		toast.error('Failed to connect to server. Please try again.');
	}

	function testWarning() {
		toast.warning('Your session will expire in 5 minutes.');
	}

	function testInfo() {
		toast.info('New features have been added. Check them out!');
	}

	function testMultiple() {
		toast.success('First notification');
		setTimeout(() => toast.info('Second notification'), 500);
		setTimeout(() => toast.warning('Third notification'), 1000);
		setTimeout(() => toast.error('Fourth notification'), 1500);
	}

	// Modal handlers
	function handleModalAction() {
		toast.success('Action confirmed!');
		showBasicModal = false;
		showSmallModal = false;
		showLargeModal = false;
		showNoBackdropModal = false;
		showNoCloseModal = false;
	}

	// Confirmation dialog tests
	async function testBasicConfirm() {
		const result = await confirm.confirm({
			title: 'Delete Item',
			message: 'Are you sure you want to delete this item?'
		});

		if (result) {
			toast.success('Item deleted!');
		} else {
			toast.info('Deletion cancelled');
		}
	}

	async function testDestructiveConfirm() {
		const result = await confirm.confirm({
			title: 'Delete Account',
			message: 'This action cannot be undone. Are you sure you want to delete your account?',
			confirmText: 'Delete Account',
			cancelText: 'Keep Account',
			destructive: true
		});

		if (result) {
			toast.error('Account deleted');
		} else {
			toast.info('Account deletion cancelled');
		}
	}

	async function testCustomConfirm() {
		const result = await confirm.confirm({
			title: 'Save Changes',
			message: 'Do you want to save your changes before leaving?',
			confirmText: 'Save',
			cancelText: "Don't Save"
		});

		if (result) {
			toast.success('Changes saved!');
		} else {
			toast.warning('Changes discarded');
		}
	}

	function handleComponentConfirm() {
		showComponentConfirm = false;
		toast.success('Action confirmed via component!');
	}

	function handleComponentCancel() {
		showComponentConfirm = false;
		toast.info('Action cancelled via component');
	}

	// Loading spinner tests
	function testOverlaySpinner() {
		showOverlaySpinner = true;
		setTimeout(() => {
			showOverlaySpinner = false;
			toast.success('Loading complete!');
		}, 3000);
	}

	function testOverlaySpinnerWithLabel() {
		showOverlaySpinner = true;
		setTimeout(() => {
			showOverlaySpinner = false;
			toast.success('Data loaded!');
		}, 3000);
	}
</script>

<svelte:head>
	<title>MAGE - Magic: The Gathering Online</title>
</svelte:head>

<div class="container" style="background-image: url('/images/{selectedBackground}')">
	<div class="card page-card">
		<div class="hero">
			<h1>MAGE</h1>
			<p class="flavor-text">Magic: The Gathering Online — play, compete, collect.</p>
		</div>

		<nav class="nav-grid" aria-label="Primary navigation">
			<a href="/login" class="nav-card">
				<div class="icon" aria-hidden="true">🔐</div>
				<h2>Login</h2>
				<p>Sign in to your account</p>
			</a>

			<a href="/register" class="nav-card">
				<div class="icon" aria-hidden="true">📝</div>
				<h2>Register</h2>
				<p>Create a new account</p>
			</a>

			<a href="/lobby" class="nav-card">
				<div class="icon" aria-hidden="true">🎯</div>
				<h2>Lobby</h2>
				<p>Find and join games</p>
			</a>

			<a href="/decks" class="nav-card">
				<div class="icon" aria-hidden="true">🎴</div>
				<h2>My Decks</h2>
				<p>Manage your deck collection</p>
			</a>

			<a href="/profile" class="nav-card">
				<div class="icon" aria-hidden="true">👤</div>
				<h2>Profile</h2>
				<p>View your stats and settings</p>
			</a>

			<a href="/table/demo" class="nav-card">
				<div class="icon" aria-hidden="true">🪑</div>
				<h2>Table (Demo)</h2>
				<p>Pre-game table lobby</p>
			</a>

			<a href="/game/demo" class="nav-card">
				<div class="icon" aria-hidden="true">⚔️</div>
				<h2>Game (Demo)</h2>
				<p>Active game view</p>
			</a>
		</nav>

		<div class="footer">
			<p class="flavor-text">All routes are placeholder pages — full functionality coming soon.</p>
		</div>

		<!-- Toast Testing Section -->
		<div class="test-panel toast-test">
			<h3>Toast Notification Test</h3>
			<div class="test-buttons">
				<button class="btn-secondary test-btn" data-tone="success" on:click={testSuccess}>
					Success Toast
				</button>
				<button class="btn-secondary test-btn" data-tone="error" on:click={testError}
					>Error Toast</button
				>
				<button class="btn-secondary test-btn" data-tone="warning" on:click={testWarning}>
					Warning Toast
				</button>
				<button class="btn-secondary test-btn" data-tone="info" on:click={testInfo}
					>Info Toast</button
				>
				<button class="btn-secondary test-btn" data-tone="multiple" on:click={testMultiple}>
					Multiple Toasts
				</button>
			</div>
		</div>

		<!-- Modal Testing Section -->
		<div class="test-panel modal-test">
			<h3>Modal Dialog Test</h3>
			<div class="test-buttons">
				<button
					class="btn-secondary test-btn"
					data-tone="info"
					on:click={() => (showBasicModal = true)}
				>
					Basic Modal
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="success"
					on:click={() => (showSmallModal = true)}
				>
					Small Modal
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="warning"
					on:click={() => (showLargeModal = true)}
				>
					Large Modal
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="error"
					on:click={() => (showNoBackdropModal = true)}
				>
					No Backdrop Close
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="multiple"
					on:click={() => (showNoCloseModal = true)}
				>
					No Close Button
				</button>
			</div>
		</div>

		<!-- Confirmation Dialog Testing Section -->
		<div class="test-panel confirm-test">
			<h3>Confirmation Dialog Test</h3>
			<div class="test-buttons">
				<button class="btn-secondary test-btn" data-tone="info" on:click={testBasicConfirm}>
					Basic Confirm
				</button>
				<button class="btn-secondary test-btn" data-tone="error" on:click={testDestructiveConfirm}>
					Destructive Confirm
				</button>
				<button class="btn-secondary test-btn" data-tone="success" on:click={testCustomConfirm}>
					Custom Text
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="multiple"
					on:click={() => (showComponentConfirm = true)}
				>
					Component-based
				</button>
			</div>
		</div>

		<!-- Loading Spinner Testing Section -->
		<div class="test-panel loading-test">
			<h3>Loading Spinner Test</h3>
			<div class="test-buttons">
				<button class="btn-secondary test-btn" data-tone="info" on:click={testOverlaySpinner}>
					Overlay Spinner
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="success"
					on:click={testOverlaySpinnerWithLabel}
				>
					With Label
				</button>
				<button
					class="btn-secondary test-btn"
					data-tone="warning"
					on:click={() => (showInlineSpinners = !showInlineSpinners)}
				>
					Toggle Inline
				</button>
			</div>

			{#if showInlineSpinners}
				<div class="inline-spinner-demo">
					<div class="spinner-row">
						<div class="spinner-item">
							<p>Small</p>
							<LoadingSpinner size="small" />
						</div>
						<div class="spinner-item">
							<p>Medium</p>
							<LoadingSpinner size="medium" />
						</div>
						<div class="spinner-item">
							<p>Large</p>
							<LoadingSpinner size="large" />
						</div>
					</div>
					<div class="spinner-row">
						<div class="spinner-item">
							<p>With Label</p>
							<LoadingSpinner size="medium" label="Loading data..." />
						</div>
						<div class="spinner-item">
							<p>Custom Color</p>
							<LoadingSpinner size="medium" color="#10b981" label="Processing..." />
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Modal Components -->
<Modal bind:open={showBasicModal} title="Basic Modal" size="medium">
	<p>This is a basic modal dialog with a title and close button.</p>
	<p>You can close it by:</p>
	<ul>
		<li>Clicking the X button</li>
		<li>Clicking outside the modal</li>
		<li>Pressing the ESC key</li>
	</ul>

	<div slot="footer">
		<button class="btn-secondary" on:click={() => (showBasicModal = false)}>Cancel</button>
		<button class="btn-primary" on:click={handleModalAction}>Confirm</button>
	</div>
</Modal>

<Modal bind:open={showSmallModal} title="Small Modal" size="small">
	<p>This is a small modal, perfect for quick confirmations or short messages.</p>

	<div slot="footer">
		<button class="btn-primary" on:click={() => (showSmallModal = false)}>Got it</button>
	</div>
</Modal>

<Modal bind:open={showLargeModal} title="Large Modal" size="large">
	<p>This is a large modal with more content space.</p>
	<p>
		Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut
		labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
		laboris nisi ut aliquip ex ea commodo consequat.
	</p>
	<p>
		Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
		pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
		mollit anim id est laborum.
	</p>

	<div slot="footer">
		<button class="btn-secondary" on:click={() => (showLargeModal = false)}>Close</button>
	</div>
</Modal>

<Modal bind:open={showNoBackdropModal} title="No Backdrop Close" closeOnBackdrop={false}>
	<p>This modal cannot be closed by clicking the backdrop.</p>
	<p>You must use the close button or press ESC.</p>

	<div slot="footer">
		<button class="btn-primary" on:click={() => (showNoBackdropModal = false)}>Close</button>
	</div>
</Modal>

<Modal bind:open={showNoCloseModal} title="Controlled Closing" showCloseButton={false}>
	<p>This modal doesn't have a close button in the header.</p>
	<p>Use the button below to close it.</p>

	<div slot="footer">
		<button class="btn-primary" on:click={() => (showNoCloseModal = false)}> Close Modal </button>
	</div>
</Modal>

<!-- Confirmation Dialog Component Example -->
<ConfirmDialog
	bind:open={showComponentConfirm}
	title="Component Confirmation"
	message="This is using the ConfirmDialog component directly. Do you want to proceed?"
	confirmText="Yes, Proceed"
	cancelText="No, Cancel"
	onConfirm={handleComponentConfirm}
	onCancel={handleComponentCancel}
/>

<!-- Loading Spinner Overlay Example -->
{#if showOverlaySpinner}
	<LoadingSpinner overlay={true} size="large" label="Loading..." />
{/if}

<style>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		padding: 2rem;
		position: relative;
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
	}

	.container::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(11, 12, 16, 0.7);
		z-index: 0;
	}

	.page-card {
		width: 100%;
		max-width: 1100px;
		position: relative;
		z-index: 1;
	}

	.hero {
		text-align: center;
		margin-bottom: 2rem;
	}

	h1 {
		font-size: 2.25rem;
		margin: 0 0 0.75rem 0;
	}

	.hero :global(.flavor-text) {
		margin: 0;
		font-size: 1rem;
		text-align: center;
	}

	.nav-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1rem;
		margin-top: 1.25rem;
	}

	.nav-card {
		background: rgba(166, 154, 168, 0.12);
		border-radius: var(--radius-xl);
		padding: 1.25rem 1.25rem;
		text-align: center;
		text-decoration: none;
		color: var(--ci-scroll-parchment);
		transition: all var(--transition-base);
		border: 1px solid rgba(59, 130, 246, 0.18);
		backdrop-filter: blur(6px);
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
	}

	.nav-card:hover {
		transform: translateY(-2px);
		border-color: rgba(59, 130, 246, 0.35);
		box-shadow:
			0 12px 28px rgba(0, 0, 0, 0.45),
			0 0 0 1px rgba(59, 130, 246, 0.12);
	}

	.icon {
		font-size: 2.25rem;
		margin-bottom: 0.75rem;
	}

	.nav-card h2 {
		margin: 0 0 0.35rem 0;
		font-size: 1.25rem;
	}

	.nav-card p {
		margin: 0;
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.9375rem;
		line-height: 1.4;
	}

	.footer {
		text-align: center;
		margin-top: 1.75rem;
	}

	.footer :global(.flavor-text) {
		margin: 0;
		text-align: center;
	}

	/* Testing Sections */
	.toast-test,
	.modal-test,
	.confirm-test,
	.loading-test {
		margin-top: 1.25rem;
	}

	.test-panel {
		background: rgba(11, 12, 16, 0.55);
		border: 1px solid rgba(59, 130, 246, 0.18);
		border-radius: var(--radius-xl);
		padding: 1.25rem;
	}

	.toast-test h3,
	.modal-test h3,
	.confirm-test h3,
	.loading-test h3 {
		margin: 0 0 1rem 0;
		text-align: left;
		font-size: 1.125rem;
	}

	.test-buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	.test-btn {
		padding: 0.625rem 1rem;
		font-size: 0.875rem;
		text-transform: none;
		letter-spacing: 0.02em;
	}

	.test-btn[data-tone='success'] {
		border-color: rgba(81, 207, 102, 0.35);
	}

	.test-btn[data-tone='error'] {
		border-color: rgba(255, 77, 77, 0.45);
	}

	.test-btn[data-tone='warning'] {
		border-color: rgba(245, 158, 11, 0.45);
	}

	.test-btn[data-tone='info'] {
		border-color: rgba(59, 130, 246, 0.45);
	}

	.test-btn[data-tone='multiple'] {
		border-color: rgba(139, 92, 246, 0.45);
	}

	:global(.modal-content p) {
		margin: 0 0 1rem 0;
		line-height: 1.5;
		color: var(--ci-scroll-parchment);
	}

	:global(.modal-content ul) {
		margin: 0 0 1rem 0;
		padding-left: 1.5rem;
		color: var(--ci-scroll-parchment);
	}

	:global(.modal-content li) {
		margin-bottom: 0.5rem;
	}

	/* Inline Spinner Demo */
	.inline-spinner-demo {
		margin-top: 1rem;
		padding: 1rem;
		background: rgba(166, 154, 168, 0.08);
		border: 1px solid rgba(59, 130, 246, 0.15);
		border-radius: var(--radius-xl);
	}

	.spinner-row {
		display: flex;
		gap: 1rem;
		margin-bottom: 1rem;
		flex-wrap: wrap;
		justify-content: flex-start;
	}

	.spinner-row:last-child {
		margin-bottom: 0;
	}

	.spinner-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		background: rgba(11, 12, 16, 0.55);
		border: 1px solid rgba(59, 130, 246, 0.18);
		border-radius: var(--radius-xl);
		min-width: 120px;
	}

	.spinner-item p {
		margin: 0;
		font-weight: 600;
		color: var(--ci-scroll-parchment);
		font-size: 0.875rem;
	}

	@media (max-width: 768px) {
		h1 {
			font-size: 2rem;
		}

		.nav-grid {
			grid-template-columns: 1fr;
		}

		.test-buttons {
			flex-direction: column;
		}

		.test-btn {
			width: 100%;
		}
	}
</style>
