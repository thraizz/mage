package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exocrine", NewExocrine)
}

// NewExocrine creates a Exocrine
// {X}{2}{R} - CREATURE
func NewExocrine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exocrine")
	card.ManaCost = "{X}{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TYRANID"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
