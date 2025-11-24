package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Raphaels Technique", NewRaphaelsTechnique)
}

// NewRaphaelsTechnique creates a Raphaels Technique
// {4}{R}{R} - INSTANT
func NewRaphaelsTechnique(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Raphaels Technique")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}