package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferi Timeless Voyager", NewTeferiTimelessVoyager)
}

// NewTeferiTimelessVoyager creates a Teferi Timeless Voyager
// {4}{U}{U} - PLANESWALKER
func NewTeferiTimelessVoyager(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferi Timeless Voyager")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEFERI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
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
