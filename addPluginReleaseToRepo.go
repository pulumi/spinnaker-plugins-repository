package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type pluginReleaseEvent struct {
	Org      string
	Repo     string
	Released plugin
}

type plugin struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Provider    string    `json:"provider"`
	Releases    []release `json:"releases"`
}

type release struct {
	Version   string `json:"version"`
	Date      string `json:"date"`
	Requires  string `json:"requires"`
	Sha512sum string `json:"sha512sum"`
	State     string `json:"state"`
	Url       string `json:"url"`
}

func main() {

	pluginReleaseJson := []byte(os.Args[1])
	var pluginReleaseEvent pluginReleaseEvent
	pluginReleaseErr := json.Unmarshal(pluginReleaseJson, &pluginReleaseEvent)
	check(pluginReleaseErr)

	pluginsJson, pluginsJsonReadErr := ioutil.ReadFile("plugins.json")
	check(pluginsJsonReadErr)
	var plugins []plugin
	pluginsErr := json.Unmarshal(pluginsJson, &plugins)
	check(pluginsErr)

	updatedPlugins := addReleaseToPlugins(pluginReleaseEvent, plugins)

	encodeBuffer := new(bytes.Buffer)
	enc := json.NewEncoder(encodeBuffer)
	enc.SetEscapeHTML(false)
	encodeErr := enc.Encode(updatedPlugins)
	check(encodeErr)
	var indentBuffer bytes.Buffer
	indentErr := json.Indent(&indentBuffer, encodeBuffer.Bytes(), "  ", "  ")
	check(indentErr)

	pluginsJsonWriteErr := ioutil.WriteFile("plugins.json", indentBuffer.Bytes(), 0644)
	check(pluginsJsonWriteErr)
}

func addReleaseToPlugins(releaseEvent pluginReleaseEvent, existingPlugins []plugin) []plugin {
	releasedPlugin := releaseEvent.Released
	rel := releasedPlugin.Releases[0]
	rel.Url = "https://github.com/" + releaseEvent.Org + "/" + releaseEvent.Repo + "/releases/download/" + rel.Version + "/" + releaseEvent.Repo + "-" + rel.Version + ".zip"
	if strings.HasPrefix(rel.Version, "v") {
		// the plugins version is supplied with a v from the bundler, but fails when update manager compares versions
		rel.Version = rel.Version[1:]
	}
	releasedPlugin.Releases = []release{
		rel,
	}

	for ip, existingPlugin := range existingPlugins {
		if existingPlugin.Id == releasedPlugin.Id {
			for ir, existingRelease := range existingPlugin.Releases {
				if existingRelease.Version == rel.Version {
					existingPlugin.Releases = append(existingPlugin.Releases[:ir], existingPlugin.Releases[ir+1:]...)
				}
			}
			releasedPlugin.Releases = append(releasedPlugin.Releases, existingPlugin.Releases...)
			existingPlugins[ip] = releasedPlugin
			return existingPlugins
		}
	}

	return append(existingPlugins, releasedPlugin)
}
