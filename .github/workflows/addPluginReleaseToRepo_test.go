package main

import (
	"reflect"
	"testing"
)

func TestAddingReleaseToEmptyRepo(t *testing.T) {

    release := []Release {
        {
            "v1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "",
        }}
    pluginReleaseEvent := PluginReleaseEvent{"org1", "repo1",
        Plugin {
            "pluginId1",
            "plugin description",
            "provider1",
            release,
        }}

    var existingPlugins []Plugin

    result := addReleaseToPlugins(pluginReleaseEvent, existingPlugins)

    expectedReleases := []Release {
        {
            "1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
        }}
    expectedPlugins := []Plugin {
        {
            "pluginId1",
            "plugin description",
            "provider1",
            expectedReleases,
        },
    }

    if !reflect.DeepEqual(result, expectedPlugins) {
        t.Errorf("Release was not added correctly: %s", result)
    }

}

func TestAddingReleaseToRepoWithOtherPlugins(t *testing.T) {

    release := []Release {
        {
            "v1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "",
        }}
    pluginReleaseEvent := PluginReleaseEvent{"org1", "repo1",
        Plugin {
            "pluginId1",
            "plugin description",
            "provider1",
            release,
        }}

    existingReleases := []Release {
        {
            "1.0.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo0/releases/download/v1.0.0/repo0-v1.0.0.zip",
        }}
    existingPlugins := []Plugin {
        {
            "pluginId0",
                "plugin description",
                "provider1",
            existingReleases,
        }}

    result := addReleaseToPlugins(pluginReleaseEvent, existingPlugins)

    expectedReleases := []Release {
        {
            "1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
        }}
    expectedPlugins := []Plugin {
        existingPlugins[0],
        {
            "pluginId1",
            "plugin description",
            "provider1",
            expectedReleases,
        },
    }
    if !reflect.DeepEqual(result, expectedPlugins) {
        t.Errorf("Release was not added correctly: %s", result)
    }

}

func TestAddingReleaseToRepoWithExistingReleases(t *testing.T) {

    release := []Release {
        {
            "v1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "",
        }}
    pluginReleaseEvent := PluginReleaseEvent{"org1", "repo1",
        Plugin {
            "pluginId1",
            "new plugin description",
            "provider1",
            release,
        }}

    existingReleases := []Release {
        {
            "1.0.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo0/releases/download/v1.0.0/repo1-v1.0.0.zip",
        }}
    existingPlugins := []Plugin {
        {
            "pluginId1",
            "plugin description",
            "provider1",
            existingReleases,
        }}

    result := addReleaseToPlugins(pluginReleaseEvent, existingPlugins)

    expectedReleases := []Release {
        {
            "1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
        },
        existingReleases[0],
    }
    expectedPlugins := []Plugin {
        {
            "pluginId1",
            "new plugin description",
            "provider1",
            expectedReleases,
        },
    }

    if !reflect.DeepEqual(result, expectedPlugins) {
        t.Errorf("Release was not added correctly: %s", result)
    }

}


func TestAddingReleaseToRepoWithTheSameRelease(t *testing.T) {

    release := []Release {
        {
            "v1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "",
        }}
    pluginReleaseEvent := PluginReleaseEvent{"org1", "repo1",
        Plugin {
            "pluginId1",
            "plugin description",
            "provider1",
            release,
        }}

    existingReleases := []Release {
        {
            "1.2.0",
            "2020-01-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo0/releases/download/v1.0.0/repo1-v1.0.0.zip",
        },
        {
            "1.0.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo0/releases/download/v1.0.0/repo1-v1.0.0.zip",
        }}
    existingPlugins := []Plugin {
        {
            "pluginId1",
            "new plugin description",
            "provider1",
            existingReleases,
        }}

    result := addReleaseToPlugins(pluginReleaseEvent, existingPlugins)

    expectedReleases := []Release {
        {
            "1.2.0",
            "2020-02-24T20:46:40.585Z",
            "orca>=0.0.0",
            "asdf",
            "RELEASE",
            "https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
        },
        existingReleases[1],
    }
    expectedPlugins := []Plugin {
        {
            "pluginId1",
            "plugin description",
            "provider1",
            expectedReleases,
        },
    }


    if !reflect.DeepEqual(result, expectedPlugins) {
        t.Errorf("Release was not added correctly: %s", result)
    }

}
