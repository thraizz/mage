package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The First Tyrannic War", NewTheFirstTyrannicWar)
}

// NewTheFirstTyrannicWar creates a The First Tyrannic War
// {2}{G}{U}{R} - ENCHANTMENT
func NewTheFirstTyrannicWar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The First Tyrannic War")
	card.ManaCost = "{2}{G}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
