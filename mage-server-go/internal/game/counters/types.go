package counters

// CounterType represents a type of counter.
// Mirrors Java CounterType enum (simplified - using string constants for now).
type CounterType string

const (
	// Common counter types
	CounterTypeLoyalty   CounterType = "loyalty"
	CounterTypePoison    CounterType = "poison"
	CounterTypeEnergy    CounterType = "energy"
	CounterTypeExperience CounterType = "experience"

	// Power/toughness boost counters
	CounterTypeP1P1 CounterType = "+1/+1"
	CounterTypeM1M1 CounterType = "-1/-1"
	CounterTypeP2P2 CounterType = "+2/+2"
	CounterTypeM2M2 CounterType = "-2/-2"
	CounterTypeP1P0 CounterType = "+1/+0"
	CounterTypeP0P1 CounterType = "+0/+1"
	CounterTypeM1M0 CounterType = "-1/+0"
	CounterTypeM0M1 CounterType = "+0/-1"

	// Ability counters (simplified - full implementation would reference abilities)
	CounterTypeFlying      CounterType = "flying"
	CounterTypeFirstStrike  CounterType = "first strike"
	CounterTypeDoubleStrike CounterType = "double strike"
	CounterTypeDeathtouch   CounterType = "deathtouch"
	CounterTypeLifelink     CounterType = "lifelink"
	CounterTypeTrample      CounterType = "trample"
	CounterTypeVigilance    CounterType = "vigilance"
	CounterTypeHaste        CounterType = "haste"
	CounterTypeHexproof     CounterType = "hexproof"
	CounterTypeIndestructible CounterType = "indestructible"
	CounterTypeReach        CounterType = "reach"
	CounterTypeMenace        CounterType = "menace"
	CounterTypeShadow        CounterType = "shadow"

	// Other counter types (expanded list from Java)
	CounterTypeAge        CounterType = "age"
	CounterTypeArrow      CounterType = "arrow"
	CounterTypeBlaze      CounterType = "blaze"
	CounterTypeBlood      CounterType = "blood"
	CounterTypeBounty      CounterType = "bounty"
	CounterTypeBrick       CounterType = "brick"
	CounterTypeCharge      CounterType = "charge"
	CounterTypeCoin        CounterType = "coin"
	CounterTypeCorpse      CounterType = "corpse"
	CounterTypeCredit      CounterType = "credit"
	CounterTypeCrystal     CounterType = "crystal"
	CounterTypeDeath       CounterType = "death"
	CounterTypeDefense     CounterType = "defense"
	CounterTypeDepletion   CounterType = "depletion"
	CounterTypeDoom        CounterType = "doom"
	CounterTypeDream       CounterType = "dream"
	CounterTypeEgg         CounterType = "egg"
	CounterTypeElixir      CounterType = "elixir"
	CounterTypeFate        CounterType = "fate"
	CounterTypeFeather     CounterType = "feather"
	CounterTypeFire        CounterType = "fire"
	CounterTypeFlame       CounterType = "flame"
	CounterTypeFungus      CounterType = "fungus"
	CounterTypeGem         CounterType = "gem"
	CounterTypeGold        CounterType = "gold"
	CounterTypeGrowth      CounterType = "growth"
	CounterTypeHour        CounterType = "hour"
	CounterTypeHourglass   CounterType = "hourglass"
	CounterTypeIce         CounterType = "ice"
	CounterTypeInfection   CounterType = "infection"
	CounterTypeInfluence   CounterType = "influence"
	CounterTypeKi          CounterType = "ki"
	CounterTypeKnowledge   CounterType = "knowledge"
	CounterTypeLevel       CounterType = "level"
	CounterTypeLore        CounterType = "lore"
	CounterTypeLuck        CounterType = "luck"
	CounterTypeMine        CounterType = "mine"
	CounterTypeMining      CounterType = "mining"
	CounterTypeMusic       CounterType = "music"
	CounterTypeMuster      CounterType = "muster"
	CounterTypeNight       CounterType = "night"
	CounterTypeOil          CounterType = "oil"
	CounterTypeOmen        CounterType = "omen"
	CounterTypeOre         CounterType = "ore"
	CounterTypePage        CounterType = "page"
	CounterTypePain        CounterType = "pain"
	CounterTypePetrification CounterType = "petrification"
	CounterTypePhylactery  CounterType = "phylactery"
	CounterTypePhyresis    CounterType = "phyresis"
	CounterTypePlague      CounterType = "plague"
	CounterTypePlot        CounterType = "plot"
	CounterTypePoint       CounterType = "point"
	CounterTypePressure    CounterType = "pressure"
	CounterTypePupa        CounterType = "pupa"
	CounterTypeQuest       CounterType = "quest"
	CounterTypeRad         CounterType = "rad"
	CounterTypeRally       CounterType = "rally"
	CounterTypeRitual      CounterType = "ritual"
	CounterTypeRope        CounterType = "rope"
	CounterTypeRust        CounterType = "rust"
	CounterTypeSilver      CounterType = "silver"
	CounterTypeScream      CounterType = "scream"
	CounterTypeShell       CounterType = "shell"
	CounterTypeShield      CounterType = "shield"
	CounterTypeSleep       CounterType = "sleep"
	CounterTypeSlime       CounterType = "slime"
	CounterTypeSlumber     CounterType = "slumber"
	CounterTypeSoot        CounterType = "soot"
	CounterTypeSoul        CounterType = "soul"
	CounterTypeSpore       CounterType = "spore"
	CounterTypeStash       CounterType = "stash"
	CounterTypeStorage     CounterType = "storage"
	CounterTypeStory       CounterType = "story"
	CounterTypeStrife      CounterType = "strife"
	CounterTypeStudy       CounterType = "study"
	CounterTypeStun        CounterType = "stun"
	CounterTypeSupply      CounterType = "supply"
	CounterTypeSuspect     CounterType = "suspect"
	CounterTypeTask        CounterType = "task"
	CounterTypeTheft       CounterType = "theft"
	CounterTypeTide        CounterType = "tide"
	CounterTypeTime        CounterType = "time"
	CounterTypeTower       CounterType = "tower"
	CounterTypeTraining    CounterType = "training"
	CounterTypeTrap        CounterType = "trap"
	CounterTypeTreasure    CounterType = "treasure"
	CounterTypeUnity       CounterType = "unity"
	CounterTypeUnlock      CounterType = "unlock"
	CounterTypeValor       CounterType = "valor"
	CounterTypeVelocity    CounterType = "velocity"
	CounterTypeVerse       CounterType = "verse"
	CounterTypeVitality    CounterType = "vitality"
	CounterTypeVoid        CounterType = "void"
	CounterTypeVortex      CounterType = "vortex"
	CounterTypeVow         CounterType = "vow"
	CounterTypeVoyage      CounterType = "voyage"
	CounterTypeWage        CounterType = "wage"
	CounterTypeWinch       CounterType = "winch"
	CounterTypeWind        CounterType = "wind"
	CounterTypeWish        CounterType = "wish"
)

// String returns the string representation of the counter type.
func (ct CounterType) String() string {
	return string(ct)
}

// CreateInstance creates a counter instance of this type with the given amount.
func (ct CounterType) CreateInstance(amount int) *Counter {
	if amount <= 0 {
		amount = 1
	}

	// Handle boost counters
	switch ct {
	case CounterTypeP1P1:
		return NewBoostCounter(1, 1, amount).Counter
	case CounterTypeM1M1:
		return NewBoostCounter(-1, -1, amount).Counter
	case CounterTypeP2P2:
		return NewBoostCounter(2, 2, amount).Counter
	case CounterTypeM2M2:
		return NewBoostCounter(-2, -2, amount).Counter
	case CounterTypeP1P0:
		return NewBoostCounter(1, 0, amount).Counter
	case CounterTypeP0P1:
		return NewBoostCounter(0, 1, amount).Counter
	case CounterTypeM1M0:
		return NewBoostCounter(-1, 0, amount).Counter
	case CounterTypeM0M1:
		return NewBoostCounter(0, -1, amount).Counter
	default:
		return NewCounter(string(ct), amount)
	}
}

// CreateBoostCounter creates a boost counter with the given power/toughness deltas.
func CreateBoostCounter(power, toughness, amount int) *BoostCounter {
	return NewBoostCounter(power, toughness, amount)
}
