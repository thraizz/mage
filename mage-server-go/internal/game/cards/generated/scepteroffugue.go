package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scepter Of Fugue", NewScepterOfFugue)
}

// NewScepterOfFugue creates a Scepter Of Fugue
// {B}{B} - ARTIFACT
func NewScepterOfFugue(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scepter Of Fugue")
	card.ManaCost = "{B}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
