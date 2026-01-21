# Game Engine Architecture: Unified Interface for Playtest & Online Modes

**Problem:** Playtest and game pages share components but use different backends:
- **Playtest**: Client-only store (`playtestGameStore`), no network, localStorage persistence
- **Game**: Server-authoritative API (`$lib/api/game.ts`, `$lib/api/direct-actions.ts`), WebSocket updates

**Current Issue:** Same action has different implementations:
```typescript
// Playtest mode
playtestGameStore.moveCardToZone(cardId, 'BATTLEFIELD')

// Game mode
moveCardToZone(gameId, cardId, 'BATTLEFIELD')
```

**Goal:** Create a unified interface so components/utilities don't need to know which mode they're in.

---

## Proposed Architecture: Game Engine Interface

### Layer 1: Core Engine Interface

Define a mode-agnostic interface that both backends implement:

**File: `$lib/engine/game-engine.ts`**

```typescript
import type { CardView, PlayerView, ManaPoolView } from '$lib/generated/mage/v1/models';

/**
 * Zone types matching both playtest and game
 */
export type GameZone =
  | 'BATTLEFIELD'
  | 'HAND'
  | 'GRAVEYARD'
  | 'EXILE'
  | 'LIBRARY'
  | 'LIBRARY_TOP'
  | 'LIBRARY_BOTTOM'
  | 'COMMAND'
  | 'STACK';

/**
 * Counter types for permanents
 */
export type CounterType = 'poison' | 'energy' | string;

/**
 * Unified game engine interface.
 * Both playtest and online game modes implement this.
 */
export interface GameEngine {
  // Mode identification
  readonly mode: 'playtest' | 'game';
  readonly gameId: string;

  // Card operations
  moveCard(cardId: string, zone: GameZone): Promise<void>;
  tapCard(cardId: string, tapped: boolean): Promise<void>;
  flipCard(cardId: string): Promise<void>;

  // Player operations
  modifyLife(playerId: string, delta: number): Promise<void>;
  setCounter(playerId: string, type: CounterType, value: number): Promise<void>;
  drawCards(playerId: string, count: number): Promise<void>;

  // Library operations
  shuffleLibrary(playerId: string): Promise<void>;
  searchLibrary(playerId: string): Promise<void>;

  // Game flow (playtest-only for some)
  untapAll?(playerId: string): Promise<void>;
  nextTurn?(): Promise<void>;
  passPriority?(): Promise<void>;

  // Playtest-only operations (optional interface members)
  createToken?(name: string, types: string[], power: string, toughness: string, color: string): Promise<void>;
  addCounter?(cardId: string, name: string, amount: number): Promise<void>;
  removeCounter?(cardId: string, name: string, amount: number): Promise<void>;
  scryCards?(playerId: string, count: number): Promise<void>;
  millCards?(playerId: string, count: number): Promise<void>;
}

/**
 * Extended engine for playtest-specific features
 */
export interface PlaytestEngine extends GameEngine {
  // Playtest has ALL operations
  createToken(name: string, types: string[], power: string, toughness: string, color: string): Promise<void>;
  addCounter(cardId: string, name: string, amount: number): Promise<void>;
  removeCounter(cardId: string, name: string, amount: number): Promise<void>;
  scryCards(playerId: string, count: number): Promise<void>;
  millCards(playerId: string, count: number): Promise<void>;
  nextTurn(): Promise<void>;
  untapAll(playerId: string): Promise<void>;

  // Session management
  saveSession(): Promise<void>;
  restoreSession(sessionId: string): Promise<void>;
}

/**
 * Check if engine supports playtest features
 */
export function isPlaytestEngine(engine: GameEngine): engine is PlaytestEngine {
  return engine.mode === 'playtest';
}
```

---

### Layer 2: Implementations

#### Playtest Engine Implementation

**File: `$lib/engine/playtest-engine.ts`**

```typescript
import { playtestGameStore } from '$lib/stores/playtest-game';
import type { PlaytestEngine, GameZone, CounterType } from './game-engine';

/**
 * Client-only game engine for playtest mode.
 * Wraps playtestGameStore with unified interface.
 */
export class PlaytestGameEngine implements PlaytestEngine {
  readonly mode = 'playtest' as const;
  readonly gameId: string;

  constructor(gameId: string) {
    this.gameId = gameId;
  }

  // Card operations
  async moveCard(cardId: string, zone: GameZone): Promise<void> {
    playtestGameStore.moveCardToZone(cardId, zone);
  }

  async tapCard(cardId: string, tapped: boolean): Promise<void> {
    playtestGameStore.tapCard(cardId, tapped);
  }

  async flipCard(cardId: string): Promise<void> {
    playtestGameStore.flipCard(cardId);
  }

  // Player operations
  async modifyLife(playerId: string, delta: number): Promise<void> {
    playtestGameStore.modifyLife(playerId, delta);
  }

  async setCounter(playerId: string, type: CounterType, value: number): Promise<void> {
    playtestGameStore.setPlayerCounter(playerId, type, value);
  }

  async drawCards(playerId: string, count: number): Promise<void> {
    playtestGameStore.drawCards(playerId, count);
  }

  // Library operations
  async shuffleLibrary(playerId: string): Promise<void> {
    playtestGameStore.shuffleLibrary(playerId);
  }

  async searchLibrary(playerId: string): Promise<void> {
    // In playtest, this is handled by UI dialog
    // Engine just needs to support the interface
    return Promise.resolve();
  }

  // Game flow
  async untapAll(playerId: string): Promise<void> {
    playtestGameStore.untapAll(playerId);
  }

  async nextTurn(): Promise<void> {
    playtestGameStore.nextTurn();
  }

  async passPriority(): Promise<void> {
    // Not applicable in playtest mode
    return Promise.resolve();
  }

  // Playtest-specific operations
  async createToken(
    name: string,
    types: string[],
    power: string,
    toughness: string,
    color: string
  ): Promise<void> {
    playtestGameStore.createToken(name, types, power, toughness, color);
  }

  async addCounter(cardId: string, name: string, amount: number): Promise<void> {
    playtestGameStore.addCounter(cardId, name, amount);
  }

  async removeCounter(cardId: string, name: string, amount: number): Promise<void> {
    playtestGameStore.removeCounter(cardId, name, amount);
  }

  async scryCards(playerId: string, count: number): Promise<void> {
    playtestGameStore.scryCards(playerId, count);
  }

  async millCards(playerId: string, count: number): Promise<void> {
    playtestGameStore.millCards(playerId, count);
  }

  // Session management
  async saveSession(): Promise<void> {
    // Store auto-saves, but we can expose explicit save
    return Promise.resolve();
  }

  async restoreSession(sessionId: string): Promise<void> {
    playtestGameStore.restoreSession(sessionId);
  }
}
```

#### Server Engine Implementation

**File: `$lib/engine/server-engine.ts`**

```typescript
import type { GameEngine, GameZone, CounterType } from './game-engine';
import {
  moveCardToZone,
  tapUntap,
  drawCards,
  shuffleLibrary,
  searchLibrary,
  modifyLife,
  setPlayerCounter,
  passPriority
} from '$lib/api/direct-actions';

/**
 * Server-authoritative game engine for online play.
 * Wraps API calls with unified interface.
 */
export class ServerGameEngine implements GameEngine {
  readonly mode = 'game' as const;
  readonly gameId: string;

  constructor(gameId: string) {
    this.gameId = gameId;
  }

  // Card operations
  async moveCard(cardId: string, zone: GameZone): Promise<void> {
    await moveCardToZone(this.gameId, cardId, zone);
  }

  async tapCard(cardId: string, tapped: boolean): Promise<void> {
    await tapUntap(this.gameId, cardId, tapped);
  }

  async flipCard(cardId: string): Promise<void> {
    // Check if flip API exists, otherwise use generic card action
    // await flipCard(this.gameId, cardId);
    throw new Error('Flip not yet implemented in server API');
  }

  // Player operations
  async modifyLife(playerId: string, delta: number): Promise<void> {
    await modifyLife(this.gameId, playerId, delta);
  }

  async setCounter(playerId: string, type: CounterType, value: number): Promise<void> {
    await setPlayerCounter(this.gameId, playerId, type, value);
  }

  async drawCards(playerId: string, count: number): Promise<void> {
    await drawCards(this.gameId, playerId, count);
  }

  // Library operations
  async shuffleLibrary(playerId: string): Promise<void> {
    await shuffleLibrary(this.gameId, playerId);
  }

  async searchLibrary(playerId: string): Promise<void> {
    await searchLibrary(this.gameId, playerId);
  }

  // Game flow
  async passPriority(): Promise<void> {
    await passPriority(this.gameId);
  }

  // Playtest-only operations (not supported in server mode)
  // These remain undefined to maintain type safety
}
```

---

### Layer 3: Engine Provider (Svelte Context)

**File: `$lib/engine/engine-context.svelte.ts`**

```typescript
import { getContext, setContext } from 'svelte';
import type { GameEngine } from './game-engine';

const ENGINE_KEY = Symbol('game-engine');

/**
 * Provide game engine to component tree
 */
export function provideEngine(engine: GameEngine): void {
  setContext(ENGINE_KEY, engine);
}

/**
 * Get game engine from context
 */
export function useEngine(): GameEngine {
  const engine = getContext<GameEngine>(ENGINE_KEY);

  if (!engine) {
    throw new Error('Game engine not provided. Did you forget to call provideEngine()?');
  }

  return engine;
}

/**
 * Optional: Get engine if available, return null otherwise
 */
export function useEngineOptional(): GameEngine | null {
  return getContext<GameEngine>(ENGINE_KEY) ?? null;
}
```

---

### Layer 4: Updated Utilities Using Engine

Now utilities can use the engine interface instead of mode-specific code:

**File: `$lib/utils/game-actions.ts` (Updated)**

```typescript
import { toast } from '$lib/stores/toast';
import type { GameEngine } from '$lib/engine/game-engine';

export interface ActionOptions {
  engine: GameEngine;
  loadingState?: { set: (value: boolean) => void };
  successMessage?: string;
  errorPrefix?: string;
  silent?: boolean;
  onSuccess?: () => void;
  onError?: (error: string) => void;
}

/**
 * Execute game action through engine interface.
 * Works with both playtest and server engines.
 */
export async function executeGameAction(
  action: (engine: GameEngine) => Promise<void>,
  options: ActionOptions
): Promise<boolean> {
  const {
    engine,
    loadingState,
    successMessage,
    errorPrefix = 'Failed',
    silent = false,
    onSuccess,
    onError
  } = options;

  if (loadingState) {
    loadingState.set(true);
  }

  try {
    await action(engine);

    if (!silent && successMessage) {
      toast.success(successMessage);
    }

    onSuccess?.();
    return true;
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    console.error(`${errorPrefix}:`, err);

    if (!silent) {
      toast.error(`${errorPrefix}: ${message}`);
    }

    onError?.(message);
    return false;
  } finally {
    if (loadingState) {
      loadingState.set(false);
    }
  }
}
```

**File: `$lib/utils/keyboard-shortcuts.ts` (Updated)**

```typescript
import type { GameEngine } from '$lib/engine/game-engine';

export interface KeyboardShortcut {
  key: string;
  handler: (engine: GameEngine) => void | Promise<void>;
  requiresHoveredCard?: boolean;
  description: string;
  category?: 'general' | 'card' | 'game' | 'ui';
}

export function createKeyboardHandler(
  engine: GameEngine,
  shortcuts: KeyboardShortcut[],
  getHoveredCard?: () => string | null
): (event: KeyboardEvent) => void {
  return (event: KeyboardEvent) => {
    if (
      event.target instanceof HTMLInputElement ||
      event.target instanceof HTMLTextAreaElement
    ) {
      return;
    }

    const key = event.key.toLowerCase();
    const hoveredCard = getHoveredCard?.();

    for (const shortcut of shortcuts) {
      if (shortcut.key === key) {
        if (shortcut.requiresHoveredCard && !hoveredCard) {
          continue;
        }

        event.preventDefault();
        shortcut.handler(engine);
        return;
      }
    }
  };
}
```

---

### Layer 5: Page Integration

#### Playtest Page

```typescript
import { provideEngine } from '$lib/engine/engine-context.svelte';
import { PlaytestGameEngine } from '$lib/engine/playtest-engine';
import { executeGameAction } from '$lib/utils/game-actions';
import { createKeyboardHandler, type KeyboardShortcut } from '$lib/utils/keyboard-shortcuts';

// Create and provide engine
const engine = new PlaytestGameEngine(gameId);
provideEngine(engine);

// Define keyboard shortcuts using engine
const shortcuts: KeyboardShortcut[] = [
  {
    key: 'c',
    handler: async (engine) => {
      if (!me) return;
      await executeGameAction(
        (eng) => eng.drawCards(me.playerId, 1),
        { engine, successMessage: 'Drew a card' }
      );
    },
    description: 'Draw card',
    category: 'game'
  },
  {
    key: 'v',
    handler: async (engine) => {
      if (!me) return;
      await executeGameAction(
        (eng) => eng.shuffleLibrary(me.playerId),
        { engine, successMessage: 'Shuffled library' }
      );
    },
    description: 'Shuffle library',
    category: 'game'
  },
  {
    key: 'd',
    handler: async (engine) => {
      if (hoveredCardId) {
        await engine.moveCard(hoveredCardId, 'GRAVEYARD');
      }
    },
    requiresHoveredCard: true,
    description: 'Move to graveyard',
    category: 'card'
  }
];

const handleGlobalKeydown = createKeyboardHandler(engine, shortcuts, () => hoveredCardId);
```

#### Game Page

```typescript
import { provideEngine } from '$lib/engine/engine-context.svelte';
import { ServerGameEngine } from '$lib/engine/server-engine';
import { executeGameAction } from '$lib/utils/game-actions';
import { createKeyboardHandler, type KeyboardShortcut } from '$lib/utils/keyboard-shortcuts';

// Create and provide engine
const engine = new ServerGameEngine(gameId);
provideEngine(engine);

// Same shortcuts work! Engine handles the difference
const shortcuts: KeyboardShortcut[] = [
  {
    key: 'c',
    handler: async (engine) => {
      const me = $localPlayer;
      if (!me) return;
      await executeGameAction(
        (eng) => eng.drawCards(me.playerId, 1),
        { engine, loadingState: { set: (v) => isActionLoading = v }, successMessage: 'Drew a card' }
      );
    },
    description: 'Draw card',
    category: 'game'
  },
  // ... same shortcuts as playtest!
];

const handleGlobalKeydown = createKeyboardHandler(engine, shortcuts, () => hoveredCardId);
```

---

## Migration Strategy

### Phase 0: Prepare Engine Layer (FIRST)

1. Create `$lib/engine/` directory
2. Implement `game-engine.ts` interface
3. Implement `playtest-engine.ts` adapter
4. Implement `server-engine.ts` adapter
5. Implement `engine-context.svelte.ts` provider
6. Write tests for both engines

### Phase 1: Migrate Playtest Page

1. Instantiate `PlaytestGameEngine`
2. Provide engine to component tree
3. Update keyboard shortcuts to use engine
4. Update action handlers to use engine
5. Test all features work identically

### Phase 2: Migrate Game Page

1. Instantiate `ServerGameEngine`
2. Provide engine to component tree
3. Update keyboard shortcuts to use engine
4. Update action handlers to use engine
5. Test all features work identically

### Phase 3: Extract Safe Patterns

Now that both pages use engine interface, extract shared utilities:
1. Keyboard shortcuts (Phase 1 from previous plan)
2. Async action wrapper (Phase 2)
3. Drop zone helpers (Phase 3)
4. Clipboard utilities (Phase 4)
5. Derived state patterns (Phase 6)
6. Dialog state management (Phase 7)

**Skip drag threshold extraction** - it's complex and risky.

---

## Testing Strategy

### Unit Tests for Engines

**File: `$lib/engine/playtest-engine.test.ts`**

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { PlaytestGameEngine } from './playtest-engine';
import { playtestGameStore } from '$lib/stores/playtest-game';

describe('PlaytestGameEngine', () => {
  let engine: PlaytestGameEngine;

  beforeEach(() => {
    engine = new PlaytestGameEngine('test-game');
    playtestGameStore.initialize('test-game', [
      { playerId: 'p1', name: 'Player 1', /* ... */ }
    ]);
  });

  it('should move card to zone', async () => {
    await engine.moveCard('card-1', 'GRAVEYARD');
    // Assert card moved
  });

  it('should modify life', async () => {
    await engine.modifyLife('p1', -5);
    // Assert life changed
  });

  // ... more tests
});
```

**File: `$lib/engine/server-engine.test.ts`**

```typescript
import { describe, it, expect, vi } from 'vitest';
import { ServerGameEngine } from './server-engine';
import * as directActions from '$lib/api/direct-actions';

describe('ServerGameEngine', () => {
  let engine: ServerGameEngine;

  beforeEach(() => {
    engine = new ServerGameEngine('game-123');
  });

  it('should call API to move card', async () => {
    const spy = vi.spyOn(directActions, 'moveCardToZone').mockResolvedValue(undefined);

    await engine.moveCard('card-1', 'GRAVEYARD');

    expect(spy).toHaveBeenCalledWith('game-123', 'card-1', 'GRAVEYARD');
  });

  // ... more tests
});
```

### Integration Tests

Test that both engines work with same utility code:

```typescript
describe('Keyboard shortcuts with engines', () => {
  it('works with playtest engine', () => {
    const engine = new PlaytestGameEngine('test');
    // Test shortcuts
  });

  it('works with server engine', () => {
    const engine = new ServerGameEngine('game-123');
    // Test shortcuts
  });
});
```

---

## Benefits of This Architecture

1. **Type Safety**: TypeScript enforces interface compliance
2. **Testability**: Mock engines for component tests
3. **Separation of Concerns**: Components don't know about stores/APIs
4. **Future-Proof**: Can add new modes (e.g., replay viewer)
5. **Safe Extraction**: Utilities use engine interface, work with both modes
6. **Fail-Fast**: Missing API methods caught at compile time

## Anti-Patterns This Prevents

❌ **Accidentally calling playtest store in game mode**
```typescript
// OLD - ERROR PRONE
playtestGameStore.moveCardToZone(cardId, zone) // Oops! In game mode!

// NEW - TYPE SAFE
engine.moveCard(cardId, zone) // Works in both modes
```

❌ **Duplicating action logic**
```typescript
// OLD
function handleDrawCard() {
  if (mode === 'playtest') {
    playtestGameStore.drawCards(playerId, 1);
  } else {
    await drawCards(gameId, playerId, 1);
  }
}

// NEW
function handleDrawCard() {
  await engine.drawCards(playerId, 1);
}
```

❌ **Breaking one mode when refactoring the other**
```typescript
// With interface, both modes MUST implement same methods
// TypeScript compiler ensures consistency
```

---

## Summary

**Architecture Layers:**
1. **Interface** (`game-engine.ts`) - Defines contract
2. **Implementations** (`playtest-engine.ts`, `server-engine.ts`) - Mode-specific adapters
3. **Context** (`engine-context.svelte.ts`) - Dependency injection
4. **Utilities** (keyboard, actions, etc.) - Use engine interface
5. **Pages** - Provide correct engine, everything else is shared

**Migration Order:**
1. Build engine layer + tests
2. Migrate playtest page to use engine
3. Migrate game page to use engine
4. Extract shared utilities (now safe!)

**Result:**
- Both modes guaranteed to work
- Shared code extraction is safe
- Type safety prevents mode-mixing bugs
- Components are mode-agnostic
- Easy to add new modes (replay, spectator, etc.)
