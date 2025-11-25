package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Kayas Guile", NewKayasGuile)
}

// NewKayasGuile creates a Kayas Guile
// {1}{W}{B} - INSTANT
func NewKayasGuile(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kayas Guile")
	card.ManaCost = "{1}{W}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ExileGraveyardAllPlayersEffect(                 StaticFilters.FILTER_CARD_CARDS, ...)
	// card.AddAbility(ability0)
	return card, nil
}
