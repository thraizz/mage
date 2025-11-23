package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Airdrop Aeronauts", NewAirdropAeronauts)
}

// NewAirdropAeronauts creates a Airdrop Aeronauts
// {3}{W}{W} - CREATURE
// Flying
func NewAirdropAeronauts(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Airdrop Aeronauts")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "SCOUT"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(5)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
