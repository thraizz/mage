package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vedalken Infiltrator", NewVedalkenInfiltrator)
}

// NewVedalkenInfiltrator creates a Vedalken Infiltrator
// {1}{U} - CREATURE
func NewVedalkenInfiltrator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vedalken Infiltrator")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VEDALKEN", "ROGUE"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
