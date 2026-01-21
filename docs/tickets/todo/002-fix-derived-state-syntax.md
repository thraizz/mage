# Fix Incorrect $derived Syntax in Svelte 5 Components

## Problem

In Svelte 5, when using `$derived` with a function that has a block body (curly braces), you must use `$derived.by(() => { ... })` instead of `$derived(() => { ... })`. 

Using `$derived(() => { ... })` creates a derived state that **returns a function** rather than evaluating the function and returning its result. This causes the function code to leak into production builds, appearing as raw JavaScript in the rendered output.

### Example Bug

**Incorrect:**
```svelte
const nextHandSize = $derived(() => {
    if (nextMulliganCount <= freeMulligans) {
        return 7;
    } else {
        const penaltyMulligans = nextMulliganCount - freeMulligans;
        return Math.max(0, 7 - penaltyMulligans);
    }
});
```

**Correct:**
```svelte
const nextHandSize = $derived.by(() => {
    if (nextMulliganCount <= freeMulligans) {
        return 7;
    } else {
        const penaltyMulligans = nextMulliganCount - freeMulligans;
        return Math.max(0, 7 - penaltyMulligans);
    }
});
```

### When to Use Each Syntax

- **`$derived(expression)`** - For simple expressions (no block body)
  ```svelte
  const count = $derived(items.length);
  const isActive = $derived(status === 'active');
  ```

- **`$derived(() => expression)`** - For simple arrow function expressions (no block body)
  ```svelte
  const filtered = $derived(() => items.filter(i => i.active));
  ```

- **`$derived.by(() => { ... })`** - For complex logic with block bodies (if/else, loops, multiple statements)
  ```svelte
  const result = $derived.by(() => {
      if (condition) return value1;
      return value2;
  });
  ```

## Affected Files

### High Priority (Already Causing Production Issues)

1. **`mage-client-web/src/lib/components/game/MulliganDialog.svelte`** (Line 30)
   - **Issue:** `nextHandSize` - Shows raw function code in button text
   - **Fix:** Change `$derived(() => { ... })` to `$derived.by(() => { ... })`

### All Files Requiring Fixes

#### Game Components

2. **`mage-client-web/src/routes/(protected)/game/[id]/+page.svelte`**
   - **Line 170:** `combatPromptOptions` - Complex logic with multiple if statements
   - **Line 189:** `damageAssignmentPrompt` - Conditional logic with type checking
   - **Line 287:** `selectedOpponent` - Multiple if statements
   - **Line 295:** `opponentBattlefield` - Conditional with filter
   - **Line 308:** `opponentCommandCards` - Conditional with filter

3. **`mage-client-web/src/lib/components/game/XManaSelector.svelte`** (Line 28)
   - **Issue:** `quickAmounts` - Multiple if statements building an array

4. **`mage-client-web/src/lib/components/game/ManaPool.svelte`** (Line 49)
   - **Issue:** `sizeClass` - Switch statement

5. **`mage-client-web/src/lib/components/game/Card.svelte`** (Line 280)
   - **Issue:** `sizeClasses` - Switch statement

6. **`mage-client-web/src/lib/components/game/Graveyard.svelte`** (Line 150)
   - **Issue:** `filteredCards` - Conditional logic with filter

7. **`mage-client-web/src/lib/components/game/ExileZone.svelte`** (Line 152)
   - **Issue:** `filteredCards` - Conditional logic with filter

8. **`mage-client-web/src/lib/components/game/LibrarySearch.svelte`** (Line 63)
   - **Issue:** `filteredCards` - Multiple filter conditions

9. **`mage-client-web/src/lib/components/game/PlaytestLibrarySearch.svelte`** (Line 45)
   - **Issue:** `filteredCards` - Conditional logic with filter

10. **`mage-client-web/src/lib/components/game/DeclareBlockers.svelte`**
    - **Line 48:** `cardNames` - For loop building a Map
    - **Line 57:** `attackersById` - For loop building a Map
    - **Line 66:** `availableBlockerIds` - For loop building a Set
    - **Line 197:** `blocksByAttacker` - For loop building a Map

11. **`mage-client-web/src/lib/components/game/DeclareAttackers.svelte`**
    - **Line 44:** `cardNames` - For loop building a Map
    - **Line 161:** `availableAttackerIds` - For loop building a Set

12. **`mage-client-web/src/lib/components/game/AssignDamage.svelte`** (Line 59)
    - **Issue:** `totalAssigned` - For loop calculating sum

13. **`mage-client-web/src/lib/components/GameStartCountdown.svelte`** (Line 77)
    - **Issue:** `countdownMessage` - Multiple if/else statements

#### Modal Components

14. **`mage-client-web/src/lib/components/CreateTableModal.svelte`** (Line 55)
    - **Issue:** `isValid` - Simple validation (could potentially be simplified to expression)

15. **`mage-client-web/src/lib/components/PlaytestModal.svelte`** (Line 138)
    - **Issue:** `isValid` - Multiple validation checks

#### Page Components

16. **`mage-client-web/src/routes/(protected)/lobby/+page.svelte`** (Line 156)
    - **Issue:** `filteredTables` - Multiple filter conditions

17. **`mage-client-web/src/routes/(protected)/playtest/+page.svelte`**
    - **Line 155:** `selectedOpponent` - Multiple if statements
    - **Line 166:** `opponentBattlefield` - Conditional with filter
    - **Line 193:** `selectedCardForCountersData` - Multiple conditional lookups
    - **Line 216:** `activePlayerName` - Find operation with fallback
    - **Line 222:** `turnNumber` - Complex calculation with Math operations
    - **Line 242:** `opponentCommandCards` - Conditional with filter

## Fix Pattern

For each affected line, change:

```svelte
const variableName = $derived(() => {
    // ... block body with statements
});
```

To:

```svelte
const variableName = $derived.by(() => {
    // ... block body with statements
});
```

## Testing Checklist

After fixing each file, verify:

- [ ] No raw JavaScript code appears in rendered output
- [ ] Derived values evaluate correctly in the browser
- [ ] Reactivity still works (values update when dependencies change)
- [ ] No console errors related to derived state
- [ ] Production build works correctly (test with `npm run build && npm run preview`)

## Priority

1. **Critical:** `MulliganDialog.svelte` - Already causing visible production bug
2. **High:** All game components that render user-facing text
3. **Medium:** Modal validation logic
4. **Low:** Internal derived state that doesn't render to UI

## Notes

- Simple arrow functions without block bodies (e.g., `$derived(() => items.filter(...))`) are **correct** and don't need changes
- Only functions with block bodies `{ ... }` need `.by()`
- This is a Svelte 5 specific requirement - Svelte 4 used different syntax
