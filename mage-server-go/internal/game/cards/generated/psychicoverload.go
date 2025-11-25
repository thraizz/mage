package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Psychic Overload", NewPsychicOverload)
}

// NewPsychicOverload creates a Psychic Overload
// {3}{U} - ENCHANTMENT - AURA
// Enchant permanent
// When Psychic Overload enters the battlefield, tap enchanted permanent.
// Enchanted permanent doesn't untap during its controller's untap step.
// Enchanted permanent has "Discard two artifact cards: Untap this permanent."
func NewPsychicOverload(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Psychic Overload")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA"}
	card.SetCode = "DST"
	card.Rarity = "uncommon"

	// Enchant permanent
	ability0 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter()))
	card.AddAbility(ability0)

	// Spell ability: attach to target permanent
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomeDetriment)).
		AddTarget(abilities.NewPermanentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)

	// When Psychic Overload enters the battlefield, tap enchanted permanent.
	ability2 := abilities.NewTriggeredAbilityBuilder(card.ID).
		SetTrigger(abilities.NewEntersBattlefieldTrigger(card.ID)).
		AddEffect(abilities.NewTapEnchantedEffect()).
		Build()
	card.AddAbility(ability2)

	// Enchanted permanent doesn't untap during its controller's untap step.
	ability3 := abilities.NewSimpleStaticAbility(card.ID, abilities.ZoneBattlefield).
		AddEffect(abilities.NewDontUntapInControllersUntapStepEnchantedEffect()).
		Build()
	card.AddAbility(ability3)

	// Enchanted permanent has "Discard two artifact cards: Untap this permanent."
	// First, create the activated ability that will be granted
	grantedAbility := abilities.NewActivatedAbilityBuilder(card.ID).
		AddCost(abilities.NewDiscardTargetCost(2, abilities.NewArtifactCardFilter())).
		AddEffect(abilities.NewUntapSourceEffect()).
		Build()

	// Then create the static ability that grants this to the enchanted permanent
	ability4 := abilities.NewSimpleStaticAbility(card.ID, abilities.ZoneBattlefield).
		AddEffect(abilities.NewGainAbilityAttachedEffect(
			grantedAbility,
			abilities.AttachmentTypeAura,
			abilities.DurationWhileOnBattlefield,
			"Enchanted permanent has \"Discard two artifact cards: Untap this permanent.\"",
		)).
		Build()
	card.AddAbility(ability4)

	return card, nil
}
