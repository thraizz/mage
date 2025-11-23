package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Liberated Livestock", NewLiberatedLivestock)
}

// NewLiberatedLivestock creates a Liberated Livestock
// {5}{W} - CREATURE
func NewLiberatedLivestock(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Liberated Livestock")
	card.ManaCost = "{5}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "BIRD", "OX"}
	card.Power = "4"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
