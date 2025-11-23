package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Finale Of Promise", NewFinaleOfPromise)
}

// NewFinaleOfPromise creates a Finale Of Promise
// {X}{R}{R} - SORCERY
func NewFinaleOfPromise(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Finale Of Promise")
	card.ManaCost = "{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
