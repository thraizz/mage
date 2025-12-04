package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nightsky Mimic", NewNightskyMimic)
}

// NewNightskyMimic creates a Nightsky Mimic
// {1}{W/B} - CREATURE
func NewNightskyMimic(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nightsky Mimic")
	card.ManaCost = "{1}{W/B}"
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
