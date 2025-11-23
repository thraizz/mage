package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaldra Compleat", NewKaldraCompleat)
}

// NewKaldraCompleat creates a Kaldra Compleat
// {7} - ARTIFACT
// Indestructible
func NewKaldraCompleat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaldra Compleat")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(5, 5)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FirstStrikeAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("IndestructibleAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("HasteAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewExileTargetEffect()).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsDamageToACreatureTriggeredAbility( // if a creature is dealt lethal damage, it is dies as a state-based action and can't be found to exile new ExileTargetEffect(null, Zone.BATTLEFIELD).setToSourceExileZone(true), true, false, true, ), AttachmentType.EQUIPMENTWhenever creature deals combat damage to a creature, exile that creature.\"")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}