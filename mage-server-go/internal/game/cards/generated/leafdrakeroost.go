package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Leafdrake Roost", NewLeafdrakeRoost)
}

// NewLeafdrakeRoost creates a Leafdrake Roost
// {3}{G}{U} - ENCHANTMENT
func NewLeafdrakeRoost(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leafdrake Roost")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomePutCreatureInPlay)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("LeafdrakeRoostDrakeToken")
	if err != nil {
		return nil, err
	}
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilityToGain, AttachmentType.AURA{G}{U}, {T}: Create a 2/2 green and blue Drake creature token with flying.\"")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}