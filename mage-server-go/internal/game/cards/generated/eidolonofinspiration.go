package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eidolon Of Inspiration", NewEidolonOfInspiration)
}

// NewEidolonOfInspiration creates a Eidolon Of Inspiration
// {1}{W}{W} - ENCHANTMENT CREATURE
func NewEidolonOfInspiration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eidolon Of Inspiration")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
