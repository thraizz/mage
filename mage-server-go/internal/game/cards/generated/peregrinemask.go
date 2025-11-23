package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Peregrine Mask", NewPeregrineMask)
}

// NewPeregrineMask creates a Peregrine Mask
// {1} - ARTIFACT
func NewPeregrineMask(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Peregrine Mask")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("DefenderAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FlyingAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("FirstStrikeAbility", abilities.AttachmentTypeEquipment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
