package specs

import (
	. "commonsyllabi/logger"
	"fmt"
)

func (m *Manifest) ResolveItems() {
	if m.Organizations.Organization.Item.Identifier != "" {
		m.traverseItemModules(m.Organizations.Organization.Item.Item)
	}
}

func (m *Manifest) traverseItemModules(itemModules []Item) {
	for i, _ := range itemModules {
		m.traverseItems(itemModules[i].Item)
	}
}

func (m *Manifest) traverseItems(items []Item) {
	Log.Debug().Int("- items traversed", len(items))
	for i := range items {

		Log.Debug().Str("- - idref", items[i].Identifierref)

		err := m.resolveItem(items[i])

		if err != nil {
			fmt.Errorf(err.Error())
		}

		if len(items[i].Item) > 0 {
			m.traverseItems(items[i].Item)
		}

	}
}

//-- given an item, resolves its relationship to a resource
//-- since items are just folders with stuff inside
func (m *Manifest) resolveItem(item Item) error {
	if item.Identifierref == "" {
		return fmt.Errorf("nope, no identifierref on item, skipping...")
	}

	for _, resource := range m.Resources.Resource {
		if resource.Identifier == item.Identifierref {
			// Log.Debug().Msg("- - mathched resource type %s\n", resource.Type)

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
	return nil
}

func (m *Manifest) PrettyPrint() {
	Log.Debug().Str("Cartridge:", m.Metadata.Lom.General.Title.String)
	Log.Debug().Int("Modules:", len(m.Organizations.Organization.Item.Item))
	for _, i := range m.Organizations.Organization.Item.Item {
		Log.Debug().Int("- items:", len(i.Item))
		for _, v := range i.Item {
			Log.Debug().Str("- - id", v.Identifierref)
		}
	}

	Log.Debug().Int("Resources (%d):\n", len(m.Resources.Resource))
	for _, r := range m.Resources.Resource {
		Log.Debug().Str(" - - type", r.Type).Str(" - - id", r.Identifier)
	}
}
