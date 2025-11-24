package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rusted Relic", NewRustedRelic)
}

// NewRustedRelic creates a Rusted Relic
// {4} - ARTIFACT
func NewRustedRelic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rusted Relic")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"GOLEM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}