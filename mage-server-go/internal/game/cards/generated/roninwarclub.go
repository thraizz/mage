package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ronin Warclub", NewRoninWarclub)
}

// NewRoninWarclub creates a Ronin Warclub
// {3} - ARTIFACT
func NewRoninWarclub(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ronin Warclub")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(2, 1)).
		Build()
	if err != nil {
		return nil, err
	}
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
