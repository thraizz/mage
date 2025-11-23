package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Nurturing Presence", NewNurturingPresence)
}

// NewNurturingPresence creates a Nurturing Presence
// {1}{W} - ENCHANTMENT
func NewNurturingPresence(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nurturing Presence")
	card.ManaCost = "{1}{W}"
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
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new EntersBattlefieldControlledTriggeredAbility( new BoostSourceEffect(1, 1) .setText(" creature gets +1/+1 until end of turn"), ), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	token3_0, err := token.GetToken("SpiritWhiteToken")
	if err != nil {
		return nil, err
	}
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token3_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}