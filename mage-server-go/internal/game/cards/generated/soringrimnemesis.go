package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sorin Grim Nemesis", NewSorinGrimNemesis)
}

// NewSorinGrimNemesis creates a Sorin Grim Nemesis
// {4}{W}{B} - PLANESWALKER
func NewSorinGrimNemesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sorin Grim Nemesis")
	card.ManaCost = "{4}{W}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SORIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(GetXValue.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
