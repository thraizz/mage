package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Muraganda", NewInvasionOfMuraganda)
}

// NewInvasionOfMuraganda creates a Invasion Of Muraganda
// {4}{G} - BATTLE
func NewInvasionOfMuraganda(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Muraganda")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}