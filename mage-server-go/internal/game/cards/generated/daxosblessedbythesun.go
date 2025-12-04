package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daxos Blessed By The Sun", NewDaxosBlessedByTheSun)
}

// NewDaxosBlessedByTheSun creates a Daxos Blessed By The Sun
// {W}{W} - ENCHANTMENT CREATURE
func NewDaxosBlessedByTheSun(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daxos Blessed By The Sun")
	card.ManaCost = "{W}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"DEMIGOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
