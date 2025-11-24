package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warden Of The Eye", NewWardenOfTheEye)
}

// NewWardenOfTheEye creates a Warden Of The Eye
// {2}{U}{R}{W} - CREATURE
func NewWardenOfTheEye(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warden Of The Eye")
	card.ManaCost = "{2}{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DJINN", "WIZARD"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnFromGraveyardToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
