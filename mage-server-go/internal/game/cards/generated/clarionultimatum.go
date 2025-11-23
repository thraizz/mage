package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clarion Ultimatum", NewClarionUltimatum)
}

// NewClarionUltimatum creates a Clarion Ultimatum
// {G}{G}{W}{W}{W}{U}{U} - SORCERY
func NewClarionUltimatum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clarion Ultimatum")
	card.ManaCost = "{G}{G}{W}{W}{W}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
