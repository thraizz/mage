package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Suntail Hawk", NewSuntailHawk)
}

// NewSuntailHawk creates a Suntail Hawk
// {W} - CREATURE
// Flying
func NewSuntailHawk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Suntail Hawk")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
