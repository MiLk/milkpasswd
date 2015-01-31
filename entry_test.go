package milkpasswd

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEntry(t *testing.T) {
	Convey("Should create, get and delete an entry", t, func() {
		path := "test"
		name := "rainbow"
		fullpath := "/" + path + "/" + name
		username := "unicorn"
		password := []byte("binary")
		website := "http://www.test.com"
		description := "Description"

		entry, _ := CreateEntry(path,
			name,
			username,
			password,
			website,
			description)
		entry.Save()

		get, _ := GetEntry(fullpath)
		So(get.FullPath(), ShouldEqual, fullpath)
		So(get.Path, ShouldEqual, "/"+path+"/")
		So(get.Name, ShouldEqual, name)
		So(get.Username, ShouldEqual, username)
		So(get.Password, ShouldResemble, password)
		So(get.Website, ShouldEqual, website)
		So(get.Description, ShouldEqual, description)

		So(get.String(), ShouldEqual,
			"Path: "+fullpath+
				" - Username: "+username+
				" - Website: "+website+
				" - Description: "+description)

		data, _ := SearchEntries(fullpath)
		So(len(data), ShouldEqual, 1)
		So(data[fullpath].FullPath(), ShouldEqual, fullpath)
		So(data[fullpath].Path, ShouldEqual, "/"+path+"/")
		So(data[fullpath].Name, ShouldEqual, name)
		So(data[fullpath].Username, ShouldEqual, username)
		So(data[fullpath].Password, ShouldResemble, password)
		So(data[fullpath].Website, ShouldEqual, website)
		So(data[fullpath].Description, ShouldEqual, description)

		data2, _ := ListEntries()
		found := false
		for k, v := range data2 {
			if k == fullpath {
				found = true
				So(v.FullPath(), ShouldEqual, fullpath)
				So(v.Path, ShouldEqual, "/"+path+"/")
				So(v.Name, ShouldEqual, name)
				So(v.Username, ShouldEqual, username)
				So(v.Password, ShouldResemble, password)
				So(v.Website, ShouldEqual, website)
				So(v.Description, ShouldEqual, description)
			}
		}
		So(found, ShouldEqual, true)

		DeleteEntry(fullpath)

		_, err := GetEntry(fullpath)
		So(err, ShouldNotEqual, nil)
	})
}
