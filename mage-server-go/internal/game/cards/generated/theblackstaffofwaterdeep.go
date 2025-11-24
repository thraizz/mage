package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Blackstaff Of Waterdeep", NewTheBlackstaffOfWaterdeep)
}

// NewTheBlackstaffOfWaterdeep creates a The Blackstaff Of Waterdeep
// {U} - ARTIFACT
func NewTheBlackstaffOfWaterdeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Blackstaff Of Waterdeep")
	card.ManaCost = "{U}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
