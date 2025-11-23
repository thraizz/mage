package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phoenix Warden Of Fire", NewPhoenixWardenOfFire)
}

// NewPhoenixWardenOfFire creates a Phoenix Warden Of Fire
//   - ENCHANTMENT CREATURE
//
// Flying, Lifelink
func NewPhoenixWardenOfFire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phoenix Warden Of Fire")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SAGA", "PHOENIX"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
