package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dollhouse Of Horrors", NewDollhouseOfHorrors)
}

// NewDollhouseOfHorrors creates a Dollhouse Of Horrors
// {5} - ARTIFACT
func NewDollhouseOfHorrors(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dollhouse Of Horrors")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
