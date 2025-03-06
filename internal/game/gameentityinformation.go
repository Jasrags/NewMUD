package game

const (
	SexMale      = "Male"
	SexFemale    = "Female"
	SexNonBinary = "Non-Binary"
)

type GameEntityInformation struct {
	ID                 string    `yaml:"id"`
	Name               string    `yaml:"name"`
	Title              string    `yaml:"title"`
	Description        string    `yaml:"description"`
	LongDescription    string    `yaml:"long_description"`
	MetatypeID         string    `yaml:"metatype_id"`
	Metatype           *Metatype `yaml:"-"`
	Age                int       `yaml:"age"`
	Sex                string    `yaml:"sex"`
	Height             int       `yaml:"height"`
	Weight             int       `yaml:"weight"`
	StreetCred         int       `yaml:"street_cred"`
	Notoriety          int       `yaml:"notoriety"`
	PublicAwareness    int       `yaml:"public_awareness"`
	ProfessionalRating int       `yaml:"professional_rating,omitempty"`
	GeneralDisposition string    `yaml:"general_disposition,omitempty"`
	Tags               []string  `yaml:"tags"`
}

func (g *GameEntityInformation) GetGeneralDisposition() string {
	if g.GeneralDisposition == "" {
		return DispositionNeutral
	}

	return g.GeneralDisposition
}

// Validate the game entity information
// TODO: Implement validation
func (g *GameEntityInformation) Validate() error {
	return nil
}
