package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Farewell", NewFarewell)
}

// NewFarewell creates a Farewell
// {4}{W}{W} - SORCERY
func NewFarewell(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Farewell")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect()
	// card.AddAbility(ability0)
	return card, nil
}
