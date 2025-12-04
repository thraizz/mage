package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yore Tiller Nephilim", NewYoreTillerNephilim)
}

// NewYoreTillerNephilim creates a Yore Tiller Nephilim
// {W}{U}{B}{R} - CREATURE
func NewYoreTillerNephilim(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yore Tiller Nephilim")
	card.ManaCost = "{W}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NEPHILIM"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldTargetEffect(true, true)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect(true, true)
	// card.AddAbility(ability1)
	return card, nil
}
