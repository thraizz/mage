package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leyline Of Punishment", NewLeylineOfPunishment)
}

// NewLeylineOfPunishment creates a Leyline Of Punishment
// {2}{R}{R} - ENCHANTMENT
func NewLeylineOfPunishment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leyline Of Punishment")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
