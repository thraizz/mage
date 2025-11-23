package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ice Tunnel", NewIceTunnel)
}

// NewIceTunnel creates a Ice Tunnel
//   - LAND
func NewIceTunnel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ice Tunnel")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"ISLAND", "SWAMP"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	return card, nil
}
