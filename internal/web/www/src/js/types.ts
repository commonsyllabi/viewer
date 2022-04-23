export interface ManifestType {
    Metadata: {
        Lom: {
            General: {
                Title: {
                    String: {
                        Text: ""
                    }
                }
            }
        }
    }
}

export interface ResourceType {

    XMLName: {
        Local: string
    },
    Title: string,
    Type: string,
    Identifier: string,
    File: Array<{
        Href: string
    }>

}

export interface ItemType {
    Item: {
        Identifier: string,
        Title: string,
    },
    Children: Array<ItemType>,
    Resources: Array<ResourceType>
}