package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Herald Of Anguish", NewHeraldOfAnguish)
}

// NewHeraldOfAnguish creates a Herald Of Anguish
// {5}{B}{B} - CREATURE
// Flying
func NewHeraldOfAnguish(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Herald Of Anguish")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(-2, -2)).
		Build()
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(TargetController.OPPONENT)
	// card.AddAbility(ability2)
	return card, nil
}
