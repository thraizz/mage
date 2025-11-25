package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Grim Captains Locker", NewTheGrimCaptainsLocker)
}

// NewTheGrimCaptainsLocker creates a The Grim Captains Locker
// {3}{B} - ARTIFACT
func NewTheGrimCaptainsLocker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Grim Captains Locker")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
