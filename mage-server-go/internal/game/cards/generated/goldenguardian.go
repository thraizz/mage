package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Golden Guardian", NewGoldenGuardian)
}

// NewGoldenGuardian creates a Golden Guardian
// {4} - ARTIFACT CREATURE
// Defender
func NewGoldenGuardian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Golden Guardian")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - FightTargetSourceEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	// card.AddAbility(ability1)
	return card, nil
}
