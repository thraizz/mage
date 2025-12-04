package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Precipice Of Mortis", NewPrecipiceOfMortis)
}

// NewPrecipiceOfMortis creates a Precipice Of Mortis
// {G}{U}{W} - ENCHANTMENT
func NewPrecipiceOfMortis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Precipice Of Mortis")
	card.ManaCost = "{G}{U}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
