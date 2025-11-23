package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sanity Gnawers", NewSanityGnawers)
}

// NewSanityGnawers creates a Sanity Gnawers
// {1}{B}{R} - CREATURE
func NewSanityGnawers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sanity Gnawers")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1, true)
	// card.AddAbility(ability0)
	return card, nil
}
