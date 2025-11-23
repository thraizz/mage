package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mishras Research Desk", NewMishrasResearchDesk)
}

// NewMishrasResearchDesk creates a Mishras Research Desk
// {1} - ARTIFACT
func NewMishrasResearchDesk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mishras Research Desk")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
