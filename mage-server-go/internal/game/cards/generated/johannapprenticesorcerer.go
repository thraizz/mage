package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Johann Apprentice Sorcerer", NewJohannApprenticeSorcerer)
}

// NewJohannApprenticeSorcerer creates a Johann Apprentice Sorcerer
// {2}{U}{R} - CREATURE
func NewJohannApprenticeSorcerer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Johann Apprentice Sorcerer")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
