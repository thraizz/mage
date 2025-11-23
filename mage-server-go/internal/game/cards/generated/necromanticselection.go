package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Necromantic Selection", NewNecromanticSelection)
}

// NewNecromanticSelection creates a Necromantic Selection
// {4}{B}{B}{B} - SORCERY
func NewNecromanticSelection(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Necromantic Selection")
	card.ManaCost = "{4}{B}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
