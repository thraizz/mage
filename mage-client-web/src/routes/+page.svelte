<script lang="ts">
	import { toast } from '$lib/stores/toast';
	import { confirm } from '$lib/stores/confirm';
	import Modal from '$lib/components/Modal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';

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

	<!-- Toast Testing Section -->
	<div class="toast-test">
		<h3>Toast Notification Test</h3>
		<div class="test-buttons">
			<button class="test-btn success" on:click={testSuccess}>Success Toast</button>
			<button class="test-btn error" on:click={testError}>Error Toast</button>
			<button class="test-btn warning" on:click={testWarning}>Warning Toast</button>
			<button class="test-btn info" on:click={testInfo}>Info Toast</button>
			<button class="test-btn multiple" on:click={testMultiple}>Multiple Toasts</button>
		</div>
	</div>

	<!-- Modal Testing Section -->
	<div class="modal-test">
		<h3>Modal Dialog Test</h3>
		<div class="test-buttons">
			<button class="test-btn info" on:click={() => (showBasicModal = true)}> Basic Modal </button>
			<button class="test-btn success" on:click={() => (showSmallModal = true)}>
				Small Modal
			</button>
			<button class="test-btn warning" on:click={() => (showLargeModal = true)}>
				Large Modal
			</button>
			<button class="test-btn error" on:click={() => (showNoBackdropModal = true)}>
				No Backdrop Close
			</button>
			<button class="test-btn multiple" on:click={() => (showNoCloseModal = true)}>
				No Close Button
			</button>
		</div>
	</div>

	<!-- Confirmation Dialog Testing Section -->
	<div class="confirm-test">
		<h3>Confirmation Dialog Test</h3>
		<div class="test-buttons">
			<button class="test-btn info" on:click={testBasicConfirm}>Basic Confirm</button>
			<button class="test-btn error" on:click={testDestructiveConfirm}>
				Destructive Confirm
			</button>
			<button class="test-btn success" on:click={testCustomConfirm}>Custom Text</button>
			<button class="test-btn multiple" on:click={() => (showComponentConfirm = true)}>
				Component-based
			</button>
		</div>
	</div>

	<!-- Loading Spinner Testing Section -->
	<div class="loading-test">
		<h3>Loading Spinner Test</h3>
		<div class="test-buttons">
			<button class="test-btn info" on:click={testOverlaySpinner}>Overlay Spinner</button>
			<button class="test-btn success" on:click={testOverlaySpinnerWithLabel}> With Label </button>
			<button class="test-btn warning" on:click={() => (showInlineSpinners = !showInlineSpinners)}>
				Toggle Inline
			</button>
		</div>

		{#if showInlineSpinners}
			<div class="inline-spinner-demo">
				<div class="spinner-row">
					<div class="spinner-item">
						<p>Small:</p>
						<LoadingSpinner size="small" />
					</div>
					<div class="spinner-item">
						<p>Medium:</p>
						<LoadingSpinner size="medium" />
					</div>
					<div class="spinner-item">
						<p>Large:</p>
						<LoadingSpinner size="large" />
					</div>
				</div>
				<div class="spinner-row">
					<div class="spinner-item">
						<p>With Label:</p>
						<LoadingSpinner size="medium" label="Loading data..." />
					</div>
					<div class="spinner-item">
						<p>Custom Color:</p>
						<LoadingSpinner size="medium" color="#10b981" label="Processing..." />
					</div>
				</div>
			</div>
		{/if}
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

	/* Testing Sections */
	.toast-test,
	.modal-test,
	.confirm-test,
	.loading-test {
		max-width: 800px;
		margin: 3rem auto 0;
		background: rgba(255, 255, 255, 0.95);
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.modal-test,
	.confirm-test,
	.loading-test {
		margin-top: 2rem;
	}

	.toast-test h3,
	.modal-test h3,
	.confirm-test h3,
	.loading-test h3 {
		margin: 0 0 1.5rem 0;
		color: #667eea;
		text-align: center;
		font-size: 1.5rem;
	}

	.test-buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 0.75rem;
		justify-content: center;
	}

	.test-btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		color: white;
	}

	.test-btn.success {
		background-color: #10b981;
	}

	.test-btn.success:hover {
		background-color: #059669;
	}

	.test-btn.error {
		background-color: #ef4444;
	}

	.test-btn.error:hover {
		background-color: #dc2626;
	}

	.test-btn.warning {
		background-color: #f59e0b;
	}

	.test-btn.warning:hover {
		background-color: #d97706;
	}

	.test-btn.info {
		background-color: #3b82f6;
	}

	.test-btn.info:hover {
		background-color: #2563eb;
	}

	.test-btn.multiple {
		background-color: #8b5cf6;
	}

	.test-btn.multiple:hover {
		background-color: #7c3aed;
	}

	/* Modal Button Styles */
	:global(.btn-primary),
	:global(.btn-secondary) {
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: none;
	}

	:global(.btn-primary) {
		background-color: #667eea;
		color: white;
	}

	:global(.btn-primary:hover) {
		background-color: #5568d3;
	}

	:global(.btn-secondary) {
		background-color: #e5e7eb;
		color: #374151;
	}

	:global(.btn-secondary:hover) {
		background-color: #d1d5db;
	}

	:global(.modal-content p) {
		margin: 0 0 1rem 0;
		line-height: 1.5;
		color: #374151;
	}

	:global(.modal-content ul) {
		margin: 0 0 1rem 0;
		padding-left: 1.5rem;
		color: #374151;
	}

	:global(.modal-content li) {
		margin-bottom: 0.5rem;
	}

	/* Inline Spinner Demo */
	.inline-spinner-demo {
		margin-top: 2rem;
		padding: 2rem;
		background: #f9fafb;
		border-radius: 0.5rem;
	}

	.spinner-row {
		display: flex;
		gap: 2rem;
		margin-bottom: 2rem;
		flex-wrap: wrap;
		justify-content: center;
	}

	.spinner-row:last-child {
		margin-bottom: 0;
	}

	.spinner-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
		padding: 1rem;
		background: white;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		min-width: 120px;
	}

	.spinner-item p {
		margin: 0;
		font-weight: 600;
		color: #374151;
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

		.test-buttons {
			flex-direction: column;
		}

		.test-btn {
			width: 100%;
		}
	}
</style>
