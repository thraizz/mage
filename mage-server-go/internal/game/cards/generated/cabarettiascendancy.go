package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cabaretti Ascendancy", NewCabarettiAscendancy)
}

// NewCabarettiAscendancy creates a Cabaretti Ascendancy
// {R}{G}{W} - ENCHANTMENT
func NewCabarettiAscendancy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cabaretti Ascendancy")
	card.ManaCost = "{R}{G}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
