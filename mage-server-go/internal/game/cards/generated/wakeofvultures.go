package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wake Of Vultures", NewWakeOfVultures)
}

// NewWakeOfVultures creates a Wake Of Vultures
// {3}{B} - CREATURE
// Flying
func NewWakeOfVultures(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wake Of Vultures")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
