package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashnods Cylix", NewAshnodsCylix)
}

// NewAshnodsCylix creates a Ashnods Cylix
// {2} - ARTIFACT
func NewAshnodsCylix(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashnods Cylix")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
