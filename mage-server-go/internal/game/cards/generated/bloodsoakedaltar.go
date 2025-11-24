package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bloodsoaked Altar", NewBloodsoakedAltar)
}

// NewBloodsoakedAltar creates a Bloodsoaked Altar
// {4}{B}{B} - ARTIFACT
func NewBloodsoakedAltar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bloodsoaked Altar")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}