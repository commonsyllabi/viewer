package commoncartridge

type Cartridge interface {
	// AsObject returns a serialized JSON representation - TODO: check how to make the CC implement Marshal
	AsObject() ([]byte, error)

	//-- ParseManifest finds the imsmanifest.xml in the ZipReader and marshals it into a struct
	ParseManifest() (Manifest, error)

	// Dump returns all the contents of the cartridge as a string
	Dump() []string

	// Title returns the title of the loaded cartridge
	Title() string

	// Metadata returns the metadata fields of the cartridge in a structured fashion
	Metadata() (string, error)

	// Resources returns an array of structs which include the resource and, if found, the item in which the resource appears
	Resources() ([]FullResource, error)
	Weblinks() ([]WebLink, error)
	Assignments() ([]Assignment, error)
	LTIs() ([]CartridgeBasicltiLink, error)
	QTIs() ([]Questestinterop, error)
	Topics() ([]Topic, error)

	Find(string) (Resource, error)
	FindFile(string) ([]byte, error)

	// Tobi's weird ideas
	// GetWeblink(string) (Weblink, error)
	// Webcontents() ([]os.File, error)
	// GetWebcontent(string) (os.File, error)
}
