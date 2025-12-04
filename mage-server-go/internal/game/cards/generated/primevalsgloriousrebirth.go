package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Primevals Glorious Rebirth", NewPrimevalsGloriousRebirth)
}

// NewPrimevalsGloriousRebirth creates a Primevals Glorious Rebirth
// {5}{W}{B} - SORCERY
func NewPrimevalsGloriousRebirth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Primevals Glorious Rebirth")
	card.ManaCost = "{5}{W}{B}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
