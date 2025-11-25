package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Acklay Of The Arena", NewAcklayOfTheArena)
}

// NewAcklayOfTheArena creates a Acklay Of The Arena
// {1}{R}{G}{W} - CREATURE
func NewAcklayOfTheArena(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Acklay Of The Arena")
	card.ManaCost = "{1}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "BEAST"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BecomesMonstrousTriggeredAbility
	//   - Effect: FightTargetsEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
