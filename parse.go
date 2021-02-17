package main

import "gopkg.in/yaml.v3"

type lockFile struct {
	Packages map[string]interface{}
}

func parsePackages(lockFileContent []byte) ([]string, error) {
	l := lockFile{}
	err := yaml.Unmarshal(lockFileContent, &l)
	if err != nil {
		return nil, err
	}

	var packageNames []string
	for name, _ := range l.Packages {
		packageNames = append(packageNames, name)
	}
	return packageNames, nil
}
