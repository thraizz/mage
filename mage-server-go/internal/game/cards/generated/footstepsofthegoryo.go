package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Footsteps Of The Goryo", NewFootstepsOfTheGoryo)
}

// NewFootstepsOfTheGoryo creates a Footsteps Of The Goryo
// {2}{B} - SORCERY
func NewFootstepsOfTheGoryo(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Footsteps Of The Goryo")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice that creature at the beginning of next ...)
	// card.AddAbility(ability0)
	return card, nil
}
