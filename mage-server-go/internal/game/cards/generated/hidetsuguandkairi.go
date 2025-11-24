package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hidetsugu And Kairi", NewHidetsuguAndKairi)
}

// NewHidetsuguAndKairi creates a Hidetsugu And Kairi
// {2}{U}{U}{B} - CREATURE
// Flying
func NewHidetsuguAndKairi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidetsugu And Kairi")
	card.ManaCost = "{2}{U}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "DEMON", "DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}