package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Epic Experiment", NewEpicExperiment)
}

// NewEpicExperiment creates a Epic Experiment
// {X}{U}{R} - SORCERY
func NewEpicExperiment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Epic Experiment")
	card.ManaCost = "{X}{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
