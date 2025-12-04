package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gallifrey Falls No More", NewGallifreyFallsNoMore)
}

// NewGallifreyFallsNoMore creates a Gallifrey Falls No More
// {4}{R}{R} - INSTANT
func NewGallifreyFallsNoMore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gallifrey Falls No More")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(4, new FilterCreaturePermanent())
	//   - PhaseOutTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
