package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haunted Cloak", NewHauntedCloak)
}

// NewHauntedCloak creates a Haunted Cloak
// {3} - ARTIFACT
func NewHauntedCloak(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haunted Cloak")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("VigilanceAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("HasteAbility", abilities.AttachmentTypeEquipment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
