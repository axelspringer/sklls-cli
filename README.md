# sklls
> sklls is a POC and not actively maintained

sklls is a simple CLI tool that aggregates git-blames by committer and generates line-counts for every committer based on file-extension and NPM dependencies.

The purpose of sklls is to automatically generate tech-usage profiles of contributors (which can serve as a proxy for tech sklls - get it?) working on any number of repositories.

> Check out how we used sklls at Axel Springer to build an automatically generated team-directory that helps developers find each other based on tech-stack: `TODO <insert blogpost here>`

# How to use
* Install `sklls` (see below)
* Run `$ sklls -out sklls-example-data` 

After installing sklls (see below) just run `sklls -out sklls-example-data` in any folder that contains git-repos.

sklls will run through all sub-folders and analyze any git-repos it can find (that have an origin set - aka which are hosted somewhere).  

After sklls has run through successfully, you will find JSON files in `/sklls-example-data` similar to the following example output:

## sklls output
```json
{
    "Ext": {
        "": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 38617,
            "Wed, 21 Dec 2022 16:45:57 +0100": 21
        },
        ".ds_store": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 2
        },
        ".gitignore": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 1
        },
        ".go": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 2929
        },
        ".json": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 668310
        },
        ".md": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 51
        },
        ".mod": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 22
        },
        ".sum": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 41
        },
        ".txt": {
            "Thu, 29 Dec 2022 16:19:58 +0100": 1
        }
    },
    "Dep": {
        "EcmaScript": {
            "@jest/globals@^29.1.2": {
                "Mon, 05 Sep 2022 12:12:24 +0200": 63,
                "Sun, 06 Nov 2022 14:12:51 +0100": 520,
                "Sun, 06 Nov 2022 15:51:31 +0100": 8
            },
            "@mui/icons-material@^5.10.9": {
                "Sun, 06 Nov 2022 14:12:51 +0100": 728
            },
            "@mui/material@^5.10.12": {
                "Sun, 06 Nov 2022 14:12:51 +0100": 843
            },
            "@octokit/graphql@^5.0.1": {
                "Sun, 06 Nov 2022 14:12:51 +0100": 404
            },
            "@storybook/react@^6.5.12": {
                "Mon, 19 Sep 2022 16:11:56 +0200": 92
            },
            "@storybook/testing-library@^0.0.13": {
                "Mon, 19 Sep 2022 16:11:56 +0200": 26
            },
            "file-extension-icon-js@^1.1.6": {
                "Sun, 06 Nov 2022 14:12:51 +0100": 728
            },
            "glob@^8.0.3": {
                "Mon, 05 Sep 2022 12:12:24 +0200": 115,
                "Sun, 06 Nov 2022 14:12:51 +0100": 51,
                "Thu, 08 Sep 2022 13:18:28 +0200": 4
            },
            "http-proxy@^1.18.1": {
                "Mon, 19 Sep 2022 16:11:56 +0200": 12
            },
            "react-dom@^18.2.0": {
                "Mon, 19 Sep 2022 10:19:53 +0200": 13
            },
            "react-router-dom@^6.4.0": {
                "Mon, 19 Sep 2022 10:19:53 +0200": 5,
                "Mon, 19 Sep 2022 16:11:56 +0200": 13,
                "Sun, 06 Nov 2022 14:12:51 +0100": 145
            },
            "react@^18.2.0": {
                "Mon, 19 Sep 2022 10:19:53 +0200": 13,
                "Mon, 19 Sep 2022 16:11:56 +0200": 269,
                "Sun, 06 Nov 2022 14:12:51 +0100": 843
            },
            "swr@^1.3.0": {
                "Mon, 19 Sep 2022 10:19:53 +0200": 5,
                "Mon, 19 Sep 2022 16:11:56 +0200": 25,
                "Sun, 06 Nov 2022 14:12:51 +0100": 30
            }
        }
    },
    "Usernames": [
        "aGuyNamedJonas"
    ]
}
```

# Usage
```
$ sklls -help

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

# Install
* Make sure go is installed
* Clone this repo
* Run `make install` from the project root
* After installation is complete, you can run `sklls` from anywhere on your machine

# License
MIT