package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Planetarium Of Wan Shi Tong", NewPlanetariumOfWanShiTong)
}

// NewPlanetariumOfWanShiTong creates a Planetarium Of Wan Shi Tong
// {6} - ARTIFACT
func NewPlanetariumOfWanShiTong(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Planetarium Of Wan Shi Tong")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddEffect(abilities.NewScryEffect(1)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
