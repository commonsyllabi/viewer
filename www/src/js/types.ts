export interface SyllabusType {
  ID: number,
  description: string,
  title: string,
  Attachments: Array<AttachmentType>,

}

export interface AttachmentType {
  ID: number,
  Name: string,
  Type: string,
  File: string,
  SyllabusID: number,
  Syllabus: SyllabusType
}

export interface ContributorType {
  ID: number,
  Name: string,
  Email: string,
  Syllabi: Array<SyllabusType>
}

export interface ManifestType {
  Metadata: {
    Lom: {
      General: {
        Language: ""
        Title: {
          String: {
            Text: "";
          };
        },
        Description: {
          String: {
            Text: "";
          }
        }
      };
      LifeCycle: {
        Contribute: {
          Date: {
            DateTime: ""
          };
          Entity: {
            Text: "",
            String: ""
          },
          Role: {
            Text: "",
            String: ""
          }
        }
      };
      Rights: {
        CopyrightAndOtherRestrictions: {
          Value: ""
        };
        Description: {
        String: ""
        }
      }
    };
    Schema: "",
    Schemaversion: ""
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
