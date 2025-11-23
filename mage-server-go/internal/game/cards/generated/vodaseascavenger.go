package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Voda Sea Scavenger", NewVodaSeaScavenger)
}

// NewVodaSeaScavenger creates a Voda Sea Scavenger
// {2}{U} - CREATURE
func NewVodaSeaScavenger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Voda Sea Scavenger")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 DomainValue.REGULAR, 1, PutCards....)
	// card.AddAbility(ability0)
	return card, nil
}
