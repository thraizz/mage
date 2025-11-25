package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaervek The Punisher", NewKaervekThePunisher)
}

// NewKaervekThePunisher creates a Kaervek The Punisher
// {1}{B}{B} - CREATURE
func NewKaervekThePunisher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaervek The Punisher")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARLOCK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: CommittedCrimeTriggeredAbility
	//   - Effect: KaervekThePunisherEffect()
	// card.AddAbility(ability0)
	return card, nil
}
