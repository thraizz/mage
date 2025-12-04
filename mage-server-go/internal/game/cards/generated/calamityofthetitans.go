package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Calamity Of The Titans", NewCalamityOfTheTitans)
}

// NewCalamityOfTheTitans creates a Calamity Of The Titans
// {4}{C}{C} - SORCERY
func NewCalamityOfTheTitans(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Calamity Of The Titans")
	card.ManaCost = "{4}{C}{C}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
