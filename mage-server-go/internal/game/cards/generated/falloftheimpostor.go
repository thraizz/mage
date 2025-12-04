package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fall Of The Impostor", NewFallOfTheImpostor)
}

// NewFallOfTheImpostor creates a Fall Of The Impostor
// {1}{G}{W} - ENCHANTMENT
func NewFallOfTheImpostor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fall Of The Impostor")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
