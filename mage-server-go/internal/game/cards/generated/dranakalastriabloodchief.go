package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drana Kalastria Bloodchief", NewDranaKalastriaBloodchief)
}

// NewDranaKalastriaBloodchief creates a Drana Kalastria Bloodchief
// {3}{B}{B} - CREATURE
// Flying
func NewDranaKalastriaBloodchief(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drana Kalastria Bloodchief")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "SHAMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(StaticValue.get(0), new SignInversionDynamicValue(GetXValue.instance))).
		AddEffect(abilities.NewBoostEffect(GetXValue.instance, StaticValue.get(0))).
		Build()
	card.AddAbility(ability1)
	return card, nil
}