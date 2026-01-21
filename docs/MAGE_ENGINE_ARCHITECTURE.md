# Mage Engine Architecture Documentation

> **File**: `mage-server-go/internal/game/direct_engine.go`  
> **Purpose**: Rules-light Magic: The Gathering game engine - assists players without enforcing rules

## Overview

The game engine should be a rules-light game engine that **assists** players rather than **enforces** rules. It provides:

- **State tracking**: Zones, life totals, counters, combat
- **Direct manipulation**: Players control all game state through UI
- **Action logging**: Every change is logged for review
- **Rollback support**: Any action can be undone
- **Real-time sync**: WebSocket notifications keep all clients updated

### What This Engine Does NOT Do

Unlike traditional MTG engines, we intentionally **do not**:

- Validate spell timing or mana costs
- Enforce priority or turn structure
- Check creature abilities (flying, menace, etc.)
- Auto-resolve state-based actions
- Process triggered abilities
- Validate targets
- Calculate combat damage automatically

Players are trusted to play correctly. The engine is a shared game board, not a referee.

---

## Design Philosophy

### Rules-Light Approach

Inspired by platforms like Untap.in, this engine follows the principle:

> **"Assist, don't enforce"**

| Traditional Engine     | Rules-Light Engine               |
| ---------------------- | -------------------------------- |
| Validates every action | Logs every action                |
| Rejects illegal plays  | Allows any play, can rollback    |
| Auto-resolves triggers | Players manually resolve         |
| Enforces priority      | Suggests phases, players control |
| Complex combat rules   | Simple attacker/blocker tracking |

### Benefits

1. **Simplicity**: We don't enforce rules, so we don't need to write them
2. **Flexibility**: House rules, casual play, testing
3. **Speed**: No validation overhead
4. **Reliability**: Fewer edge cases = fewer bugs
5. **Maintainability**: Easy to understand and modify
