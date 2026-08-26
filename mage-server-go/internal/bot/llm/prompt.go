package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/magefree/mage-server-go/internal/bot"
)

// prompt.go holds the frozen system prompt and the per-decision user message.
//
// SystemPrompt is adapted from reference/system-prompt.md (mage-bench
// puppeteer/prompts/default.md, dual MIT -- see reference/LICENSE-mage-bench.txt).
// It is ADAPTED, not copied, because their tool surface is not ours: there is
// no blocking pass_priority that returns the next decision, no server-side
// legality filter, and no mana_plan, because none of those subsystems exist in
// this engine (§0.5). Promising the model a capability the harness lacks is how
// a bot ends up issuing moves that silently do nothing (anti-pattern 8).
//
// IT IS A CONSTANT AND MUST STAY ONE. It is the second-largest cacheable block
// in every request; a per-game name, a persona line or a timestamp spliced in
// here costs a full cache rebuild every turn, for every seat, for the whole game.
// Per-seat voice belongs in Pacing/ChatSource, which is not part of the prefix.
func SystemPrompt() string { return systemPrompt }

const systemPrompt = `You are a competitive Magic: The Gathering player at a four-player Commander table. Your goal is to WIN. Play to maximise your win rate -- make optimal strategic decisions, not flashy ones. Think about sequencing, card evaluation and board development.

## Game Flow

The core loop is: read the decision, take one action, repeat. Each of your messages contains the current board state and a numbered list of the options available to you RIGHT NOW.

- ` + "`choose_action`" + ` -- take one of the listed options.
- ` + "`pass_priority`" + ` -- decline to act. If a "Pass the turn" option is listed, that is what happens.

Take exactly one action per decision. The options are pre-filtered: everything listed is something you can legally do this instant, and anything not listed is not available, however much you would like it.

## Reading a Decision

- The board shows life totals, hands, battlefields, graveyards and exile.
- Your hand is shown in full. Opponents' hands show only a count. You cannot see any library, including your own.
- A Card Reference section gives oracle text for non-basic cards the first time they appear. It is not repeated -- call ` + "`get_oracle_text`" + ` if you need a reminder.
- Every game object has a short id like "p3". Ids are stable as a card changes zones.
- Every OPTION has an id like "m3". Take it with ` + "`choose_action(choice=\"m3\")`" + `.

## Mulligans

Before the game starts you will be asked to keep or mulligan.

- ` + "`choose_action(choice=\"yes\")`" + ` = MULLIGAN -- shuffle back and draw a new hand.
- ` + "`choose_action(choice=\"no\")`" + ` = KEEP this hand.

The listed options say the same thing in words; either form works.

## Chat

Use ` + "`send_chat_message`" + ` to talk to the table. Chat at least once every two turn cycles (a cycle is one turn for each player). React to big plays, comment on the board, and always reply when someone talks to you. Keep it to one short line. Do not play in silence -- but do not let chat replace your action either: a turn where you only chatted is a wasted turn.

## Constraints

- This table runs a rules-light engine. There is no stack, no priority window and no combat damage step: an option listed is executed the moment you pick it.
- Never invent an option. If nothing listed is good, pass.
`

// decisionMessage renders the user turn: the board, the Card Reference, and the
// enumerated options.
//
// The board half is bot.Serializer -- the Phase 2 port of mage-bench's decision
// renderer -- so the prompt the model sees is the format the golden tests
// already cover. The options half is this package's own: upstream's choices
// come from an XMage decision, ours from bot.LegalMoves.
func decisionMessage(ctx context.Context, s *bot.Serializer, v *bot.SafeView, moves []bot.Macro, index int, recentChat []string) string {
	d := &bot.Decision{
		Index:         index,
		SnapshotIndex: index,
		Turn:          v.Turn,
		Phase:         phaseOf(v),
		Player:        seatName(v),
		Message:       decisionPrompt(v),
		Choices:       choicesFor(moves),
		RespondWith:   "choice=mN to take that option, or call pass_priority",
		RecentChat:    recentChat,
	}
	return s.Render(ctx, v, d)
}

// choicesFor maps macros onto the renderer's Choice rows.
//
// The id is "mN", 1-based, and it is the ONLY handle the model gets on a macro.
// It is not the card's short id: a macro is an action, several macros can name
// the same card ("Tap Forest" / "Untap Forest"), and reusing pN would make the
// two indistinguishable. Action carries the macro kind so the model can tell a
// land drop from a cast without parsing the label.
func choicesFor(moves []bot.Macro) []bot.Choice {
	out := make([]bot.Choice, 0, len(moves))
	for i, m := range moves {
		out = append(out, bot.Choice{
			Name:   m.Label,
			ID:     macroID(i),
			Action: string(m.KindOf()),
		})
	}
	return out
}

func macroID(i int) string { return fmt.Sprintf("m%d", i+1) }

// decisionPrompt is the one-line instruction above the choices.
func decisionPrompt(v *bot.SafeView) string {
	if v.Me != nil && !v.Me.KeptHand {
		return "Mulligan or keep this hand?"
	}
	if v.ActivePlayerID == v.ViewerID {
		return "Your turn. Choose one action."
	}
	return "Choose one action."
}

// phaseOf reports what this engine can actually say about the phase. There is
// no phase structure (§0.5: the turn is a counter), so it is PREGAME until
// every seat has kept, and MAIN afterwards -- honest, and stable enough to
// cache.
func phaseOf(v *bot.SafeView) string {
	if v.Me != nil && !v.Me.KeptHand {
		return ""
	}
	return "MAIN"
}

func seatName(v *bot.SafeView) string {
	if v.Me == nil {
		return v.ViewerID
	}
	if v.Me.Name != "" {
		return v.Me.Name
	}
	return v.Me.PlayerID
}

// boardSummary is the digest carried by the context-bridge message: the state a
// model needs to keep playing after the middle of its transcript was
// summarised away. reference/pilot_rendering.py:341-369 builds the same thing
// from get_game_state; we have the view in hand, so no round trip is needed.
func boardSummary(v *bot.SafeView) string {
	if v == nil {
		return ""
	}
	var parts []string
	parts = append(parts, fmt.Sprintf("Turn %d", v.Turn))
	if v.Me != nil {
		parts = append(parts, fmt.Sprintf("%s: %dhp, %d cards in hand, %d in library",
			seatName(v), v.Me.Life, len(v.Me.Hand), v.Me.LibraryCount))
	}
	for _, o := range v.Opponents {
		if o == nil {
			continue
		}
		name := o.Name
		if name == "" {
			name = o.PlayerID
		}
		parts = append(parts, fmt.Sprintf("%s: %dhp, hand=%d", name, o.Life, o.HandCount))
	}
	return "Board: " + strings.Join(parts, " | ")
}
