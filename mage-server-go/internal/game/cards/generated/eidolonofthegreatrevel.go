package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eidolon Of The Great Revel", NewEidolonOfTheGreatRevel)
}

// NewEidolonOfTheGreatRevel creates a Eidolon Of The Great Revel
// {R}{R} - ENCHANTMENT CREATURE
func NewEidolonOfTheGreatRevel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eidolon Of The Great Revel")
	card.ManaCost = "{R}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
