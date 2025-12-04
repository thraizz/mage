package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Meeting Of The Five", NewMeetingOfTheFive)
}

// NewMeetingOfTheFive creates a Meeting Of The Five
// {3}{W}{U}{B}{R}{G} - SORCERY
func NewMeetingOfTheFive(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Meeting Of The Five")
	card.ManaCost = "{3}{W}{U}{B}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
