package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Summons Of Saruman", NewSummonsOfSaruman)
}

// NewSummonsOfSaruman creates a Summons Of Saruman
// {X}{U}{R} - SORCERY
func NewSummonsOfSaruman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Summons Of Saruman")
	card.ManaCost = "{X}{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
