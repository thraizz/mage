package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yawgmoths Agenda", NewYawgmothsAgenda)
}

// NewYawgmothsAgenda creates a Yawgmoths Agenda
// {3}{B}{B} - ENCHANTMENT
func NewYawgmothsAgenda(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yawgmoths Agenda")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
