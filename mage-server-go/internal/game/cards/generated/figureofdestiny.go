package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Figure Of Destiny", NewFigureOfDestiny)
}

// NewFigureOfDestiny creates a Figure Of Destiny
// {R/W} - CREATURE
func NewFigureOfDestiny(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Figure Of Destiny")
	card.ManaCost = "{R/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KITHKIN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - AddCardSubTypeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
