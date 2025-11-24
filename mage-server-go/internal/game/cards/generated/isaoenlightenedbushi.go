package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Isao Enlightened Bushi", NewIsaoEnlightenedBushi)
}

// NewIsaoEnlightenedBushi creates a Isao Enlightened Bushi
// {2}{G} - CREATURE
func NewIsaoEnlightenedBushi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Isao Enlightened Bushi")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SAMURAI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	// card.AddAbility(ability0)
	return card, nil
}
