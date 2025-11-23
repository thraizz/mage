package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Untaidake The Cloud Keeper", NewUntaidakeTheCloudKeeper)
}

// NewUntaidakeTheCloudKeeper creates a Untaidake The Cloud Keeper
//   - LAND
func NewUntaidakeTheCloudKeeper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Untaidake The Cloud Keeper")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
