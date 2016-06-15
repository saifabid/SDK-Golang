package recast

import (
	"encoding/json"
	"errors"
)

// Response handles the response from Recast.Ai and provides utility to get informations
// about data
type Response struct {
	status    int
	source    string
	version   string
	intents   []string
	sentences []*Sentence
	language  string
	timestamp string
}

func newResponse(jsonString string) (*Response, error) {
	r := &Response{}
	var temp map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &temp)
	if err != nil {
		return nil, err
	}
	resultsMap := temp["results"].(map[string]interface{})
	r.status = int(resultsMap["status"].(float64))
	r.source = resultsMap["source"].(string)
	r.version = resultsMap["version"].(string)
	r.timestamp = resultsMap["timestamp"].(string)
	r.language = resultsMap["language"].(string)
	r.intents = make([]string, len(resultsMap["intents"].([]interface{})))
	for i, intent := range resultsMap["intents"].([]interface{}) {
		r.intents[i] = intent.(string)
	}
	r.sentences = make([]*Sentence, len(resultsMap["sentences"].([]interface{})))
	for i, sentence := range resultsMap["sentences"].([]interface{}) {
		r.sentences[i] = newSentence(sentence.(map[string]interface{}))
	}
	return r, nil
}

// Language returns the language of the processed text
func (r *Response) Language() string {
	return r.language
}

// Status returns the status of the request
func (r *Response) Status() int {
	return r.status
}

// Source returns the original text sent to Recast
func (r *Response) Source() string {
	return r.source
}

// Timestamp returns the timestamp of the request formatted following the ISO 8061 standard
func (r *Response) Timestamp() string {
	return r.timestamp
}

// Version returns the Recast version that processes the input
func (r *Response) Version() string {
	return r.version
}

//Intents returns a slice of strings representing the matched intents, order by probability
func (r *Response) Intents() []string {
	return r.intents
}

// Intent returns the main intent matched by Recast or an error if no intent where found
func (r *Response) Intent() (string, error) {
	if len(r.intents) > 0 {
		return r.intents[0], nil
	}
	return "", errors.New("No intent found")
}

// Sentences returns a slice of Sentence
func (r *Response) Sentences() []*Sentence {
	return r.sentences
}

// Sentence returns the first sentence of the input
func (r *Response) Sentence() *Sentence {
	return r.sentences[0]
}

// Entity returns the first entity matching the name parameter
func (r *Response) Entity(name string) *Entity {
	for _, sentence := range r.sentences {
		if ent := sentence.Entity(name); ent != nil {
			return ent
		}
	}
	return nil
}

// Entities returns a slice of Entity containing all entities matching with name
func (r *Response) Entities(name string) []*Entity {
	var entities []*Entity
	for _, sentence := range r.sentences {
		if sentenceEntities := sentence.Entities(name); len(sentenceEntities) > 0 {
			entities = append(entities, sentenceEntities...)
		}
	}
	if len(entities) > 0 {
		return entities
	}
	return nil
}

// AllEntities returns a map containing slices of entities matching names, with names as keys
// AllEntities called with no arguments returns a map of all entities detected in the input
func (r *Response) AllEntities(names ...string) map[string][]*Entity {
	entities := make(map[string][]*Entity)
	for _, sentence := range r.sentences {
		if sentenceEntities := sentence.AllEntities(names...); len(sentenceEntities) > 0 {
			for entityName, ents := range sentenceEntities {
				entities[entityName] = append(entities[entityName], ents...)
			}
		}
	}
	return entities
}
