package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Summon G F Cerberus", NewSummonGFCerberus)
}

// NewSummonGFCerberus creates a Summon G F Cerberus
// {2}{R}{R} - ENCHANTMENT CREATURE
func NewSummonGFCerberus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Summon G F Cerberus")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SAGA", "DOG"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
