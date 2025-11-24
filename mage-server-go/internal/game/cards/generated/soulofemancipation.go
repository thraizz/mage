package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soul Of Emancipation", NewSoulOfEmancipation)
}

// NewSoulOfEmancipation creates a Soul Of Emancipation
// {4}{G}{W}{U} - CREATURE
func NewSoulOfEmancipation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soul Of Emancipation")
	card.ManaCost = "{4}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "5"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
