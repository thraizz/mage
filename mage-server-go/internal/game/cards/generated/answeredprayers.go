package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Answered Prayers", NewAnsweredPrayers)
}

// NewAnsweredPrayers creates a Answered Prayers
// {1}{W}{W} - ENCHANTMENT
// Flying
func NewAnsweredPrayers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Answered Prayers")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"ANGEL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
