package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mask Of Intolerance", NewMaskOfIntolerance)
}

// NewMaskOfIntolerance creates a Mask Of Intolerance
// {2} - ARTIFACT
func NewMaskOfIntolerance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mask Of Intolerance")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(3, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
