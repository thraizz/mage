package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Pilfered Proof", NewCaseOfThePilferedProof)
}

// NewCaseOfThePilferedProof creates a Case Of The Pilfered Proof
// {1}{W} - ENCHANTMENT
func NewCaseOfThePilferedProof(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Pilfered Proof")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}