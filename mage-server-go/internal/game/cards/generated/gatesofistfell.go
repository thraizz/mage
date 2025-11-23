package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gates Of Istfell", NewGatesOfIstfell)
}

// NewGatesOfIstfell creates a Gates Of Istfell
//   - LAND
func NewGatesOfIstfell(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gates Of Istfell")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewGainLifeEffect(2)).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
