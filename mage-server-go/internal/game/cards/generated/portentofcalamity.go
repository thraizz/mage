package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Portent Of Calamity", NewPortentOfCalamity)
}

// NewPortentOfCalamity creates a Portent Of Calamity
// {X}{U} - SORCERY
func NewPortentOfCalamity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Portent Of Calamity")
	card.ManaCost = "{X}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}