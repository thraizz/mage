package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Excalibur Sword Of Eden", NewExcaliburSwordOfEden)
}

// NewExcaliburSwordOfEden creates a Excalibur Sword Of Eden
// {12} - ARTIFACT
func NewExcaliburSwordOfEden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Excalibur Sword Of Eden")
	card.ManaCost = "{12}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(10, 0)).
		AddEffect(abilities.NewGainAbilityAttachedEffect("VigilanceAbility", abilities.AttachmentTypeEquipment)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
