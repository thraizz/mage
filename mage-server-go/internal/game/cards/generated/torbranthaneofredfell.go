package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Torbran Thane Of Red Fell", NewTorbranThaneOfRedFell)
}

// NewTorbranThaneOfRedFell creates a Torbran Thane Of Red Fell
// {1}{R}{R}{R} - CREATURE
func NewTorbranThaneOfRedFell(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Torbran Thane Of Red Fell")
	card.ManaCost = "{1}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}