package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ugin Eye Of The Storms", NewUginEyeOfTheStorms)
}

// NewUginEyeOfTheStorms creates a Ugin Eye Of The Storms
// {7} - PLANESWALKER
func NewUginEyeOfTheStorms(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ugin Eye Of The Storms")
	card.ManaCost = "{7}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"UGIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileTargetEffect()).
		AddEffect(abilities.NewExileTargetEffect()).
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
