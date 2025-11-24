<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import { websocketStore } from '$lib/stores/websocket';
	import type { Table, GameFormat } from '$lib/types/table';
	import type { OnlinePlayer } from '$lib/types/player';
	import { fetchTables, fetchOnlinePlayers, getGameFormats } from '$lib/api/lobby';
	import {
		subscribeLobbyUpdates,
		connectLobbyUpdates,
		disconnectLobbyUpdates,
		type TableUpdateEvent
	} from '$lib/services/lobby-updates';
	import { getMageClient } from '$lib/grpc/client';
	import TableCard from '$lib/components/TableCard.svelte';
	import LoadingSpinner from '$lib/components/LoadingSpinner.svelte';
	import CreateTableModal from '$lib/components/CreateTableModal.svelte';
	import JoinTableModal from '$lib/components/JoinTableModal.svelte';
	import OnlinePlayersList from '$lib/components/OnlinePlayersList.svelte';
	import LobbyChat from '$lib/components/LobbyChat.svelte';

	// State
	let tables = $state<Table[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	// Online players state
	let onlinePlayers = $state<OnlinePlayer[]>([]);
	let playersListOpen = $state(true);

	// Filter state
	let searchQuery = $state('');
	let selectedFormat = $state<GameFormat | 'All'>('All');
	let openOnly = $state(false);

	// Modal state
	let showCreateModal = $state(false);
	let showJoinModal = $state(false);
	let joiningTable = $state<Table | null>(null);

	// Available formats
	const formats = getGameFormats();

	// WebSocket connection state
	let wsState = $derived($websocketStore.state);

	/**
	 * Load tables from API
	 */
	async function loadTables(): Promise<void> {
		loading = true;
		error = null;

		try {
			tables = await fetchTables();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load tables';
			console.error('Failed to load tables:', err);
		} finally {
			loading = false;
		}
	}

	/**
	 * Load online players
	 */
	async function loadOnlinePlayers(): Promise<void> {
		try {
			const currentUsername = $auth.user?.username;
			onlinePlayers = await fetchOnlinePlayers(currentUsername);
		} catch (err) {
			console.error('Failed to load online players:', err);
		}
	}

	/**
	 * Filter tables based on current filters
	 */
	const filteredTables = $derived(() => {
		let result = tables;

		// Filter by format
		if (selectedFormat !== 'All') {
			result = result.filter((table) => table.format === selectedFormat);
		}

		// Filter by search query
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase().trim();
			result = result.filter(
				(table) =>
					table.name.toLowerCase().includes(query) ||
					table.hostUsername.toLowerCase().includes(query) ||
					table.format.toLowerCase().includes(query)
			);
		}

		// Filter by open only
		if (openOnly) {
			result = result.filter(
				(table) => table.status === 'waiting' && table.players.length < table.maxPlayers
			);
		}

		return result;
	});

	/**
	 * Handle table click - show join modal
	 */
	function handleTableClick(table: Table): void {
		joiningTable = table;
		showJoinModal = true;
	}

	/**
	 * Handle successful table join
	 */
	function handleTableJoined(tableId: string): void {
		showJoinModal = false;
		joiningTable = null;
		// Navigate to table page
		window.location.href = `/table/${tableId}`;
	}

	/**
	 * Refresh tables list
	 */
	async function handleRefresh(): Promise<void> {
		await loadTables();
	}

	/**
	 * Clear all filters
	 */
	function clearFilters(): void {
		searchQuery = '';
		selectedFormat = 'All';
		openOnly = false;
	}

	/**
	 * Open create table modal
	 */
	function openCreateModal(): void {
		showCreateModal = true;
	}

	/**
	 * Close create table modal
	 */
	function closeCreateModal(): void {
		showCreateModal = false;
	}

	/**
	 * Handle successful table creation
	 */
	function handleTableCreated(tableId: string): void {
		// Refresh table list
		loadTables();
		// Navigate to the new table
		window.location.href = `/table/${tableId}`;
	}

	// WebSocket subscription cleanup function
	let unsubscribeLobby: (() => void) | null = null;

	/**
	 * Handle table updates from WebSocket
	 */
	function handleTableUpdate(event: TableUpdateEvent): void {
		console.log('[Lobby] Table update:', event.type, event.table.id);

		switch (event.type) {
			case 'created':
			case 'updated': {
				// Find existing table
				const existingIndex = tables.findIndex((t) => t.id === event.table.id);

				if (existingIndex >= 0) {
					// Update existing table
					tables[existingIndex] = event.table;
					// Trigger reactivity
					tables = [...tables];
				} else {
					// Add new table with animation
					tables = [...tables, event.table];
				}
				break;
			}
			case 'deleted': {
				// Remove table
				tables = tables.filter((t) => t.id !== event.table.id);
				break;
			}
		}
	}

	/**
	 * Connect to WebSocket for real-time updates
	 */
	async function connectWebSocket(): Promise<void> {
		try {
			const client = getMageClient();
			const sessionId = client.getSessionId();

			if (!sessionId) {
				console.warn('[Lobby] No session ID available for WebSocket connection');
				return;
			}

			// Connect to WebSocket
			await connectLobbyUpdates(sessionId);

			// Subscribe to table updates
			unsubscribeLobby = subscribeLobbyUpdates(handleTableUpdate);

			console.log('[Lobby] WebSocket connected and subscribed to updates');
		} catch (err) {
			console.error('[Lobby] Failed to connect WebSocket:', err);
		}
	}

	// Load tables and online players on mount
	// Wait for auth to be ready before making API calls
	onMount(() => {
		// Wait for auth to be authenticated and sessionId to be available
		const checkAuth = () => {
			if ($auth.isAuthenticated) {
				// Give a small delay to ensure sessionId is restored
				setTimeout(async () => {
					// Load initial data
					await loadTables();
					await loadOnlinePlayers();

					// Connect WebSocket for real-time updates
					await connectWebSocket();
				}, 100);
			} else {
				// Wait a bit and check again
				setTimeout(checkAuth, 100);
			}
		};

		checkAuth();
	});

	// Cleanup on unmount
	onDestroy(() => {
		// Unsubscribe from lobby updates
		if (unsubscribeLobby) {
			unsubscribeLobby();
			unsubscribeLobby = null;
		}

		// Disconnect WebSocket
		disconnectLobbyUpdates();

		console.log('[Lobby] Cleaned up WebSocket subscriptions');
	});
</script>

<svelte:head>
	<title>Lobby - Mage</title>
</svelte:head>

<div class="lobby-page">
	<!-- Header -->
	<div class="lobby-header">
		<div class="header-content">
			<div class="header-left">
				<h1 class="page-title">Game Lobby</h1>
				{#if !loading && tables.length > 0}
					<span class="table-count">
						{filteredTables().length}
						{#if filteredTables().length !== tables.length}
							<span class="count-divider">/</span>
							<span class="total-count">{tables.length}</span>
						{/if}
						{filteredTables().length === 1 ? 'table' : 'tables'}
					</span>
				{/if}
				{#if wsState === 'connected'}
					<span class="ws-status ws-connected" title="Real-time updates active">
						<span class="ws-dot"></span>
						Live
					</span>
				{:else if wsState === 'reconnecting'}
					<span class="ws-status ws-reconnecting" title="Reconnecting...">
						<span class="ws-dot"></span>
						Reconnecting
					</span>
				{/if}
			</div>

			<div class="header-right">
				<button
					class="refresh-button"
					onclick={handleRefresh}
					disabled={loading}
					title="Refresh tables"
				>
					<svg
						class="refresh-icon"
						class:spinning={loading}
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path
							d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"
						/>
					</svg>
					<span>Refresh</span>
				</button>

				<button class="create-button" onclick={openCreateModal} title="Create new table">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="20"
						height="20"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<line x1="12" y1="5" x2="12" y2="19"></line>
						<line x1="5" y1="12" x2="19" y2="12"></line>
					</svg>
					<span>Create Table</span>
				</button>
			</div>
		</div>

		<!-- Filters Bar -->
		{#if !loading}
			<div class="filters-bar">
				<!-- Search Input -->
				<div class="search-input-wrapper">
					<svg
						class="search-icon"
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="11" cy="11" r="8"></circle>
						<path d="m21 21-4.35-4.35"></path>
					</svg>
					<input
						type="text"
						class="search-input"
						placeholder="Search tables..."
						bind:value={searchQuery}
					/>
					{#if searchQuery}
						<button class="clear-search" onclick={() => (searchQuery = '')} title="Clear search">
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
								<line x1="18" y1="6" x2="6" y2="18"></line>
								<line x1="6" y1="6" x2="18" y2="18"></line>
							</svg>
						</button>
					{/if}
				</div>

				<!-- Format Filter -->
				<select class="format-select" bind:value={selectedFormat}>
					<option value="All">All Formats</option>
					{#each formats as format}
						<option value={format}>{format}</option>
					{/each}
				</select>

				<!-- Open Only Toggle -->
				<label class="checkbox-label">
					<input type="checkbox" class="checkbox-input" bind:checked={openOnly} />
					<span class="checkbox-text">Open Only</span>
				</label>

				<!-- Clear Filters Button -->
				{#if searchQuery || selectedFormat !== 'All' || openOnly}
					<button class="clear-filters-button" onclick={clearFilters}> Clear Filters </button>
				{/if}
			</div>
		{/if}
	</div>

	<!-- Main Content with Sidebar -->
	<div class="lobby-main">
		<!-- Sidebar: Online Players -->
		<aside class="sidebar sidebar-left">
			<OnlinePlayersList
				bind:players={onlinePlayers}
				bind:isOpen={playersListOpen}
				currentUsername={$auth.user?.username || ''}
			/>
		</aside>

		<!-- Content -->
		<div class="lobby-content">
			{#if loading}
				<!-- Loading State -->
				<div class="loading-container">
					<LoadingSpinner size="large" />
					<p class="loading-text">Loading tables...</p>
				</div>
			{:else if error}
				<!-- Error State -->
				<div class="error-container">
					<div class="error-icon">⚠️</div>
					<h2 class="error-title">Failed to Load Tables</h2>
					<p class="error-message">{error}</p>
					<button class="retry-button" onclick={handleRefresh}>Try Again</button>
				</div>
			{:else if tables.length === 0}
				<!-- Empty State - No Tables -->
				<div class="empty-container">
					<div class="empty-icon">🎮</div>
					<h2 class="empty-title">No Active Tables</h2>
					<p class="empty-message">
						There are no active games right now. Be the first to create one!
					</p>
					<button class="create-table-button" onclick={openCreateModal}>Create a Table</button>
				</div>
			{:else if filteredTables().length === 0}
				<!-- Empty State - No Results -->
				<div class="empty-container">
					<div class="empty-icon">🔍</div>
					<h2 class="empty-title">No Tables Found</h2>
					<p class="empty-message">
						No tables match your current filters. Try adjusting your search criteria.
					</p>
					<button class="clear-filters-button" onclick={clearFilters}>Clear Filters</button>
				</div>
			{:else}
				<!-- Tables Grid -->
				<div class="tables-grid">
					{#each filteredTables() as table (table.id)}
						<TableCard {table} onClick={handleTableClick} />
					{/each}
				</div>
			{/if}
		</div>

		<!-- Sidebar: Chat -->
		<aside class="sidebar sidebar-right">
			<LobbyChat />
		</aside>
	</div>
</div>

<!-- Create Table Modal -->
<CreateTableModal
	bind:open={showCreateModal}
	onClose={closeCreateModal}
	onSuccess={handleTableCreated}
/>

<!-- Join Table Modal -->
<JoinTableModal bind:open={showJoinModal} table={joiningTable} onSuccess={handleTableJoined} />

<style>
	.lobby-page {
		display: flex;
		flex-direction: column;
		height: 100%;
		background-color: #f9fafb;
	}

	/* Header */
	.lobby-header {
		background: white;
		border-bottom: 1px solid #e5e7eb;
		padding: 1.5rem 2rem;
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.page-title {
		font-size: 1.875rem;
		font-weight: 700;
		color: #111827;
		margin: 0;
	}

	.table-count {
		background-color: #667eea;
		color: white;
		padding: 0.25rem 0.75rem;
		border-radius: 1rem;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.count-divider {
		opacity: 0.7;
		margin: 0 0.125rem;
	}

	.total-count {
		opacity: 0.8;
	}

	.ws-status {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.75rem;
		border-radius: 1rem;
		font-size: 0.8125rem;
		font-weight: 600;
		transition: all 0.2s;
	}

	.ws-connected {
		background-color: #10b981;
		color: white;
	}

	.ws-reconnecting {
		background-color: #f59e0b;
		color: white;
	}

	.ws-dot {
		width: 0.5rem;
		height: 0.5rem;
		border-radius: 50%;
		background-color: currentColor;
	}

	.ws-connected .ws-dot {
		animation: pulse-dot 2s ease-in-out infinite;
	}

	.ws-reconnecting .ws-dot {
		animation: blink-dot 1s ease-in-out infinite;
	}

	@keyframes pulse-dot {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}

	@keyframes blink-dot {
		0%,
		50%,
		100% {
			opacity: 1;
		}
		25%,
		75% {
			opacity: 0.3;
		}
	}

	.header-right {
		display: flex;
		gap: 0.75rem;
	}

	.refresh-button,
	.create-button {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1.25rem;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
		border: 1px solid #d1d5db;
		background: white;
		color: #374151;
	}

	.refresh-button:hover:not(:disabled),
	.create-button:hover {
		background-color: #f9fafb;
		border-color: #9ca3af;
	}

	.refresh-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.create-button {
		background-color: #667eea;
		color: white;
		border-color: #667eea;
	}

	.create-button:hover {
		background-color: #5568d3;
		border-color: #5568d3;
	}

	.refresh-icon {
		transition: transform 0.6s ease;
	}

	.refresh-icon.spinning {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	/* Main Layout */
	.lobby-main {
		display: flex;
		flex: 1;
		overflow: hidden;
		gap: 1.5rem;
	}

	.sidebar {
		width: 300px;
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
	}

	.sidebar-left {
		padding: 1.5rem 0 1.5rem 1.5rem;
	}

	.sidebar-right {
		padding: 1.5rem 1.5rem 1.5rem 0;
	}

	/* Content */
	.lobby-content {
		flex: 1;
		padding: 2rem;
		overflow-y: auto;
		min-width: 0; /* Allow flexbox to shrink below content size */
	}

	/* Loading State */
	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 1.5rem;
	}

	.loading-text {
		color: #6b7280;
		font-size: 1.125rem;
		font-weight: 500;
		margin: 0;
	}

	/* Error State */
	.error-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 1rem;
		text-align: center;
	}

	.error-icon {
		font-size: 4rem;
		margin-bottom: 0.5rem;
	}

	.error-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #111827;
		margin: 0;
	}

	.error-message {
		color: #6b7280;
		font-size: 1rem;
		margin: 0;
		max-width: 28rem;
	}

	.retry-button {
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background-color: #667eea;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.retry-button:hover {
		background-color: #5568d3;
	}

	/* Empty State */
	.empty-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		gap: 1rem;
		text-align: center;
	}

	.empty-icon {
		font-size: 5rem;
		margin-bottom: 0.5rem;
		opacity: 0.5;
	}

	.empty-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: #111827;
		margin: 0;
	}

	.empty-message {
		color: #6b7280;
		font-size: 1rem;
		margin: 0;
		max-width: 28rem;
	}

	.create-table-button {
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background-color: #667eea;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.create-table-button:hover {
		background-color: #5568d3;
	}

	/* Filters Bar */
	.filters-bar {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding-top: 1.5rem;
		flex-wrap: wrap;
	}

	.search-input-wrapper {
		position: relative;
		flex: 1;
		min-width: 250px;
		max-width: 400px;
	}

	.search-icon {
		position: absolute;
		left: 0.875rem;
		top: 50%;
		transform: translateY(-50%);
		color: #9ca3af;
		pointer-events: none;
	}

	.search-input {
		width: 100%;
		padding: 0.625rem 2.75rem 0.625rem 2.75rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		transition: all 0.2s;
	}

	.search-input:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.clear-search {
		position: absolute;
		right: 0.625rem;
		top: 50%;
		transform: translateY(-50%);
		background: transparent;
		border: none;
		color: #9ca3af;
		cursor: pointer;
		padding: 0.25rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 0.25rem;
		transition: all 0.2s;
	}

	.clear-search:hover {
		color: #6b7280;
		background-color: #f3f4f6;
	}

	.format-select {
		padding: 0.625rem 2.5rem 0.625rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		background-color: white;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
		background-position: right 0.5rem center;
		background-repeat: no-repeat;
		background-size: 1.5em 1.5em;
		cursor: pointer;
		transition: all 0.2s;
		appearance: none;
	}

	.format-select:focus {
		outline: none;
		border-color: #667eea;
		box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		user-select: none;
		padding: 0.625rem 1rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		background: white;
		transition: all 0.2s;
	}

	.checkbox-label:hover {
		background-color: #f9fafb;
		border-color: #9ca3af;
	}

	.checkbox-input {
		width: 1.125rem;
		height: 1.125rem;
		border: 2px solid #d1d5db;
		border-radius: 0.25rem;
		cursor: pointer;
		transition: all 0.2s;
		accent-color: #667eea;
	}

	.checkbox-text {
		font-size: 0.875rem;
		font-weight: 500;
		color: #374151;
	}

	.clear-filters-button {
		padding: 0.625rem 1.25rem;
		background-color: #ef4444;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	.clear-filters-button:hover {
		background-color: #dc2626;
	}

	/* Tables Grid */
	.tables-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 1.5rem;
		padding-bottom: 2rem;
	}

	/* Responsive */
	@media (max-width: 1024px) {
		.tables-grid {
			grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
			gap: 1.25rem;
		}
	}

	@media (max-width: 768px) {
		.lobby-header {
			padding: 1rem;
		}

		.header-content {
			flex-direction: column;
			align-items: stretch;
		}

		.header-left {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}

		.page-title {
			font-size: 1.5rem;
		}

		.header-right {
			width: 100%;
		}

		.refresh-button,
		.create-button {
			flex: 1;
			justify-content: center;
		}

		.filters-bar {
			flex-direction: column;
			align-items: stretch;
			gap: 0.75rem;
		}

		.search-input-wrapper {
			max-width: none;
		}

		.format-select,
		.checkbox-label {
			width: 100%;
		}

		.lobby-main {
			flex-direction: column;
			gap: 1rem;
		}

		.sidebar {
			width: 100%;
		}

		.sidebar-left,
		.sidebar-right {
			padding: 1rem;
		}

		.lobby-content {
			padding: 1rem;
		}

		.tables-grid {
			grid-template-columns: 1fr;
			gap: 1rem;
		}
	}

	@media (max-width: 640px) {
		.refresh-button span,
		.create-button span {
			display: none;
		}

		.refresh-button,
		.create-button {
			padding: 0.625rem;
		}

		.filters-bar {
			padding-top: 1rem;
		}
	}
</style>
