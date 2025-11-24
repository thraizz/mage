package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ritual Of Subdual", NewRitualOfSubdual)
}

// NewRitualOfSubdual creates a Ritual Of Subdual
// {4}{G}{G} - ENCHANTMENT
func NewRitualOfSubdual(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ritual Of Subdual")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}