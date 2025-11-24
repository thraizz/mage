package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Envoy Of Okinec Ahau", NewEnvoyOfOkinecAhau)
}

// NewEnvoyOfOkinecAhau creates a Envoy Of Okinec Ahau
// {2}{W} - CREATURE
func NewEnvoyOfOkinecAhau(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Envoy Of Okinec Ahau")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "ADVISOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
