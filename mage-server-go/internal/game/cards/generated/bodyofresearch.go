package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Body Of Research", NewBodyOfResearch)
}

// NewBodyOfResearch creates a Body Of Research
// {G}{G}{G}{U}{U}{U} - SORCERY
func NewBodyOfResearch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Body Of Research")
	card.ManaCost = "{G}{G}{G}{U}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
