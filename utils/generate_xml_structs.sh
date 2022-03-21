#!/bin/bash

echo "generating structs from manifest..."

declare -a arr=("manifest", "item", "resource")

touch viewer/specs/autogen_manifest.go && zek -P specs -t manifest -o viewer/specs/autogen_manifest.go tmp/canvas_large_1.3/imsmanifest.xml
touch viewer/specs/autogen_item.go && zek -P specs -t item -o viewer/specs/autogen_item.go tmp/canvas_large_1.3/imsmanifest.xml
touch viewer/specs/autogen_resource.go && zek -P specs -t resource -o viewer/specs/autogen_resource.go tmp/canvas_large_1.3/imsmanifest.xml

echo "done!"

echo "WARNING! you still need to remove anonymous structs from the manifest"

# example with go generate //go:generate bash -c "zek -e < repo.xml > repo.go.stub"
# example with xsdgen: xsdgen -f -o autogen_assignment.go cc_extresource_assignmentv1p0_v1p0.xsd :