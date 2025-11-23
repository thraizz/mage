package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Plargg And Nassari", NewPlarggAndNassari)
}

// NewPlarggAndNassari creates a Plargg And Nassari
// {3}{R}{R} - CREATURE
func NewPlarggAndNassari(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plargg And Nassari")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "EFREET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
