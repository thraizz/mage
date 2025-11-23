package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Henrika Infernal Seer", NewHenrikaInfernalSeer)
}

// NewHenrikaInfernalSeer creates a Henrika Infernal Seer
//   - CREATURE
//
// Flying, Deathtouch, Lifelink
func NewHenrikaInfernalSeer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Henrika Infernal Seer")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability2)
	return card, nil
}
