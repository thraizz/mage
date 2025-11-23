package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Radha Heir To Keld", NewRadhaHeirToKeld)
}

// NewRadhaHeirToKeld creates a Radha Heir To Keld
// {R}{G} - CREATURE
func NewRadhaHeirToKeld(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Radha Heir To Keld")
	card.ManaCost = "{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	return card, nil
}