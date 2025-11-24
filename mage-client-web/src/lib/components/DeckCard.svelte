<script lang="ts">
	import type { Deck } from '$lib/types/deck';

	export let deck: Deck;

	// Format the last modified date
	function formatDate(timestamp: number): string {
		const date = new Date(timestamp);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

		if (diffDays === 0) {
			return 'Today';
		} else if (diffDays === 1) {
			return 'Yesterday';
		} else if (diffDays < 7) {
			return `${diffDays} days ago`;
		} else if (diffDays < 30) {
			const weeks = Math.floor(diffDays / 7);
			return `${weeks} week${weeks > 1 ? 's' : ''} ago`;
		} else {
			return date.toLocaleDateString();
		}
	}

	// Get format color class for visual distinction
	function getFormatColorClass(format: string): string {
		const colors: Record<string, string> = {
			Standard: 'format-badge-blue',
			Modern: 'format-badge-purple',
			Commander: 'format-badge-yellow',
			Legacy: 'format-badge-red',
			Vintage: 'format-badge-pink',
			Pioneer: 'format-badge-green',
			Pauper: 'format-badge-gray',
			Historic: 'format-badge-indigo'
		};
		return colors[format] || 'format-badge-gray';
	}
</script>

<div
	class="deck-card"
	role="button"
	tabindex="0"
	on:click
	on:keydown={(e) => {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			e.currentTarget.click();
		}
	}}
>
	<!-- Deck Name -->
	<h3 class="deck-name" title={deck.name}>
		{deck.name}
	</h3>

	<!-- Format Badge -->
	<div class="format-badge-container">
		<span class="format-badge {getFormatColorClass(deck.format)}">
			{deck.format}
		</span>
	</div>

	<!-- Deck Stats -->
	<div class="deck-stats">
		<!-- Card Count -->
		<div class="stat-row">
			<span class="stat-label">
				<svg class="stat-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
					/>
				</svg>
				Cards
			</span>
			<span class="stat-value">{deck.cardCount}</span>
		</div>

		<!-- Main Deck Count -->
		<div class="stat-row">
			<span>Main Deck</span>
			<span class="stat-value">{deck.mainDeck.length} cards</span>
		</div>

		<!-- Sideboard Count -->
		{#if deck.sideboard && deck.sideboard.length > 0}
			<div class="stat-row">
				<span>Sideboard</span>
				<span class="stat-value">{deck.sideboard.length} cards</span>
			</div>
		{/if}

		<!-- Validity Badge -->
		<div class="stat-row stat-row-border">
			<span>Status</span>
			{#if deck.isValid}
				<span class="status-badge status-valid">
					<svg class="status-icon" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
							clip-rule="evenodd"
						/>
					</svg>
					Valid
				</span>
			{:else}
				<span class="status-badge status-invalid">
					<svg class="status-icon" fill="currentColor" viewBox="0 0 20 20">
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
							clip-rule="evenodd"
						/>
					</svg>
					Invalid
				</span>
			{/if}
		</div>

		<!-- Last Modified -->
		<div class="stat-row stat-row-small">
			<span>Modified</span>
			<span>{formatDate(deck.updatedAt)}</span>
		</div>
	</div>
</div>

<style>
	.deck-card {
		background-color: white;
		border-radius: 0.5rem;
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
		transition: box-shadow 0.2s, border-color 0.2s;
		padding: 1rem;
		cursor: pointer;
		border: 1px solid #e5e7eb;
	}

	.deck-card:hover {
		box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
		border-color: #60a5fa;
	}

	.deck-name {
		font-size: 1.125rem;
		font-weight: 600;
		color: #111827;
		margin: 0 0 0.5rem 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.format-badge-container {
		margin-bottom: 0.75rem;
	}

	.format-badge {
		display: inline-block;
		padding: 0.25rem 0.5rem;
		font-size: 0.75rem;
		font-weight: 500;
		border-radius: 0.25rem;
	}

	/* Format badge colors */
	.format-badge-blue {
		background-color: #dbeafe;
		color: #1e40af;
	}

	.format-badge-purple {
		background-color: #f3e8ff;
		color: #6b21a8;
	}

	.format-badge-yellow {
		background-color: #fef3c7;
		color: #92400e;
	}

	.format-badge-red {
		background-color: #fee2e2;
		color: #991b1b;
	}

	.format-badge-pink {
		background-color: #fce7f3;
		color: #9f1239;
	}

	.format-badge-green {
		background-color: #d1fae5;
		color: #065f46;
	}

	.format-badge-gray {
		background-color: #f3f4f6;
		color: #1f2937;
	}

	.format-badge-indigo {
		background-color: #e0e7ff;
		color: #3730a3;
	}

	.deck-stats {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-size: 0.875rem;
		color: #4b5563;
	}

	.stat-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.stat-row-border {
		padding-top: 0.5rem;
		border-top: 1px solid #f3f4f6;
	}

	.stat-row-small {
		font-size: 0.75rem;
		color: #6b7280;
		padding-top: 0.5rem;
	}

	.stat-label {
		display: flex;
		align-items: center;
	}

	.stat-icon {
		width: 1rem;
		height: 1rem;
		margin-right: 0.25rem;
	}

	.stat-value {
		font-weight: 500;
		color: #111827;
	}

	.status-badge {
		display: inline-flex;
		align-items: center;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-icon {
		width: 1rem;
		height: 1rem;
		margin-right: 0.25rem;
	}

	.status-valid {
		color: #15803d;
	}

	.status-invalid {
		color: #b91c1c;
	}
</style>
