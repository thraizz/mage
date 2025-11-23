package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Illicit Auction", NewIllicitAuction)
}

// NewIllicitAuction creates a Illicit Auction
// {3}{R}{R} - SORCERY
// FirstStrike, DoubleStrike, Trample
func NewIllicitAuction(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Illicit Auction")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
