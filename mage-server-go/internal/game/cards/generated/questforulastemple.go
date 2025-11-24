package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quest For Ulas Temple", NewQuestForUlasTemple)
}

// NewQuestForUlasTemple creates a Quest For Ulas Temple
// {U} - ENCHANTMENT
func NewQuestForUlasTemple(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quest For Ulas Temple")
	card.ManaCost = "{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}