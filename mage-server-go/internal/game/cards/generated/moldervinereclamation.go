package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Moldervine Reclamation", NewMoldervineReclamation)
}

// NewMoldervineReclamation creates a Moldervine Reclamation
// {3}{B}{G} - ENCHANTMENT
func NewMoldervineReclamation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moldervine Reclamation")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
