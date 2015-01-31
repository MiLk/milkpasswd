package milkpasswd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

var _ = fmt.Println

type Entry struct {
	Path        string
	Name        string
	Username    string
	Password    []byte
	Website     string
	Description string

	oldPath string
	oldName string
}

func (e *Entry) SetPath(value string) {
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	if !strings.HasSuffix(value, "/") {
		value = value + "/"
	}

	if value == e.Path {
		return
	}

	if e.oldPath == "" {
		e.oldPath = e.Path
	}
	e.Path = value
	return
}

func (e *Entry) SetName(value string) {
	if value == e.Name {
		return
	}

	if e.oldName == "" {
		e.oldName = e.Name
	}
	e.Name = value
	return
}

func (e *Entry) FullPath() string {
	return e.Path + e.Name
}

func (e *Entry) oldFullPath() string {
	if e.oldPath != "" && e.oldName != "" {
		return e.oldPath + e.oldName
	} else if e.oldPath != "" {
		return e.oldPath + e.Name
	} else if e.oldName != "" {
		return e.Path + e.oldName
	} else {
		return ""
	}
}

func (e *Entry) Bytes() []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(e)
	return buffer.Bytes()
}

func (e *Entry) String() string {
	return "Path: " + e.FullPath() +
		" - Username: " + e.Username +
		" - Website: " + e.Website +
		" - Description: " + e.Description
}

func (e *Entry) Save() (err error) {
	// TODO: Transaction

	// Remove old record
	if len(e.oldFullPath()) > 0 {
		err = deleteRecord(e.oldFullPath())
		if err != nil {
			return
		}
		e.oldPath = ""
		e.oldName = ""
	}

	// Insert new record
	err = setRecord(e.FullPath(), e.Bytes())
	return
}

func CreateEntry(path, name, username string, password []byte, website, description string) (*Entry, error) {
	entry := new(Entry)
	entry.SetPath(path)
	entry.SetName(name)
	entry.Username = username
	entry.Password = password
	entry.Website = website
	entry.Description = description
	return entry, nil
}

func GetEntry(fullPath string) (entry *Entry, err error) {
	data, err := getRecord(fullPath)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	dec := gob.NewDecoder(&buffer)

	_, err = buffer.Write(data)
	if err != nil {
		return
	}

	dec.Decode(&entry)
	return
}

func DeleteEntry(fullPath string) (err error) {
	err = deleteRecord(fullPath)
	return
}

func decodeRecords(data map[string][]byte) (entries map[string]*Entry, err error) {
	entries = make(map[string]*Entry)
	for path, value := range data {
		var buffer bytes.Buffer
		dec := gob.NewDecoder(&buffer)

		_, err = buffer.Write(value)
		if err != nil {
			return
		}
		var entry *Entry
		dec.Decode(&entry)

		entries[path] = entry
	}

	return
}
func ListEntries() (entries map[string]*Entry, err error) {
	data, err := listRecords()
	if err != nil {
		return
	}

	entries, err = decodeRecords(data)
	return
}

func SearchEntries(prefix string) (entries map[string]*Entry, err error) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	data, err := searchRecords(prefix)
	if err != nil {
		return
	}

	entries, err = decodeRecords(data)
	return
}
