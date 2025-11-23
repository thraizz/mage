package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ink Eyes Servant Of Oni", NewInkEyesServantOfOni)
}

// NewInkEyesServantOfOni creates a Ink Eyes Servant Of Oni
// {4}{B}{B} - CREATURE
func NewInkEyesServantOfOni(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ink Eyes Servant Of Oni")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT", "NINJA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
