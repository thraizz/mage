package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Throne Of The God Pharaoh", NewThroneOfTheGodPharaoh)
}

// NewThroneOfTheGodPharaoh creates a Throne Of The God Pharaoh
// {2} - ARTIFACT
func NewThroneOfTheGodPharaoh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Throne Of The God Pharaoh")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
