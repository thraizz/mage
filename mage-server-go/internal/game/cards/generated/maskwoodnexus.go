package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Maskwood Nexus", NewMaskwoodNexus)
}

// NewMaskwoodNexus creates a Maskwood Nexus
// {4} - ARTIFACT
func NewMaskwoodNexus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maskwood Nexus")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("ShapeshifterBlueToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
