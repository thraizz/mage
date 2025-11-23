package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Snow Covered Mountain", NewSnowCoveredMountain)
}

// NewSnowCoveredMountain creates a Snow Covered Mountain
//   - LAND
func NewSnowCoveredMountain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Snow Covered Mountain")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"MOUNTAIN"}
	card.Supertypes = []string{"BASIC", "SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	return card, nil
}
