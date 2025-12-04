package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Garruks Harbinger", NewGarruksHarbinger)
}

// NewGarruksHarbinger creates a Garruks Harbinger
// {1}{G}{G} - CREATURE
func NewGarruksHarbinger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Garruks Harbinger")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(SavedDamageValue.MANY, 1, filter, PutCards.HAND, P...)
	// card.AddAbility(ability0)
	return card, nil
}
