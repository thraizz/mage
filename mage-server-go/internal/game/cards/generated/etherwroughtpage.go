package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Etherwrought Page", NewEtherwroughtPage)
}

// NewEtherwroughtPage creates a Etherwrought Page
// {1}{W}{U}{B} - ARTIFACT
func NewEtherwroughtPage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Etherwrought Page")
	card.ManaCost = "{1}{W}{U}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}