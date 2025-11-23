package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Terraformer", NewTerraformer)
}

// NewTerraformer creates a Terraformer
// {2}{U} - CREATURE
func NewTerraformer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Terraformer")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "R")
	card.AddAbility(ability2)
	ability3 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability3)
	ability4 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability4)
	return card, nil
}
