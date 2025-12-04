package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Moritte Of The Frost", NewMoritteOfTheFrost)
}

// NewMoritteOfTheFrost creates a Moritte Of The Frost
// {2}{G}{U}{U} - CREATURE
func NewMoritteOfTheFrost(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Moritte Of The Frost")
	card.ManaCost = "{2}{G}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Supertypes = []string{"LEGENDARY", "SNOW"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(filter, new MoritteOfTheFrostCopyApplier())
	// card.AddAbility(ability0)
	return card, nil
}
