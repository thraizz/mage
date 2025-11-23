package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ascend From Avernus", NewAscendFromAvernus)
}

// NewAscendFromAvernus creates a Ascend From Avernus
// {X}{W}{W}{W} - SORCERY
func NewAscendFromAvernus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ascend From Avernus")
	card.ManaCost = "{X}{W}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
