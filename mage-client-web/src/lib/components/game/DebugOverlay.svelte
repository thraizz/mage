<script lang="ts">
	import { websocketStore } from '$lib/stores/websocket';
	import type { GameStoreState } from '$lib/stores/game.legacy';
	import type { PlayerView, CardView } from '$lib/generated/mage/v1/models';

	interface Props {
		open: boolean;
		gameId: string;
		localPlayerId: string;
		gameState: GameStoreState;
		allPlayers: PlayerView[];
		battlefieldCards: CardView[];
		stackCards: CardView[];
		commandCards: CardView[];
		turn: number;
		phase: string;
		havePriority: boolean;
		isMulliganPhase: boolean;
		gameFormat: string;
		isGameOver: boolean;
		gameWinner: string | null;
		activePlayerName: string;
		prompt: { type: string; message: string } | null;
		error: string | null;
		onClose: () => void;
	}

	let {
		open = $bindable(),
		gameId,
		localPlayerId,
		gameState,
		allPlayers,
		battlefieldCards,
		stackCards,
		commandCards,
		turn,
		phase,
		havePriority,
		isMulliganPhase,
		gameFormat,
		isGameOver,
		gameWinner,
		activePlayerName,
		prompt,
		error,
		onClose
	}: Props = $props();

	// Local state
	let updateCount = $state(0);
	let lastUpdate = $state<Date | null>(null);
	let expandedSections = $state<Record<string, boolean>>({
		gameState: true,
		players: true,
		zones: true,
		websocket: false,
		clientState: true,
		messages: false,
		combat: true,
		rawJson: false
	});

	// Track state changes
	let previousStateJson = $state('');

	$effect(() => {
		if (open) {
			const currentJson = JSON.stringify(gameState.gameView);
			if (currentJson !== previousStateJson && gameState.gameView) {
				previousStateJson = currentJson;
				lastUpdate = new Date();
				updateCount++;
			}
		}
	});

	/**
	 * Format JSON for display
	 */
	function formatJson(obj: unknown): string {
		try {
			return JSON.stringify(obj, null, 2);
		} catch {
			return String(obj);
		}
	}

	/**
	 * Toggle section expansion
	 */
	function toggleSection(section: string) {
		expandedSections[section] = !expandedSections[section];
	}

	/**
	 * Handle close
	 */
	function handleClose() {
		onClose();
	}

	/**
	 * Handle backdrop click
	 */
	function handleBackdropKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			handleClose();
		}
	}
</script>

{#if open}
	<div
		class="debug-overlay"
		role="dialog"
		aria-modal="true"
		onkeydown={handleBackdropKeydown}
		tabindex="0"
	>
		<div class="debug-modal">
			<header class="debug-header">
				<div class="debug-header-left">
					<h2>🔧 Debug View</h2>
					<div class="debug-status" class:connected={gameState.isConnected}>
						{gameState.isConnected ? '● Connected' : '○ Disconnected'}
					</div>
				</div>
				<div class="debug-header-right">
					<span class="debug-updates">Updates: {updateCount}</span>
					{#if lastUpdate}
						<span class="debug-time">Last: {lastUpdate.toLocaleTimeString()}</span>
					{/if}
					<button class="debug-close" onclick={handleClose}>✕</button>
				</div>
			</header>

			<main class="debug-content">
				<!-- Game State Overview -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('gameState')}>
						<span class="debug-icon">{expandedSections.gameState ? '▼' : '▶'}</span>
						<span>Game State Overview</span>
					</button>
					{#if expandedSections.gameState}
						<div class="debug-code">
							<pre><code
									>{@html `<span class="dk">gameId:</span> <span class="ds">"${gameId}"</span>
<span class="dk">localPlayerId:</span> <span class="ds">"${localPlayerId}"</span>
<span class="dk">state:</span> <span class="ds">"${gameState.gameView?.state || 'N/A'}"</span>
<span class="dk">turn:</span> <span class="dn">${turn}</span>
<span class="dk">phase:</span> <span class="ds">"${phase}"</span>
<span class="dk">step:</span> <span class="ds">"${gameState.gameView?.step || 'N/A'}"</span>
<span class="dk">activePlayerId:</span> <span class="ds">"${gameState.gameView?.activePlayerId || 'N/A'}"</span>
<span class="dk">activePlayerName:</span> <span class="ds">"${activePlayerName}"</span>
<span class="dk">priorityPlayerId:</span> <span class="ds">"${gameState.gameView?.priorityPlayerId || 'N/A'}"</span>
<span class="dk">hasPriority:</span> <span class="db">${havePriority}</span>
<span class="dk">isMulliganPhase:</span> <span class="db">${isMulliganPhase}</span>
<span class="dk">gameFormat:</span> <span class="ds">"${gameFormat}"</span>
<span class="dk">isGameOver:</span> <span class="db">${isGameOver}</span>
<span class="dk">winner:</span> <span class="ds">${gameWinner ? `"${gameWinner}"` : 'null'}</span>`}</code
								></pre>
						</div>
					{/if}
				</section>

				<!-- Players -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('players')}>
						<span class="debug-icon">{expandedSections.players ? '▼' : '▶'}</span>
						<span>Players ({allPlayers.length})</span>
					</button>
					{#if expandedSections.players}
						{#each allPlayers as player (player.playerId)}
							<div class="debug-player">
								<div class="debug-player-header">
									<span class="debug-badge" class:local={player.playerId === localPlayerId}>
										{player.playerId === localPlayerId ? '👤 You' : '👥 Opponent'}
									</span>
									<span>{player.name}</span>
								</div>
								<div class="debug-code">
									<pre><code
											>{@html `<span class="dk">playerId:</span> <span class="ds">"${player.playerId}"</span>
<span class="dk">life:</span> <span class="dn">${player.life}</span>
<span class="dk">poison:</span> <span class="dn">${player.poison}</span>
<span class="dk">libraryCount:</span> <span class="dn">${player.libraryCount}</span>
<span class="dk">handCount:</span> <span class="dn">${player.handCount}</span>
<span class="dk">hand:</span> [${player.hand?.map((c) => `\n  <span class="ds">"${c.name}"</span> <span class="dc">// ${c.id}</span>`).join(',') || ''}
]
<span class="dk">graveyard:</span> [${player.graveyard?.map((c) => `\n  <span class="ds">"${c.name}"</span>`).join(',') || ''}
]
<span class="dk">manaPool:</span> { W:${player.manaPool?.white || 0}, U:${player.manaPool?.blue || 0}, B:${player.manaPool?.black || 0}, R:${player.manaPool?.red || 0}, G:${player.manaPool?.green || 0}, C:${player.manaPool?.colorless || 0} }`}</code
										></pre>
								</div>
							</div>
						{/each}
					{/if}
				</section>

				<!-- Zones -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('zones')}>
						<span class="debug-icon">{expandedSections.zones ? '▼' : '▶'}</span>
						<span>Zones</span>
					</button>
					{#if expandedSections.zones}
						<div class="debug-zones-grid">
							<div class="debug-zone">
								<h4>🏟️ Battlefield ({battlefieldCards.length})</h4>
								<div class="debug-code small">
									<pre><code
											>{battlefieldCards.length > 0
												? formatJson(
														battlefieldCards.map((c) => ({
															id: c.id,
															name: c.name,
															type: c.type,
															controllerId: c.controllerId,
															tapped: c.tapped
														}))
													)
												: '[]'}</code
										></pre>
								</div>
							</div>
							<div class="debug-zone">
								<h4>📚 Stack ({stackCards.length})</h4>
								<div class="debug-code small">
									<pre><code
											>{stackCards.length > 0
												? formatJson(
														stackCards.map((c) => ({
															id: c.id,
															name: c.name,
															controllerId: c.controllerId
														}))
													)
												: '[]'}</code
										></pre>
								</div>
							</div>
							<div class="debug-zone">
								<h4>⚔️ Command ({commandCards.length})</h4>
								<div class="debug-code small">
									<pre><code
											>{commandCards.length > 0
												? formatJson(commandCards.map((c) => ({ id: c.id, name: c.name })))
												: '[]'}</code
										></pre>
								</div>
							</div>
							<div class="debug-zone">
								<h4>🚫 Exile ({gameState.gameView?.exile?.length || 0})</h4>
								<div class="debug-code small">
									<pre><code
											>{gameState.gameView?.exile?.length
												? formatJson(
														gameState.gameView.exile.map((c) => ({ id: c.id, name: c.name }))
													)
												: '[]'}</code
										></pre>
								</div>
							</div>
						</div>
					{/if}
				</section>

				<!-- Client State -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('clientState')}>
						<span class="debug-icon">{expandedSections.clientState ? '▼' : '▶'}</span>
						<span>Client Store State</span>
					</button>
					{#if expandedSections.clientState}
						<div class="debug-code">
							<pre><code
									>{@html `<span class="dc">// Store meta</span>
<span class="dk">isConnected:</span> <span class="db">${gameState.isConnected}</span>
<span class="dk">isLoading:</span> <span class="db">${gameState.isLoading}</span>
<span class="dk">error:</span> <span class="ds">${error ? `"${error}"` : 'null'}</span>
<span class="dk">selectedCardIds:</span> [${gameState.selectedCardIds.map((id) => `<span class="ds">"${id}"</span>`).join(', ')}]
<span class="dk">showStack:</span> <span class="db">${gameState.showStack}</span>
<span class="dk">gameOver:</span> <span class="db">${gameState.gameOver}</span>
<span class="dk">pendingPrompt:</span> ${prompt ? `{ type: "${prompt.type}", message: "${prompt.message}" }` : 'null'}`}</code
								></pre>
						</div>
					{/if}
				</section>

				<!-- WebSocket State -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('websocket')}>
						<span class="debug-icon">{expandedSections.websocket ? '▼' : '▶'}</span>
						<span>WebSocket State</span>
					</button>
					{#if expandedSections.websocket}
						<div class="debug-code">
							<pre><code
									>{@html `<span class="dk">state:</span> <span class="ds">"${$websocketStore.state}"</span>
<span class="dk">error:</span> <span class="ds">${$websocketStore.error ? `"${$websocketStore.error}"` : 'null'}</span>
<span class="dk">lastConnected:</span> <span class="dn">${$websocketStore.lastConnected ? new Date($websocketStore.lastConnected).toISOString() : 'null'}</span>
<span class="dk">reconnectAttempts:</span> <span class="dn">${$websocketStore.reconnectAttempts}</span>`}</code
								></pre>
						</div>
					{/if}
				</section>

				<!-- Game Messages -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('messages')}>
						<span class="debug-icon">{expandedSections.messages ? '▼' : '▶'}</span>
						<span>Game Messages ({gameState.gameView?.messages?.length || 0})</span>
					</button>
					{#if expandedSections.messages}
						<div class="debug-code messages">
							{#if gameState.gameView?.messages?.length}
								{#each gameState.gameView.messages.slice(-50) as msg, idx}
									<div class="debug-log-line">
										<span class="debug-log-idx">{idx + 1}.</span>
										<span class="debug-log-msg">{msg}</span>
									</div>
								{/each}
							{:else}
								<pre><code><span class="dc">// No messages</span></code></pre>
							{/if}
						</div>
					{/if}
				</section>

				<!-- Raw JSON -->
				<section class="debug-section">
					<button class="debug-section-header" onclick={() => toggleSection('rawJson')}>
						<span class="debug-icon">{expandedSections.rawJson ? '▼' : '▶'}</span>
						<span>Raw GameView JSON</span>
					</button>
					{#if expandedSections.rawJson}
						<div class="debug-code raw">
							<pre><code>{formatJson(gameState.gameView)}</code></pre>
						</div>
					{/if}
				</section>
			</main>
		</div>
	</div>
{/if}

<style>
	.debug-overlay {
		position: fixed;
		inset: 0;
		top: 9rem;
		background: #2d2d2d;
		z-index: 2000;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		font-family: 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
	}

	.debug-modal {
		flex: 1;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.debug-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1rem;
		background: #1a1a1a;
		border-bottom: 1px solid #444;
		flex-shrink: 0;
	}

	.debug-header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.debug-header h2 {
		font-size: 1rem;
		font-weight: 600;
		margin: 0;
		color: #00ff00;
	}

	.debug-status {
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 999px;
		background: #3a0000;
		color: #ff6b6b;
		border: 1px solid #ff6b6b;
	}

	.debug-status.connected {
		background: #003a00;
		color: #00ff00;
		border-color: #00ff00;
	}

	.debug-header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.debug-updates {
		font-size: 0.75rem;
		color: #ffff00;
	}

	.debug-time {
		font-size: 0.75rem;
		color: #888;
	}

	.debug-close {
		width: 32px;
		height: 32px;
		border-radius: 6px;
		background: #333;
		border: 1px solid #555;
		color: #fff;
		font-size: 1rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.debug-close:hover {
		background: #ef4444;
		border-color: #ef4444;
	}

	.debug-content {
		flex: 1;
		overflow-y: auto;
		padding: 1rem;
	}

	.debug-section {
		margin-bottom: 1rem;
		background: #1e1e1e;
		border: 1px solid #444;
		border-radius: 8px;
		overflow: hidden;
	}

	.debug-section-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
		padding: 0.75rem 1rem;
		background: #252525;
		border: none;
		color: #00ff00;
		cursor: pointer;
		font-family: inherit;
		font-size: 0.875rem;
		text-align: left;
		transition: background 0.2s;
	}

	.debug-section-header:hover {
		background: #2a2a2a;
	}

	.debug-icon {
		font-size: 0.625rem;
		color: #888;
	}

	.debug-code {
		background: #0d0d0d;
		padding: 0.75rem 1rem;
		overflow-x: auto;
		max-height: 300px;
		overflow-y: auto;
	}

	.debug-code.small {
		max-height: 150px;
	}

	.debug-code.messages {
		max-height: 250px;
	}

	.debug-code.raw {
		max-height: 400px;
	}

	.debug-code pre {
		margin: 0;
		font-size: 0.75rem;
		line-height: 1.4;
	}

	.debug-code code {
		color: #00ff00;
	}

	/* Syntax highlighting */
	:global(.debug-code .dk) {
		color: #9cdcfe;
	}
	:global(.debug-code .ds) {
		color: #ce9178;
	}
	:global(.debug-code .dn) {
		color: #b5cea8;
	}
	:global(.debug-code .db) {
		color: #569cd6;
	}
	:global(.debug-code .dc) {
		color: #6a9955;
		font-style: italic;
	}

	.debug-player {
		border-top: 1px solid #333;
	}

	.debug-player-header {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
		background: #1a1a1a;
	}

	.debug-badge {
		font-size: 0.625rem;
		padding: 0.125rem 0.375rem;
		border-radius: 4px;
		background: #333;
		color: #888;
	}

	.debug-badge.local {
		background: #003a00;
		color: #00ff00;
	}

	.debug-zones-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1px;
		background: #333;
	}

	.debug-zone {
		background: #1e1e1e;
	}

	.debug-zone h4 {
		font-size: 0.75rem;
		font-weight: 600;
		margin: 0;
		padding: 0.5rem 0.75rem;
		background: #252525;
		color: #ffff00;
		border-bottom: 1px solid #333;
	}

	.debug-log-line {
		display: flex;
		gap: 0.5rem;
		padding: 0.125rem 0;
		font-size: 0.6875rem;
		border-bottom: 1px solid #222;
	}

	.debug-log-idx {
		color: #666;
		min-width: 2rem;
		text-align: right;
	}

	.debug-log-msg {
		color: #00ff00;
	}

	/* Scrollbars */
	.debug-code::-webkit-scrollbar,
	.debug-content::-webkit-scrollbar {
		width: 6px;
		height: 6px;
	}

	.debug-code::-webkit-scrollbar-track,
	.debug-content::-webkit-scrollbar-track {
		background: #1a1a1a;
	}

	.debug-code::-webkit-scrollbar-thumb,
	.debug-content::-webkit-scrollbar-thumb {
		background: #444;
		border-radius: 3px;
	}

	.debug-code::-webkit-scrollbar-thumb:hover,
	.debug-content::-webkit-scrollbar-thumb:hover {
		background: #555;
	}
</style>
