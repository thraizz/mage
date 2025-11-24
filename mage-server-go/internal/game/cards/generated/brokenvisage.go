package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Broken Visage", NewBrokenVisage)
}

// NewBrokenVisage creates a Broken Visage
// {4}{B} - INSTANT
func NewBrokenVisage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Broken Visage")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("Sacrifice the token at the beginning of the next ...)
	// card.AddAbility(ability0)
	return card, nil
}
