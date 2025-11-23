package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Orah Skyclave Hierophant", NewOrahSkyclaveHierophant)
}

// NewOrahSkyclaveHierophant creates a Orah Skyclave Hierophant
// {2}{W}{B} - CREATURE
// Lifelink
func NewOrahSkyclaveHierophant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Orah Skyclave Hierophant")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOR", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
