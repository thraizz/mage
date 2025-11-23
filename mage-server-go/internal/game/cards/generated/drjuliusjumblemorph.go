package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dr Julius Jumblemorph", NewDrJuliusJumblemorph)
}

// NewDrJuliusJumblemorph creates a Dr Julius Jumblemorph
// {2}{G}{W} - CREATURE
func NewDrJuliusJumblemorph(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dr Julius Jumblemorph")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
