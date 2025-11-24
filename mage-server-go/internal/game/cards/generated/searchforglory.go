package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Search For Glory", NewSearchForGlory)
}

// NewSearchForGlory creates a Search For Glory
// {2}{W} - SORCERY
func NewSearchForGlory(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Search For Glory")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}