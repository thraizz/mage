package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vish Kal Blood Arbiter", NewVishKalBloodArbiter)
}

// NewVishKalBloodArbiter creates a Vish Kal Blood Arbiter
// {4}{W}{B}{B} - CREATURE
// Flying, Lifelink
func NewVishKalBloodArbiter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vish Kal Blood Arbiter")
	card.ManaCost = "{4}{W}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(VishKalBloodArbiterDynamicValue.instance, VishKalBloodArbiterDynamicValue.instance)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}