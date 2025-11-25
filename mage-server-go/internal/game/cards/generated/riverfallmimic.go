package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Riverfall Mimic", NewRiverfallMimic)
}

// NewRiverfallMimic creates a Riverfall Mimic
// {1}{U/R} - CREATURE
func NewRiverfallMimic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Riverfall Mimic")
	card.ManaCost = "{1}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastControllerTriggeredAbility
	//   - Effect: SetBasePowerToughnessSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
