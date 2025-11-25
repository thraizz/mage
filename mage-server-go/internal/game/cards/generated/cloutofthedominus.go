package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clout Of The Dominus", NewCloutOfTheDominus)
}

// NewCloutOfTheDominus creates a Clout Of The Dominus
// {U/R} - ENCHANTMENT
func NewCloutOfTheDominus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clout Of The Dominus")
	card.ManaCost = "{U/R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeAddAbility)).
		AddEffect(abilities.NewBoostEnchantedEffect(1, 1)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordShroud), abilities.AttachmentTypeAura, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordShroud), abilities.AttachmentTypeAura, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste), abilities.AttachmentTypeAura, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
