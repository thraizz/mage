package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hidetsugu Devouring Chaos", NewHidetsuguDevouringChaos)
}

// NewHidetsuguDevouringChaos creates a Hidetsugu Devouring Chaos
// {3}{B} - CREATURE
func NewHidetsuguDevouringChaos(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidetsugu Devouring Chaos")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "DEMON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}