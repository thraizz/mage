package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nahiris Lithoforming", NewNahirisLithoforming)
}

// NewNahirisLithoforming creates a Nahiris Lithoforming
// {X}{R}{R} - SORCERY
func NewNahirisLithoforming(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nahiris Lithoforming")
	card.ManaCost = "{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
