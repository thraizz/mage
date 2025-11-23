package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Waltz Of Rage", NewWaltzOfRage)
}

// NewWaltzOfRage creates a Waltz Of Rage
// {3}{R}{R} - SORCERY
func NewWaltzOfRage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Waltz Of Rage")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
