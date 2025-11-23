package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sunastian Falconer", NewSunastianFalconer)
}

// NewSunastianFalconer creates a Sunastian Falconer
// {3}{R}{G} - CREATURE
func NewSunastianFalconer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sunastian Falconer")
	card.ManaCost = "{3}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SHAMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
