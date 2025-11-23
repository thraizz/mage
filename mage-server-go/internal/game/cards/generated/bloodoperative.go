package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blood Operative", NewBloodOperative)
}

// NewBloodOperative creates a Blood Operative
// {1}{B}{B} - CREATURE
// Lifelink
func NewBloodOperative(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blood Operative")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "ASSASSIN"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new ReturnSourceFromGraveyardToHa...)
	// card.AddAbility(ability2)
	return card, nil
}
