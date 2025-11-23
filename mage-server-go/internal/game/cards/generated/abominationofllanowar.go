package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abomination Of Llanowar", NewAbominationOfLlanowar)
}

// NewAbominationOfLlanowar creates a Abomination Of Llanowar
// {1}{B}{G} - CREATURE
// Vigilance, Menace
func NewAbominationOfLlanowar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abomination Of Llanowar")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability1)
	return card, nil
}
