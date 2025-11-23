package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexias Core", NewPhyrexiasCore)
}

// NewPhyrexiasCore creates a Phyrexias Core
//   - LAND
func NewPhyrexiasCore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexias Core")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddManaCost("{1}").
		AddEffect(abilities.NewGainLifeEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
