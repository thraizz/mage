package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Overwhelming Stampede", NewOverwhelmingStampede)
}

// NewOverwhelmingStampede creates a Overwhelming Stampede
// {3}{G}{G} - SORCERY
func NewOverwhelmingStampede(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Overwhelming Stampede")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
