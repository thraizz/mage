package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gelatinous Genesis", NewGelatinousGenesis)
}

// NewGelatinousGenesis creates a Gelatinous Genesis
// {X}{X}{G} - SORCERY
func NewGelatinousGenesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gelatinous Genesis")
	card.ManaCost = "{X}{X}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
