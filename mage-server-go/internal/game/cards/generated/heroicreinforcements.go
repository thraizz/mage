package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Heroic Reinforcements", NewHeroicReinforcements)
}

// NewHeroicReinforcements creates a Heroic Reinforcements
// {2}{R}{W} - SORCERY
func NewHeroicReinforcements(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heroic Reinforcements")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
