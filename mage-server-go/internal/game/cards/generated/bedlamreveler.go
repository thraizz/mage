package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bedlam Reveler", NewBedlamReveler)
}

// NewBedlamReveler creates a Bedlam Reveler
// {6}{R}{R} - CREATURE
func NewBedlamReveler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bedlam Reveler")
	card.ManaCost = "{6}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardHandControllerEffect()
	// card.AddAbility(ability0)
	return card, nil
}
