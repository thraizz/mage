package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Secret Passage", NewSecretPassage)
}

// NewSecretPassage creates a Secret Passage
//   - LAND
func NewSecretPassage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Secret Passage")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - InvestigateEffect()
	//
	// Costs:
	//   - AddManaCost("{4}")
	//   - AddTapCost()
	// card.AddAbility(ability2)
	return card, nil
}
