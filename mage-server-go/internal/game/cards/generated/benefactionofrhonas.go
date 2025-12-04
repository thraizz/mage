package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Benefaction Of Rhonas", NewBenefactionOfRhonas)
}

// NewBenefactionOfRhonas creates a Benefaction Of Rhonas
// {2}{G} - SORCERY
func NewBenefactionOfRhonas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Benefaction Of Rhonas")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
