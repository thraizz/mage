package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vorinclex Voice Of Hunger", NewVorinclexVoiceOfHunger)
}

// NewVorinclexVoiceOfHunger creates a Vorinclex Voice Of Hunger
// {6}{G}{G} - CREATURE
// Trample
func NewVorinclexVoiceOfHunger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vorinclex Voice Of Hunger")
	card.ManaCost = "{6}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "PRAETOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
