package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mouth Of Ronom", NewMouthOfRonom)
}

// NewMouthOfRonom creates a Mouth Of Ronom
//   - LAND
func NewMouthOfRonom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mouth Of Ronom")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(4)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
