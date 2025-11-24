package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nissa Worldsoul Speaker", NewNissaWorldsoulSpeaker)
}

// NewNissaWorldsoulSpeaker creates a Nissa Worldsoul Speaker
// {3}{G} - CREATURE
func NewNissaWorldsoulSpeaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nissa Worldsoul Speaker")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
