package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Relentless Pursuit", NewRelentlessPursuit)
}

// NewRelentlessPursuit creates a Relentless Pursuit
// {2}{G} - SORCERY
func NewRelentlessPursuit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Relentless Pursuit")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
