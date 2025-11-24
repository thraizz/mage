package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Buzzing Whack A Doodle", NewBuzzingWhackADoodle)
}

// NewBuzzingWhackADoodle creates a Buzzing Whack A Doodle
// {4} - ARTIFACT
func NewBuzzingWhackADoodle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Buzzing Whack A Doodle")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
