package shared

type RuleSource string

const (
	RuleSourceSR5Core RuleSource = "SR5:Core"
	RuleSourceSR5CF   RuleSource = "SR5:ChromeFlesh"
	RuleSourceSR5RG   RuleSource = "SR5:RunAndGun"
	RuleSourceSR5SG   RuleSource = "SR5:StreetGrimoire"
	RuleSourceSR5HT   RuleSource = "SR5:HardTargets"
	RuleSourceSR5R5   RuleSource = "SR5:Rigger5"
	RuleSourceSR5DT   RuleSource = "SR5:DataTrails"
	RuleSourceSR5CA   RuleSource = "SR5:CuttingAces"
	RuleSourceSR5SASS RuleSource = "SR5:SailAwaySweetSister"
	RuleSourceSR5GH3  RuleSource = "SR5:GunH(e)aven3"
	RuleSourceSR5BB   RuleSource = "SR5:BulletsAndBandages"

	DispositionFriendly   = "Friendly"
	DispositionNeutral    = "Neutral"
	DispositionAggressive = "Aggressive"
)
