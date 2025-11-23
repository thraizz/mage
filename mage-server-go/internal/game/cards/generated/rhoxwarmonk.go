package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rhox War Monk", NewRhoxWarMonk)
}

// NewRhoxWarMonk creates a Rhox War Monk
// {G}{W}{U} - CREATURE
// Lifelink
func NewRhoxWarMonk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rhox War Monk")
	card.ManaCost = "{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RHINO", "MONK"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	return card, nil
}
