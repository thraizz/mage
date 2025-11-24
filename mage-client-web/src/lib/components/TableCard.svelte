<script lang="ts">
	import type { Table } from '$lib/types/table';

	// Props
	export let table: Table;
	// eslint-disable-next-line no-unused-vars
	export let onClick: ((table: Table) => void) | undefined = undefined;

	function getStatusColor(): string {
		switch (table.status) {
			case 'waiting':
				return '#f59e0b'; // Orange
			case 'ready':
				return '#10b981'; // Green
			case 'playing':
				return '#3b82f6'; // Blue
			case 'finished':
				return '#6b7280'; // Gray
			default:
				return '#6b7280';
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
		on:click={handleClick}
		on:keypress={(e) => e.key === 'Enter' && handleClick()}
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
		background: white;
		border-radius: 0.75rem;
		padding: 1.25rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		transition: all 0.2s;
		position: relative;
		overflow: hidden;
	}

	.table-card.clickable {
		cursor: pointer;
	}

	.table-card.clickable:hover {
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		transform: translateY(-2px);
	}

	.table-card.clickable:active {
		transform: translateY(0);
	}

	.table-card.full {
		opacity: 0.7;
	}

	.table-card.full.clickable:hover {
		opacity: 0.8;
	}

	/* Header */
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.75rem;
	}

	.format-badge {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background-color: #667eea;
		color: white;
		padding: 0.375rem 0.75rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.format-icon {
		font-size: 1rem;
	}

	.format-name {
		line-height: 1;
	}

	.password-badge {
		font-size: 1.25rem;
		opacity: 0.7;
	}

	/* Table Name */
	.table-name {
		font-size: 1.125rem;
		font-weight: 700;
		color: #111827;
		margin: 0 0 0.75rem 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* Table Info */
	.table-info {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 0.75rem;
	}

	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 0.875rem;
	}

	.info-label {
		color: #6b7280;
		font-weight: 500;
	}

	.info-value {
		color: #111827;
		font-weight: 600;
	}

	.players-count.full {
		color: #ef4444;
	}

	/* Status Bar */
	.status-bar {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding-top: 0.75rem;
		border-top: 1px solid #e5e7eb;
	}

	.status-indicator {
		width: 0.75rem;
		height: 0.75rem;
		border-radius: 50%;
	}

	.status-text {
		font-size: 0.8125rem;
		color: #6b7280;
		font-weight: 500;
	}

	/* Join Overlay */
	.join-overlay {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
		padding: 2rem 1.25rem 1.25rem;
		opacity: 0;
		transition: opacity 0.2s;
		pointer-events: none;
	}

	.table-card.clickable:hover .join-overlay {
		opacity: 1;
	}

	.join-button {
		width: 100%;
		padding: 0.625rem 1rem;
		background-color: #667eea;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 600;
		font-size: 0.875rem;
		cursor: pointer;
		transition: background-color 0.2s;
		pointer-events: auto;
	}

	.join-button:hover {
		background-color: #5568d3;
	}

	/* Responsive */
	@media (max-width: 640px) {
		.table-card {
			padding: 1rem;
		}

		.table-name {
			font-size: 1rem;
		}
	}
</style>
