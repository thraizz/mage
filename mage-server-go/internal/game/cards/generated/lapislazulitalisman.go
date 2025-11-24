package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lapis Lazuli Talisman", NewLapisLazuliTalisman)
}

// NewLapisLazuliTalisman creates a Lapis Lazuli Talisman
// {2} - ARTIFACT
func NewLapisLazuliTalisman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lapis Lazuli Talisman")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new UntapTargetEffect(), new ManaCostsImpl<>("{3}"...)
	// card.AddAbility(ability0)
	return card, nil
}