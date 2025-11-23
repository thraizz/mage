package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quandrix Campus", NewQuandrixCampus)
}

// NewQuandrixCampus creates a Quandrix Campus
//   - LAND
func NewQuandrixCampus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quandrix Campus")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
