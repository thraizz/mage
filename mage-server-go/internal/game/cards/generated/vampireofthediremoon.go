package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vampire Of The Dire Moon", NewVampireOfTheDireMoon)
}

// NewVampireOfTheDireMoon creates a Vampire Of The Dire Moon
// {B} - CREATURE
// Deathtouch, Lifelink
func NewVampireOfTheDireMoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vampire Of The Dire Moon")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
