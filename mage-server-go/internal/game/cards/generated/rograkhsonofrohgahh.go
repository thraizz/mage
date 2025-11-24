package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rograkh Son Of Rohgahh", NewRograkhSonOfRohgahh)
}

// NewRograkhSonOfRohgahh creates a Rograkh Son Of Rohgahh
// {0} - CREATURE
// FirstStrike, Trample
func NewRograkhSonOfRohgahh(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rograkh Son Of Rohgahh")
	card.ManaCost = "{0}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOBOLD", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}