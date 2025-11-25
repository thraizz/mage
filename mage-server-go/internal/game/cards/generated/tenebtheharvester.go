package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teneb The Harvester", NewTenebTheHarvester)
}

// NewTenebTheHarvester creates a Teneb The Harvester
// {3}{W}{B}{G} - CREATURE
// Flying
func NewTenebTheHarvester(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teneb The Harvester")
	card.ManaCost = "{3}{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DealsCombatDamageToAPlayerTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldTargetEffect()
	//   - Effect: DoIfCostPaid(new ReturnFromGraveyardToBattlefieldTargetEffect()...)
	// card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	return card, nil
}
