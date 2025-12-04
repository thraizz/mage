package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warrior En Kor", NewWarriorEnKor)
}

// NewWarriorEnKor creates a Warrior En Kor
// {W}{W} - CREATURE
func NewWarriorEnKor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warrior En Kor")
	card.ManaCost = "{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOR", "WARRIOR", "KNIGHT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RedirectDamageFromSourceToTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{0}")
	// card.AddAbility(ability0)
	return card, nil
}
