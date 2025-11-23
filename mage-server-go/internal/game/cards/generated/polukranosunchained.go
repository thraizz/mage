package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Polukranos Unchained", NewPolukranosUnchained)
}

// NewPolukranosUnchained creates a Polukranos Unchained
// {2}{B}{G} - CREATURE
func NewPolukranosUnchained(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Polukranos Unchained")
	card.ManaCost = "{2}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "HYDRA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
