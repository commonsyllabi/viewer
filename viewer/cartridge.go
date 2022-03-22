package viewer

type Cartridge interface {
	Load(string) (Cartridge, error)
	AsObject() []byte
	Dump()
}
