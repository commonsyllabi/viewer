package viewer

import "commonsyllabi/viewer/specs"

type Cartridge interface {
	Load(string) (Cartridge, error)
	AsObject() ([]byte, error)
	ParseManifest() (specs.Manifest, error)
	Dump() []string
}
