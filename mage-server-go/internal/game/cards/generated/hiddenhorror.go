package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hidden Horror", NewHiddenHorror)
}

// NewHiddenHorror creates a Hidden Horror
// {1}{B}{B} - CREATURE
func NewHiddenHorror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidden Horror")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HORROR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
