package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakalite", NewRakalite)
}

// NewRakalite creates a Rakalite
// {6} - ARTIFACT
func NewRakalite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakalite")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PreventDamageToTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	// card.AddAbility(ability0)
	return card, nil
}
