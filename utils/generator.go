package utils

//go:generate bash -c "zek -P specs -t manifest -o ../viewer/specs/autogen_manifest.go ../viewer/specs/types/examples/manifest.xml"
//go:generate bash -c "zek -P specs -t item -o ../viewer/specs/autogen_item.go ../viewer/specs/types/examples/manifest.xml"
//-- when extracting Item, it needs to look very simple (xmlname, text, identifier, identifierref, []item)
//go:generate bash -c "zek -P specs -t resource -o ../viewer/specs/autogen_resource.go ../viewer/specs/types/examples/manifest.xml"

//go:generate bash -c "zek -P specs -t topic -o ../viewer/specs/autogen_topic.go ../viewer/specs/types/examples/topic.xml"
//go:generate bash -c "zek -P specs -t lti -o ../viewer/specs/autogen_lti.go ../viewer/specs/types/examples/lti.xml"
//go:generate bash -c "zek -P specs -t qti -o ../viewer/specs/autogen_qti.go ../viewer/specs/types/examples/qti.xml"
//go:generate bash -c "zek -P specs -t weblink -o ../viewer/specs/autogen_weblink.go ../viewer/specs/types/examples/weblink.xml"
