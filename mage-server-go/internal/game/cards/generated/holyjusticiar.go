package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Holy Justiciar", NewHolyJusticiar)
}

// NewHolyJusticiar creates a Holy Justiciar
// {3}{W} - CREATURE
func NewHolyJusticiar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Holy Justiciar")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "CLERIC"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
