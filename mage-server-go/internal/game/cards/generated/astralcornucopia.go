package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Astral Cornucopia", NewAstralCornucopia)
}

// NewAstralCornucopia creates a Astral Cornucopia
// {X}{X}{X} - ARTIFACT
func NewAstralCornucopia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Astral Cornucopia")
	card.ManaCost = "{X}{X}{X}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
