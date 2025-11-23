package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slaughter The Strong", NewSlaughterTheStrong)
}

// NewSlaughterTheStrong creates a Slaughter The Strong
// {1}{W}{W} - SORCERY
func NewSlaughterTheStrong(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slaughter The Strong")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
