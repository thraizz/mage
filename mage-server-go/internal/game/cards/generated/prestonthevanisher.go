package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Preston The Vanisher", NewPrestonTheVanisher)
}

// NewPrestonTheVanisher creates a Preston The Vanisher
//
//	-
func NewPrestonTheVanisher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Preston The Vanisher")
	card.ManaCost = ""
	card.Subtypes = []string{"RABBIT", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: ExileTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
