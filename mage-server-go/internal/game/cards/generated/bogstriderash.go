package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bog Strider Ash", NewBogStriderAsh)
}

// NewBogStriderAsh creates a Bog Strider Ash
// {3}{G} - CREATURE
func NewBogStriderAsh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bog Strider Ash")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new GainLifeEffect(2), new ManaCostsImpl<>("{G}"))
	// card.AddAbility(ability0)
	return card, nil
}