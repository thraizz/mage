package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Guardian Of Vitu Ghazi", NewGuardianOfVituGhazi)
}

// NewGuardianOfVituGhazi creates a Guardian Of Vitu Ghazi
// {6}{G}{W} - CREATURE
// Vigilance
func NewGuardianOfVituGhazi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Guardian Of Vitu Ghazi")
	card.ManaCost = "{6}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "4"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
