package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abandoned Campground", NewAbandonedCampground)
}

// NewAbandonedCampground creates a Abandoned Campground
//   - LAND
func NewAbandonedCampground(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abandoned Campground")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	return card, nil
}
