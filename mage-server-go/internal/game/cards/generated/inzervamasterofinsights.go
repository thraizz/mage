package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Inzerva Master Of Insights", NewInzervaMasterOfInsights)
}

// NewInzervaMasterOfInsights creates a Inzerva Master Of Insights
// {1}{2/U}{2/R} - PLANESWALKER
func NewInzervaMasterOfInsights(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Inzerva Master Of Insights")
	card.ManaCost = "{1}{2/U}{2/R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"INZERVA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
