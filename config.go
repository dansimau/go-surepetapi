package surepetapi

type Config struct {
	AuthToken    string `yaml:"authToken"`
	EmailAddress string `yaml:"emailAddress"`
	Password     string `yaml:"password"`
}
