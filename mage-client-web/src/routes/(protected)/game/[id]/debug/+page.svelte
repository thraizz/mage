<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import type { CardView } from '$lib/generated/mage/v1/models';
  import { auth } from '$lib/stores/auth';
  import {
    multiplayerBattlefield,
    multiplayerGameStore,
    multiplayerIsInitialized,
    multiplayerLocalPlayer,
    multiplayerOpponents,
    multiplayerPlayers,
    multiplayerStack
  } from '$lib/stores/multiplayer-game';
  import { websocketStore } from '$lib/stores/websocket';
  import { getSessionIdFromToken } from '$lib/utils/jwt';
  import { onDestroy, onMount } from 'svelte';

  // Game ID from route params
  const gameId = $derived(page.params.id);

  // Local state
  let initialized = $state(false);
  let lastUpdateTime = $state<Date | null>(null);
  let updateCount = $state(0);
  let rawWebsocketEvents = $state<Array<{ time: Date; type: string; data: unknown }>>([]);
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

  // Get local player ID from auth
  const localPlayerId = $derived($auth.user?.username || '');

  // Derived state from multiplayer game store
  const gameState = $derived($multiplayerGameStore);
  const allPlayers = $derived($multiplayerPlayers);
  const me = $derived($multiplayerLocalPlayer);
  const otherPlayers = $derived($multiplayerOpponents);
  const myCards = $derived(me?.hand || []);
  const myGrave = $derived(me?.graveyard || []);
  const myMana = $derived(me?.manaPool);
  const havePriority = $derived(false); // Not available in multiplayer store
  const phase = $derived('N/A'); // Not available in multiplayer store
  const turn = $derived(gameState.turn);
  const battlefieldCards = $derived($multiplayerBattlefield);
  const stackCards = $derived($multiplayerStack);
  const commandCards = $derived(gameState.command);
  const prompt = $derived(null); // Not available in multiplayer store
  const isGameOver = $derived(false); // Not available in multiplayer store
  const gameWinner = $derived(null); // Not available in multiplayer store
  const error = $derived(null); // Not available in multiplayer store
  const loading = $derived(!$multiplayerIsInitialized);

  /**
   * Initialize game connection
   */
  async function initializeGame() {
    if (!localPlayerId || !gameId) {
      console.error('[DebugPage] Missing player ID or game ID');
      return;
    }

    try {
      console.log('[DebugPage] Starting game initialization...', { gameId, localPlayerId });

      const wsState = $websocketStore;
      if (wsState.state !== 'connected') {
        const token = $auth.token;
        const sessionId = token ? getSessionIdFromToken(token) : null;
        if (sessionId) {
          console.log('[DebugPage] Connecting to WebSocket...');
          await websocketStore.connect(sessionId);
          console.log('[DebugPage] WebSocket connected');
        } else {
          throw new Error('No session ID available');
        }
      }

      console.log('[DebugPage] Initializing game store...');
      await multiplayerGameStore.initialize(gameId);

      console.log('[DebugPage] Game store initialized');

      initialized = true;
      console.log('[DebugPage] Game initialization complete');
    } catch (err) {
      console.error('[DebugPage] Failed to initialize game:', err);
      goto('/lobby');
    }
  }

  /**
   * Format JSON with syntax highlighting
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

  // Initialize on mount
  onMount(() => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    initializeGame();
  });

  // Cleanup on destroy
  onDestroy(() => {
    multiplayerGameStore.reset();
  });
</script>

<svelte:head>
  <title>Debug - Game {gameId} - MAGE</title>
</svelte:head>

<div class="debug-container">
  <header class="debug-header">
    <div class="header-left">
      <a href="/game/{gameId}" class="back-link">← Back to Game</a>
      <a href="/admin/games" class="admin-link">🛠️ All Games</a>
      <h1>🔧 Game Debug View</h1>
    </div>
    <div class="header-right">
      <div class="status-pill" class:connected={gameState.isConnected}>
        {gameState.isConnected ? '● Connected' : '○ Disconnected'}
      </div>
    </div>
  </header>

  <main class="debug-content">
    <!-- Game State Overview -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('gameState')}>
        <span class="section-icon">{expandedSections.gameState ? '▼' : '▶'}</span>
        <h2>Game State Overview</h2>
      </button>
      {#if expandedSections.gameState}
        <div class="code-block">
          <pre><code
              >{@html `<span class="key">gameId:</span> <span class="string">"${gameId}"</span>
<span class="key">localPlayerId:</span> <span class="string">"${localPlayerId}"</span>
<span class="key">state:</span> <span class="string">"${gameState.gameView?.state || 'N/A'}"</span>
<span class="key">turn:</span> <span class="number">${turn}</span>
<span class="key">phase:</span> <span class="string">"${phase}"</span>
<span class="key">step:</span> <span class="string">"${gameState.gameView?.step || 'N/A'}"</span>
<span class="key">activePlayerId:</span> <span class="string">"${gameState.gameView?.activePlayerId || 'N/A'}"</span>
<span class="key">activePlayerName:</span> <span class="string">"${gameState.gameView?.activePlayerName || 'N/A'}"</span>
<span class="key">priorityPlayerId:</span> <span class="string">"${gameState.gameView?.priorityPlayerId || 'N/A'}"</span>
<span class="key">hasPriority:</span> <span class="boolean">${havePriority}</span>
<span class="key">isMulliganPhase:</span> <span class="boolean">${gameState.gameView?.isMulliganPhase || false}</span>
<span class="key">gameFormat:</span> <span class="string">"${gameState.gameView?.gameFormat || 'N/A'}"</span>
<span class="key">isGameOver:</span> <span class="boolean">${isGameOver}</span>
<span class="key">winner:</span> <span class="string">${gameWinner ? `"${gameWinner}"` : 'null'}</span>`}</code
            ></pre>
        </div>
      {/if}
    </section>

    <!-- Players -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('players')}>
        <span class="section-icon">{expandedSections.players ? '▼' : '▶'}</span>
        <h2>Players ({allPlayers.length})</h2>
      </button>
      {#if expandedSections.players}
        {#each allPlayers as player, idx (player.playerId)}
          <div class="player-block">
            <div class="player-header">
              <span class="player-badge" class:local={player.playerId === localPlayerId}>
                {player.playerId === localPlayerId ? '👤 You' : '👥 Opponent'}
              </span>
              <span class="player-name">{player.name}</span>
            </div>
            <div class="code-block">
              <pre><code
                  >{@html `<span class="key">playerId:</span> <span class="string">"${player.playerId}"</span>
<span class="key">name:</span> <span class="string">"${player.name}"</span>
<span class="key">life:</span> <span class="number">${player.life}</span>
<span class="key">poison:</span> <span class="number">${player.poison}</span>
<span class="key">libraryCount:</span> <span class="number">${player.libraryCount}</span>
<span class="key">handCount:</span> <span class="number">${player.handCount}</span>
<span class="key">hand:</span> [${player.hand?.map((c: CardView) => `\n  <span class="string">"${c.name}"</span> <span class="comment">// ${c.id}</span>`).join(',') || ''}
]
<span class="key">graveyard:</span> [${player.graveyard?.map((c: CardView) => `\n  <span class="string">"${c.name}"</span>`).join(',') || ''}
]
<span class="key">manaPool:</span> {
  <span class="key">white:</span> <span class="number">${player.manaPool?.white || 0}</span>,
  <span class="key">blue:</span> <span class="number">${player.manaPool?.blue || 0}</span>,
  <span class="key">black:</span> <span class="number">${player.manaPool?.black || 0}</span>,
  <span class="key">red:</span> <span class="number">${player.manaPool?.red || 0}</span>,
  <span class="key">green:</span> <span class="number">${player.manaPool?.green || 0}</span>,
  <span class="key">colorless:</span> <span class="number">${player.manaPool?.colorless || 0}</span>
}`}</code
                ></pre>
            </div>
          </div>
        {/each}
      {/if}
    </section>

    <!-- Zones -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('zones')}>
        <span class="section-icon">{expandedSections.zones ? '▼' : '▶'}</span>
        <h2>Zones</h2>
      </button>
      {#if expandedSections.zones}
        <div class="zone-grid">
          <div class="zone-block">
            <h3>🏟️ Battlefield ({battlefieldCards.length})</h3>
            <div class="code-block small">
              <pre><code
                  >{battlefieldCards.length > 0
                    ? formatJson(
                        battlefieldCards.map((c: CardView) => ({
                          id: c.id,
                          name: c.name,
                          type: c.type,
                          controllerId: c.controllerId,
                          tapped: c.tapped,
                          power: c.power,
                          toughness: c.toughness
                        }))
                      )
                    : '[]'}</code
                ></pre>
            </div>
          </div>

          <div class="zone-block">
            <h3>📚 Stack ({stackCards.length})</h3>
            <div class="code-block small">
              <pre><code
                  >{stackCards.length > 0
                    ? formatJson(
                        stackCards.map((c: CardView) => ({
                          id: c.id,
                          name: c.name,
                          type: c.type,
                          controllerId: c.controllerId
                        }))
                      )
                    : '[]'}</code
                ></pre>
            </div>
          </div>

          <div class="zone-block">
            <h3>⚔️ Command ({commandCards.length})</h3>
            <div class="code-block small">
              <pre><code
                  >{commandCards.length > 0
                    ? formatJson(
                        commandCards.map((c: CardView) => ({
                          id: c.id,
                          name: c.name,
                          type: c.type
                        }))
                      )
                    : '[]'}</code
                ></pre>
            </div>
          </div>

          <div class="zone-block">
            <h3>🚫 Exile ({gameState.exile?.length || 0})</h3>
            <div class="code-block small">
              <pre><code
                  >{gameState.exile && gameState.exile.length > 0
                    ? formatJson(
                        gameState.exile.map((c: CardView) => ({
                          id: c.id,
                          name: c.name
                        }))
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
      <button class="section-header" onclick={() => toggleSection('clientState')}>
        <span class="section-icon">{expandedSections.clientState ? '▼' : '▶'}</span>
        <h2>Client Store State</h2>
      </button>
      {#if expandedSections.clientState}
        <div class="code-block">
          <pre><code
              >{@html `<span class="comment">// Store meta state</span>
<span class="key">isConnected:</span> <span class="boolean">${gameState.isConnected}</span>
<span class="key">isLoading:</span> <span class="boolean">${gameState.isLoading}</span>
<span class="key">error:</span> <span class="string">${error ? `"${error}"` : 'null'}</span>
<span class="key">selectedCardIds:</span> [${gameState.selectedCardIds.map((id) => `<span class="string">"${id}"</span>`).join(', ')}]
<span class="key">showStack:</span> <span class="boolean">${gameState.showStack}</span>
<span class="key">gameOver:</span> <span class="boolean">${gameState.gameOver}</span>
<span class="key">winner:</span> <span class="string">${gameState.winner ? `"${gameState.winner}"` : 'null'}</span>

<span class="comment">// Pending prompt</span>
<span class="key">pendingPrompt:</span> ${
                prompt
                  ? `{
  <span class="key">type:</span> <span class="string">"${prompt.type}"</span>,
  <span class="key">message:</span> <span class="string">"${prompt.message}"</span>,
  <span class="key">data:</span> ${formatJson(prompt.data)}
}`
                  : 'null'
              }`}</code
            ></pre>
        </div>
      {/if}
    </section>

    <!-- WebSocket State -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('websocket')}>
        <span class="section-icon">{expandedSections.websocket ? '▼' : '▶'}</span>
        <h2>WebSocket State</h2>
      </button>
      {#if expandedSections.websocket}
        <div class="code-block">
          <pre><code
              >{@html `<span class="comment">// WebSocket connection</span>
<span class="key">state:</span> <span class="string">"${$websocketStore.state}"</span>
<span class="key">error:</span> <span class="string">${$websocketStore.error ? `"${$websocketStore.error}"` : 'null'}</span>
<span class="key">lastConnected:</span> <span class="number">${$websocketStore.lastConnected ? new Date($websocketStore.lastConnected).toISOString() : 'null'}</span>
<span class="key">reconnectAttempts:</span> <span class="number">${$websocketStore.reconnectAttempts}</span>`}</code
            ></pre>
        </div>
      {/if}
    </section>

    <!-- Game Messages / Action Log -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('messages')}>
        <span class="section-icon">{expandedSections.messages ? '▼' : '▶'}</span>
        <h2>Game Messages ({gameState.gameView?.messages?.length || 0})</h2>
      </button>
      {#if expandedSections.messages}
        <div class="code-block messages-log">
          {#if gameState.gameView?.messages && gameState.gameView.messages.length > 0}
            {#each gameState.gameView.messages.slice(-50) as msg, idx}
              <div class="log-line">
                <span class="log-idx">{idx + 1}.</span>
                <span class="log-msg">{msg}</span>
              </div>
            {/each}
          {:else}
            <pre><code><span class="comment">// No messages yet</span></code></pre>
          {/if}
        </div>
      {/if}
    </section>

    <!-- Combat State -->
    {#if gameState.gameView?.combat}
      <section class="debug-section">
        <button class="section-header" onclick={() => toggleSection('combat')}>
          <span class="section-icon">{expandedSections.combat ? '▼' : '▶'}</span>
          <h2>⚔️ Combat State</h2>
        </button>
        {#if expandedSections.combat}
          <div class="code-block">
            <pre><code>{formatJson(gameState.gameView.combat)}</code></pre>
          </div>
        {/if}
      </section>
    {/if}

    <!-- Raw Game View JSON -->
    <section class="debug-section">
      <button class="section-header" onclick={() => toggleSection('rawJson')}>
        <span class="section-icon">{expandedSections.rawJson ? '▼' : '▶'}</span>
        <h2>Raw GameView JSON</h2>
      </button>
      {#if expandedSections.rawJson}
        <div class="code-block raw-json">
          <pre><code>{formatJson(gameState.gameView)}</code></pre>
        </div>
      {/if}
    </section>
  </main>
</div>

<style>
  .debug-container {
    min-height: 100vh;
    background: #2d2d2d;
    color: #00ff00;
    font-family: 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
  }

  .debug-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background: #1a1a1a;
    border-bottom: 1px solid #444;
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .back-link {
    color: #888;
    text-decoration: none;
    font-size: 0.875rem;
    transition: color 0.2s;
  }

  .back-link:hover {
    color: #00ff00;
  }

  .admin-link {
    color: #ffff00;
    text-decoration: none;
    font-size: 0.875rem;
    padding: 0.25rem 0.5rem;
    background: rgba(255, 255, 0, 0.1);
    border: 1px solid rgba(255, 255, 0, 0.3);
    border-radius: 4px;
    transition: all 0.2s;
  }

  .admin-link:hover {
    background: rgba(255, 255, 0, 0.2);
    border-color: #ffff00;
  }

  .debug-header h1 {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
    color: #00ff00;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .status-pill {
    padding: 0.375rem 0.75rem;
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 600;
    background: #3a0000;
    color: #ff6b6b;
    border: 1px solid #ff6b6b;
  }

  .status-pill.connected {
    background: #003a00;
    color: #00ff00;
    border-color: #00ff00;
  }

  .update-info {
    display: flex;
    gap: 1rem;
    font-size: 0.75rem;
    color: #888;
  }

  .update-count {
    color: #ffff00;
  }

  .debug-content {
    padding: 1.5rem;
    max-width: 1400px;
    margin: 0 auto;
  }

  .debug-section {
    margin-bottom: 1.5rem;
    background: #1e1e1e;
    border: 1px solid #444;
    border-radius: 8px;
    overflow: hidden;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    padding: 0.875rem 1rem;
    background: #252525;
    border: none;
    color: #00ff00;
    cursor: pointer;
    font-family: inherit;
    font-size: 1rem;
    text-align: left;
    transition: background 0.2s;
  }

  .section-header:hover {
    background: #2a2a2a;
  }

  .section-icon {
    font-size: 0.75rem;
    color: #888;
  }

  .section-header h2 {
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
  }

  .code-block {
    background: #0d0d0d;
    padding: 1rem;
    overflow-x: auto;
  }

  .code-block.small {
    max-height: 200px;
    overflow-y: auto;
  }

  .code-block.raw-json {
    max-height: 500px;
    overflow-y: auto;
  }

  .code-block pre {
    margin: 0;
    font-size: 0.8125rem;
    line-height: 1.5;
  }

  .code-block code {
    color: #00ff00;
  }

  /* JSON syntax highlighting */
  :global(.code-block .key) {
    color: #9cdcfe;
  }

  :global(.code-block .string) {
    color: #ce9178;
  }

  :global(.code-block .number) {
    color: #b5cea8;
  }

  :global(.code-block .boolean) {
    color: #569cd6;
  }

  :global(.code-block .comment) {
    color: #6a9955;
    font-style: italic;
  }

  .player-block {
    border-top: 1px solid #333;
  }

  .player-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: #1a1a1a;
  }

  .player-badge {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.6875rem;
    font-weight: 600;
    background: #333;
    color: #888;
  }

  .player-badge.local {
    background: #003a00;
    color: #00ff00;
  }

  .player-name {
    font-weight: 600;
    color: #fff;
  }

  .zone-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1px;
    background: #333;
  }

  .zone-block {
    background: #1e1e1e;
  }

  .zone-block h3 {
    font-size: 0.875rem;
    font-weight: 600;
    margin: 0;
    padding: 0.75rem 1rem;
    background: #252525;
    border-bottom: 1px solid #333;
    color: #ffff00;
  }

  /* Scrollbar styling */
  .code-block::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }

  .code-block::-webkit-scrollbar-track {
    background: #1a1a1a;
  }

  .code-block::-webkit-scrollbar-thumb {
    background: #444;
    border-radius: 4px;
  }

  .code-block::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  /* Messages log */
  .messages-log {
    max-height: 400px;
    overflow-y: auto;
  }

  .log-line {
    display: flex;
    gap: 0.75rem;
    padding: 0.25rem 0;
    border-bottom: 1px solid #222;
    font-size: 0.75rem;
  }

  .log-line:last-child {
    border-bottom: none;
  }

  .log-idx {
    color: #666;
    min-width: 2.5rem;
    text-align: right;
  }

  .log-msg {
    color: #00ff00;
    word-break: break-word;
  }

  /* Terminal cursor blink effect */
  @keyframes blink {
    0%,
    50% {
      opacity: 1;
    }
    51%,
    100% {
      opacity: 0;
    }
  }

  .debug-header h1::after {
    content: '▋';
    animation: blink 1s infinite;
    margin-left: 4px;
  }
</style>
