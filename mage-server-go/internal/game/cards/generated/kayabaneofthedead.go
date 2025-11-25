package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaya Bane Of The Dead", NewKayaBaneOfTheDead)
}

// NewKayaBaneOfTheDead creates a Kaya Bane Of The Dead
// {3}{W/B}{W/B}{W/B} - PLANESWALKER
func NewKayaBaneOfTheDead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaya Bane Of The Dead")
	card.ManaCost = "{3}{W/B}{W/B}{W/B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
