package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rooting Kavu", NewRootingKavu)
}

// NewRootingKavu creates a Rooting Kavu
// {2}{G}{G} - CREATURE
func NewRootingKavu(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rooting Kavu")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KAVU"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new RootingKavuEffect(), new ExileSourceFromGraveC...)
	// card.AddAbility(ability0)
	return card, nil
}
