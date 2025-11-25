package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dovin Architect Of Law", NewDovinArchitectOfLaw)
}

// NewDovinArchitectOfLaw creates a Dovin Architect Of Law
// {4}{W}{U} - PLANESWALKER
func NewDovinArchitectOfLaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dovin Architect Of Law")
	card.ManaCost = "{4}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOVIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(2)).
		AddEffect(abilities.NewTapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
