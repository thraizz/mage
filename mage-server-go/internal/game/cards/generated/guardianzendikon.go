package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Guardian Zendikon", NewGuardianZendikon)
}

// NewGuardianZendikon creates a Guardian Zendikon
// {2}{W} - ENCHANTMENT
// Defender
func NewGuardianZendikon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Guardian Zendikon")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"AURA", "WALL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	ability1 := abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAttachEffect(abilities.OutcomePutCreatureInPlay)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
