package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mimeoplasm Revered One", NewMimeoplasmReveredOne)
}

// NewMimeoplasmReveredOne creates a Mimeoplasm Revered One
// {X}{B}{G}{U} - CREATURE
func NewMimeoplasmReveredOne(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mimeoplasm Revered One")
	card.ManaCost = "{X}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OOZE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}