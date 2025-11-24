package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Genesis Hydra", NewGenesisHydra)
}

// NewGenesisHydra creates a Genesis Hydra
// {X}{G}{G} - CREATURE
func NewGenesisHydra(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Genesis Hydra")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PLANT", "HYDRA"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
