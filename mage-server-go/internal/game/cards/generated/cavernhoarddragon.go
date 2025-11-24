package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Cavern Hoard Dragon", NewCavernHoardDragon)
}

// NewCavernHoardDragon creates a Cavern Hoard Dragon
// {7}{R}{R} - CREATURE
// Flying, Trample, Haste
func NewCavernHoardDragon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cavern Hoard Dragon")
	card.ManaCost = "{7}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	token3_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: CreateTokenEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
