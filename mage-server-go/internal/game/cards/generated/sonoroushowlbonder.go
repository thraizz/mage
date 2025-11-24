package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sonorous Howlbonder", NewSonorousHowlbonder)
}

// NewSonorousHowlbonder creates a Sonorous Howlbonder
// {1}{B/R}{B/R} - CREATURE
func NewSonorousHowlbonder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sonorous Howlbonder")
	card.ManaCost = "{1}{B/R}{B/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
