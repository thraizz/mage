package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Olag Ludevics Hubris", NewOlagLudevicsHubris)
}

// NewOlagLudevicsHubris creates a Olag Ludevics Hubris
//   - CREATURE
func NewOlagLudevicsHubris(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Olag Ludevics Hubris")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
