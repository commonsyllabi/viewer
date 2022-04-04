package commoncartridge

import (
	zero "commonsyllabi/internals/logger" //-- TODO remove this dependency because internals should not be imported in the public package, and because zerolog is an opinion and libraries shouldn't import opinions
	"fmt"
)

func (m *Manifest) ResolveItems() {
	if m.Organizations.Organization.Item.Identifier != "" {
		m.traverseItemModules(m.Organizations.Organization.Item.Item)
	}
}

func (m *Manifest) traverseItemModules(itemModules []Item) {
	for i := range itemModules {
		m.traverseItems(itemModules[i].Item)
	}
}

func (m *Manifest) traverseItems(items []Item) {
	zero.Log.Debug().Int("- items traversed", len(items))
	for i := range items {

		zero.Log.Debug().Str("- - idref", items[i].Identifierref)

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
			// Log.Debug().Msg("- - matched resource id %s\n", resource.Identifier)

			switch resource.Type {
			case "imsdt_xmlv1p1":
				//-- topic
				zero.Log.Info().Msgf("found topic %v", resource.File)

			case "webcontent":
				//-- webcontent
				zero.Log.Info().Msgf("found webcontent %v", resource.File)

			case "imswl_xmlv1p1":
				//-- weblink
				zero.Log.Info().Msgf("found weblink %v", resource.File)

			case "assignment_xmlv1p0":
				//-- assignment
				zero.Log.Info().Msgf("found assignment %v", resource.File)

			case "assessment":
				//-- qti
				zero.Log.Info().Msgf("found question bank %v", resource.File)

			default:
				return fmt.Errorf("[resolveItem] No matching type found: %s\n", resource.Type)
			}
		}
	}
	return nil
}

func (m *Manifest) PrettyPrint() {
	zero.Log.Debug().Str("Cartridge:", m.Metadata.Lom.General.Title.String.Text)
	zero.Log.Debug().Int("Modules:", len(m.Organizations.Organization.Item.Item))
	for _, i := range m.Organizations.Organization.Item.Item {
		zero.Log.Debug().Int("- items:", len(i.Item))
		for _, v := range i.Item {
			zero.Log.Debug().Str("- - id", v.Identifierref)
		}
	}

	zero.Log.Debug().Int("Resources (%d):\n", len(m.Resources.Resource))
	for _, r := range m.Resources.Resource {
		zero.Log.Debug().Str(" - - type", r.Type).Str(" - - id", r.Identifier)
	}
}
