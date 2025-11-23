package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ivorytusk Fortress", NewIvorytuskFortress)
}

// NewIvorytuskFortress creates a Ivorytusk Fortress
// {2}{W}{B}{G} - CREATURE
func NewIvorytuskFortress(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ivorytusk Fortress")
	card.ManaCost = "{2}{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT"}
	card.Power = "5"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
