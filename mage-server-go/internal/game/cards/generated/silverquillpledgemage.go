package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Silverquill Pledgemage", NewSilverquillPledgemage)
}

// NewSilverquillPledgemage creates a Silverquill Pledgemage
// {1}{W/B}{W/B} - CREATURE
// Flying, Lifelink
func NewSilverquillPledgemage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silverquill Pledgemage")
	card.ManaCost = "{1}{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "CLERIC"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
