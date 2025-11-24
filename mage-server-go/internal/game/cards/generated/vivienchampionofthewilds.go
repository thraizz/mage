package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vivien Champion Of The Wilds", NewVivienChampionOfTheWilds)
}

// NewVivienChampionOfTheWilds creates a Vivien Champion Of The Wilds
// {2}{G} - PLANESWALKER
func NewVivienChampionOfTheWilds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vivien Champion Of The Wilds")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VIVIEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}