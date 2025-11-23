package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Snowfield Sinkhole", NewSnowfieldSinkhole)
}

// NewSnowfieldSinkhole creates a Snowfield Sinkhole
//   - LAND
func NewSnowfieldSinkhole(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Snowfield Sinkhole")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"PLAINS", "SWAMP"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	return card, nil
}
