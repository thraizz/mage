package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ruthless Invasion", NewRuthlessInvasion)
}

// NewRuthlessInvasion creates a Ruthless Invasion
// {3}{R/P} - SORCERY
func NewRuthlessInvasion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ruthless Invasion")
	card.ManaCost = "{3}{R/P}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
