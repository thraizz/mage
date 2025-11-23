package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cyclops Of One Eyed Pass", NewCyclopsOfOneEyedPass)
}

// NewCyclopsOfOneEyedPass creates a Cyclops Of One Eyed Pass
// {2}{R}{R} - CREATURE
func NewCyclopsOfOneEyedPass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cyclops Of One Eyed Pass")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CYCLOPS"}
	card.Power = "5"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
