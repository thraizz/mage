package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lonis Genetics Expert", NewLonisGeneticsExpert)
}

// NewLonisGeneticsExpert creates a Lonis Genetics Expert
// {1}{G/U}{G/U} - CREATURE
func NewLonisGeneticsExpert(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lonis Genetics Expert")
	card.ManaCost = "{1}{G/U}{G/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "ELF", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
