package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boneyard Wurm", NewBoneyardWurm)
}

// NewBoneyardWurm creates a Boneyard Wurm
// {1}{G} - CREATURE
func NewBoneyardWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boneyard Wurm")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
