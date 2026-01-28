# gitstats

[go-git-contributions](https://thevalleyofcode.com/go-git-contributions/) was the starting point for this project, so that I could learn some GO. I've modified the original version quite a bit, as it's an older code base and some things did not work and there were other things that did not work as I expected so I changed them. Originaly you could just see the green boxes of your local git repos, but I added an option to see yours or any other users green github boxes.


# Installation

Clone the repo:

```bash
https://github.com/Helland369/gitstats.git
```

Build the project:

```bash
cd gitstats

go build
```

# Add local repos to scan

To get the local git stats/green boxes (color may vary on your terminal theme)

```bash
./gitstats -add "path/to/your/local/repo"
```

# See your local git stats

```bash
./gitstats -email "your@mail.com"
```

# See your github stats

```bash
./gitstats -user "githubUserName"
```
