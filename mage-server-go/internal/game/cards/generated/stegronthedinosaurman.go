package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stegron The Dinosaur Man", NewStegronTheDinosaurMan)
}

// NewStegronTheDinosaurMan creates a Stegron The Dinosaur Man
// {4}{R} - CREATURE
// Menace
func NewStegronTheDinosaurMan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stegron The Dinosaur Man")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "VILLAIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(3, 1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
