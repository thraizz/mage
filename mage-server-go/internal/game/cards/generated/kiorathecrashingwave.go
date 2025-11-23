package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kiora The Crashing Wave", NewKioraTheCrashingWave)
}

// NewKioraTheCrashingWave creates a Kiora The Crashing Wave
// {2}{G}{U} - PLANESWALKER
func NewKioraTheCrashingWave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kiora The Crashing Wave")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KIORA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
