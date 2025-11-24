package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nicol Bolas Dragon God", NewNicolBolasDragonGod)
}

// NewNicolBolasDragonGod creates a Nicol Bolas Dragon God
// {U}{B}{B}{B}{R} - PLANESWALKER
func NewNicolBolasDragonGod(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nicol Bolas Dragon God")
	card.ManaCost = "{U}{B}{B}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"BOLAS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
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