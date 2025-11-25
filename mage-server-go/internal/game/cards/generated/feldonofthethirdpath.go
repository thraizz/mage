package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Feldon Of The Third Path", NewFeldonOfTheThirdPath)
}

// NewFeldonOfTheThirdPath creates a Feldon Of The Third Path
// {1}{R}{R} - CREATURE
func NewFeldonOfTheThirdPath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Feldon Of The Third Path")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice the token at the beginning of the next ...)
	//   - CreateTokenCopyTargetEffect(source.getControllerId(), CardType.ARTIFACT, true)
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - FeldonOfTheThirdPathEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
