package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soulfire Eruption", NewSoulfireEruption)
}

// NewSoulfireEruption creates a Soulfire Eruption
// {6}{R}{R}{R} - SORCERY
func NewSoulfireEruption(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soulfire Eruption")
	card.ManaCost = "{6}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
