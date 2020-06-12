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

type PluginReleaseEvent struct {
    Org string
    Repo string
    Released Plugin
}

type Plugin struct {
    Id string           `json:"id"`
    Description string  `json:"description"`
    Provider string     `json:"provider"`
    Releases []Release  `json:"releases"`
}

type Release struct {
    Version string      `json:"version"`
    Date string         `json:"date"`
    Requires string     `json:"requires"`
    Sha512sum string    `json:"sha512sum"`
    State string        `json:"state"`
    Url string          `json:"url"`
}

func main() {

    pluginReleaseJson := []byte(os.Args[1])
    var pluginReleaseEvent PluginReleaseEvent
    pluginReleaseErr := json.Unmarshal(pluginReleaseJson, &pluginReleaseEvent)
    check(pluginReleaseErr)

    pluginsJson, pluginsJsonReadErr := ioutil.ReadFile("plugins.json")
    check(pluginsJsonReadErr)
    var plugins []Plugin
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


func addReleaseToPlugins(releaseEvent PluginReleaseEvent, existingPlugins []Plugin) []Plugin {
    releasedPlugin := releaseEvent.Released
    release := releasedPlugin.Releases[0]
    release.Url = "https://github.com/" + releaseEvent.Org + "/" + releaseEvent.Repo + "/releases/download/" + release.Version + "/" + releaseEvent.Repo + "-" + release.Version + ".zip"
    if strings.HasPrefix(release.Version,"v") {
        // the plugins version is supplied with a v from the bundler, but fails when update manager compares versions
        release.Version = release.Version[1:]
    }
    releasedPlugin.Releases = []Release {
        release,
    }

    for ip, existingPlugin := range existingPlugins {
        if existingPlugin.Id == releasedPlugin.Id {
            for ir, existingRelease := range existingPlugin.Releases {
                if existingRelease.Version == release.Version {
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



