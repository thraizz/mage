package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fleetfoot Dancer", NewFleetfootDancer)
}

// NewFleetfootDancer creates a Fleetfoot Dancer
// {1}{R}{G}{W} - CREATURE
// Trample, Lifelink, Haste
func NewFleetfootDancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fleetfoot Dancer")
	card.ManaCost = "{1}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "DRUID"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	return card, nil
}
