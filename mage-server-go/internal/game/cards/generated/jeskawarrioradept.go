package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jeska Warrior Adept", NewJeskaWarriorAdept)
}

// NewJeskaWarriorAdept creates a Jeska Warrior Adept
// {2}{R}{R} - CREATURE
// FirstStrike, Haste
func NewJeskaWarriorAdept(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jeska Warrior Adept")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "BARBARIAN", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDamageEffect(1)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
