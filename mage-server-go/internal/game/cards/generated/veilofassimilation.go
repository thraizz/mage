package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Veil Of Assimilation", NewVeilOfAssimilation)
}

// NewVeilOfAssimilation creates a Veil Of Assimilation
// {1}{W} - ARTIFACT
func NewVeilOfAssimilation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Veil Of Assimilation")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}