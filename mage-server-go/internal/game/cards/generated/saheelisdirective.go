package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Saheelis Directive", NewSaheelisDirective)
}

// NewSaheelisDirective creates a Saheelis Directive
// {X}{R}{R}{R} - SORCERY
func NewSaheelisDirective(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Saheelis Directive")
	card.ManaCost = "{X}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
