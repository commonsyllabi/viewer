// 1. set up app state
// 2. take an action
// 3. make an assertion about the resulting state

describe('Cartridge page', () => {
  it('uploads a cartridge', () => {
    cy.intercept('POST', '/api/upload', { fixture: 'test_01_upload_response.json' }).as('uploadCartridge')
    cy.visit('/cartridge.html')

    cy.get("#upload-file").selectFile("cypress/fixtures/test_01.imscc");
    cy.get("#upload-submit").click();

    cy.get("#log").contains("uploading");

    cy.wait("@uploadCartridge");
    cy.get("#log").contains("uploaded");
  });

  it('fills in the email', () => {
    cy.get("#email").type("pierre.depaz@gmail.com")
    cy.get("#email-conf").type("pierre.depaz@gmail.com")
  });

  it('submits a new course', () => {
    cy.get("#course-submit")
  })
});
