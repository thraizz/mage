package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Feast Of The Victorious Dead", NewFeastOfTheVictoriousDead)
}

// NewFeastOfTheVictoriousDead creates a Feast Of The Victorious Dead
// {W}{B} - ENCHANTMENT
func NewFeastOfTheVictoriousDead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Feast Of The Victorious Dead")
	card.ManaCost = "{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
