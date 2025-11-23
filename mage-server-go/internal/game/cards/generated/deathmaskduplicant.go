package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Death Mask Duplicant", NewDeathMaskDuplicant)
}

// NewDeathMaskDuplicant creates a Death Mask Duplicant
// {7} - ARTIFACT CREATURE
func NewDeathMaskDuplicant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Death Mask Duplicant")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewExileTargetEffect()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
