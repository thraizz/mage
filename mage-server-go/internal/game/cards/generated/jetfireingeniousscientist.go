package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jetfire Ingenious Scientist", NewJetfireIngeniousScientist)
}

// NewJetfireIngeniousScientist creates a Jetfire Ingenious Scientist
// {4}{U} - ARTIFACT CREATURE
// Flying
func NewJetfireIngeniousScientist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jetfire Ingenious Scientist")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - JetfireIngeniousScientistEffect()
	//   - TransformSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
