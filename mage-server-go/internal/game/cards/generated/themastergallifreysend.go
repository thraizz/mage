package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Master Gallifreys End", NewTheMasterGallifreysEnd)
}

// NewTheMasterGallifreysEnd creates a The Master Gallifreys End
// {2}{B}{R} - CREATURE
func NewTheMasterGallifreysEnd(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Master Gallifreys End")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TIME_LORD", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
