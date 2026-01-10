<script lang="ts">
	/**
	 * LibraryZone - Compact button to view library/deck contents
	 * Styled like Graveyard and ExileZone buttons.
	 * Clicking opens the LibrarySearch modal via the searchLibrary() API.
	 */

	// Props
	let {
		libraryCount = 0,
		playerName = 'Player',
		isOpponent = false,
		onSearch = () => {}
	}: {
		libraryCount?: number;
		playerName?: string;
		isOpponent?: boolean;
		onSearch?: () => void;
	} = $props();

	// Derived values
	const isEmpty = $derived(libraryCount === 0);

	/**
	 * Handle click - trigger library search
	 */
	function handleClick(): void {
		if (!isEmpty && !isOpponent) {
			onSearch();
		}
	}
</script>

<button
	class="library-compact"
	class:has-cards={!isEmpty}
	class:opponent={isOpponent}
	onclick={handleClick}
	disabled={isOpponent}
	title="{playerName}'s Library ({libraryCount} cards){isEmpty || isOpponent ? '' : ' - Click to view'}"
>
	<span class="library-icon">📚</span>
	<span class="library-label">Deck</span>
	<span class="card-count" class:zero={isEmpty}>{libraryCount}</span>
</button>

<style>
	/* Compact Library Button - styled like Graveyard */
	.library-compact {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.625rem;
		background: rgba(26, 31, 46, 0.6);
		border: 1px solid rgba(34, 197, 94, 0.3);
		border-radius: 6px;
		min-height: 32px;
		cursor: default;
		transition: all 0.15s;
		color: inherit;
	}

	.library-compact.has-cards:not(.opponent) {
		cursor: pointer;
		background: rgba(26, 31, 46, 0.9);
		border-color: rgba(34, 197, 94, 0.4);
	}

	.library-compact.has-cards:not(.opponent):hover {
		background: rgba(42, 52, 65, 0.9);
		border-color: rgba(34, 197, 94, 0.6);
	}

	.library-compact.opponent {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.library-compact:disabled {
		cursor: not-allowed;
	}

	.library-icon {
		font-size: 0.875rem;
		opacity: 0.7;
	}

	.library-compact.has-cards .library-icon {
		opacity: 1;
	}

	.library-label {
		font-size: 0.6875rem;
		color: #22c55e;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.card-count {
		font-size: 0.75rem;
		font-weight: 700;
		color: #22c55e;
		background: rgba(34, 197, 94, 0.2);
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		min-width: 1.25rem;
		text-align: center;
	}

	.card-count.zero {
		color: #4b5563;
		background: transparent;
	}
</style>

