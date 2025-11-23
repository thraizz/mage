package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elsewhere Flask", NewElsewhereFlask)
}

// NewElsewhereFlask creates a Elsewhere Flask
// {2} - ARTIFACT
func NewElsewhereFlask(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elsewhere Flask")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
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
	ability5, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability5)
	return card, nil
}
