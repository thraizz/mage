package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eerie Ultimatum", NewEerieUltimatum)
}

// NewEerieUltimatum creates a Eerie Ultimatum
// {W}{W}{B}{B}{B}{G}{G} - SORCERY
func NewEerieUltimatum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eerie Ultimatum")
	card.ManaCost = "{W}{W}{B}{B}{B}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}