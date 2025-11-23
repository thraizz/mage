package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Prosthetic Injector", NewProstheticInjector)
}

// NewProstheticInjector creates a Prosthetic Injector
// {1} - ARTIFACT
func NewProstheticInjector(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prosthetic Injector")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{1}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(0, 2)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new ToxicAbility(1), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}