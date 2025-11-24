package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Coax From The Blind Eternities", NewCoaxFromTheBlindEternities)
}

// NewCoaxFromTheBlindEternities creates a Coax From The Blind Eternities
// {2}{U} - SORCERY
func NewCoaxFromTheBlindEternities(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Coax From The Blind Eternities")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
