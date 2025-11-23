package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Doran The Siege Tower", NewDoranTheSiegeTower)
}

// NewDoranTheSiegeTower creates a Doran The Siege Tower
// {W}{B}{G} - CREATURE
func NewDoranTheSiegeTower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Doran The Siege Tower")
	card.ManaCost = "{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TREEFOLK", "SHAMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
