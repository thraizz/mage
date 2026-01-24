<script lang="ts">
	import type { Table } from '$lib/types/table';

	// Props
	export let table: Table;
	 
	export let onClick: ((table: Table) => void) | undefined = undefined;

	function getStatusColor(): string {
		switch (table.status) {
			case 'waiting':
				return 'var(--table-waiting)'; // Orange - waiting for players
			case 'ready':
				return 'var(--ci-forest-emerald)'; // Green - ready to start
			case 'playing':
				return 'var(--ci-jace-cloak)'; // Blue - in progress
			case 'finished':
				return 'var(--ci-swamp-obsidian)'; // Gray - finished
			default:
				return 'var(--ci-swamp-obsidian)';
		}
	}

	function getStatusText(): string {
		switch (table.status) {
			case 'waiting':
				return 'Waiting for players';
			case 'ready':
				return 'Ready to start';
			case 'playing':
				return 'In progress';
			case 'finished':
				return 'Finished';
			default:
				return 'Unknown';
		}
	}

	function getFormatIcon(): string {
		switch (table.format) {
			case 'Commander':
				return '👑';
			case 'Standard':
				return '⚔️';
			case 'Modern':
				return '🔥';
			case 'Legacy':
				return '💎';
			case 'Vintage':
				return '🏺';
			case 'Pauper':
				return '💰';
			case 'Pioneer':
				return '🌟';
			case 'Limited':
			case 'Draft':
			case 'Sealed':
				return '📦';
			default:
				return '🎮';
		}
	}

	function handleClick() {
		if (onClick) {
			onClick(table);
		}
	}

	$: isFull = table.players.length >= table.maxPlayers;
	$: isOpen = table.status === 'waiting' && !isFull;
</script>

{#if onClick}
	<div
		class="table-card clickable"
		class:full={isFull}
		onclick={handleClick}
		onkeypress={(e) => e.key === 'Enter' && handleClick()}
		role="button"
		tabindex="0"
	>
		<div class="card-header">
			<div class="format-badge">
				<span class="format-icon">{getFormatIcon()}</span>
				<span class="format-name">{table.format}</span>
			</div>
			{#if table.hasPassword}
				<div class="password-badge" title="Password protected">🔒</div>
			{/if}
		</div>

		<!-- Table Name -->
		<h3 class="table-name">{table.name}</h3>

		<!-- Host Info -->
		<div class="table-info">
			<div class="info-row">
				<span class="info-label">Host:</span>
				<span class="info-value">{table.hostUsername}</span>
			</div>

			<!-- Player Count -->
			<div class="info-row">
				<span class="info-label">Players:</span>
				<span class="info-value players-count" class:full={isFull}>
					{table.players.length}/{table.maxPlayers}
				</span>
			</div>
		</div>

		<!-- Status -->
		<div class="status-bar">
			<div class="status-indicator" style="background-color: {getStatusColor()}"></div>
			<span class="status-text">{getStatusText()}</span>
		</div>

		<!-- Join Button Overlay (only if clickable and open) -->
		{#if isOpen}
			<div class="join-overlay">
				<button class="join-button">Join Table</button>
			</div>
		{/if}
	</div>
{:else}
	<div class="table-card" class:full={isFull} role="article">
		<div class="card-header">
			<div class="format-badge">
				<span class="format-icon">{getFormatIcon()}</span>
				<span class="format-name">{table.format}</span>
			</div>
			{#if table.hasPassword}
				<div class="password-badge" title="Password protected">🔒</div>
			{/if}
		</div>

		<!-- Table Name -->
		<h3 class="table-name">{table.name}</h3>

		<!-- Host Info -->
		<div class="table-info">
			<div class="info-row">
				<span class="info-label">Host:</span>
				<span class="info-value">{table.hostUsername}</span>
			</div>

			<!-- Player Count -->
			<div class="info-row">
				<span class="info-label">Players:</span>
				<span class="info-value players-count" class:full={isFull}>
					{table.players.length}/{table.maxPlayers}
				</span>
			</div>
		</div>

		<!-- Status -->
		<div class="status-bar">
			<div class="status-indicator" style="background-color: {getStatusColor()}"></div>
			<span class="status-text">{getStatusText()}</span>
		</div>
	</div>
{/if}

<style>
	.table-card {
		background: var(--bg-obsidian);
		border: 1px solid var(--border-subtle);
		border-radius: var(--radius-lg);
		padding: var(--space-5);
		transition: all var(--transition-base);
		position: relative;
		overflow: hidden;
	}

	.table-card.clickable {
		cursor: pointer;
	}

	.table-card.clickable:hover {
		border-color: var(--ci-jace-cloak);
		box-shadow:
			0 0 20px rgba(59, 130, 246, 0.2),
			0 4px 12px rgba(0, 0, 0, 0.3);
		transform: translateY(-2px);
	}

	.table-card.clickable:active {
		transform: translateY(0) scale(0.99);
	}

	.table-card.full {
		opacity: 0.6;
	}

	.table-card.full.clickable:hover {
		opacity: 0.75;
		border-color: var(--border-default);
	}

	/* Header */
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--space-3);
	}

	.format-badge {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
		color: var(--ci-scroll-parchment);
		padding: var(--space-1) var(--space-3);
		border-radius: var(--radius-md);
		font-size: var(--text-sm);
		font-weight: var(--weight-semibold);
		box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
	}

	.format-icon {
		font-size: var(--text-base);
	}

	.format-name {
		line-height: 1;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.password-badge {
		font-size: var(--text-xl);
		opacity: 0.8;
	}

	/* Table Name */
	.table-name {
		font-family: var(--font-display);
		font-size: var(--text-lg);
		font-weight: var(--weight-bold);
		color: var(--ci-scroll-parchment);
		margin: 0 0 var(--space-3) 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		letter-spacing: 0.02em;
	}

	/* Table Info */
	.table-info {
		display: flex;
		flex-direction: column;
		gap: var(--space-2);
		margin-bottom: var(--space-3);
	}

	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: var(--text-sm);
	}

	.info-label {
		color: var(--ci-swamp-obsidian);
		font-weight: var(--weight-medium);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		font-size: var(--text-xs);
	}

	.info-value {
		color: var(--ci-scroll-parchment);
		font-weight: var(--weight-semibold);
	}

	.players-count.full {
		color: var(--ci-mountain-ember);
	}

	/* Status Bar */
	.status-bar {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding-top: var(--space-3);
		border-top: 1px solid var(--border-subtle);
	}

	.status-indicator {
		width: 0.625rem;
		height: 0.625rem;
		border-radius: var(--radius-full);
		box-shadow: 0 0 6px currentColor;
	}

	.status-text {
		font-size: var(--text-xs);
		color: var(--ci-swamp-obsidian);
		font-weight: var(--weight-medium);
		font-style: italic;
	}

	/* Join Overlay */
	.join-overlay {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		background: linear-gradient(to top, rgba(11, 12, 16, 0.95), transparent);
		padding: var(--space-8) var(--space-5) var(--space-5);
		opacity: 0;
		transition: opacity var(--transition-base);
		pointer-events: none;
	}

	.table-card.clickable:hover .join-overlay {
		opacity: 1;
	}

	.join-button {
		width: 100%;
		padding: var(--space-3) var(--space-4);
		background: linear-gradient(135deg, var(--ci-jace-cloak) 0%, #2563eb 100%);
		color: var(--ci-scroll-parchment);
		border: none;
		border-radius: var(--radius-md);
		font-weight: var(--weight-bold);
		font-size: var(--text-sm);
		font-family: var(--font-display);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		cursor: pointer;
		transition: all var(--transition-base);
		pointer-events: auto;
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
	}

	.join-button:hover {
		background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
		box-shadow: 0 6px 16px rgba(59, 130, 246, 0.5);
		transform: translateY(-1px);
	}

	/* Responsive */
	@media (max-width: 640px) {
		.table-card {
			padding: var(--space-4);
		}

		.table-name {
			font-size: var(--text-base);
		}
	}
</style>
