<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { initializePlaytest, validateDeckIds } from '$lib/playtest/initializer';
	import { gameStore } from '$lib/stores/game.legacy';
	import {
		playtestActiveControlSeat,
		playtestBattlefield,
		playtestExile,
		playtestGameStore,
		playtestIsInitialized,
		playtestLocalPlayer,
		playtestOpponents,
		playtestPlayers,
		type PlaytestSessionMeta
	} from '$lib/stores/playtest-game';
	import { toast } from '$lib/stores/toast';
	import { onMount, untrack } from 'svelte';
	// Game components
	import type { MenuAction } from '$lib/components/game/DeckContextMenu.svelte';
	import type { ScrySession } from '$lib/stores/playtest-game';
	import { dragDropStore, getAllValidDropZones, type SourceZone } from '$lib/utils/drag-drop';
	import Copy from '@lucide/svelte/icons/copy';

	// Game log
	const gameLog = $derived($playtestGameStore.log || []);

	/**
	 * Copy game log to clipboard
	 */
	async function handleCopyLog(): Promise<void> {
		const logText = playtestGameStore.buildLogText($playtestGameStore);
		try {
			if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
				await navigator.clipboard.writeText(logText);
				toast.success('Game log copied to clipboard!');
				return;
			}
			// Fallback for older browsers
			const textarea = document.createElement('textarea');
			textarea.value = logText;
			textarea.style.position = 'fixed';
			textarea.style.top = '0';
			textarea.style.left = '0';
			textarea.style.opacity = '0';
			document.body.appendChild(textarea);
			textarea.focus();
			textarea.select();
			const ok = document.execCommand('copy');
			document.body.removeChild(textarea);
			if (ok) {
				toast.success('Game log copied to clipboard!');
			} else {
				toast.error('Failed to copy log');
			}
		} catch (err) {
			console.error('Failed to copy log to clipboard:', err);
			toast.error('Failed to copy log');
		}
	}
</script>

<section class="debug-section">
	<div class="debug-section-header">
		<span>Game State Log ({gameLog.length} events)</span>
		<button
			class="debug-copy-btn"
			onclick={handleCopyLog}
			title="Copy log to clipboard"
			aria-label="Copy log to clipboard"
		>
			<Copy size={16} aria-hidden="true" />
			<span>Copy</span>
		</button>
	</div>
	<div class="debug-log-container">
		{#if gameLog.length === 0}
			<div class="debug-log-empty">No events logged yet</div>
		{:else}
			<div class="debug-log-entries">
				{#each gameLog.slice().reverse() as entry (entry.id)}
					<div class="debug-log-entry">
						<span class="debug-log-time">
							{new Date(entry.at).toLocaleTimeString([], {
								hour: '2-digit',
								minute: '2-digit',
								second: '2-digit'
							})}
						</span>
						<span class="debug-log-turn">T{entry.turn}</span>
						<span class="debug-log-kind">{entry.kind}</span>
						<span class="debug-log-message">{entry.message}</span>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</section>
