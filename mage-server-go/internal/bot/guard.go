package bot

import (
	"sync"

	"github.com/magefree/mage-server-go/internal/game"
)

// guard.go closes the unsynchronized-read window that GameEngine.GetGameView
// leaves open.
//
// THE PROBLEM. GetGameView (internal/game/game_engine.go:270) takes the engine
// read lock, calls buildGameView, and returns -- releasing the lock. But
// buildGameView shares Battlefield, Exile, Stack, Command, hands, graveyards,
// mana pools and the game log with the authoritative GameState BY POINTER
// (internal/game/view.go:74-81, 102, 119-120). So the caller is handed live
// engine memory and then reads it with no lock held, while another goroutine
// mutates the very same *game.Card structs under the engine's write lock.
//
// Redact's deep copy is exactly the read that races. `go test -race` reports
// it as, for example:
//
//	Write at 0x... by goroutine 18: game.(*GameEngine).Mulligan()
//	                                game.(*EngineAdapter).ProcessGameActions()
//	Previous read at 0x... by goroutine 19: bot.copyCard()
//	                                        bot.RedactErr()
//
// This is a PRE-EXISTING ENGINE DEFECT, not something the bot runner
// introduced. The production websocket path has it too: internal/server calls
// GetGameView and then protojson-marshals the result, reading the same live
// pointers outside the lock. A bot simulation under -race is simply the first
// thing in this repo that exercises it hard enough to be caught.
//
// THE FIX, AND WHY IT LIVES HERE. The right repair is for GetGameView to deep
// copy under the read lock, but that is engine surgery that belongs with the
// people who own internal/game. Phase 3 needs correctness now and needs it
// without touching shared engine code, so ViewGuard serializes the two sides
// from outside:
//
//   - the bot runner holds ViewGuard's read lock across GetGameView AND the
//     Redact deep copy, so the copy is made from a stable snapshot;
//   - every mutation is funnelled through WrapEngine, which takes ViewGuard's
//     write lock around ProcessAction.
//
// Lock ordering is always ViewGuard before the engine's own mutex, in both
// paths, so this cannot deadlock. Once the deep copy is done, the SafeView owns
// everything it holds and the guard is irrelevant.
type ViewGuard struct {
	mu sync.RWMutex
}

// WrapEngine returns e with ProcessAction serialized against view reads.
//
// Hand the result to game.NewEngineAdapter; EngineAdapter.ProcessGameActions
// then drains the action queue through the guarded ProcessAction with no
// changes on the engine side.
//
// NOTE: the returned value is not a *game.GameEngine, so
// EngineAdapter.SetNotificationCallback -- which type-asserts to *GameEngine
// unchecked (anti-pattern 5) -- would panic on an adapter built over it. Bots
// never call it; they poll, for the reasons documented at the top of runner.go.
func (g *ViewGuard) WrapEngine(e game.EngineInterface) game.EngineInterface {
	return &guardedEngine{EngineInterface: e, guard: g}
}

type guardedEngine struct {
	game.EngineInterface
	guard *ViewGuard
}

func (g *guardedEngine) ProcessAction(gameID string, action game.PlayerAction) error {
	g.guard.mu.Lock()
	defer g.guard.mu.Unlock()
	return g.EngineInterface.ProcessAction(gameID, action)
}

// Snapshot runs fn -- a GetGameView call plus its Redact -- with mutations held
// off.
func (g *ViewGuard) Snapshot(fn func() error) error {
	if g == nil {
		return fn()
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	return fn()
}
