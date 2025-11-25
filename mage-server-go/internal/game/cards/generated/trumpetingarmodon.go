package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trumpeting Armodon", NewTrumpetingArmodon)
}

// NewTrumpetingArmodon creates a Trumpeting Armodon
// {3}{G} - CREATURE
func NewTrumpetingArmodon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trumpeting Armodon")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - MustBeBlockedByTargetSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
