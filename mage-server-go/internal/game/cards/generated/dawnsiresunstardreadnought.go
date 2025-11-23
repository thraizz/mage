package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dawnsire Sunstar Dreadnought", NewDawnsireSunstarDreadnought)
}

// NewDawnsireSunstarDreadnought creates a Dawnsire Sunstar Dreadnought
// {5} - ARTIFACT
// Flying
func NewDawnsireSunstarDreadnought(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dawnsire Sunstar Dreadnought")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPACECRAFT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(100)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
