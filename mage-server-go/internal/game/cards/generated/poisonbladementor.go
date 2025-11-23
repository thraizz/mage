package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Poison Blade Mentor", NewPoisonBladeMentor)
}

// NewPoisonBladeMentor creates a Poison Blade Mentor
// {1}{B} - CREATURE
// Deathtouch
func NewPoisonBladeMentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Poison Blade Mentor")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ASSASSIN"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect("DeathtouchAbility")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
