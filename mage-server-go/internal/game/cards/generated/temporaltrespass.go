package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Temporal Trespass", NewTemporalTrespass)
}

// NewTemporalTrespass creates a Temporal Trespass
// {8}{U}{U}{U} - SORCERY
func NewTemporalTrespass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Temporal Trespass")
	card.ManaCost = "{8}{U}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
