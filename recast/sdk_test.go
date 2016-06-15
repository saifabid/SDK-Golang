package recast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	//Basic test
	c := NewClient("0496a4db6a50b7ec02e4ce101d3f8613", "en")
	require.True(t, c != nil)
	r, err := c.TextRequest("Hello !", nil)
	require.True(t, err == nil)
	require.True(t, r != nil)
	intent, err := r.Intent()
	require.True(t, err == nil)
	require.True(t, intent == "hello-greetings")
	require.Equal(t, r.Source(), "Hello !")
	require.Equal(t, r.Version(), "1.3.0")
	require.Equal(t, r.Status(), 200)
	require.True(t, r.Timestamp() != "")
	require.True(t, len(r.Intents()) == 1)

	//Test no intent
	r, err = c.TextRequest("abcd", nil)
	require.True(t, err == nil)
	require.True(t, len(r.Intents()) == 0)
	_, err = r.Intent()
	require.True(t, err != nil)

	//Test invalid token
	c = &Client{token: "foobar"}
	require.True(t, c != nil)
	r, err = c.TextRequest("Hello !", nil)
	require.True(t, err != nil)
	c.SetToken("0496a4db6a50b7ec02e4ce101d3f8613")
	r, err = c.TextRequest("Hello!", nil)
	require.True(t, err == nil)
	require.True(t, r != nil)

	c = &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613", language: "en"}
	require.True(t, c != nil)
	_, err = c.FileRequest("./test/test.wav", nil)
	require.True(t, err == nil)
	require.True(t, r != nil)
	entities := r.AllEntities()
	require.True(t, len(entities) == 0)

}

func TestOpts(t *testing.T) {
	c := &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613"}
	require.True(t, c != nil)
	_, err := c.TextRequest("I go from Paris to London.", map[string]string{
		"language": "foo",
	})
	require.True(t, err != nil)
	_, err = c.TextRequest("I go from Paris to London.", map[string]string{
		"language": "en",
		"token":    "foobar",
	})
	require.True(t, err != nil)
	_, err = c.FileRequest("test/test.wav", map[string]string{
		"language": "en",
	})
	require.True(t, err == nil)
}

func TestSentences(t *testing.T) {
	c := &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613"}
	require.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.", map[string]string{
		"token":    "c271ef3e774f72315ce6856f3bfc5876",
		"language": "en",
	})
	require.True(t, err == nil)
	require.True(t, r != nil)

	require.True(t, r.Sentence() != nil)
	require.True(t, len(r.Sentences()) == 1)
	require.True(t, r.Sentence().Source() == "I go from Paris to London.")
	s := r.Sentence()
	require.True(t, s.Type() == "assert")
	require.True(t, s.Action() == "go")
	require.True(t, s.Agent() == "i")
	require.True(t, s.Polarity() == "positive")
	require.True(t, s.Entity("location") != nil)
	require.True(t, s.AllEntities() != nil)
	require.True(t, s.AllEntities("pronoun", "location") != nil)
	require.True(t, len(s.AllEntities("pronoun", "location")) == 2)
	var locations []*Entity
	locations = s.AllEntities("location")["location"]
	require.True(t, len(locations) == 2)
	location := locations[0]
	require.True(t, location.Name() == "location")
	raw := location.Field("raw")
	require.True(t, raw != nil)
	fakeField := location.Field("aDdasdasd")
	require.True(t, fakeField == nil)
	require.True(t, raw == location.Raw())
	require.True(t, raw == "Paris")
	lat := location.Field("lat")
	require.True(t, lat != nil)
	require.True(t, lat.(float64) == 48.856614)
	pronoun := s.Entity("pronoun")
	require.True(t, pronoun != nil)
	require.True(t, pronoun.Field("person") != nil)
	require.True(t, int(pronoun.Field("person").(float64)) == 1)
	require.True(t, pronoun.Field("number").(string) == "singular")
	require.True(t, pronoun.Field("gender").(string) != "")
	require.True(t, pronoun.Field("swag") == nil)
	require.True(t, pronoun.Field("raw").(string) == "I")
}

func TestResponseEntities(t *testing.T) {
	c := &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613"}
	require.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.", nil)
	require.True(t, err == nil)
	require.True(t, r != nil)
	var entities map[string][]*Entity
	entities = r.AllEntities("pronoun", "location", "riendutout")
	require.True(t, len(entities) == 2)
	require.True(t, entities["pronoun"] != nil)
	require.True(t, entities["location"] != nil)
	require.True(t, len(entities["location"]) == 2)
	location := r.Entity("location")
	require.True(t, location.Field("lat") != nil)
	require.True(t, location.Field("foo") == nil)
	require.True(t, location.Raw() == "Paris")
	require.True(t, r.Entity("bar") == nil)
}

func TestMoar(t *testing.T) {
	c := &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613"}
	require.True(t, c != nil)
	r, err := c.TextRequest("I go from Paris to London.", nil)
	require.True(t, err == nil)
	require.True(t, r != nil)
	locations := r.Entities("location")
	require.True(t, len(locations) == 2)
	require.True(t, locations[0].Raw() == "Paris")
	require.True(t, locations[1].Raw() == "London")
	require.True(t, locations[0].Field("lat").(float64) == 48.856614)
	require.True(t, r.Entities("foo") == nil)
}

func TestFileUpload(t *testing.T) {
	c := &Client{token: "0496a4db6a50b7ec02e4ce101d3f8613"}
	require.True(t, c != nil)
	r, err := c.FileRequest("foobar.wav", nil)
	require.True(t, err != nil)
	require.True(t, r == nil)
	r, err = c.FileRequest("./test/invalid.wav", map[string]string{
		"token":    "0496a4db6a50b7ec02e4ce101d3f8613",
		"language": "en",
	})
	require.True(t, err == nil)
	require.True(t, r != nil)
}
