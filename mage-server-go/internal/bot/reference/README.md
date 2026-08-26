# mage-bench reference material

Vendored reference artifacts from the public upstream project **mage-bench** — a fork of
XMage that exposes a tool interface so LLMs can play Magic: The Gathering.

- **Upstream repo:** https://github.com/GregorStocks/mage-bench
- **Default branch:** `master`
- **Pinned commit:** `b81453e887e7935b1162ef3cbc7dd2485a3303a6`

## THE RULE: REFERENCE ONLY

**These files are REFERENCE ONLY. They are never imported, never executed, never compiled.**

They are Python, Java, JSON5 and Markdown sitting inside a Go module *by design*. Nothing in
`mage-server-go` may `import`, `exec`, embed, or otherwise depend on anything in this
directory. They exist to be *read* by humans and agents while porting behaviour into Go.

A `go build ./...` must never see them, and `grep -rn "reference/" --include=*.go` must
return nothing.

## Deliberately unmodified

Every file is byte-identical to upstream except for a leading provenance header comment.
That is intentional: the value of these artifacts is that they stay **diffable against
upstream**. Do not reformat, lint, translate, re-indent, or "clean up" any of them. When
upstream changes, re-fetch at a new SHA and diff.

`toolsets.json` is strict JSON, which does not permit comments, so it carries **no** header
and is byte-identical to upstream. Its provenance is recorded here instead.

## How these were fetched (reproducible)

```sh
SHA=b81453e887e7935b1162ef3cbc7dd2485a3303a6
gh api "repos/GregorStocks/mage-bench/contents/<upstream-path>?ref=$SHA" --jq '.content' \
  | base64 -d > <local-file>
```

To re-pin against the current head:

```sh
gh api repos/GregorStocks/mage-bench/commits/master --jq '.sha'
```

(Fallback for files over the 1 MB `contents` API limit:
`curl -sL https://raw.githubusercontent.com/GregorStocks/mage-bench/$SHA/<upstream-path>`)

## File map

| Local file | Upstream path | What it is |
|---|---|---|
| `mcp-tools.json5` | `website/src/data/mcp-tools.json5` | Generated, CI-verified wire-format schemas for all 10 tools |
| `system-prompt.md` | `puppeteer/prompts/default.md` | The complete system prompt |
| `decision_renderer.py` | `src/magebench/game/decision_renderer.py` | The `## Decision` board-state serializer |
| `pilot_rendering.py` | `src/magebench/pilot/pilot_rendering.py` | Context windowing, cache breakpoints, trailing lines |
| `pilot_recovery.py` | `src/magebench/pilot/pilot_recovery.py` | Error-recovery primitives |
| `pilot_state.py` | `src/magebench/pilot/pilot_state.py` | Loop state + context reset |
| `ShortIdRegistry.java` | `Mage/src/main/java/mage/util/ShortIdRegistry.java` | Short-ID allocation, near-zero deps |
| `harness_epoch.py` | `src/magebench/game/harness_epoch.py` | Changelog of every breaking design change — read as history |
| `toolsets.json` | `puppeteer/toolsets.json` | Named tool subsets offered to the model |
| `json5_utils.py` | `src/magebench/common/json5_utils.py` | Diff-friendly JSON5 serializer — the pattern Phase 2 golden files follow |
| `LICENSE-mage-bench.txt` | `LICENSE.txt` | Upstream license file, verbatim |

## Licensing

`LICENSE-mage-bench.txt` contains **two stacked MIT grants**:

1. MageBench code and data — © 2026 Gregor Stocks
2. Inherited XMage code (mage-bench is a fork) — © 2010 betasteward@gmail.com

Both are permissive MIT, so copying is permitted. **Both notices must travel with any copied
material** — if code or text from this directory is adapted into `mage-server-go`, the
attribution obligation comes with it. (GitHub reports the upstream repo as "NOASSERTION"
only because two license texts share one file; it is dual MIT.)

**Not covered by either grant:** Magic: The Gathering card names and oracle text are Wizards
of the Coast intellectual property. mage-bench pulls them from Scryfall at runtime, exactly
as this repo does via the `scryfall_cards` table. Do not vendor card text here.
