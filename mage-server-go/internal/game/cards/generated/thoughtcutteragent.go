package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thoughtcutter Agent", NewThoughtcutterAgent)
}

// NewThoughtcutterAgent creates a Thoughtcutter Agent
// {U}{B} - ARTIFACT CREATURE
func NewThoughtcutterAgent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thoughtcutter Agent")
	card.ManaCost = "{U}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealHandTargetEffect(TargetController.ANY)
	// card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewLoseLifeEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
