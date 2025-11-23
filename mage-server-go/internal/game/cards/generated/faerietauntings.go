package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Faerie Tauntings", NewFaerieTauntings)
}

// NewFaerieTauntings creates a Faerie Tauntings
// {2}{B} - KINDRED ENCHANTMENT
func NewFaerieTauntings(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Faerie Tauntings")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"FAERIE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
