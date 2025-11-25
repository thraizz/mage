package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shadows Lair", NewShadowsLair)
}

// NewShadowsLair creates a Shadows Lair
//   - LAND
func NewShadowsLair(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shadows Lair")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewLoseLifeEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
