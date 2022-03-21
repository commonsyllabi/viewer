#!/bin/bash

echo "generating structs..."


touch viewer/specs/autogen_manifest.go && zek -P specs -t manifest -o viewer/specs/autogen_manifest.go tmp/canvas_large_1.3/imsmanifest.xml
touch viewer/specs/autogen_item.go && zek -P specs -t item -o viewer/specs/autogen_item.go tmp/canvas_large_1.3/imsmanifest.xml
touch viewer/specs/autogen_resource.go && zek -P specs -t resource -o viewer/specs/autogen_resource.go tmp/canvas_large_1.3/imsmanifest.xml

echo "done!"

echo "WARNING! you still need to remove anonymous structs from the manifest"