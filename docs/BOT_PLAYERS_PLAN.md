# LLM Bot Players — Phased Implementation Plan

Bot opponents for `mage-server-go`, reusing mage-bench (a fork of XMage that gives LLMs a
tool interface to play Magic). MVP = local playtesting bots. v2 = online NPCs.

Each phase is self-contained and executable in a fresh chat context. Read Phase 0 first,
every time — it is the "Allowed APIs" contract.

---

## Phase 0: Discovery Findings (READ FIRST, EVERY PHASE)

Consolidated from three discovery agents. Everything below was verified against source.
Treat this section as authoritative over recollection.

### 0.1 Upstream: mage-bench

Repo `GregorStocks/mage-bench`, default branch `master`. A fork of XMage — the same upstream
project this repo descends from (`go.mod` module path is `github.com/magefree/mage-server-go`).

**License: two stacked MIT grants in `LICENSE.txt`** — mage-bench code © 2026 Gregor Stocks,
inherited XMage code © 2010 betasteward@gmail.com. GitHub reports "NOASSERTION" only because
two licenses share one file. Copying is permitted; both notices must travel with copied
material. Card oracle text is Wizards of the Coast IP under neither grant — mage-bench pulls
it from Scryfall at runtime, as this repo does via `scryfall_cards`.

**Files to vendor (highest value first):**

| Path | Lines | What it is |
|---|---|---|
| `website/src/data/mcp-tools.json5` | 1640 | Generated, CI-verified wire-format schemas for all 10 tools |
| `puppeteer/prompts/default.md` | 66 | The complete system prompt |
| `src/magebench/game/decision_renderer.py` | 767 | The `## Decision` board serializer |
| `src/magebench/pilot/pilot_rendering.py` | 369 | Context windowing, cache breakpoints, trailing lines |
| `src/magebench/pilot/pilot_recovery.py` | 193 | Error-recovery primitives |
| `src/magebench/pilot/pilot_state.py` | 97 | Loop state + context reset |
| `Mage/src/main/java/mage/util/ShortIdRegistry.java` | 203 | Short-ID allocation, near-zero deps |
| `src/magebench/common/json5_utils.py` | 69 | Diff-friendly JSON5 serializer for golden files |
| `src/magebench/game/harness_epoch.py` | ~80 | 60-entry changelog of every breaking design change — read as history |

Fetch with `gh api repos/GregorStocks/mage-bench/contents/<path> --jq '.content' | base64 -d`.

**Read for reference, do not copy:** `src/magebench/pilot/pilot.py` (too coupled; mine the
constants block at `:64-81`), `doc/golden-prompts.md` (~15% stale), `tests/golden_helpers.py`
(2029 lines; copy only `_strip_volatile`, `_normalize_prompt_for_golden`, `_json_diff`,
`assert_golden_prompt` from `:1617-1890`).

**Known-stale claims in their docs** — do not follow: `doc/golden-prompts.md:254` says short
IDs start at p3 with players at p1/p2 (contradicted by `GameView.java:377-386`, where players
are assigned *last*); `doc/golden-prompts.md:188` describes a `[redacted]` normalization that
no longer exists in `golden_helpers.py` (it moved to Java).

### 0.2 The board-state format (the thing we are porting)

Produced by `decision_renderer.py::_render_decision_block`. Complete field list, in emission
order — this is the full contract, not an example:

```
## Card Reference
- <Name> <mana_cost> -- <type_line> <P/T>: <oracle text with \n replaced by " / ">

## Decision

[Decision {index}, snapshot={n}] Turn {turn} {phase} - {player}
  Board: {player_line} | {player_line} | ...
  Stack: [{item}, {item}]                     (only when stack non-empty)
  Combat: {groups}                            (only when combat groups exist)
  Combat Phase: {combat_phase}
  Incoming Attackers: {...}                   (only at declare-blockers, non-empty)
  Untapped lands: {n}, Land drops remaining: {n}
  Message: {message}                          (always, may be empty)
  Choices ({n}): {name} [id=p3, {action}, {mana_cost}], ...
    -- OR --
  Items ({n}): total={x} | total_min=..,total_max=..
    {i}: {desc} [min=.., max=..]
  NOTE: ...                                   (only when message contains "Pick triggered ability")
  Respond: {respond_with}                     (or "Response type: {t}")
  Mana pool: {C}={n}, ...                     (only when non-zero)
  Recent chat: {msg} | {msg}
```

Player line: `{name}: {life}hp hand=[{cards}] lib={n}[ {counter}={v}][ bf=[..]][ gy=[..]][ exile=[..]]`
— own seat gets `hand=[...]`, opponents get `hand={count}` (omitted when zero).

Permanent display: `Name P/T (tapped, sick, face_down, loyalty=N, {counter}={n}, copy of X, token)`.

`phase` renders as the literal `PREGAME` when empty.

Card Reference excludes basic lands (an 11-name frozenset incl. Wastes and Snow-Covered
basics) and dedupes: oracle text appears on a card's **first** appearance only, tracked in a
`seen_oracle_cards` set. That dedup is the single largest token saving in the design.

### 0.3 The tool surface

10 tools exist; only 7 are exposed to the model (`puppeteer/toolsets.json`):
`pass_priority`, `get_action_choices`, `choose_action`, `get_game_state`, `get_game_log`,
`get_oracle_text`, `send_chat_message`. (`concede`, `get_game_history`, `join_table` are not.)

Only `send_chat_message.message` is a **required** parameter anywhere. Every other parameter
on every tool is optional — the model is never forced to supply a field it cannot reason about.

`choose_action` takes nine optional params: `choice` (string — `"p3"` | `"0"` | `"yes"` |
`"no"`), `amount` (int), `amounts` (int[]), `pile` (int), `text` (string), `mana_plan`
(string, e.g. `"p1,p5:1"`), `auto_tap` (bool, default true), `attackers` (string, comma IDs
or `"all"`), `blockers` (string, `"blocker:attacker"` pairs). `choice` parsing:
`yes|true` → answer true; `no|false` → answer false; parseable int → index; else → id.

`pass_priority` / `get_action_choices` take `until` (enum of 11 step names) and `board_cursor`
(int — omits the board from the response when unchanged, a token optimization).

Results carry ~35 optional fields; the ones that drive the loop are `action_pending`,
`response_type`, `respond_with`, `stop_reason`, `game_over`, `player_dead`. `respond_with`
is a per-decision-type instruction string generated in
`processor/BridgePublishedQueryBuilder.java` — 13 distinct strings, e.g.
`"choice=pN to play, or choice=no to pass"`, `"amounts=[N,N,...] — one per item, sum between total_min and total_max"`.

Schemas are **reflected from Java annotations** (`@Tool`, `@Param`, `@ResultField`) by
`McpToolRegistry.java`, always emitting `additionalProperties: false`, and CI-verified stale
by `make verify-mcp-tools`. We are porting the *output*, not the reflection machinery.

### 0.4 Their loop, and its measured constants

From `pilot.py:64-81` — these are empirical, arrived at over 60 harness epochs:

```
MAX_TOKENS = 20_000
LLM_REQUEST_TIMEOUT_SECS = 120       # epoch #15 raised this from 45s; 45 was not enough
MAX_CONSECUTIVE_TIMEOUTS = 3
MAX_CONSECUTIVE_EMPTY_CHOICES = 5
MAX_GAME_DURATION_SECS = 3 * 3600
MAX_TURNS_WITHOUT_PROGRESS = 20
MAX_CONSECUTIVE_PASS_ERRORS = 3
MAX_CONSECUTIVE_TRUNCATIONS = 3
MAX_CONSECUTIVE_EMPTY_ERRORS = 10
MAX_EMPTY_RESPONSES = 10
MAX_CHAT_MESSAGES_PER_TURN = 2
INFO_ONLY_TOOLS = {"get_game_state", "get_oracle_text", "send_chat_message"}
```

**There is no client-side timeout on the blocking tool call**:
`httpx.Timeout(30.0, read=None)` (`bridge_transport.py:152`). The Java side uses virtual
threads so a blocked `pass_priority` does not stall concurrent requests. The 120s timeout
applies to the *LLM request*, not the game-side block. These are two different timers and
conflating them is the main design trap.

**Context is windowed, not unbounded** (`pilot_rendering.py:20-25, 275-343`): last 40 messages
verbatim, the 20 before that with tool results >200 chars summarized, and a synthetic
"state bridge" user message between them carrying a `get_game_state` summary refreshed every
5 renders. Two `cache_control` breakpoints — one on the system prompt, one located by scanning
backwards for the literal marker `"All cards listed are playable right now."`

Degradation is graceful and in-character: on repeated failure the pilot chats
*"My brain is fried... going on autopilot for the rest of this game. GG!"* and switches to an
auto-pass loop rather than hanging or crashing the table.

### 0.5 Local engine reality — the constraint that shapes everything

`ProcessAction` (`internal/game/game_engine.go:150`) is a thin dispatcher over 20 direct state
mutations. **No legality checks, no priority, no stack resolution, no combat, no mana payment,
no turn structure beyond a counter.** Verified: `internal/game/actions.go` imports only `fmt`,
`math/rand`, `strings`, `time`.

Every subsystem that could supply legality is **orphaned** — zero non-test importers outside
its own package:

| Package | LOC | External importers |
|---|---|---|
| `internal/game/mana` | 1098 | 0 |
| `internal/game/targeting` | 364 | 0 |
| `internal/game/cards` | 579 | 0 (2 test-only) |
| `internal/game/token` | 13462 | 0 |
| `internal/game/abilities` | ~6100 | 1 (`ability_registry.go`, itself unreferenced) |
| `internal/game/counters` | 463 | 1 (`card.go` `LegacyCard`, deprecated) |

Three mutually incompatible type hierarchies exist for the same concepts: two mana-cost
parsers (`mana.ParseCost` vs `abilities.ParseManaCost`), two targeting systems
(`targeting.TargetRequirement` vs `abilities.TargetRequirement`), three card types
(`game.Card`, `game.LegacyCard`, `abilities.internalCard`). `abilities.SpellAbility.CanActivate`
returns `true` unconditionally (`abilities/spell.go:39-44`).

Combat and event tracking are `.disabled` against packages that **no longer exist**
(`internal/game/rules`, `internal/game/effects`). `internal/game/watchers/` contains only
`.disabled` files — an empty Go package.

**Conclusion: there is no legality layer, and no partially-wired one either.** A bot plays the
same rules-light sandbox a human playtester does. This is a Phase 7 problem, deliberately.

### 0.6 Allowed APIs — local

Verified signatures. Do not invent others.

```go
// internal/game/manager.go:292 — needs nothing but a username string.
func (m *Manager) SendPlayerAction(gameID, playerID, actionType string, data interface{}) error

// internal/game/manager.go:482 — public, takes the engine read lock. Pull-based, zero source changes.
func (ea *EngineAdapter) GetGameView(gameID, playerID string) (interface{}, error)

// internal/game/game_engine.go:270
func (e *GameEngine) GetGameView(gameID, playerID string) (interface{}, error)

// internal/game/game_engine.go:414
func (e *GameEngine) ParseAndExecuteStringCommand(gameID, playerID, command string) error

// internal/game/game_engine.go:74
func (e *GameEngine) StartGameWithDecks(gameID string, players []string, gameType string, decks map[string]DeckList) error

// internal/game/manager.go:41
type PlayerAction struct { PlayerID, ActionType string; Data interface{}; Timestamp time.Time }

// internal/game/manager.go:343
type DeckList struct { MainDeck, Sideboard, Commanders []string }

// internal/game/game_engine.go:26
type NotificationHandler interface {
    NotifyGameStateChange(playerID string, gameView interface{})
    NotifyGameEvent(gameID string, eventType string, data interface{})
}
```

`ProcessAction` action types: `SEND_STRING`, `DRAW`, `PLAY`, `MOVE`, `TAP`, `UNTAP_ALL`,
`FLIP`, `MODIFY_LIFE`, `SET_COUNTER`, `SHUFFLE`, `CREATE_TOKEN`, `ADD_COUNTER`,
`REMOVE_COUNTER`, `SET_CARD_COUNTER`, `MILL`, `SCRY`, `SET_REVEALED_TOP`, `NEXT_TURN`,
`MULLIGAN`, `KEEP_HAND`. Anything else errors.

`ParseAndExecuteStringCommand` verbs, colon-delimited: `TAP:<id>`, `UNTAP:<id>`,
`MOVE:<id>:<zone>`, `FLIP:<id>:<bool>`, `DRAW:<player>:<n>`, `MODIFY_LIFE:<player>:<delta>`,
`SET_COUNTER:<player>:<type>:<v>`, `SHUFFLE:<player>`,
`CREATE_TOKEN:<name>:<types>:<p>:<t>:<color>`, `ADD_COUNTER:<id>:<name>:<n>`,
`REMOVE_COUNTER:<id>:<name>:<n>`, `SET_CARD_COUNTER:<id>:<name>:<n>`, `MILL:<player>:<n>`,
`SCRY:<player>:<n>`, `REVEAL_TOP:<player>:<bool>`, `NEXT_TURN:<player>`, `MULLIGAN:<player>`,
`KEEP_HAND:<player>`.

Zones: `LIBRARY HAND BATTLEFIELD GRAVEYARD EXILE STACK COMMAND` (`game_state.go:113-121`).

### 0.7 Allowed APIs — Anthropic Go SDK

Not currently a dependency (`go.mod` has no `anthropic` entry). Add in Phase 5.

```go
// Manual loop — USE THIS, not the tool runner. See anti-patterns.
func (r *BetaMessageService) New(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (*BetaMessage, error)
func (r BetaMessage) ToParam() BetaMessageParam

func NewBetaToolResultBlock(toolUseID string, content string, isError bool) BetaContentBlockParamUnion
func NewBetaUserMessage(blocks ...BetaContentBlockParamUnion) BetaMessageParam
func NewBetaTextBlock(text string) BetaContentBlockParamUnion

type BetaToolParam struct {
    InputSchema BetaToolInputSchemaParam
    Name        string
    Description param.Opt[string]
    Strict      param.Opt[bool]
    CacheControl BetaCacheControlEphemeralParam
    // ...
}
type BetaToolInputSchemaParam struct {
    Properties  any
    Required    []string
    ExtraFields map[string]any   // additionalProperties goes here — no typed field
}

func NewBetaCacheControlEphemeralParam() BetaCacheControlEphemeralParam
// 1h TTL: BetaCacheControlEphemeralParam{TTL: BetaCacheControlEphemeralTTLTTL1h}
// verify: resp.Usage.CacheCreationInputTokens / resp.Usage.CacheReadInputTokens

func BetaToolChoiceParamOfTool(name string) BetaToolChoiceUnionParam
// DisableParallelToolUse param.Opt[bool] exists on Auto/Any/Tool variants

func option.WithRequestTimeout(dur time.Duration) RequestOption   // PER ATTEMPT
func option.WithMaxRetries(retries int) RequestOption             // default 2
```

Model constants exist in SDK `main`: `ModelClaudeOpus5`, `ModelClaudeSonnet5`,
`ModelClaudeHaiku4_5`, `ModelClaudeHaiku4_5_20251001`. **Verify against the pinned version in
`go.mod` before use** — the bundled skill doc claims Opus 5 has no constant; source disagrees.

| Model | ID | $/MTok in | $/MTok out | Context | Max out | Min cacheable prefix |
|---|---|---|---|---|---|---|
| Opus 5 | `claude-opus-5` | 5 | 25 | 1M | 128K | 512 |
| Sonnet 5 | `claude-sonnet-5` | 2 | 10 | 1M | 128K | 1024 |
| Haiku 4.5 | `claude-haiku-4-5` | 1 | 5 | 200K | 64K | **4096** |

Thinking/effort per model:
- Opus 5 / Sonnet 5 — `Thinking: ThinkingConfigParamUnion{OfAdaptive: &ThinkingConfigAdaptiveParam{}}`,
  and `OutputConfig.Effort` from `low`|`medium`|`high`|`xhigh`|`max`.
- **Haiku 4.5 — no adaptive thinking and NO effort support at all.** Use
  `ThinkingConfigParamOfEnabled(N)` (budget < MaxTokens, min 1024) or omit thinking. Do not
  set `OutputConfig.Effort`; it is not offered for this model.

Error handling: `errors.As(err, &apierr)` into `*anthropic.Error`, switch on
`apierr.StatusCode`. Retryable: 429, 500, 529 (SDK's own predicate also covers 408, 409, all
≥500). Not retryable: 400, 401, 403, 404, 413.

### 0.8 Anti-patterns — do NOT do these

1. **Do not use `BetaToolRunner`.** It overwrites `params.Tools` with definitions built only
   from `Name()/Description()/InputSchema()`, so `CacheControl`, `Strict`, and `DeferLoading`
   cannot be set on a runner-managed tool. Cached tool definitions are exactly what a
   game-long conversation needs. Its handlers also run concurrently via `errgroup`, and turns
   ending in `MaxTokens`/`Refusal` silently skip pending tool calls. Use the manual loop.

2. **Do not call `GetGameView` or `SendPlayerAction` synchronously from a notification
   handler.** `broadcast` (`game_engine.go:342`) runs while `e.mu` is held for **writing** —
   every `actions.go` method does `Lock(); defer Unlock()` then `broadcast()`. A synchronous
   call deadlocks. The codebase already documents this hazard at `grpc.go:188-190` and
   `grpc.go:301-303`. Hand off to a goroutine or a queue.

3. **Do not treat the view as a private copy.** `buildGameView` shares `Battlefield`, `Exile`,
   `Stack`, `Command`, graveyards, and mana pools **by pointer** (`view.go:74-81, 102, 119-120`).
   A bot holding a view can mutate live engine state through those slices — a hazard a
   websocket client does not have, because protojson copies. `Redact()` must deep-copy.

4. **Do not put an LLM API key in `config/config.yaml`.** It is a **hard link** to
   `config.dev.yaml` (same inode) and is the file that loads by default
   (`cmd/server/main.go:35`). Also, `Load()` calls `v.AutomaticEnv()` with no prefix and **no
   `SetEnvKeyReplacer`** (`config.go:166`), so nested keys like `bot.api_key` are not reachable
   as `BOT_API_KEY` without adding `v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))`.

5. **Do not assume the notification handler is available.** `GameEngine.notifyFn` is a single
   field; `SetNotificationHandler` overwrites it, and `EngineAdapter.SetNotificationCallback`
   already holds the only registration (`grpc.go:122`). `SetNotificationCallback` also
   type-asserts `ea.engine.(*GameEngine)` **unchecked** — it panics on any other engine.

6. **Do not use `plugin.*` expecting it to do anything.** The entire package is referenced
   only by a blank import at `cmd/server/main.go:19`. `grep -rn "plugin\."` outside the package
   returns zero hits. Starting life is hardcoded to **20**, not 40 (`game_state.go:149`);
   commanders are shuffled into the **library**, not the command zone (`game_engine.go:94`);
   commander damage is never tracked and `Card` has no `Damage` field.

7. **Do not add a field to a golden test's volatile-strip list to make it pass.** mage-bench's
   rule, worth adopting verbatim (`doc/golden-prompts.md:197-204`): the strip list is a
   whitelist of fields verified non-semantic. Adding to it to silence a flake is forbidden.

8. **Do not send action types the engine does not handle.** The gRPC layer already sends
   `SEND_UUID`, `SEND_BOOLEAN`, `SEND_INTEGER`, `SEND_MANA_TYPE`, `PLAYER_ACTION`,
   `SPECIAL_ACTION`, `ACTIVATE_ABILITY` — none are in `ProcessAction`'s switch, so all land in
   `default` and are silently swallowed by `ProcessGameActions`, which logs errors and never
   returns them. Bot actions must be verified by reading back state, not by trusting success.

---

## Phase 1: Vendor mage-bench Artifacts

**Goal:** get the reference material into the repo, unmodified, with attribution. No Go code yet.

### Tasks

1. `mkdir -p mage-server-go/internal/bot/reference/`
2. Fetch and commit verbatim, each with a header comment naming the source path and commit SHA:
   - `mcp-tools.json5` ← `website/src/data/mcp-tools.json5`
   - `system-prompt.md` ← `puppeteer/prompts/default.md`
   - `decision_renderer.py` ← `src/magebench/game/decision_renderer.py`
   - `pilot_rendering.py` ← `src/magebench/pilot/pilot_rendering.py`
   - `pilot_recovery.py`, `pilot_state.py`
   - `ShortIdRegistry.java`
   - `harness_epoch.py`
3. Create `mage-server-go/internal/bot/reference/LICENSE-mage-bench.txt` containing both MIT
   notices from upstream `LICENSE.txt` verbatim.
4. Create `mage-server-go/internal/bot/reference/README.md` recording: upstream repo, commit
   SHA pinned, fetch command, and the rule that these files are **reference only — never
   imported, never executed**. Note they are Python/Java in a Go module by design.
5. Record the pinned upstream SHA: `gh api repos/GregorStocks/mage-bench/commits/master --jq '.sha'`

### Verification

- [ ] `go build ./...` still clean (the vendored files must not be in any Go package)
- [ ] `grep -rn "reference/" --include=*.go mage-server-go/` returns nothing
- [ ] Both MIT copyright lines present in `LICENSE-mage-bench.txt`
- [ ] `README.md` records a concrete 40-char SHA, not `master`

### Anti-pattern guards

- Do not "clean up", reformat, or translate the vendored files. Their value is being diffable
  against upstream.
- Do not add them to `.gitignore`.

---

## Phase 2: `Redact()` + Serializer + Golden Tests

**Goal:** turn `*PlaytestGameView` into the mage-bench board text, provably leaking nothing.
**No LLM in this phase.**

### Tasks

1. **`internal/bot/redact.go`** — `func Redact(v *game.PlaytestGameView, viewerID string) *SafeView`
   - Deep-copy every slice and pointer. See anti-pattern 3 — the input aliases live engine state.
   - Drop `PlaytestPlayerView.Library` entirely. It is the viewer's own library in full
     (`view.go:100`) and the client needs it for search/scry UI, but a bot seeing it draws
     perfectly. Reintroduce a revealed slice only for an explicit scry/search decision.
   - Assert, do not trust, that every `PlaytestOpponentView.Hand` and `.Library` is empty.
   - `SafeView` is a new type in `internal/bot`. Do not reuse the engine's view type — the type
     boundary is the security boundary.

2. **`internal/bot/shortid.go`** — port `ShortIdRegistry.java`.
   - `map[string]string` UUID→short, monotonic counter from 1, prefix `p`.
   - Stability across zones is inherent (keyed by card ID, and `game.Card.ID` is stable).
   - Deterministic assignment order — sort by `(name, sequence)`, **never** by card ID.
   - Old IDs stay resolvable as aliases after reassignment.

3. **`internal/bot/serialize.go`** — port `_render_decision_block`. Implement the complete
   field list in §0.2, in that exact order. Include Card Reference with `seenOracleCards`
   dedup and the basic-land exclusion set.
   - Oracle text source: the existing `scryfall_cards` repo (`internal/repository/scryfall_cards.go`).
     `game.Card` carries only `Name`/`DisplayName` after `StartGameWithDecks` (`game_engine.go:94-114`),
     so type line, mana cost and P/T must be joined from Scryfall, not read off the card.

4. **`internal/bot/golden_test.go`** + `internal/bot/testdata/golden/*.json5`
   - First golden files in this repo — no existing pattern to follow
     (`find . -name testdata -o -name "*.golden"` returns nothing).
   - Port `dumps_json5` semantics from `json5_utils.py`: 2-space indent, sorted keys, trailing
     commas, and `\n` inside strings expanded to a line continuation so long prompts diff
     line-by-line. This is what makes review possible.
   - `UPDATE_GOLDEN=1 go test ./internal/bot/...` regenerates.
   - Scenarios to cover, mirroring upstream: opening mulligan, land-drop main phase, stack
     non-empty, declare-attackers, declare-blockers with `Incoming Attackers`, multi-amount
     items, oracle-dedup on second appearance.

5. **`internal/bot/leak_test.go`** — the security regression test.
   - Build a `SafeView` for a bot seat in a fixture game, walk the **serialized string**, and
     fail if it contains any card name that is not in that seat's legitimate knowledge set
     (own hand, public zones, revealed cards).
   - Also assert `SafeView` shares no pointer with the source view — mutate the source after
     redaction and confirm the `SafeView` is unchanged.

### Documentation references

- Field list and format strings: §0.2 above, and `reference/decision_renderer.py:125-248`
  (`_render_decision_block`), `:251-291` (`_render_board`), `:301-342` (`permanent_display`),
  `:451-473` (`format_choice`), `:705-767` (`_render_card_reference`).
- Trailing lines (`Respond:` / `Mana pool:` / `Recent chat:`): `reference/pilot_rendering.py:66-90`.
- Oracle dedup: `reference/pilot_rendering.py:53-56`.
- Short IDs: `reference/ShortIdRegistry.java:13-34` (class doc states the invariants).
- Test style: stdlib `t.Errorf`/`t.Fatalf` (`internal/server/grpc_table_test.go:5-13`) or
  testify (`internal/game/abilities/costs_test.go:303`) — both exist in-repo; testify is a
  direct dep at `go.mod:11`. Logger: `zaptest.NewLogger(t)`.

### Verification

- [ ] `go test ./internal/bot/...` passes; `make test` (`go test -race ./...`) still passes
- [ ] Golden files exist and are human-readable — open one and confirm the board renders
      line-by-line, not as one escaped blob
- [ ] `leak_test.go` fails when `Redact` is deliberately broken (delete the library-strip line
      and confirm red) — an untested guard is not a guard
- [ ] Pointer-aliasing assertion passes
- [ ] `grep -n "Library" internal/bot/redact.go` shows the field is explicitly dropped, with a
      comment explaining why

### Anti-pattern guards

- Do not skip the deep copy "because we only read it".
- Do not make `SafeView` an alias of `*game.PlaytestGameView`.
- Do not put the basic-land exclusion list inline in three places — one package-level set.

---

## Phase 3: `LegalMoves` + `BotRunner` (Random Policy)

**Goal:** a headless 4-bot Commander game that completes, driven by random choice. Still no LLM.
This is the real milestone — if random bots cannot finish a game, an LLM will not fix that.

### Tasks

1. **Fix `NextTurn` determinism first.** `actions.go:655-670` builds player order by ranging
   over `state.Players`, a `map[string]*Player` — Go randomizes map iteration, so turn order
   differs every call. This makes bot games unreproducible and must be fixed before any test
   depends on sequencing. Sort the player IDs, or store an explicit turn-order slice on
   `GameState`. Add a test asserting turn order is stable across 100 runs.

2. **`internal/bot/moves.go`** — `func LegalMoves(v *SafeView) []Macro`
   - Honor-system: the manual action set a human playtester has. Not real legality — see §0.5.
   - `type Macro struct { Label string; Steps []string }` — `Label` is what the model sees,
     `Steps` are `ParseAndExecuteStringCommand` strings.
   - Cover: play land, "cast" (move hand→stack→battlefield with lands tapped), tap/untap,
     move between zones, modify life, draw, mulligan/keep, pass turn.
   - Mana payment is a **solver**, not a prompt: subset-sum over untapped sources. Never ask
     the model which lands to tap.

3. **`internal/bot/runner.go`** — `BotRunner`
   - One goroutine per bot seat.
   - Integration via **polling `EngineAdapter.GetGameView`** — anti-pattern 2 makes the
     push path a deadlock hazard, and polling requires **zero source changes** to the engine.
     Revisit in Phase 7 if latency demands it.
   - `Policy` interface: `Pick(ctx, *SafeView, []Macro) (Macro, error)`. Phase 3 ships
     `RandomPolicy`. Phase 5 adds `LLMPolicy` behind the same interface.
   - Execute a macro by sending each step via `SendPlayerAction(gameID, botID, "SEND_STRING", step)`.
   - **Verify by reading back state, not by trusting success** — anti-pattern 8.

4. **`internal/bot/harness_test.go`** — headless 4-player Commander game.
   ```go
   logger := zaptest.NewLogger(t)
   engine := game.NewGameEngine(logger)
   adapter := game.NewEngineAdapter(engine, logger)
   mgr := game.NewManager(logger)
   g := mgr.CreateGame("table-1", "Commander Free For All", []string{"p1","p2","p3","p4"})
   decks := map[string]game.DeckList{ /* ... */ }
   _ = adapter.StartGameWithDecks(g, decks)
   go adapter.ProcessGameActions(g)
   ```
   No DB, no session manager, no table manager, no gRPC, no websocket. This bypasses
   `MatchStart`'s two blockers (`len(players) < 2` and the websocket-session requirement at
   `grpc_game.go:93-110`).

   Known inherited defects to assert around, **not** fix here: 20 life not 40, commanders in
   library not command zone, libraries unshuffled, cards name-only.

5. Add a `make bot-sim` target running N headless games and reporting completion rate.

### Verification

- [ ] 4-bot game reaches a terminal state in ≥95% of 100 headless runs
- [ ] Turn order deterministic across 100 runs
- [ ] No deadlock under `-race` — this is the phase where anti-pattern 2 would surface
- [ ] `make test` clean
- [ ] Every `LegalMoves` macro's steps execute without an engine error (assert by reading state back)

### Anti-pattern guards

- Do not wire into `SetNotificationCallback` in this phase. Poll.
- Do not try to make the bot play *well*. Random is the point — it isolates loop bugs from
  policy bugs.
- Do not fix the 20-life / commander-zone defects here; they belong to a Commander-rules task,
  and conflating them hides which layer broke.

---

## Phase 4: Pacer (Human-Like Timing)

**Goal:** bot actions look human in the real client. Judged by eye, before any LLM exists.

### Tasks

1. **`internal/bot/pace.go`**
   - Stagger macro steps: emit each with 200–600ms jittered delay, so lands tap one at a time
     and the card visibly slides to the stack. Each step broadcasts
     (`game_engine.go:342`), so the client animates naturally.
   - Weight the pre-action delay by decision type: land drop / untap 0.5–1.5s; routine cast
     1.5–4s; combat or targeted removal 4–10s; mulligan or complex stack 8–20s.
   - **Log-normal, not uniform** — humans have a long tail. ~5% of decisions "tank" past 25s.
   - Optional hesitation: occasionally `TAP` then `UNTAP` before doing something else.

2. **Chat.** Wire `internal/chat`. Adopt mage-bench's cadence instruction (≥1 message per 2
   turn cycles) and their `MAX_CHAT_MESSAGES_PER_TURN = 2` cap. In Phase 4 use canned lines;
   Phase 5 replaces them with the model's `why`.

3. **Presence.** Whatever activity/heartbeat signal the client shows for humans, bots emit too.

### Verification

- [ ] Run `make bot-sim` with the real web client attached, watch a full game, and judge
      whether it reads as human. This is a subjective gate and it is the point of the phase.
- [ ] No step-stagger delay blocks the engine lock (delays happen in the bot goroutine, between
      `SendPlayerAction` calls, never inside a handler)
- [ ] Pacing is configurable and can be set to zero for fast headless sims — Phase 3's
      completion-rate test must not take hours

### Anti-pattern guards

- Do not `time.Sleep` while holding any engine reference in a way that could stall broadcast.
- Do not make delays uniform-random; it reads as robotic in a different way.

---

## Phase 5: LLM Policy (Manual Loop)

**Goal:** replace `RandomPolicy` with `LLMPolicy`. Same interface, same runner.

### Tasks

1. `go get github.com/anthropics/anthropic-sdk-go`. Record the resolved version and **verify
   `ModelClaudeSonnet5` / `ModelClaudeHaiku4_5` exist in that version** (§0.7).

2. **`internal/bot/llm/client.go`** — manual loop, per anti-pattern 1.
   - Tools: port `choose_action`, `pass_priority`, `get_oracle_text`, `send_chat_message` from
     `reference/mcp-tools.json5`. Skip `get_game_state` / `get_game_log` initially — our
     `SafeView` is small enough to inline.
   - `Strict: true` + `ExtraFields: {"additionalProperties": false}` on every tool.
   - Keep tool order **deterministic** (sort by name) — tools render at position 0 and any
     reorder invalidates the whole cache.
   - `CacheControl` on the last tool definition and on the system prompt.

3. **`internal/bot/llm/context.go`** — context management. This is its own component, not a
   detail; two independent constraints bite here.
   - Window as upstream does: last 40 messages verbatim, previous 20 with tool results >200
     chars summarized, synthetic state-bridge message between, board summary refreshed every 5
     renders (`reference/pilot_rendering.py:275-343`).
   - **The 20-block cache lookback**: each `cache_control` breakpoint walks back at most 20
     content blocks to find a prior entry. An agentic turn with many tool_use/tool_result pairs
     exceeds this and silently rebuilds the cache. Place an intermediate breakpoint every ~15
     blocks. Max 4 breakpoints per request.
   - Never mutate the system prompt or the tool list mid-game.

4. **`internal/bot/llm/policy.go`** — `LLMPolicy` implementing `Policy`.
   - Model: start with `claude-sonnet-5`. **Do not assume Haiku is cheaper** — its minimum
     cacheable prefix is 4096 tokens vs Sonnet 5's 1024, and it supports neither adaptive
     thinking nor `effort`, which removes the main cost dial. Measure both in Phase 6 and let
     data decide. If Haiku is used, `ThinkingConfigParamOfEnabled(N)` and **no**
     `OutputConfig.Effort`.
   - Two separate timers, per §0.4:
     - **LLM request timeout: 120s** (upstream's measured value; 45s was empirically too short).
     - **Table-stall guard**: total wall-clock the bot may hold priority, 45–60s, force-pass past it.
     - The pacing delay from Phase 4 **overlaps** the LLM latency — fire the request
       immediately, pad up to the persona's desired pause only if the model returns early.
       Latency is only visible when it exceeds the pause we wanted anyway.
   - Bound the whole loop with `context.WithTimeout`; `WithRequestTimeout` is **per attempt**,
     so worst case is timeout × (retries+1).

5. **`internal/bot/llm/recovery.go`** — port the failure matrix from
   `reference/pilot_recovery.py` and the constants in §0.4. Include the in-character
   degradation: chat, then auto-pass, rather than hanging or crashing the table.

6. **Config.** Add `Bot BotConfig` to `Config` (`config.go:21`), the struct after
   `MetricsConfig` (`:149`), defaults after the metrics block (`:245`), enum validation before
   `return nil` (`:282`, precedent at `:267`). **API key from `ANTHROPIC_API_KEY` env only** —
   never the YAML (anti-pattern 4). If a config path for the key is wanted, add
   `v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))` at `config.go:166` first.

### Verification

- [ ] A 4-bot game completes with `LLMPolicy` under `make bot-sim`
- [ ] `resp.Usage.CacheReadInputTokens > 0` on the second and every later request. **If it is
      zero, stop and find the invalidator** — a timestamp, an unsorted map, a mutated tool list.
      Costs will be ~10× otherwise and nothing else will report it.
- [ ] Malformed tool output retries once, then falls back — assert with a stubbed client
- [ ] Every `Policy` implementation is interchangeable: the same harness test passes with
      `RandomPolicy` and `LLMPolicy`
- [ ] `grep -rn "ANTHROPIC_API_KEY\|api_key" config/` returns nothing

### Anti-pattern guards

- No `BetaToolRunner` (anti-pattern 1).
- No `OutputConfig.Effort` on Haiku 4.5.
- No `budget_tokens` on Sonnet 5 / Opus 5 — it returns 400.
- Do not let the model emit raw engine commands. It picks a macro; Go expands it.

---

## Phase 6: Metrics + Golden Prompt Tests

**Goal:** make cost, quality, and prompt drift observable. This is what catches a silent
regression, and it is where the Haiku-vs-Sonnet question gets answered with data.

### Tasks

1. **`internal/bot/metrics.go`** — adopt upstream's field set (`leaderboard/stats.py:84-99`):
   `gamesPlayed`, `wins`, `totalCostUsd`, `totalToolCallsOk`, `totalToolCallsFailed`,
   `totalPromptTokens`, `totalCompletionTokens`, `totalCachedTokens`, `successfulResponses`,
   `errors{}`, `contextResets`, `latencyP50`, `latencyP95`.
   Derived: cost/game, cache %, tool-fail %, timeout %.

2. **Golden prompt tests.** Extend Phase 2's harness to capture the complete wire-format
   `messages` array for scripted scenarios and diff it in CI. This catches the failure mode
   that has no error path: a serializer refactor drops a field from the board line, bots
   quietly get worse, nothing logs.
   - Port `_strip_volatile` from `golden_helpers.py:1793-1852` — and adopt their rule that the
     strip list is a verified whitelist, never a flake silencer (anti-pattern 7).
   - Add a `make regen-golden` target mirroring upstream's.

3. **Harness epoch.** Copy the idea from `harness_epoch.py`: a single integer plus a changelog,
   bumped whenever tools, prompt, or loop semantics change enough to make results
   non-comparable. Without it, cost and win-rate numbers from different weeks get averaged
   together and mean nothing.

4. **Run the model comparison.** N games each on `claude-haiku-4-5` and `claude-sonnet-5`.
   Report cost/game and completion rate. Decide the MVP default from the result, not from
   the sticker price.

### Verification

- [ ] `make regen-golden` produces reviewable line-by-line diffs
- [ ] CI fails when the board serializer changes without a golden update
- [ ] Metrics land for a full game; cost/game is within 2× of a hand-computed estimate
- [ ] The Haiku-vs-Sonnet comparison is recorded in this document with real numbers

---

## Phase 7: Legality Layer + Combat (v2 Unlock)

**Goal:** real rules enforcement, so bots can face humans online without being a bug report.
Weeks, not days. Needed for online NPCs regardless of whether bots ever ship — see §0.5.

### Tasks

1. **Pick one hierarchy per concept and delete the other.** Currently two mana-cost parsers,
   two targeting systems, three card types (§0.5). This decision blocks everything else.
2. **Rebuild `internal/game/rules`** (EventBus, Event, Watcher/BaseWatcher) — the package the
   `.disabled` watchers import. Nothing else revives without it.
3. **Combat.** `combat_test_harness.go.disabled` (349 lines) is a sound design whose engine was
   removed. Reviving it needs ten `GameEngine` methods (`ResetCombat`, `SetAttacker`,
   `SetDefenders`, `DeclareAttacker`, `DeclareBlocker`, `AcceptBlockers`, `AssignCombatDamage`,
   `ApplyCombatDamage`, `EndCombat`, `HasFirstOrDoubleStrike`) plus a `Damage` field on `Card`
   and an int-typed zone. Add matching engine action types and string commands.
4. **`CanPlay(card, state, player)`** wiring `mana` + `abilities` into `ProcessAction`, and a
   real `LegalMoves` behind Phase 3's interface. `abilities.GameContext`
   (`abilities/ability.go:75`) is the intended bridge; the only implementation that ever
   existed is `game_context.go.disabled` (1676 lines).
5. **Commander rules.** Wire `internal/plugin` for real: 40 life, commanders to the command
   zone, commander tax, 21-damage. Both the behavior (`plugin/commander_behavior.go`) and the
   watcher (`watchers/commander_damage.go.disabled`, which cites rule 903.10a) exist and
   neither is connected.
6. **Fix the client/server command mismatch.** The web client sends ~11 commands the server
   does not implement (`UNTAP_ALL`, `TRANSFORM:`, `SET_LIFE:`, `CLEAR_COMBAT`, `STACK_ADD:`, …)
   and sends `NEXT_TURN`/`SHUFFLE` without the required `playerId`. Reconcile both directions.
7. **FIX THE `GetGameView` DATA RACE (pre-existing, affects production).** Found in Phase 3 by
   running the bot sim under `-race`. `GameEngine.GetGameView` (`game_engine.go:270`) takes
   `e.mu.RLock()` with `defer RUnlock()` — so the lock is released when it **returns**, but the
   view it returns aliases live engine state by pointer (`view.go:74-81`). Every caller then
   reads that live state with no lock held:
   - `internal/server/grpc_game.go:387` (`GameGetView`)
   - `internal/server/grpc.go:314` — hands the aliased view to protojson

   Observed race: `GameEngine.Mulligan` writing under the write lock vs. a concurrent read of
   the same memory through the returned view. **This is not a bot bug** — the websocket path
   has always had it; bots are just the first workload that hits it hard enough to catch.

   Real fix: deep-copy inside `GetGameView` while the read lock is held, so the returned view
   owns its memory. Phase 3 works around it locally with `internal/bot/guard.go` (`ViewGuard`),
   which holds a read lock across `GetGameView` + `Redact`'s copy and wraps `ProcessAction` in a
   write lock. **Delete `guard.go` once `GetGameView` is fixed** — the workaround only protects
   the bot path, not the server's.

8. **Notification fan-out.** Convert `NotificationHandler` from a single field to a slice so
   bots can subscribe push-style, and replace polling. `broadcast` already builds a per-player
   filtered view, so a bot handler receives exactly what a human does. Keep the goroutine
   hand-off (anti-pattern 2).

### Verification

- [ ] Bots cannot make an illegal play — property test over random game states
- [ ] Combat harness revived and passing
- [ ] A Commander game starts at 40 life with commanders in the command zone
- [ ] Human-vs-bot game over the real network stack completes

---

## Phase 8: Final Verification

- [ ] `make test` (`go test -race -coverprofile ./...`) clean
- [ ] `make lint`, `make vet`, `make fmt` clean
- [ ] `grep -rn "BetaToolRunner" --include=*.go .` → nothing (anti-pattern 1)
- [ ] `grep -rniE "anthropic|llm_api_key|bot\.api_key" mage-server-go/config/` → nothing (anti-pattern 4).
      Note: `config.example.yaml:63` has a pre-existing **Mailgun** `api_key: "key-changeme"`
      placeholder. That is unrelated and expected — do not "fix" it, and do not use a bare
      `api_key` grep, which matches it and produces a permanent false positive.
- [ ] Leak test still fails when `Redact` is deliberately broken
- [ ] `CacheReadInputTokens > 0` in a live game
- [ ] Both MIT notices present and the pinned upstream SHA recorded
- [ ] Every macro the bot can emit is a command the engine actually implements — cross-check
      `LegalMoves` output against the §0.6 verb list
- [ ] The Haiku-vs-Sonnet decision is recorded with real numbers, not assumed

---

## Open Questions

1. **Poll vs push.** Phase 3 polls because push is a deadlock hazard and needs an engine
   change. If poll latency proves visible in Phase 4, bring the fan-out from Phase 7 forward.
2. **Where bots live.** In-process with the server (simplest, chosen here) vs a separate
   process speaking gRPC like a real client (more faithful, and would let us reuse
   mage-bench's harness directly). In-process is right for the MVP; revisit for v2.
3. **Whether to fix Commander rules before Phase 5.** Bots playing at 20 life with commanders
   in the library are testing the loop, not the format. Fine for MVP; wrong for evaluating
   play quality.
