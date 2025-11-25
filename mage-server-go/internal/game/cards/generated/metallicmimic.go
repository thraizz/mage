package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Metallic Mimic", NewMetallicMimic)
}

// NewMetallicMimic creates a Metallic Mimic
// {2} - ARTIFACT CREATURE
func NewMetallicMimic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Metallic Mimic")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AsEntersBattlefieldAbility
	//   - Effect: ChooseCreatureTypeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
