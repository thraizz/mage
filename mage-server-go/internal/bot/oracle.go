package bot

import (
	"context"
	"database/sql"

	"github.com/magefree/mage-server-go/internal/repository"
)

// OracleCard is the subset of a card's printed characteristics the Card
// Reference section renders. Field names follow the keys mage-bench's
// oracle_texts dict uses (reference/decision_renderer.py:745-752).
type OracleCard struct {
	Name       string
	ManaCost   string
	TypeLine   string
	OracleText string
	Power      string
	Toughness  string
	Loyalty    string
}

// PowerToughness renders "P/T", or "" when the card is not a creature.
// reference/decision_renderer.py:750-752.
func (o *OracleCard) PowerToughness() string {
	if o == nil || o.Power == "" {
		return ""
	}
	return o.Power + "/" + o.Toughness
}

// OracleLookup resolves a card name to its printed characteristics.
//
// This is a narrow interface on purpose. game.Card carries only Name and
// DisplayName after StartGameWithDecks (internal/game/game_engine.go:94-114) --
// ManaCost, Type, Power and Toughness are all left at their zero values -- so
// mana cost, type line, P/T and oracle text must be joined in from Scryfall.
// Depending on the concrete repository here would put a live Postgres in the
// path of every serializer unit test, so the serializer depends on this
// interface and the tests supply MapOracle.
type OracleLookup interface {
	// Oracle returns the characteristics of the named card. The bool is false
	// when the name is unknown; that is not an error, it just means the card
	// gets no Card Reference entry (reference/decision_renderer.py:753-755).
	Oracle(ctx context.Context, name string) (*OracleCard, bool)
}

// MapOracle is an in-memory OracleLookup, for tests and for fixtures.
type MapOracle map[string]*OracleCard

// Oracle implements OracleLookup.
func (m MapOracle) Oracle(_ context.Context, name string) (*OracleCard, bool) {
	c, ok := m[name]
	if !ok || c == nil {
		return nil, false
	}
	return c, true
}

// scryfallCardGetter is the one method ScryfallOracle needs from
// *repository.ScryfallCardRepository (internal/repository/scryfall_cards.go:142).
type scryfallCardGetter interface {
	GetByNameCaseInsensitive(ctx context.Context, name string) ([]*repository.ScryfallCard, error)
}

// ScryfallOracle adapts the existing Scryfall card repository to OracleLookup.
// It caches by name for the life of the value: a game re-renders the same
// couple of hundred card names hundreds of times, and printed characteristics
// do not change mid-game.
type ScryfallOracle struct {
	repo  scryfallCardGetter
	cache map[string]*OracleCard
}

// NewScryfallOracle wraps a Scryfall card repository.
func NewScryfallOracle(repo *repository.ScryfallCardRepository) *ScryfallOracle {
	return &ScryfallOracle{repo: repo, cache: make(map[string]*OracleCard)}
}

// Oracle implements OracleLookup. A name with several printings resolves to the
// first row the repository returns; printed characteristics are identical
// across printings for everything the Card Reference renders.
func (s *ScryfallOracle) Oracle(ctx context.Context, name string) (*OracleCard, bool) {
	if cached, ok := s.cache[name]; ok {
		return cached, cached != nil
	}
	rows, err := s.repo.GetByNameCaseInsensitive(ctx, name)
	if err != nil || len(rows) == 0 || rows[0] == nil {
		// A miss is cached too: an unknown name stays unknown, and a failing
		// lookup must not be retried once per decision for the rest of a game.
		s.cache[name] = nil
		return nil, false
	}
	c := rows[0]
	out := &OracleCard{
		Name:       c.Name,
		ManaCost:   nullString(c.ManaCost),
		TypeLine:   c.TypeLine,
		OracleText: nullString(c.OracleText),
		Power:      nullString(c.Power),
		Toughness:  nullString(c.Toughness),
		Loyalty:    nullString(c.Loyalty),
	}
	s.cache[name] = out
	return out, true
}

func nullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}
