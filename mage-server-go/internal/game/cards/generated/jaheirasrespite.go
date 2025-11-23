package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jaheiras Respite", NewJaheirasRespite)
}

// NewJaheirasRespite creates a Jaheiras Respite
// {4}{G} - INSTANT
func NewJaheirasRespite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jaheiras Respite")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
