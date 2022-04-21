package commoncartridge

import (
	"io/fs"

	"github.com/commonsyllabi/viewer/pkg/commoncartridge/types"
)

type Cartridge interface {
	// MarshalJSON returns a serialized JSON representation
	MarshalJSON() ([]byte, error)

	//-- ParseManifest finds the imsmanifest.xml in the ZipReader and marshals it into a struct
	ParseManifest() (types.Manifest, error)

	// Dump returns all the contents of the cartridge as a string
	Dump() []string

	// Title returns the title of the loaded cartridge
	Title() string

	// Metadata returns the metadata fields of the cartridge in a structured fashion
	Metadata() (string, error)

	// Items returns an array of structs which include the Item, the Resources and the children Items it might have
	Items() ([]FullItem, error)

	// Resources returns an array of structs which include the resource and, if found, the item in which the resource appears
	Resources() ([]FullResource, error)
	Weblinks() ([]types.WebLink, error)
	Assignments() ([]types.Assignment, error)
	LTIs() ([]types.CartridgeBasicltiLink, error)
	QTIs() ([]types.Questestinterop, error)
	Topics() ([]types.Topic, error)

	Find(string) (interface{}, error)
	FindFile(string) (fs.File, error)

	// Tobi's weird ideas
	// Webcontents() ([]os.File, error)
	// GetWebcontent(string) (os.File, error)
}
