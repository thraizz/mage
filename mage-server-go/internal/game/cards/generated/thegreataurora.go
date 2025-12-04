package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Great Aurora", NewTheGreatAurora)
}

// NewTheGreatAurora creates a The Great Aurora
// {6}{G}{G}{G} - SORCERY
func NewTheGreatAurora(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Great Aurora")
	card.ManaCost = "{6}{G}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
