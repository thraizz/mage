package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tinder Farm", NewTinderFarm)
}

// NewTinderFarm creates a Tinder Farm
//   - LAND
func NewTinderFarm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tinder Farm")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	return card, nil
}
