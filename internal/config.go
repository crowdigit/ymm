package internal

type ConfigExec struct {
	Path string
	Args []string
}

func (c ConfigExec) ReplacePlaceholder(placeholder, value string) {
	for i := range c.Args {
		if c.Args[i] == placeholder {
			c.Args[i] = value
		}
	}
}

type Config struct {
	Command struct {
		Metadata struct {
			Youtube ConfigExec
			JSON    ConfigExec
		}
		Download struct {
			Youtube ConfigExec
		}
	}
}
