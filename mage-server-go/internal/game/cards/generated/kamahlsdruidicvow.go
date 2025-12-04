package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kamahls Druidic Vow", NewKamahlsDruidicVow)
}

// NewKamahlsDruidicVow creates a Kamahls Druidic Vow
// {X}{G}{G} - SORCERY
func NewKamahlsDruidicVow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kamahls Druidic Vow")
	card.ManaCost = "{X}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
