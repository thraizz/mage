package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goblin Surveyor", NewGoblinSurveyor)
}

// NewGoblinSurveyor creates a Goblin Surveyor
// {2}{R} - CREATURE
// Trample
func NewGoblinSurveyor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goblin Surveyor")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SCOUT"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
