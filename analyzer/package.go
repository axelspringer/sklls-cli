package analyzer

type FilePath string
type PackageVersion string
type PackageName string

type Package interface {
	// Called by the file crawler (that searches for git repos initially)
	IsPackageFile(filePath FilePath) bool

	// Used to add package files
	AddPackageFile(packageFileContents string) error

	// Returns whether dependency was found in packages, and what version it has
	GetDepFromPackages(dep Dependency) (bool, Dependency, PackageVersion)
}
