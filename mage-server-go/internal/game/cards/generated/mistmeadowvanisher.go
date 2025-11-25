package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mistmeadow Vanisher", NewMistmeadowVanisher)
}

// NewMistmeadowVanisher creates a Mistmeadow Vanisher
// {2}{W/U} - CREATURE
func NewMistmeadowVanisher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mistmeadow Vanisher")
	card.ManaCost = "{2}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KITHKIN", "WIZARD"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BecomesTappedSourceTriggeredAbility
	//   - Effect: ExileReturnBattlefieldNextEndStepTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
