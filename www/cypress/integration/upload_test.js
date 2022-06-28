// 1. set up app state
// 2. take an action
// 3. make an assertion about the resulting state

describe('Home page', () => {
  it('clicks to upload a cartridge', () => {
    cy.visit('/')
    cy.get("#cta-upload").click()
  });

  it('scrolls through listed syllabi', () => {

  });
});

describe('Cartridge page', () => {
  it('uploads a cartridge', () => {
    cy.intercept('POST', '/api/upload', { fixture: 'test_01_upload_response.json' }).as('uploadCartridge')
    cy.visit('/cartridge.html')

    cy.get("#upload-file").selectFile("cypress/fixtures/test_01.imscc", { force: true });
    cy.wait("@uploadCartridge");

    cy.get("div.title").first().contains("Loaded Course")
  });

  it('fills in the email', () => {
    cy.get('#show-upload').click()
    cy.get("#email").type("pierre.depaz@gmail.com")
    cy.get("#email-conf").type("pierre.depaz@gmail.com")
  });

  it('submits a new course', () => {
    cy.intercept('POST', '/syllabi/', { fixture: 'test_01_submit_response.json' }).as('submitCartridge')
    cy.get("#course-submit").click() //-elaborate with more precise state checking
    cy.wait("@submitCartridge")
    cy.get("#submit-log").contains("success")
  })

  it('selects an example', () => {
    cy.visit('/cartridge.html')

    cy.get("#examples").select('0')

    cy.get("div.title").first().contains("Loaded Course")
    cy.get("#course-submit").should('not.exist')
  })

  it('resets the example', () => {
    cy.get("#reset-upload").click()
    cy.get("#examples").should('exist')
    cy.get("div.title").should('not.exist')
  })
});