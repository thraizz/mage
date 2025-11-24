package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Survivors Encampment", NewSurvivorsEncampment)
}

// NewSurvivorsEncampment creates a Survivors Encampment
//   - LAND
func NewSurvivorsEncampment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Survivors Encampment")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"DESERT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}
