package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spinewoods Armadillo", NewSpinewoodsArmadillo)
}

// NewSpinewoodsArmadillo creates a Spinewoods Armadillo
// {4}{G}{G} - CREATURE
// Reach
func NewSpinewoodsArmadillo(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spinewoods Armadillo")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ARMADILLO"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: SearchLibraryPutInHandEffect with complex parameters
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
