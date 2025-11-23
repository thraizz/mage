package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Warmaker Gunship", NewWarmakerGunship)
}

// NewWarmakerGunship creates a Warmaker Gunship
// {2}{R} - ARTIFACT
// Flying
func NewWarmakerGunship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Warmaker Gunship")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"SPACECRAFT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(ArtifactYouControlCount.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
