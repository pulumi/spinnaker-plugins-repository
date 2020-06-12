package main

import (
	"reflect"
	"testing"
)

func TestAddingReleaseToEmptyRepo(t *testing.T) {

	rel := []release{
		{
			"v1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"",
		}}
	releaseEvent := pluginReleaseEvent{"org1", "repo1",
		plugin{
			"pluginId1",
			"plugin description",
			"provider1",
			rel,
		}}

	var existingPlugins []plugin

	result := addReleaseToPlugins(releaseEvent, existingPlugins)

	expectedReleases := []release{
		{
			"1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
		}}
	expectedPlugins := []plugin{
		{
			"pluginId1",
			"plugin description",
			"provider1",
			expectedReleases,
		},
	}

	if !reflect.DeepEqual(result, expectedPlugins) {
		t.Errorf("release was not added correctly: %s", result)
	}

}

func TestAddingReleaseToRepoWithOtherPlugins(t *testing.T) {

	rel := []release{
		{
			"v1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"",
		}}
	releaseEvent := pluginReleaseEvent{"org1", "repo1",
		plugin{
			"pluginId1",
			"plugin description",
			"provider1",
			rel,
		}}

	existingReleases := []release{
		{
			"1.0.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"https://github.com/org1/repo0/releases/download/v1.0.0/repo0-v1.0.0.zip",
		}}
	existingPlugins := []plugin{
		{
			"pluginId0",
			"plugin description",
			"provider1",
			existingReleases,
		}}

	result := addReleaseToPlugins(releaseEvent, existingPlugins)

	expectedReleases := []release{
		{
			"1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"https://github.com/org1/repo1/releases/download/v1.2.0/repo1-v1.2.0.zip",
		}}
	expectedPlugins := []plugin{
		existingPlugins[0],
		{
			"pluginId1",
			"plugin description",
			"provider1",
			expectedReleases,
		},
	}
	if !reflect.DeepEqual(result, expectedPlugins) {
		t.Errorf("release was not added correctly: %s", result)
	}

}

func TestAddingReleaseToRepoWithExistingReleases(t *testing.T) {

	rel := []release{
		{
			"v1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"",
		}}
	releaseEvent := pluginReleaseEvent{"org1", "repo1",
		plugin{
			"pluginId1",
			"new plugin description",
			"provider1",
			rel,
		}}

	existingReleases := []release{
		{
			"1.0.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"https://github.com/org1/repo0/releases/download/v1.0.0/repo1-v1.0.0.zip",
		}}
	existingPlugins := []plugin{
		{
			"pluginId1",
			"plugin description",
			"provider1",
			existingReleases,
		}}

	result := addReleaseToPlugins(releaseEvent, existingPlugins)

	expectedReleases := []release{
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
	expectedPlugins := []plugin{
		{
			"pluginId1",
			"new plugin description",
			"provider1",
			expectedReleases,
		},
	}

	if !reflect.DeepEqual(result, expectedPlugins) {
		t.Errorf("release was not added correctly: %s", result)
	}

}

func TestAddingReleaseToRepoWithTheSameRelease(t *testing.T) {

	rel := []release{
		{
			"v1.2.0",
			"2020-02-24T20:46:40.585Z",
			"orca>=0.0.0",
			"asdf",
			"RELEASE",
			"",
		}}
	releaseEvent := pluginReleaseEvent{"org1", "repo1",
		plugin{
			"pluginId1",
			"plugin description",
			"provider1",
			rel,
		}}

	existingReleases := []release{
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
	existingPlugins := []plugin{
		{
			"pluginId1",
			"new plugin description",
			"provider1",
			existingReleases,
		}}

	result := addReleaseToPlugins(releaseEvent, existingPlugins)

	expectedReleases := []release{
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
	expectedPlugins := []plugin{
		{
			"pluginId1",
			"plugin description",
			"provider1",
			expectedReleases,
		},
	}

	if !reflect.DeepEqual(result, expectedPlugins) {
		t.Errorf("release was not added correctly: %s", result)
	}

}
