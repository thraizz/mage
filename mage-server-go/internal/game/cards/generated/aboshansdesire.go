package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aboshans Desire", NewAboshansDesire)
}

// NewAboshansDesire creates a Aboshans Desire
// {U} - ENCHANTMENT
func NewAboshansDesire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aboshans Desire")
	card.ManaCost = "{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeAddAbility)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	// Enchanted creature has flying
	flyingAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(flyingAbility, abilities.AttachmentTypeAura, abilities.DurationWhileOnBattlefield, "Enchanted creature has flying")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	// Threshold: Enchanted creature has shroud (simplified - full implementation would check graveyard)
	shroudAbility := abilities.NewKeywordAbility(card.ID, abilities.KeywordShroud)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(shroudAbility, abilities.AttachmentTypeAura, abilities.DurationWhileOnBattlefield, "Enchanted creature has shroud")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
