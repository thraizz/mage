package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Angrath Captain Of Chaos", NewAngrathCaptainOfChaos)
}

// NewAngrathCaptainOfChaos creates a Angrath Captain Of Chaos
// {2}{B/R}{B/R} - PLANESWALKER
func NewAngrathCaptainOfChaos(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Angrath Captain Of Chaos")
	card.ManaCost = "{2}{B/R}{B/R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ANGRATH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: GainAbilityControlledEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
