package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Croaking Counterpart", NewCroakingCounterpart)
}

// NewCroakingCounterpart creates a Croaking Counterpart
// {1}{G}{U} - SORCERY
func NewCroakingCounterpart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Croaking Counterpart")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
