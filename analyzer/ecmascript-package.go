package analyzer

import (
	"encoding/json"
	"strings"
)

type EcmaScriptPackages struct {
	Packages []PackageJson
}

type PackageJson struct {
	Dependencies    map[Dependency]PackageVersion `json:"dependencies"`
	DevDependencies map[Dependency]PackageVersion `json:"devDependencies"`
}

func NewEcmaScriptPackages() *EcmaScriptPackages {
	return &EcmaScriptPackages{
		Packages: []PackageJson{},
	}
}

func (pckg *EcmaScriptPackages) IsPackageFile(filePath FilePath) bool {
	return strings.HasSuffix(string(filePath), "package.json")
}

func (pckg *EcmaScriptPackages) AddPackageFile(packageFileContents string) error {
	var packageFile PackageJson
	err := json.Unmarshal([]byte(packageFileContents), &packageFile)
	if err != nil {
		return err
	}

	pckg.Packages = append(pckg.Packages, packageFile)
	return nil
}

func (pckg *EcmaScriptPackages) GetDepFromPackages(dep Dependency) (bool, Dependency, PackageVersion) {
	for _, dependencies := range pckg.Packages {
		for dependency, version := range dependencies.Dependencies {
			if doesFileDepMatchPackageDep(dep, dependency) {
				return true, dependency, version
			}
		}

		for devDependency, version := range dependencies.DevDependencies {
			if doesFileDepMatchPackageDep(dep, devDependency) {
				return true, devDependency, version
			}
		}
	}

	return false, Dependency(""), PackageVersion("")
}

func doesFileDepMatchPackageDep(fileDep Dependency, packageDep Dependency) bool {
	fileDepStr := string(fileDep)
	packageDepStr := string(packageDep)

	itsAMatch := strings.Split(fileDepStr, "/")[0] == packageDepStr
	if strings.HasPrefix(fileDepStr, "@") {
		// The strings.Split(fileDepStr, "/")[1] bit fails sometimes (not sure why!)
		if len(strings.Split(fileDepStr, "/")) < 2 {
			return false
		}

		filePackageWithOrg := strings.Split(fileDepStr, "/")[0] + "/" + strings.Split(fileDepStr, "/")[1]
		itsAMatch = filePackageWithOrg == packageDepStr
	}

	return itsAMatch
}
