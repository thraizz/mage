package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Skarrgan Skybreaker", NewSkarrganSkybreaker)
}

// NewSkarrganSkybreaker creates a Skarrgan Skybreaker
// {4}{R}{R}{G} - CREATURE
func NewSkarrganSkybreaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skarrgan Skybreaker")
	card.ManaCost = "{4}{R}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(SourcePermanentPowerValue.NOT_NEGATIVE)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
