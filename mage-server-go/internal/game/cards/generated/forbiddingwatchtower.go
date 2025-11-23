package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Forbidding Watchtower", NewForbiddingWatchtower)
}

// NewForbiddingWatchtower creates a Forbidding Watchtower
//   - LAND
func NewForbiddingWatchtower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forbidding Watchtower")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"SOLDIER"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	return card, nil
}
