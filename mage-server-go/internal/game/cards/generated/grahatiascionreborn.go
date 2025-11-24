package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Graha Tia Scion Reborn", NewGrahaTiaScionReborn)
}

// NewGrahaTiaScionReborn creates a Graha Tia Scion Reborn
// {W}{U}{B} - CREATURE
// Lifelink
func NewGrahaTiaScionReborn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Graha Tia Scion Reborn")
	card.ManaCost = "{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new GrahaTiaScionRebornEffect(), ...)
	// card.AddAbility(ability1)
	return card, nil
}