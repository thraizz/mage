package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Curse Of Conformity", NewCurseOfConformity)
}

// NewCurseOfConformity creates a Curse Of Conformity
// {4}{W} - ENCHANTMENT
func NewCurseOfConformity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Curse Of Conformity")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA", "CURSE"}
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
	return card, nil
}
