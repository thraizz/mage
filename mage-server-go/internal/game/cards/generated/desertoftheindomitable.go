package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Desert Of The Indomitable", NewDesertOfTheIndomitable)
}

// NewDesertOfTheIndomitable creates a Desert Of The Indomitable
//   - LAND
func NewDesertOfTheIndomitable(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Desert Of The Indomitable")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"DESERT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	return card, nil
}
