package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ceremonial Knife", NewCeremonialKnife)
}

// NewCeremonialKnife creates a Ceremonial Knife
// {1} - ARTIFACT
func NewCeremonialKnife(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ceremonial Knife")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{2}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("BloodToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(1, 0)).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsCombatDamageTriggeredAbility( new CreateTokenEffect(token1_0), false ).setTriggerPhrase("Whenever creature deals combat damage, "), AttachmentType.AURA)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}