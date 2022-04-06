package commoncartridge

//go:generate echo "Generating Manifest, Item, Resource, ..."
//go:generate bash -c "zek -P commoncartridge -t manifest -o ./autogen_manifest.go ./types/examples/manifest.xml"
//go:generate bash -c "zek -P commoncartridge -t item -o ./autogen_item.go ./types/examples/item.xml"
//go:generate echo "WARNING! when extracting Item, it needs to have a []Item field and not []struct"
//go:generate bash -c "zek -P commoncartridge -t resource -o ./autogen_resource.go ./types/examples/resource.xml"

//go:generate echo "Generating Topic, LTI, QTI, WebLink, ..."
//go:generate bash -c "zek -P commoncartridge -t topic -o ./autogen_topic.go ./types/examples/topic.xml"
//go:generate bash -c "zek -P commoncartridge -t lti -o ./autogen_lti.go ./types/examples/lti.xml"
//go:generate bash -c "zek -P commoncartridge -t qti -o ./autogen_qti.go ./types/examples/qti.xml"
//go:generate bash -c "zek -P commoncartridge -t weblink -o ./autogen_weblink.go ./types/examples/weblink.xml"
//go:generate bash -c "zek -P commoncartridge -t assignment -o ./autogen_assignment.go ./types/examples/assignment.xml"

//go:generate "...done!"
