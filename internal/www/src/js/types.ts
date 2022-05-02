export interface ManifestType {
  Metadata: {
    Lom: {
      General: {
        Title: {
          String: {
            Text: "";
          };
        };
      };
    };
  };
}

export interface ResourceType {
  //-- basic learning application resource / webcontent
  XMLName: {
    Local: string;
  };
  Type: string;
  Identifier: string;
  File: Array<{
    Href: string;
  }>;

  //-- topic
  Title: string;
  Text: {
    Text: string;
  };
  Attachments: {
    Text: string;
    Attachment: Array<{
      Text: string;
      Href: string;
    }>;
  };

  //-- assignment
  Gradable: {
    Text: string;
    PointsPossible: string;
  };
  SubmissionFormats: {
    Text: string;
    Format: Array<{
      Text: string;
      Type: string;
    }>;
  };

  //-- LTI
  Description: string;
  LaunchURL: string;
  SecureLaunchURL: string;
  Vendor: {
    Text: string;
    Name: string;
    Description: string;
    URL: string;
  };

  //-- QTI
  Assessment: {
    Title: string;
    Text: string;
  };

  //-- Weblink
  URL: {
    Text: string;
    Href: string;
  };
}

export interface ItemType {
  Item: {
    Identifier: string;
    Title: string;
  };
  Children: Array<ItemType>;
  Resources: Array<ResourceType>;
}
