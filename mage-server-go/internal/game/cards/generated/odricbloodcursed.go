package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Odric Blood Cursed", NewOdricBloodCursed)
}

// NewOdricBloodCursed creates a Odric Blood Cursed
// {1}{R}{W} - CREATURE
func NewOdricBloodCursed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Odric Blood Cursed")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
