package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brilliant Wings", NewBrilliantWings)
}

// NewBrilliantWings creates a Brilliant Wings
// {1}{W} - ENCHANTMENT
// Flash
func NewBrilliantWings(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brilliant Wings")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBoostCreature)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("HexproofAbility", abilities.AttachmentTypeAura)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new AttachEffect(Outcome....)
	// card.AddAbility(ability4)
	return card, nil
}
