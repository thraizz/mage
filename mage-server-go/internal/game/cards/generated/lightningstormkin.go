package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lightning Stormkin", NewLightningStormkin)
}

// NewLightningStormkin creates a Lightning Stormkin
// {U}{R} - CREATURE
// Flying, Haste
func NewLightningStormkin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lightning Stormkin")
	card.ManaCost = "{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	return card, nil
}
