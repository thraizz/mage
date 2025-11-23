package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scaled Wurm", NewScaledWurm)
}

// NewScaledWurm creates a Scaled Wurm
// {7}{G} - CREATURE
func NewScaledWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scaled Wurm")
	card.ManaCost = "{7}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
