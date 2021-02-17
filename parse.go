package main

import "gopkg.in/yaml.v3"

type parse struct {
	Packages map[string]struct {
		Dependency  string
		Description interface{}
		Source      string
		Version     string
	}
}

type LockFile struct {
	Packages map[string]Package
}

type PackageDescription struct {
	Name string
	URL  string
}

type Package struct {
	Dependency  string
	Description PackageDescription
	Source      string
	Version     string
}

func parseLockFile(lockFileContent []byte) (*LockFile, error) {
	internal := parse{}
	err := yaml.Unmarshal(lockFileContent, &internal)
	if err != nil {
		return nil, err
	}

	l := LockFile{Packages: map[string]Package{}}

	for k, v := range internal.Packages {
		var descr PackageDescription
		if m, ok := v.Description.(map[string]interface{}); ok {
			descr = PackageDescription{
				Name: m["name"].(string),
				URL:  m["url"].(string),
			}
		} else if s, ok := v.Description.(string); ok {
			descr = PackageDescription{
				Name: s,
				URL:  "local",
			}
		}

		l.Packages[k] = Package{
			Dependency:  v.Dependency,
			Description: descr,
			Source:      v.Source,
			Version:     v.Version,
		}
	}

	return &l, nil
}
