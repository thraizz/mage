package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sibsig Host", NewSibsigHost)
}

// NewSibsigHost creates a Sibsig Host
// {4}{B} - CREATURE
func NewSibsigHost(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sibsig Host")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE"}
	card.Power = "2"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
