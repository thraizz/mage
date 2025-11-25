package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tetzin Gnome Champion", NewTetzinGnomeChampion)
}

// NewTetzinGnomeChampion creates a Tetzin Gnome Champion
// {U}{R}{W} - ARTIFACT
func NewTetzinGnomeChampion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tetzin Gnome Champion")
	card.ManaCost = "{U}{R}{W}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"GNOME"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
