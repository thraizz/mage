package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arboria", NewArboria)
}

// NewArboria creates a Arboria
// {2}{G}{G} - ENCHANTMENT
func NewArboria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arboria")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"WORLD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
