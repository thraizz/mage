package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unexpectedly Absent", NewUnexpectedlyAbsent)
}

// NewUnexpectedlyAbsent creates a Unexpectedly Absent
// {X}{W}{W} - INSTANT
func NewUnexpectedlyAbsent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unexpectedly Absent")
	card.ManaCost = "{X}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
