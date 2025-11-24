package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fevered Suspicion", NewFeveredSuspicion)
}

// NewFeveredSuspicion creates a Fevered Suspicion
// {6}{B}{R} - SORCERY
func NewFeveredSuspicion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fevered Suspicion")
	card.ManaCost = "{6}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
