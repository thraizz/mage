package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nicol Bolas The Deceiver", NewNicolBolasTheDeceiver)
}

// NewNicolBolasTheDeceiver creates a Nicol Bolas The Deceiver
// {5}{U}{B}{R} - PLANESWALKER
func NewNicolBolasTheDeceiver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nicol Bolas The Deceiver")
	card.ManaCost = "{5}{U}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BOLAS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
