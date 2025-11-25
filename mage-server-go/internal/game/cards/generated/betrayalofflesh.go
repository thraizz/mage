package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Betrayal Of Flesh", NewBetrayalOfFlesh)
}

// NewBetrayalOfFlesh creates a Betrayal Of Flesh
// {5}{B} - INSTANT
func NewBetrayalOfFlesh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Betrayal Of Flesh")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
