package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blitzwing Cruel Tormentor", NewBlitzwingCruelTormentor)
}

// NewBlitzwingCruelTormentor creates a Blitzwing Cruel Tormentor
// {5}{B} - ARTIFACT CREATURE
func NewBlitzwingCruelTormentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blitzwing Cruel Tormentor")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: BlitzwingCruelTormentorEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
