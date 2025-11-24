package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fractal Summoning", NewFractalSummoning)
}

// NewFractalSummoning creates a Fractal Summoning
// {X}{G/U}{G/U} - SORCERY
func NewFractalSummoning(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fractal Summoning")
	card.ManaCost = "{X}{G/U}{G/U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}