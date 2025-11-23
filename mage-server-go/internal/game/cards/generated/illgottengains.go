package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ill Gotten Gains", NewIllGottenGains)
}

// NewIllGottenGains creates a Ill Gotten Gains
// {2}{B}{B} - SORCERY
func NewIllGottenGains(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ill Gotten Gains")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
