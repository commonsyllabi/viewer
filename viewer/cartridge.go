package viewer

type Cartridge interface {
	Load(string) (interface{}, error)
	AsObject() []byte
}
