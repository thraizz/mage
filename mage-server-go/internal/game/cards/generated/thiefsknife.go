package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thiefs Knife", NewThiefsKnife)
}

// NewThiefsKnife creates a Thiefs Knife
// {2}{U} - ARTIFACT
func NewThiefsKnife(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thiefs Knife")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{4}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(1, 1)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsCombatDamageToAPlayerTriggeredAbility( new DrawCardSourceControllerEffect(1) ), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}