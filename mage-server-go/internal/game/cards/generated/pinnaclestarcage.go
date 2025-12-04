package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pinnacle Starcage", NewPinnacleStarcage)
}

// NewPinnacleStarcage creates a Pinnacle Starcage
// {1}{W}{W} - ARTIFACT
func NewPinnacleStarcage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pinnacle Starcage")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PinnacleStarcageTokenEffect()
	//   - SacrificeSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
