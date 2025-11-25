package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stromgald Spy", NewStromgaldSpy)
}

// NewStromgaldSpy creates a Stromgald Spy
// {3}{B} - CREATURE
func NewStromgaldSpy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stromgald Spy")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksAndIsNotBlockedTriggeredAbility
	//   - Effect: ConditionalContinuousEffect()
	// card.AddAbility(ability0)
	return card, nil
}
