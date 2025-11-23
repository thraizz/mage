package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Gnottvold Slumbermound", NewGnottvoldSlumbermound)
}

// NewGnottvoldSlumbermound creates a Gnottvold Slumbermound
//   - LAND
func NewGnottvoldSlumbermound(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gnottvold Slumbermound")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("TrollWarriorToken")
	if err != nil {
		return nil, err
	}
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
