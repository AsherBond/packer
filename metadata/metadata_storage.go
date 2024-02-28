package metadata

import (
	"fmt"
	"sync"

	"github.com/hashicorp/packer/version"
)

type PluginComponentUsedInBuild struct {
	BuildName    string
	ComponentKey string
}

type MetadataStorage struct {
	PluginComponentUsage []*PluginComponentUsedInBuild
}

var (
	metadataStorage     *MetadataStorage
	metadataStorageOnce sync.Once
)

func GetMetadataStorage() *MetadataStorage {
	metadataStorageOnce.Do(func() {
		metadataStorage = &MetadataStorage{
			PluginComponentUsage: []*PluginComponentUsedInBuild{},
		}
	})
	return metadataStorage
}

func (ms *MetadataStorage) AddPluginUsageMetadataFor(buildName string, componentKey string) {
	ms.PluginComponentUsage = append(ms.PluginComponentUsage, &PluginComponentUsedInBuild{
		BuildName:    buildName,
		ComponentKey: componentKey,
	})
}

func (ms *MetadataStorage) GetMetadataByBuild() map[string]map[string]string {
	pluginStorage := GetAllPluginsStorage()
	resp := map[string]map[string]string{}

	fmt.Printf("<===> All Plugins: %q\n", pluginStorage.Components)
	fmt.Printf("<===> All Builds: %q\n", ms.PluginComponentUsage)
	for _, pluginUsage := range ms.PluginComponentUsage {
		pluginUsageByBuild, ok := resp[pluginUsage.BuildName]
		if !ok {
			pluginUsageByBuild = map[string]string{}
		}
		fmt.Printf("<===> Build Name: %q\n", pluginUsage.BuildName)
		fmt.Printf("<===> Plugin Component: %q\n", pluginUsage.ComponentKey)

		compDetails := pluginStorage.GetPluginDetailsFor(pluginUsage.ComponentKey)
		if compDetails == nil {
			fmt.Printf("<===> Can not find component details. \n")
			continue
		}

		fmt.Printf("<===> Plugin Name: %q\n", compDetails.Name)
		fmt.Printf("<===> Plugin Version: %q\n", compDetails.Description.Version)
		pluginUsageByBuild[compDetails.Name] = compDetails.Description.Version
		resp[pluginUsage.BuildName] = pluginUsageByBuild
	}
	return resp
}

type PluginMetadata struct {
	Name    string
	Version string
}

type PackerCoreMetadata struct {
	Version string
}

type BuildsMetadata struct {
	PluginsMetadata    map[string][]PluginMetadata
	PackerCoreMetadata PackerCoreMetadata
}

func (ms *MetadataStorage) GetMetadata() *BuildsMetadata {
	pluginStorage := GetAllPluginsStorage()

	respMetadata := BuildsMetadata{
		PackerCoreMetadata: PackerCoreMetadata{Version: version.FormattedVersion()},
	}

	for _, pluginComponentUsage := range ms.PluginComponentUsage {
		pluginUsageByBuild, ok := respMetadata.PluginsMetadata[pluginComponentUsage.BuildName]
		if !ok {
			pluginUsageByBuild = []PluginMetadata{}
		}

		fmt.Printf("<===> Build Name: %q\n", pluginComponentUsage.BuildName)
		fmt.Printf("<===> Plugin Component: %q\n", pluginComponentUsage.ComponentKey)

		compDetails := pluginStorage.GetPluginDetailsFor(pluginComponentUsage.ComponentKey)
		if compDetails == nil {
			fmt.Printf("<===> Can not find component details. \n")
			continue
		}

		fmt.Printf("<===> Plugin Name: %q\n", compDetails.Name)
		fmt.Printf("<===> Plugin Version: %q\n", compDetails.Description.Version)
		pluginUsageByBuild = append(pluginUsageByBuild, PluginMetadata{
			Name:    compDetails.Name,
			Version: compDetails.Description.Version,
		})
		respMetadata.PluginsMetadata[pluginComponentUsage.BuildName] = pluginUsageByBuild
	}

	return &respMetadata
}
