This repository is a metadata repository for Spinnaker's Halyard to discover plugins available to install in a Spinnaker instance.

```
hal plugins repository add pulumi --url https://raw.githubusercontent.com/pulumi/spinnaker-plugins-repository/master/repositories.json
```

## `repositories.json`

This JSON file contains an array that lists all the `plugins.json` files for Halyard to scan to find a plugin. Each entry contains an `id` and a `url`.

## `plugins.json`

This file contains an array of plugins available served by this plugin repository. 

## `repository_dispatch` Event

When a new release of a Spinnaker plugin is made, a `repository_dispatch` POST request should be made to this repository with the information about the plugin that was released. The type of the event must always be `onPluginRelease`. This is the event type that the `release.yml` workflow file is expecting.

For example, consider a plugin repo at https://github.com/pulumi/a-plugin-repo, then the following JSON payload should be POSTed to https://api.github.com/repos/pulumi/spinnaker-plugins-repository/dispatches.

```json
{
	"event_type": "onPluginRelease",
	"client_payload": {
		"org": "pulumi",
		"repo": "a-plugin-repo",
		"released": {
			"id": "...",
			"description": "...",
			"provider": "https://github.com/pulumi",
			"releases": [{
				"version": "vx.y.z",
				"date": "2020-06-12T23:12:39.344522Z",
				"requires": "orca>=0.0.0",
				"sha512sum": "d73f6e7887eccad920d4e404d45112898c6ef8637f18d3d1b28d20f0fcd8f6d4e6d25a6755a0e97b772834035de49f9568f34b143be7b78d990d57d6798517dc",
				"state": "RELEASE"
			}]
		}
	}
}
```

Luckily, the Spinnaker Gradle plugin for building release version of a Spinnaker plugin produces a manifest file alongside the release JAR, which could be used as the value for the `released` property in the above JSON. In that case, the plugin repo just needs to add the following step to its release workflow file:

> Note that in order to make a `repository_dispatch` POST request, you need to add a Personal Access Token with the `repo` scope to the plugin repo.
> In the following snippet, that's the `USERNAME`, and `GH_PAT` GitHub Actions secrets.

```yaml
- name: Add release to plugin repo
  run: |
    curl -XPOST -u "${{ secrets.USERNAME }}:${{ secrets.GH_PAT }}" -H "Accept: application/vnd.github.everest-preview+json" -H "Content-Type: application/json" https://api.github.com/repos/pulumi/spinnaker-plugins-repository/dispatches --data "{\"event_type\": \"onPluginRelease\", \"client_payload\": {\"org\": \"${{ github.repository_owner }}\", \"repo\": \"${{ steps.get_project_info.outputs.PROJECT }}\", \"released\": $(cat build/distributions/plugin-info.json)}}"
```
