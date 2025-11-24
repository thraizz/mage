package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Herbal Poultice", NewHerbalPoultice)
}

// NewHerbalPoultice creates a Herbal Poultice
// {0} - ARTIFACT
func NewHerbalPoultice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Herbal Poultice")
	card.ManaCost = "{0}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}