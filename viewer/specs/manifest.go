package specs

import (
	"encoding/json"
	"fmt"
)

func (m Manifest) TraverseItemOrg(itemOrg Item, root string) {
	if itemOrg.Identifier != "" {
		m.traverseItemModules(itemOrg.Item, root)
	}
}

func (m Manifest) traverseItemModules(itemModules []Item, root string) {
	for i, _ := range itemModules {
		m.traverseItems(itemModules[i].Item, root)
	}
}

func (m Manifest) traverseItems(items []Item, root string) {
	fmt.Printf("- traversing batch of %d items\n", len(items))
	for i := range items {

		fmt.Printf("- - checking for idref: %v\n", items[i].Identifierref)

		err := m.resolveItem(items[i])

		if err != nil {
			fmt.Errorf(err.Error())
		}

		if len(items[i].Item) > 0 {
			m.traverseItems(items[i].Item, root)
		}

	}
}

//-- given an item, resolves its relationship to a resource
//-- since items are just folders with stuff inside
func (m Manifest) resolveItem(item Item) error {
	if item.Identifierref == "" {
		return fmt.Errorf("nope, no identifierref on item, skipping...")
	}

	identified := false
	for _, resource := range m.Resources.Resource {
		if resource.Identifier == item.Identifierref {
			identified = true
			fmt.Printf("- - mathched resource type %s\n", resource.Type)

			//-- TODO here we could parse the resources based on type
			switch resource.Type {
			case "imsdt_xmlv1p1":
				//-- topic
				break
			case "webcontent":
				//-- webcontent
				break
			case "imswl_xmlv1p1":
				//-- weblink
				break
			case "assignment_xmlv1p0":
				//-- assignment
				break
			case "assessment":
				//-- qti
				break
			default:
				return fmt.Errorf("[resolveItem] No matching type found: %s\n", resource.Type)
			}
		}
	}

	fmt.Printf("found match? %v\n", identified)

	return nil
}

func (m Manifest) PrettyPrint() {
	s, _ := json.MarshalIndent(m, "", "\t")
	fmt.Printf("manifest struct: %s\n", s)

	fmt.Printf("title: %v\n", m.Metadata.Lom.General.Title)
}
