package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jetfire Air Guardian", NewJetfireAirGuardian)
}

// NewJetfireAirGuardian creates a Jetfire Air Guardian
//   - ARTIFACT
//
// Flying
func NewJetfireAirGuardian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jetfire Air Guardian")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - TransformSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
