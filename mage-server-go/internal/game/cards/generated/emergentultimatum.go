package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emergent Ultimatum", NewEmergentUltimatum)
}

// NewEmergentUltimatum creates a Emergent Ultimatum
// {B}{B}{G}{G}{G}{U}{U} - SORCERY
func NewEmergentUltimatum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emergent Ultimatum")
	card.ManaCost = "{B}{B}{G}{G}{G}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
