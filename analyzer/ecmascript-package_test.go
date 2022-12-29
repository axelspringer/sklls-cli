package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const packageJson = `
{
	"name": "stack-starter",
	"version": "0.1.0",
	"bin": {
	  "stack-starter": "bin/stack-starter.js"
	},
	"scripts": {
	  "build": "tsc",
	  "watch": "tsc -w",
	  "test": "jest",
	  "cdk": "cdk",
	  "destroy": "cdk destroy",
	  "deploy-dev": "cdk deploy --app 'cdk.out/' sklls-dev",
	  "deploy-production": "cdk deploy --app 'cdk.out/' sklls-production --require-approval never",
	  "diff": "cdk diff",
	  "synth": "cdk synth"
	},
	"devDependencies": {
	  "@aws-cdk/assert": "2.3.0",
	  "@types/jest": "^27.4.0",
	  "@types/node": "17.0.8",
	  "jest": "^27.4.6",
	  "ts-jest": "^27.1.2",
	  "ts-node": "^10.4.0",
	  "typescript": "~4.5.4"
	},
	"dependencies": {
	  "aws-cdk": "^2.3.0",
	  "aws-cdk-lib": "^2.3.0",
	  "constructs": "^10.0.25",
	  "dotenv": "^10.0.0",
	  "source-map-support": "^0.5.21"
	}
  }
`

const anotherPackageJson = `
{
	"dependencies": {
		"react": "1.2.3"
	}
}
`

func TestIsEcmaScriptPackageFile(t *testing.T) {
	p := NewEcmaScriptPackages()
	assert.Equal(t, p.IsPackageFile("go.mod"), false, "go.mod file is recognized as NOT a package.json file")
	assert.Equal(t, p.IsPackageFile("/Users/jpeeck/dev/projectA/package.json"), true, "package.json is correctly recognized")
}

func TestAddPackageFile(t *testing.T) {
	p := NewEcmaScriptPackages()

	const brokenJson = `{`
	err := p.AddPackageFile(brokenJson)
	assert.NotNil(t, err, "Expect faulty JSON to return error")

	err = p.AddPackageFile(packageJson)
	assert.Nil(t, err, "Expect no error to be thrown for correct json")
	assert.Equal(t, p.Packages[0].Dependencies["aws-cdk"], PackageVersion("^2.3.0"), "AddPackageFile works to parse package file correctly")
}

func TestGetDepFromPackages(t *testing.T) {
	p := NewEcmaScriptPackages()
	p.AddPackageFile(packageJson)
	p.AddPackageFile(anotherPackageJson)

	found, dep, version := p.GetDepFromPackages(Dependency("react"))
	assert.Equal(t, true, found, "Returns true if found")
	assert.Equal(t, Dependency("react"), dep, "Returns correct dependency (the one from package)")
	assert.Equal(t, PackageVersion("1.2.3"), version, "Returns the correct version")

	found, dep, version = p.GetDepFromPackages(Dependency("aws-cdk-lib/aws-s3"))
	assert.Equal(t, true, found, "Returns true also for deep imports from packages")
	assert.Equal(t, Dependency("aws-cdk-lib"), dep, "Returns correct dependency (the one from package)")
	assert.Equal(t, PackageVersion("^2.3.0"), version, "Returns the correct version")

	found, dep, version = p.GetDepFromPackages(Dependency("@types/node/random"))
	assert.Equal(t, true, found, "Returns true also for deep imports from packages")
	assert.Equal(t, Dependency("@types/node"), dep, "Returns correct dependency (the one from package)")
	assert.Equal(t, PackageVersion("17.0.8"), version, "Returns the correct version")

	found, dep, version = p.GetDepFromPackages(Dependency("DOES_NOT_EXIST"))
	assert.Equal(t, false, found, "Returns false for packages that don't exist")
	assert.Equal(t, Dependency(""), dep, "Returns correct dependency (the one from package)")
	assert.Equal(t, PackageVersion(""), version, "Returns empty package version for not found packages")

	found, dep, version = p.GetDepFromPackages(Dependency(""))
	assert.Equal(t, false, found, "Returns false for empty package")
	assert.Equal(t, Dependency(""), dep, "Returns correct dependency (the one from package)")
	assert.Equal(t, PackageVersion(""), version, "Returns empty package version for empty package")
}
