package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Security Bypass", NewSecurityBypass)
}

// NewSecurityBypass creates a Security Bypass
// {1}{U} - ENCHANTMENT
func NewSecurityBypass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Security Bypass")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeBoostCreature)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsCombatDamageToAPlayerTriggeredAbility(new ConniveSourceEffect(), false) .setTriggerPhrase("Whenever creature deals combat damage to a playerenchanted creature has " + "\"Whenever creature deals combat damage to a player, it connives.\" " + "<i>(Its controller draws a card, then discards a card. If they discarded a nonland cardthey put a +1/+1 counter on creature.)</i>")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}