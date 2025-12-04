package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tinybones Bauble Burglar", NewTinybonesBaubleBurglar)
}

// NewTinybonesBaubleBurglar creates a Tinybones Bauble Burglar
// {1}{B} - CREATURE
func NewTinybonesBaubleBurglar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tinybones Bauble Burglar")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SKELETON", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: DiscardEachPlayerEffect(                         StaticValue.get(1), false...)
	// card.AddAbility(ability0)
	return card, nil
}
