package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Diabolic Revelation", NewDiabolicRevelation)
}

// NewDiabolicRevelation creates a Diabolic Revelation
// {X}{3}{B}{B} - SORCERY
func NewDiabolicRevelation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Diabolic Revelation")
	card.ManaCost = "{X}{3}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
