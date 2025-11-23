package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scrubland", NewScrubland)
}

// NewScrubland creates a Scrubland
//   - LAND
func NewScrubland(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scrubland")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"PLAINS", "SWAMP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	return card, nil
}
