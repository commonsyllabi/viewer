package commoncartridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const singleTestFile = "./test_files/test_01.imscc"
const allTestFilesDir = "./test_files/dump"

// have setup() to create tmp files, with e.g. different schema versions
// check go test tables

func TestLoadEmpty(t *testing.T) {
	_, err := Load("")

	if err == nil {
		t.Fail()
	}
}

func TestParseManifest(t *testing.T) {
	// declaring cc as Cartridge to test that the return value implements Cartridge interface
	var cc Cartridge = load(t, singleTestFile)

	manifest, err := cc.ParseManifest()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(manifest).Kind() != reflect.Struct {
		t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(manifest).Kind())
	}
}

func TestExactCartridge(t *testing.T) {
	var cc Cartridge = load(t, singleTestFile)

	manifest, err := cc.ParseManifest()

	if err != nil {
		t.Errorf("Error parsing manifest: %v", err)
	}

	if manifest.Metadata.Lom.General.Title.String.Text != "Loaded Course" {
		t.Errorf("Expected title to be Loaded Course, got %s", manifest.Metadata.Lom.General.Title)
	}

	if manifest.Metadata.Lom.General.Description.String.Text != "Sample Description" {
		t.Errorf("Expected description to be Sample Description, got %s", manifest.Metadata.Lom.General.Description.String.Text)
	}

	if manifest.Metadata.Lom.General.Keyword.String.Text != "Test, Attempt" {
		t.Errorf("Expected description to be Test, Attempt, got %s", manifest.Metadata.Lom.General.Keyword.String.Text)
	}

	if manifest.Metadata.Lom.General.Language != "en-US" {
		t.Errorf("Expected description to be en-US, got %s", manifest.Metadata.Lom.General.Language)
	}

	if manifest.Metadata.Lom.LifeCycle.Contribute.Date.DateTime != "2014-09-08" {
		t.Errorf("Expected description to be 2014-09-08, got %s", manifest.Metadata.Lom.LifeCycle.Contribute.Date.DateTime)
	}

	if manifest.Metadata.Lom.Rights.CopyrightAndOtherRestrictions.Value != "yes" {
		t.Errorf("Expected description to be yes, got %s", manifest.Metadata.Lom.Rights.CopyrightAndOtherRestrictions.Value)
	}

	if manifest.Metadata.Lom.Rights.Description.String != "Private (Copyrighted) - http://en.wikipedia.org/wiki/Copyright" {
		t.Errorf("Expected description to be Private (Copyrighted) - http://en.wikipedia.org/wiki/Copyright, got %s", manifest.Metadata.Lom.Rights.Description.String)
	}

	if manifest.Organizations.Organization.Item.Identifier != "LearningModules" {
		t.Errorf("Expected top level item to be called Learning Module, got %v", manifest.Organizations.Organization.Item.Identifier)
	}

	public := manifest.Organizations.Organization.Item.Item[0]
	if len(public.Item) != 11 {
		t.Errorf("Expected number of public items to be 11, got %v", len(public.Item))
	}

	locked := manifest.Organizations.Organization.Item.Item[1]
	if len(locked.Item) != 1 {
		t.Errorf("Expected number of locked items to be 1, got %v", len(locked.Item))
	}

	if len(manifest.Resources.Resource) != 120 {
		t.Errorf("Expected to have 120 resources, got %v", len(manifest.Resources.Resource))
	}
}

func TestLoadCorrect(t *testing.T) {
	cc := load(t, singleTestFile)

	var empty IMSCC
	if reflect.DeepEqual(cc, empty) {
		t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
	}
}

func TestLoadAll(t *testing.T) {
	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(filepath.Join(cwd, allTestFilesDir))

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Test loading %d cartridges\n", len(files))
	var i int = 0
	for _, file := range files {
		i++
		if file.IsDir() {
			continue
		}

		cc, err := Load(filepath.Join(allTestFilesDir, file.Name()))

		if err != nil {
			t.Error(err)
		}

		var empty IMSCC
		if reflect.DeepEqual(cc, empty) {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
		}

		if cc.Title() == "" {
			t.Error("Cartridge Title should not be empty!")
		}

		fmt.Printf("Parsed %d/%d - %s\n", i, len(files), cc.Title())
	}
}

func TestParseManifestAll(t *testing.T) {
	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(filepath.Join(cwd, allTestFilesDir))

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("testing %d cartridges\n", len(files))
	var i int = 0
	for _, file := range files {
		i++
		if file.IsDir() {
			continue
		}

		cc, err := Load(filepath.Join(allTestFilesDir, file.Name()))

		if err != nil {
			t.Error(err)
		}

		var empty IMSCC
		if reflect.DeepEqual(cc, empty) {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
		}

		manifest, err := cc.ParseManifest()

		if err != nil {
			t.Errorf("Error parsing manifest: %v", err)
		}

		if reflect.TypeOf(manifest).Kind() != reflect.Struct {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(manifest).Kind())
		}
	}
}

func TestMetadata(t *testing.T) {
	cc := load(t, singleTestFile)
	meta, err := cc.Metadata()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(meta).Kind() != reflect.String {
		t.Errorf("Expecting metadata to be serialized JSON, got %v", reflect.TypeOf(meta).Kind())
	}
}

func TestResources(t *testing.T) {
	cc := load(t, singleTestFile)
	resources, err := cc.Resources()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(resources[0]) != reflect.TypeOf(FullResource{}) {
		t.Errorf("Expected Resources()[0] to be of type FullResource, got %v", reflect.TypeOf(resources[0]))
	}

	if resources[0].Resource.XMLName.Local != "resource" {
		t.Errorf("Expected FullResource to have a XMLName of 'resource', got %s", resources[0].Resource.XMLName.Local)
	}

	if len(resources) != 120 {
		t.Errorf("Expected 120 resources, got %d", len(resources))
	}
}

func TestAssignments(t *testing.T) {
	cc := load(t, singleTestFile)
	assignments, err := cc.Assignments()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(assignments) != reflect.TypeOf([]Assignment{}) {
		t.Errorf("Expected assignments to be of type []Assignment, got %v", reflect.TypeOf(assignments))
	}

	if strings.Contains(assignments[0].XMLName.Local, "assignment_xmlv1p") {
		t.Errorf("Expected assignemnt[0] to have a XMLName of 'assignment_xmlv1p', got %v", assignments[0].XMLName.Local)
	}
}

func TestLTIs(t *testing.T) {
	cc := load(t, singleTestFile)
	ltis, err := cc.LTIs()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(ltis) != reflect.TypeOf([]CartridgeBasicltiLink{}) {
		t.Errorf("Expected ltis to be of type []CartridgeBasicltiLink, got %v", reflect.TypeOf(ltis))
	}

	if strings.Contains(ltis[0].XMLName.Local, "imsbasiclti_xmlv1p") {
		t.Errorf("Expected assignemnt[0] to have a XMLName of 'imsbasiclti_xmlv1p', got %v", ltis[0].XMLName.Local)
	}
}

func TestQTIs(t *testing.T) {
	cc := load(t, singleTestFile)
	qtis, err := cc.QTIs()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(qtis) != reflect.TypeOf([]Questestinterop{}) {
		t.Errorf("Expected qtis to be of type []Questestinterop, got %v", reflect.TypeOf(qtis))
	}

	if strings.Contains(qtis[0].XMLName.Local, "imsbasicqti_xmlv1p") {
		t.Errorf("Expected assignemnt[0] to have a XMLName of 'imsbasicqti_xmlv1p', got %v", qtis[0].XMLName.Local)
	}
}

func TestTopics(t *testing.T) {
	cc := load(t, singleTestFile)
	topics, err := cc.Topics()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(topics) != reflect.TypeOf([]Topic{}) {
		t.Errorf("Expected topics to be of type []Topic, got %v", reflect.TypeOf(topics))
	}

	if strings.Contains(topics[0].XMLName.Local, "imsdt_xmlv1p") {
		t.Errorf("Expected assignemnt[0] to have a XMLName of 'imsdt_xmlv1p', got %v", topics[0].XMLName.Local)
	}
}

func TestWeblinks(t *testing.T) {
	cc := load(t, singleTestFile)
	weblinks, err := cc.Weblinks()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(weblinks) != reflect.TypeOf([]WebLink{}) {
		t.Errorf("Expected weblinks to be of type []webLink, got %v", reflect.TypeOf(weblinks))
	}

	if strings.Contains(weblinks[0].XMLName.Local, "imswl_xmlv1p") {
		t.Errorf("Expected assignment[0] to have a XMLName of 'imswl_xmlv1p', got %v", weblinks[0].XMLName.Local)
	}
}

func TestFind(t *testing.T) {
	cc := load(t, singleTestFile)

	//-- generic learing application resource
	found, err := cc.Find("ic1b5d76bd9a4bd37eb78cf0bcb5b84da")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(Resource{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- topic
	found, err = cc.Find("i528c2ce0186a758d13a9bd193bd88611")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(Topic{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- weblink
	found, err = cc.Find("ibb3ca45e774c0c487daeb9352e7a4553")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(WebLink{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- assignment
	found, err = cc.Find("ie801a403cd25e9a771ab7e3a2d6bea3a")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(WebLink{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- qti
	found, err = cc.Find("iad7e264143b9f2ec9dbc71a9d166f6f2")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(WebLink{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- lti
	found, err = cc.Find("iae0220efe8693f664806e9bfe43b6e30")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(WebLink{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}

	//-- webcontent
	found, err = cc.Find("i3755487a331b36c76cec8bbbcdb7cc66")

	if err != nil {
		t.Errorf("%v", err)
	}

	if reflect.TypeOf(found).Kind() != reflect.TypeOf(WebLink{}).Kind() {
		t.Errorf("Expected to have the returned resource to be of type Resource{}, got %v", reflect.TypeOf(found).Kind())
	}
}

func TestFindFile(t *testing.T) {
	cc := load(t, singleTestFile)

	bytes, err := cc.FindFile("i3755487a331b36c76cec8bbbcdb7cc66")

	if err != nil {
		t.Error(err)
	}

	if len(bytes) == 0 {
		t.Error("Found file should not be equal to 0")
	}
}

func TestDump(t *testing.T) {

	cc := load(t, singleTestFile)
	dump := cc.Dump()

	if reflect.TypeOf(dump).Kind() != reflect.Slice {
		t.Errorf("Expecting slice type, got: %v", reflect.TypeOf(dump).Kind())
	}

	if len(dump) == 0 {
		t.Error("Empty byte array returned!")
	}
}

func TestAsObject(t *testing.T) {

	cc := load(t, singleTestFile)

	obj, err := cc.AsObject()

	if err != nil {
		t.Errorf("Error parsing the JSON: %v\n", err)
	}

	if reflect.TypeOf(obj).Kind() != reflect.Slice {
		t.Errorf("Expecting slice type, got: %v", reflect.TypeOf(obj).Kind())
	}

	if len(obj) == 0 {
		t.Error("Empty byte array returned!")
	}

}

func load(t *testing.T, p string) Cartridge {
	cc, err := Load(p)
	if err != nil {
		t.Errorf("could not load %s: %s", p, err)
	}
	return cc
}
