package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dryad Militant", NewDryadMilitant)
}

// NewDryadMilitant creates a Dryad Militant
// {G/W} - CREATURE
func NewDryadMilitant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dryad Militant")
	card.ManaCost = "{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRYAD", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}