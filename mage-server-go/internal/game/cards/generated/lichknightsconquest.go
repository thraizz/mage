package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lich Knights Conquest", NewLichKnightsConquest)
}

// NewLichKnightsConquest creates a Lich Knights Conquest
// {4}{B} - SORCERY
func NewLichKnightsConquest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lich Knights Conquest")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
