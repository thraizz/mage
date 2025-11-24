package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Academic Probation", NewAcademicProbation)
}

// NewAcademicProbation creates a Academic Probation
// {1}{W} - SORCERY
func NewAcademicProbation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Academic Probation")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
