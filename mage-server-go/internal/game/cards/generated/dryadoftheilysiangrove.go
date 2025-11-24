package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dryad Of The Ilysian Grove", NewDryadOfTheIlysianGrove)
}

// NewDryadOfTheIlysianGrove creates a Dryad Of The Ilysian Grove
// {2}{G} - ENCHANTMENT CREATURE
func NewDryadOfTheIlysianGrove(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dryad Of The Ilysian Grove")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"NYMPH", "DRYAD"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}