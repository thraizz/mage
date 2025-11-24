package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Praetors Counsel", NewPraetorsCounsel)
}

// NewPraetorsCounsel creates a Praetors Counsel
// {5}{G}{G}{G} - SORCERY
func NewPraetorsCounsel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Praetors Counsel")
	card.ManaCost = "{5}{G}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
