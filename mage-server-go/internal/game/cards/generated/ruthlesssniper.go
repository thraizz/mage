package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ruthless Sniper", NewRuthlessSniper)
}

// NewRuthlessSniper creates a Ruthless Sniper
// {B} - CREATURE
func NewRuthlessSniper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ruthless Sniper")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARCHER"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
