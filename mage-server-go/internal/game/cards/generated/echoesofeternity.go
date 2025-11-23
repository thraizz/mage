package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Echoes Of Eternity", NewEchoesOfEternity)
}

// NewEchoesOfEternity creates a Echoes Of Eternity
// {3}{C}{C}{C} - KINDRED ENCHANTMENT
func NewEchoesOfEternity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Echoes Of Eternity")
	card.ManaCost = "{3}{C}{C}{C}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"ELDRAZI"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
