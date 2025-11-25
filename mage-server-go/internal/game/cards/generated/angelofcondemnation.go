package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Angel Of Condemnation", NewAngelOfCondemnation)
}

// NewAngelOfCondemnation creates a Angel Of Condemnation
// {2}{W}{W} - CREATURE
// Flying, Vigilance
func NewAngelOfCondemnation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Angel Of Condemnation")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - ExileReturnBattlefieldNextEndStepTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability2)
	// TODO: Implement activated ability with unmapped effects
	//   - ExileUntilSourceLeavesEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability3)
	return card, nil
}
