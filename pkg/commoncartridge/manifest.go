package commoncartridge

import (
	"fmt"
)

func (m *Manifest) ResolveItems() {
	if m.Organizations.Organization.Item.Identifier != "" {
		m.traverseItemModules(m.Organizations.Organization.Item.Item)
	}
}

func (m *Manifest) FindItem(id string) (Item, error) {
	// fmt.Printf("starting search for %s\n", id)
	var item Item
	for i := range m.Organizations.Organization.Item.Item {
		item, err := m.findItem(m.Organizations.Organization.Item.Item[i].Item, id)

		if err != nil {
			return item, nil
		}

		if item.Identifierref == id { //-- found
			return item, nil
		}
	}

	return item, nil
}

func (m *Manifest) findItem(items []Item, id string) (Item, error) {
	// fmt.Printf("looking for %s in %d items\n", id, len(items))
	var item Item
	var err error
	for i := range items {
		if items[i].Identifierref == id {
			return items[i], nil
		}

		if len(items[i].Item) > 0 {
			item, err = m.findItem(items[i].Item, id)
		}
	}

	return item, err
}

func (m *Manifest) traverseItemModules(itemModules []Item) {
	for i := range itemModules {
		m.traverseItems(itemModules[i].Item)
	}
}

func (m *Manifest) traverseItems(items []Item) {
	// fmt.Printf("- %d items traversed\n", len(items))
	for i := range items {

		// fmt.Printf("- - idref %v\n", items[i].Identifierref)

		err := m.resolveItem(items[i])

		if err != nil {
			fmt.Println(err)
		}

		if len(items[i].Item) > 0 {
			m.traverseItems(items[i].Item)
		}

	}
}

//-- resolveItem resolves its relationship to a resource (since items are just folders with stuff inside)
func (m *Manifest) resolveItem(item Item) error {
	if item.Identifierref == "" {
		return fmt.Errorf("no identifierref for item %s", item.Identifier)
	}

	for _, resource := range m.Resources.Resource {
		if resource.Identifier == item.Identifierref {
			// fmt.Printf("- - matched resource id %s\n", resource.Identifier)

			switch resource.Type {
			case "imsdt_xmlv1p1", "imsdt_xmlv1p2", "imsdt_xmlv1p3":
				//-- topic
				// fmt.Printf("found topic %v", resource.File)

			case "webcontent":
				//-- webcontent
				// fmt.Printf("found webcontent %v", resource.File)

			case "imswl_xmlv1p1":
				//-- weblink
				// fmt.Printf("found weblink %v", resource.File)

			case "assignment_xmlv1p0":
				//-- assignment
				// fmt.Printf("found assignment %v", resource.File)

			case "imsqti_xmlv1p2/imscc_xmlv1p1/assessment", "imsqti_xmlv1p2/imscc_xmlv1p2/assessment",
				"imsqti_xmlv1p2/imscc_xmlv1p3/assessment":
				//-- qti
				// fmt.Printf("found question bank %v", resource.File)

			default:
				return fmt.Errorf("[resolveItem] No matching type found: %s", resource.Type)
			}
		}
	}
	return nil
}

func (m *Manifest) PrettyPrint() {
	fmt.Printf("Cartridge: %v\n", m.Metadata.Lom.General.Title.String.Text)
	fmt.Printf("Modules: %v\n", len(m.Organizations.Organization.Item.Item))

	for _, i := range m.Organizations.Organization.Item.Item {
		fmt.Printf("- %d items\n", len(i.Item))
		for _, v := range i.Item {
			fmt.Printf("- - id %v\n", v.Identifierref)
		}
	}

	fmt.Printf("Resources (%d):\n", len(m.Resources.Resource))
	for _, r := range m.Resources.Resource {
		fmt.Printf(" - - type %v - id %v\n", r.Type, r.Identifier)
	}
}
