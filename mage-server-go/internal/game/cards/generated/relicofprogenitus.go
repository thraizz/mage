package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Relic Of Progenitus", NewRelicOfProgenitus)
}

// NewRelicOfProgenitus creates a Relic Of Progenitus
// {1} - ARTIFACT
func NewRelicOfProgenitus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Relic Of Progenitus")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability0)
	return card, nil
}
