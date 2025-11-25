package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gix Yawgmoth Praetor", NewGixYawgmothPraetor)
}

// NewGixYawgmothPraetor creates a Gix Yawgmoth Praetor
// {1}{B}{B} - CREATURE
func NewGixYawgmothPraetor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gix Yawgmoth Praetor")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "PRAETOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - GixYawgmothPraetorExileEffect()
	// card.AddAbility(ability0)
	return card, nil
}
