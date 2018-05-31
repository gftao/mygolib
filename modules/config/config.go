package config

type IConfig interface {
	SetSection(section string)
	SetOption(name, value string)
	Int(option string) (result int, found bool)
	IntDefault(option string, dfault int) int
	Bool(option string) (result, found bool)
	BoolDefault(option string, dfault bool) bool
	String(option string) (result string, found bool)
	StringDefault(option, dfault string) string
	HasSection(section string) bool
	Options(prefix string) []string
}

var instance IConfig

func SetSection(section string) {
	instance.SetSection(section)
}
func SetOption(name, value string) {
	instance.SetOption(name, value)
}
func Int(option string) (result int, found bool) {
	return instance.Int(option)
}
func IntDefault(option string, dfault int) int {
	return instance.IntDefault(option, dfault)
}
func Bool(option string) (result, found bool) {
	return instance.Bool(option)
}
func BoolDefault(option string, dfault bool) bool {
	return instance.BoolDefault(option, dfault)
}
func String(option string) (result string, found bool) {
	return instance.String(option)
}
func StringDefault(option, dfault string) string {
	return instance.StringDefault(option, dfault)
}
func HasSection(section string) bool {
	return instance.HasSection(section)
}
func Options(prefix string) []string {
	return instance.Options(prefix)
}
