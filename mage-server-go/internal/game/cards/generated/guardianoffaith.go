package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Guardian Of Faith", NewGuardianOfFaith)
}

// NewGuardianOfFaith creates a Guardian Of Faith
// {1}{W}{W} - CREATURE
// Flash, Vigilance
func NewGuardianOfFaith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Guardian Of Faith")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "KNIGHT"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	// card.AddAbility(ability2)
	return card, nil
}
