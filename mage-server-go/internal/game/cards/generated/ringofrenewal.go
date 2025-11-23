package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ring Of Renewal", NewRingOfRenewal)
}

// NewRingOfRenewal creates a Ring Of Renewal
// {5} - ARTIFACT
func NewRingOfRenewal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ring Of Renewal")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - DiscardControllerEffect(1, true)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
