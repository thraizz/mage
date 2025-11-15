# MageEngine Logical Differences: Go vs Java

This document compares the Go implementation (`mage-server-go/internal/game/mage_engine.go`) against the Java implementation (`Mage/src/main/java/mage/game/GameImpl.java`) to identify logical differences.

## Critical Differences

### 1. Priority Retention After Casting

**Java Implementation:**
- When a player casts a spell, they **retain priority** by default
- The caster can continue to cast spells or activate abilities without passing
- Priority only passes to the next player when the caster explicitly passes

**Go Implementation:**
- After casting a spell, priority is **immediately passed** to the next player
- Comment in code (line 621-623) acknowledges this simplification: "In strict MTG rules, caster retains priority, but for simplicity we pass priority"
- The caster's `Passed` flag is set to `true` immediately after casting (line 624)

**Impact:** This changes the fundamental priority flow and prevents players from casting multiple spells in sequence without passing priority.

### 2. Stack Resolution Logic

**Java Implementation:**
- Resolves **one item at a time** in a loop
- After each resolution:
  - Checks state-based actions (`checkStateBasedActions()`)
  - Handles triggered abilities (`checkTriggered()`)
  - Processes simultaneous events
  - Applies effects
- Continues until stack is empty

**Go Implementation:**
- Resolves **all items sequentially** without checking between resolutions
- Only checks legality before each resolution
- No state-based action checks between resolutions
- No triggered ability processing between resolutions

**Impact:** This can miss state-based actions and triggered abilities that should occur between stack item resolutions, potentially breaking game rules.

### 3. Pass Checking Logic

**Java Implementation (`allPassed()`):**
```java
protected boolean allPassed() {
    for (Player player : state.getPlayers().values()) {
        if (!player.isPassed() && player.canRespond()) {
            return false;
        }
    }
    return true;
}
```
- Only considers players who `canRespond()` (not lost, not left, can take actions)
- Players who have lost or left are automatically considered as "passed"

**Go Implementation (`handlePass()`):**
```go
allPassed := true
activePlayers := 0
for _, pid := range gameState.playerOrder {
    p := gameState.players[pid]
    if !p.Lost && !p.Left {
        activePlayers++
        if !p.Passed {
            allPassed = false
        }
    }
}
```
- Checks all players who haven't lost or left
- Doesn't check `canRespond()` equivalent
- Has special case: "If only one active player, they can't pass to themselves" (line 476-478)

**Impact:** The Go implementation may incorrectly handle edge cases where players can't respond but haven't explicitly passed.

### 4. State-Based Actions and Triggered Abilities

**Java Implementation:**
- Calls `checkStateAndTriggered()` **before each priority** (line 1735)
- This method:
  1. Checks state-based actions repeatedly until none occur
  2. Handles triggered abilities and puts them on the stack
  3. Repeats until no more state-based actions or triggers occur
- Ensures game state is clean before giving priority

**Go Implementation:**
- **No equivalent** to `checkStateAndTriggered()`
- No systematic state-based action checking
- Triggered abilities are created immediately when spells are cast (line 618)
- No pre-priority state cleanup

**Impact:** Missing state-based actions (like creatures dying from damage, players losing from 0 life, etc.) can cause game state inconsistencies.

### 5. Reset Passed Logic

**Java Implementation (`resetPassed()`):**
```java
public void resetPassed() {
    this.passed = this.loses || this.hasLeft();
}
```
- Sets `passed = true` if player has lost or left
- Otherwise sets `passed = false`
- Called after stack resolution or phase advancement

**Go Implementation:**
```go
for _, p := range gameState.players {
    p.Passed = false
}
```
- Simply sets all players' `Passed = false`
- Doesn't account for lost/left players

**Impact:** Lost or left players might incorrectly get priority again.

### 6. Stack Resolution Triggering

**Java Implementation:**
- Stack resolves when `allPassed()` returns true AND stack is not empty
- Happens in the main priority loop (line 1763-1772)
- After resolution, resets passed flags and continues priority loop

**Go Implementation:**
- Stack resolves when all players have passed (line 480-484, 510-512, 642-644)
- Multiple code paths can trigger resolution
- After resolution, resets passed flags and sets priority to active player

**Impact:** The Go implementation has multiple code paths that can trigger stack resolution, which could lead to inconsistent behavior.

### 7. Triggered Ability Timing

**Java Implementation:**
- Triggered abilities are queued when events occur
- They are put on the stack **before priority is given** (via `checkTriggered()`)
- Handled in APNAP order (Active Player, Non-Active Player)

**Go Implementation:**
- Triggered abilities are created **immediately** when a spell is cast (line 618)
- Pushed directly to stack without queuing
- No APNAP ordering

**Impact:** Triggered abilities may resolve in incorrect order or at incorrect times.

### 8. Priority Loop Structure

**Java Implementation (`playPriority()`):**
- Outer loop: continues until game paused, ended, or turn ends
- Inner loop: cycles through players until all pass
- For each player:
  - Checks state and triggers
  - Applies effects
  - Gives priority
  - If player acts, processes simultaneous events and applies effects
- After all pass: resolves stack if not empty, or advances step/phase

**Go Implementation (`handlePass()`):**
- Single method handles pass action
- Checks if all passed
- If yes: resolves stack or advances step/phase
- If no: passes priority to next player
- No systematic loop structure

**Impact:** The Go implementation lacks the comprehensive priority loop structure that ensures proper game state management.

### 9. Card Zone Management After Resolution

**Java Implementation:**
- Stack resolution removes items from stack
- Cards are moved to appropriate zones (graveyard, exile, etc.) by the resolve method
- Zone changes are tracked and events are fired

**Go Implementation:**
- Stack resolution pops items and calls resolve functions
- Cards are moved to zones in resolve functions
- Some zone tracking may be incomplete (e.g., cards removed from stack but zone not always updated correctly)

**Impact:** Cards might end up in incorrect zones or have inconsistent zone tracking.

### 10. Error Handling and Rollback

**Java Implementation:**
- Has rollback mechanism (`bookmarkState()`, `restoreState()`)
- Can restore to previous state on errors
- Error counting and limits
- Comprehensive error handling in priority loop

**Go Implementation:**
- No rollback mechanism
- Errors are logged but game continues
- No state restoration capability

**Impact:** Errors can leave game in inconsistent state with no recovery mechanism.

## Summary of Missing Features in Go Implementation

1. ✅ Priority retention after casting
2. ✅ State-based action checking (`checkStateBasedActions()`)
3. ✅ Systematic triggered ability handling (`checkTriggered()`)
4. ✅ Pre-priority state cleanup (`checkStateAndTriggered()`)
5. ✅ Proper `canRespond()` checking in pass logic
6. ✅ Rollback/restore state mechanism
7. ✅ APNAP ordering for triggered abilities
8. ✅ Comprehensive priority loop structure
9. ✅ Simultaneous event handling between stack resolutions
10. ✅ Proper lost/left player handling in resetPassed

## Recommendations

1. **Implement priority retention**: Allow casters to retain priority after casting spells
2. **Add state-based action checking**: Implement `checkStateBasedActions()` and call it before priority
3. **Add triggered ability queue**: Queue triggered abilities and process them before priority (APNAP order)
4. **Fix resetPassed logic**: Account for lost/left players when resetting passed flags
5. **Improve stack resolution**: Add state-based action and triggered ability checks between resolutions
6. **Add canRespond checking**: Check if players can respond before giving them priority
7. **Implement rollback mechanism**: Add state bookmarking and restoration for error recovery
8. **Unify stack resolution paths**: Consolidate multiple code paths that trigger stack resolution
