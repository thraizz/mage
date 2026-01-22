# Consolidated Magic: The Gathering Simulator Feature List

We need to check if we offer full support for:

## Zone & Card Movement
- [ ] Move cards between zones (Library, Hand, Battlefield, Graveyard, Exile, Command Zone, Stack, Sideboard)
- [ ] Draw X cards
- [ ] Discard specific card
- [ ] Discard random card
- [ ] Search library (with filters by type, cost, etc.)
- [ ] Shuffle library
- [ ] Mill X cards (move from library to graveyard)
- [ ] Scry/Surveil/Explore (look at top X, reorder to top/bottom/graveyard)
- [ ] Look at top X cards
- [ ] Reveal cards (temporary or persistent)
- [ ] Reveal hand
- [ ] Exile face-down
- [ ] Peek at face-down cards or opponent's hand

## Card States
- [ ] Tap/Untap (individual or all)
- [ ] Flip/Transform (double-faced cards, Day/Night)
- [ ] Face up/Face down (Morph)
- [ ] Phase in/Phase out
- [ ] Attach/Detach (Auras, Equipment, Fortifications)
- [ ] Meld (combine two cards)

## Counters on Permanents
- [ ] Add/remove +1/+1 counters
- [ ] Add/remove -1/-1 counters
- [ ] Add/remove loyalty counters (Planeswalkers)
- [ ] Add/remove keyword counters (Flying, Lifelink, Deathtouch, Shield, etc.)
- [ ] Add/remove utility counters (Charge, Oil, Lore, Time, etc.)

## Counters on Players
- [ ] Add/remove poison counters
- [ ] Add/remove energy counters
- [ ] Add/remove experience counters
- [ ] Add/remove rad counters
- [ ] Add/remove ticket counters

## Battlefield Actions
- [ ] Create tokens (custom or predefined)
- [ ] Create emblems
- [ ] Copy spells on stack
- [ ] Copy permanents on battlefield
- [ ] Change control of permanent
- [ ] Track ownership vs control

## Combat
- [ ] Declare attackers (choose creatures and targets)
- [ ] Declare blockers
- [ ] Assign combat damage
- [ ] Apply trample logic
- [ ] Track commander damage per opponent (21-damage tracker)
- [ ] Visual targeting indicators (arrows/lines)
- [ ] Mark creatures as attacking/blocking

## Player Resources & State
- [ ] Change life total
- [ ] Add mana to pool (W, U, B, R, G, C) via optional modal. Players can also just tap lands.
- [ ] Track floating mana
- [ ] Track monarch status
- [ ] Track initiative status
- [ ] Track City's Blessing
- [ ] Track Day/Night cycle

## Turn & Phase Management
- [ ] Progress through turn phases (Untap, Upkeep, Draw, Main 1, Combat, Main 2, End, Cleanup)
- [ ] Pass priority
- [ ] Pass turn
- [ ] Take extra turn
- [ ] Mulligan

## Stack & Spell Resolution
- [ ] Add spells/abilities to stack (LIFO)
- [ ] Resolve stack items
- [ ] Target objects or players
- [ ] Make modal choices

## Commander-Specific
- [ ] Command Zone management
- [ ] Track commander tax (+2 per cast)
- [ ] Commander damage matrix (track per commander to each player)
- [ ] Partner/Background support
- [ ] Dungeon progression tracking

## Randomization
- [ ] Roll dice (D6, D20, N-sided)
- [ ] Flip coins

## UI & Game Tools
- [ ] Game log/action history
- [ ] Priority indicator
- [ ] Turn/phase indicator
- [ ] Card preview on hover
- [ ] Undo functionality
- [ ] Spectator mode
- [ ] Ping/Point at cards
- [ ] Track temporary effects ("until end of turn")
- [ ] Track Storm count/spells cast this turn
- [ ] Concede