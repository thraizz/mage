package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gyrus Waker Of Corpses", NewGyrusWakerOfCorpses)
}

// NewGyrusWakerOfCorpses creates a Gyrus Waker Of Corpses
// {X}{B}{R}{G} - CREATURE
func NewGyrusWakerOfCorpses(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gyrus Waker Of Corpses")
	card.ManaCost = "{X}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HYDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksTriggeredAbility
	//   - Effect: GyrusWakerOfCorpsesEffect()
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(source.getControllerId(), null, true, 1, true, tru...)
	// card.AddAbility(ability1)
	return card, nil
}
