package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Radagast The Brown", NewRadagastTheBrown)
}

// NewRadagastTheBrown creates a Radagast The Brown
// {2}{G}{G} - CREATURE
func NewRadagastTheBrown(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Radagast The Brown")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                         PermanentEnteringBattlefi...)
	// card.AddAbility(ability0)
	return card, nil
}