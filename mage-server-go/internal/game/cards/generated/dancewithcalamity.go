package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dance With Calamity", NewDanceWithCalamity)
}

// NewDanceWithCalamity creates a Dance With Calamity
// {7}{R} - SORCERY
func NewDanceWithCalamity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dance With Calamity")
	card.ManaCost = "{7}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
