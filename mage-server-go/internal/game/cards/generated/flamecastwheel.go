package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flamecast Wheel", NewFlamecastWheel)
}

// NewFlamecastWheel creates a Flamecast Wheel
// {1} - ARTIFACT
func NewFlamecastWheel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flamecast Wheel")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{5}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(3)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
