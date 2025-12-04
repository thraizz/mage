package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Locket Of Yesterdays", NewLocketOfYesterdays)
}

// NewLocketOfYesterdays creates a Locket Of Yesterdays
// {1} - ARTIFACT
func NewLocketOfYesterdays(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Locket Of Yesterdays")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
