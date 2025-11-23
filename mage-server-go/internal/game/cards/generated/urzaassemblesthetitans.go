package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urza Assembles The Titans", NewUrzaAssemblesTheTitans)
}

// NewUrzaAssemblesTheTitans creates a Urza Assembles The Titans
// {3}{W}{W} - ENCHANTMENT
func NewUrzaAssemblesTheTitans(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urza Assembles The Titans")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
