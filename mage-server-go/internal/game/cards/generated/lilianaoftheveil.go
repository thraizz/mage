package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Liliana Of The Veil", NewLilianaOfTheVeil)
}

// NewLilianaOfTheVeil creates a Liliana Of The Veil
// {1}{B}{B} - PLANESWALKER
func NewLilianaOfTheVeil(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Liliana Of The Veil")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LILIANA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect()
	// card.AddAbility(ability0)
	return card, nil
}