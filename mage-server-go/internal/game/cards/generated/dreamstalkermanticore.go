package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dreamstalker Manticore", NewDreamstalkerManticore)
}

// NewDreamstalkerManticore creates a Dreamstalker Manticore
// {2}{R} - ENCHANTMENT CREATURE
func NewDreamstalkerManticore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dreamstalker Manticore")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"MANTICORE"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
