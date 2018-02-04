# mbbacklight

Utility for controlling the screen and keyboard backlight on my Macbook running GNU/Linux.

Must be run (at your own risk) using sudo, setuid or something similar.

## Usage

```
usage: mbbacklight [flags] [system] [operation]

Systems:
  -kbd
  -screen

Operations:
  -get
  -up
  -down
  -max
  -set [value]

Flags:
  -step int
        step value for command
```
