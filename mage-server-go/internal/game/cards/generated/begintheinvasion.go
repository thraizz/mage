package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Begin The Invasion", NewBeginTheInvasion)
}

// NewBeginTheInvasion creates a Begin The Invasion
// {X}{W}{U}{B}{R}{G} - SORCERY
func NewBeginTheInvasion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Begin The Invasion")
	card.ManaCost = "{X}{W}{U}{B}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
