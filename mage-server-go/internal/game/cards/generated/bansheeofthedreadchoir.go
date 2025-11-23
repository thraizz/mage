package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Banshee Of The Dread Choir", NewBansheeOfTheDreadChoir)
}

// NewBansheeOfTheDreadChoir creates a Banshee Of The Dread Choir
// {3}{B}{B} - CREATURE
func NewBansheeOfTheDreadChoir(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Banshee Of The Dread Choir")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
