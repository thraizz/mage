package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rampage Of The Clans", NewRampageOfTheClans)
}

// NewRampageOfTheClans creates a Rampage Of The Clans
// {3}{G} - INSTANT
func NewRampageOfTheClans(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rampage Of The Clans")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
