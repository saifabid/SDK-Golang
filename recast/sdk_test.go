package recast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	//Basic test
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	r, err := c.TextRequest("Hello !")
	assert.True(t, err == nil)
	if assert.True(t, r != nil) {
		intent, err := r.Intent()
		assert.True(t, err == nil)
		assert.True(t, intent == "hello-greetings")
		assert.Equal(t, r.Source(), "Hello !")
		assert.Equal(t, r.Version(), "0.1.4")
		assert.Equal(t, r.Status(), 200)
		assert.True(t, r.Timestamp() != "")
		assert.True(t, len(r.Intents()) == 1)
	}

	//Test no intent
	r, err = c.TextRequest("abcd")
	if assert.True(t, err == nil) {
		assert.True(t, len(r.Intents()) == 0)
		_, err := r.Intent()
		assert.True(t, err != nil)
	}

	//Test invalid token
	c = NewClient("foobar")
	assert.True(t, c != nil)
	r, err = c.TextRequest("Hello !")
	assert.True(t, err != nil)
	c.SetToken("0496a4db6a50b7ec02e4ce101d3f8613")
	r, err = c.TextRequest("Hello!")
	assert.True(t, err == nil)
	assert.True(t, r != nil)

	c = NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	_, err = c.FileRequest("./test/test.wav")
	assert.True(t, err == nil)
	assert.True(t, r != nil)
	entities := r.AllEntities()
	assert.True(t, len(entities) == 0)

}

func TestSentences(t *testing.T) {
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.")
	assert.True(t, err == nil)
	assert.True(t, r != nil)

	assert.True(t, r.Sentence() != nil)
	assert.True(t, len(r.Sentences()) == 1)
	assert.True(t, r.Sentence().Source() == "I go from Paris to London.")
	s := r.Sentence()
	assert.True(t, s.Type() == "assert")
	assert.True(t, s.Action() == "go")
	assert.True(t, s.Agent() == "i")
	assert.True(t, s.Polarity() == "positive")
	assert.True(t, s.Entity("location") != nil)
	assert.True(t, s.AllEntities() != nil)
	assert.True(t, s.AllEntities("pronoun", "location") != nil)
	assert.True(t, len(s.AllEntities("pronoun", "location")) == 2)
	var locations []*Entity
	locations = s.AllEntities("location")["location"]
	assert.True(t, len(locations) == 2)
	location := locations[0]
	assert.True(t, location.Name() == "location")
	raw := location.Field("raw")
	assert.True(t, raw != nil)
	fakeField := location.Field("aDdasdasd")
	assert.True(t, fakeField == nil)
	assert.True(t, raw == location.Raw())
	assert.True(t, raw == "Paris")
	lat := location.Field("lat")
	assert.True(t, lat != nil)
	assert.True(t, lat.(float64) == 48.856614)
	pronoun := s.Entity("pronoun")
	assert.True(t, pronoun != nil)
	assert.True(t, pronoun.Field("person") != nil)
	assert.True(t, int(pronoun.Field("person").(float64)) == 1)
	assert.True(t, pronoun.Field("number").(string) == "singular")
	assert.True(t, pronoun.Field("gender").(string) != "")
	assert.True(t, pronoun.Field("swag") == nil)
	assert.True(t, pronoun.Field("raw").(string) == "I")
}

func TestResponseEntities(t *testing.T) {
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.")
	assert.True(t, err == nil)
	assert.True(t, r != nil)
	var entities map[string][]*Entity
	entities = r.AllEntities("pronoun", "location", "riendutout")
	assert.True(t, len(entities) == 2)
	assert.True(t, entities["pronoun"] != nil)
	assert.True(t, entities["location"] != nil)
	assert.True(t, len(entities["location"]) == 2)
	location := r.Entity("location")
	assert.True(t, location.Field("lat") != nil)
	assert.True(t, location.Field("foo") == nil)
	assert.True(t, location.Raw() == "Paris")
	assert.True(t, r.Entity("bar") == nil)
}

func TestMoar(t *testing.T) {
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.")
	assert.True(t, err == nil)
	assert.True(t, r != nil)
	locations := r.Entities("location")
	assert.True(t, len(locations) == 2)
	assert.True(t, locations[0].Raw() == "Paris")
	assert.True(t, locations[1].Raw() == "London")
	assert.True(t, locations[0].Field("lat").(float64) == 48.856614)
	assert.True(t, r.Entities("foo") == nil)
}

func TestFileUpload(t *testing.T) {
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613")
	assert.True(t, c != nil)
	r, err := c.FileRequest("foobar.wav")
	assert.True(t, err != nil)
	assert.True(t, r == nil)
	r, err = c.FileRequest("./test/invalid.wav")
	assert.True(t, err != nil)
	assert.True(t, r == nil)
	assert.True(t, err.Error() == `Request failed: 400 Bad Request ({"results":null,"message":"Speech is not recognizable."})`)
}
