package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sorin Imperious Bloodlord", NewSorinImperiousBloodlord)
}

// NewSorinImperiousBloodlord creates a Sorin Imperious Bloodlord
// {2}{B} - PLANESWALKER
func NewSorinImperiousBloodlord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sorin Imperious Bloodlord")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SORIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: SorinImperiousBloodlordEffect()
	// card.AddAbility(ability0)
	return card, nil
}
