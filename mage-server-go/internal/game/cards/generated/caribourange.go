package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Caribou Range", NewCaribouRange)
}

// NewCaribouRange creates a Caribou Range
// {2}{W}{W} - ENCHANTMENT
func NewCaribouRange(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Caribou Range")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("CaribouToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(AttachmentType.AURA)).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeAddAbility)).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(AttachmentType.AURA)).
		AddEffect(abilities.NewGainLifeEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("CaribouToken")
	if err != nil {
		return nil, err
	}
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
