# Manga Download

### About
---
manga-dl is a cli tool build to download manga in pdf format from mangafire.to.
The binaries for this cli tool can be download from the [releases](https://github.com/xortock/manga-download/releases) page

### Install 
---

#### linux
```
curl -L https://github.com/xortock/manga-download/releases/download/v0.1.0/manga-dl-linux-amd64 -o ./manga-dl && sudo mv ./manga-dl /usr/local/bin/ && chmod +x /usr/local/bin/manga-dl
```

#### windows
```
New-Item -ItemType Directory -Force -Path "C:\Program Files\MangaDL";`
wget https://github.com/xortock/manga-download/releases/download/v0.1.0/manga-dl-windows-amd64.exe  -O "C:\Program Files\MangaDL\manga-dl.exe";`
[System.Environment]::SetEnvironmentVariable('PATH',"$env:path;C:\Program Files\MangaDL\", 'Machine')
```

### Usage
---
```
NAME:
   manga-dl - download manga in PDF or CBZ format

USAGE:
   manga-dl [global options] command [command options]

VERSION:
   development

DESCRIPTION:
   manga-dl is a cli tool build to download manga in PDF or CBZ format from mangafire.to

AUTHOR:
   xortock <bgmaduro@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --code value, -c value         the mangafire manga code
   --name value, -n value         the folder name where the manga chapters should be stored
   --output-path value, -o value  the output path for the manga folder
   --type value, -t value         the type in which the manga should be downloaded zip or cbz (default: "zip")
   --division value, -d value     the division type in which the manga should be download (does not make a difference in case of file type cbz) chapter or volume (default: "chapter")
   --help, -h                     show help
   --version, -v                  print the version

COPYRIGHT:
   (C) 2024 xortock
```