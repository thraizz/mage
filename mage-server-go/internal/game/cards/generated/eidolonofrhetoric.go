package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eidolon Of Rhetoric", NewEidolonOfRhetoric)
}

// NewEidolonOfRhetoric creates a Eidolon Of Rhetoric
// {2}{W} - ENCHANTMENT CREATURE
func NewEidolonOfRhetoric(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eidolon Of Rhetoric")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
