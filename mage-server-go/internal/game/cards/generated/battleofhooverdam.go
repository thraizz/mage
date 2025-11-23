package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Battle Of Hoover Dam", NewBattleOfHooverDam)
}

// NewBattleOfHooverDam creates a Battle Of Hoover Dam
// {3}{W} - ENCHANTMENT
func NewBattleOfHooverDam(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Battle Of Hoover Dam")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
