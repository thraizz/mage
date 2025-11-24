package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Ikoria", NewInvasionOfIkoria)
}

// NewInvasionOfIkoria creates a Invasion Of Ikoria
// {X}{G}{G} - BATTLE
func NewInvasionOfIkoria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Ikoria")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}