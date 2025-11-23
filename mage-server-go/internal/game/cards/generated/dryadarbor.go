package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dryad Arbor", NewDryadArbor)
}

// NewDryadArbor creates a Dryad Arbor
//   - LAND CREATURE
func NewDryadArbor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dryad Arbor")
	card.ManaCost = ""
	card.Types = []string{"LAND", "CREATURE"}
	card.Subtypes = []string{"FOREST", "DRYAD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	return card, nil
}
