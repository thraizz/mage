package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Immortal Servitude", NewImmortalServitude)
}

// NewImmortalServitude creates a Immortal Servitude
// {X}{W/B}{W/B}{W/B} - SORCERY
func NewImmortalServitude(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Immortal Servitude")
	card.ManaCost = "{X}{W/B}{W/B}{W/B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}