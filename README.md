# sklls

> sklls CLI makes it easy to pull git blame logs for repositories and group them by email address

```
Usage of sklls:
  -cloneDir string
    	Directory to use for cloning repos into (default: Temp dir is being used
  -concurrency int
    	Dictates how many files are analyzed in parallel (default 10)
  -dir string
    	Directory where to look for repos for. Default is CWD.
  -exclude string
    	Comma separated list of file suffixes to exclude from scanning (can be both extensions and filenames) (default "package-lock.json,yarn.lock")
  -ghpat string
    	Github Personal Access Token for cloning non-public repos
  -ghrepos string
    	Comma separated list of Github repos (org/repo - e.g. spring-media/ep-curato) to clone from Github & analyze (requires ghpat to be set as well, if repos are non-public)
  -out string
    	Output folder (default "./")
  -perfLogThreshold int
    	Performance data is logged for any file analysis that takes longer than perfLogThreshold (in ms) (default 500)
  -timeout int
    	Timeout in seconds for analyzing repositieries (default: 30 [seconds]) (default 30)
  -verbose
    	Enable verbose output
```

# Example output file
```json
// jonas.peeck@axelspringer.com
{
    "Ext": {
        ".md": {
            "Wed, 13 Jul 2019 16:57:22 +0200": 25
        }
    },
    "Dep": {},
    "Usernames": [
        "Jonas Peeck"
    ]
}
```

# Install
`make install`

# Try it yourself
`sklls -out ./sklls-data` (just run that in any folder that contains repositories)

# License
MIT