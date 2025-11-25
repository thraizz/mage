package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Balduvian Atrocity", NewBalduvianAtrocity)
}

// NewBalduvianAtrocity creates a Balduvian Atrocity
// {2}{B} - CREATURE
func NewBalduvianAtrocity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Balduvian Atrocity")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "BERSERKER"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: BalduvianAtrocityEffect()
	// card.AddAbility(ability0)
	ability1 := abilities.NewKickerAbility(card.ID, "{R}")
	card.AddAbility(ability1)
	return card, nil
}
