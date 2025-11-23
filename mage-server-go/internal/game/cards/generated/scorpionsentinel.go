package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scorpion Sentinel", NewScorpionSentinel)
}

// NewScorpionSentinel creates a Scorpion Sentinel
// {1}{U} - ARTIFACT CREATURE
func NewScorpionSentinel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scorpion Sentinel")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "SCORPION"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(3, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
