package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rush Of Inspiration", NewRushOfInspiration)
}

// NewRushOfInspiration creates a Rush Of Inspiration
//  - INSTANT
func NewRushOfInspiration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rush Of Inspiration")
	card.ManaCost = ""
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1, true)
	// card.AddAbility(ability2)
	return card, nil
}