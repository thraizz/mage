package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jace The Mind Sculptor", NewJaceTheMindSculptor)
}

// NewJaceTheMindSculptor creates a Jace The Mind Sculptor
// {2}{U}{U} - PLANESWALKER
func NewJaceTheMindSculptor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jace The Mind Sculptor")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
