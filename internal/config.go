package internal

type ConfigExec struct {
	Path string
	Args []string
}

func (c ConfigExec) ReplacePlaceholder(placeholder, value string) ConfigExec {
	args := make([]string, 0, len(c.Args))
	for _, arg := range c.Args {
		if arg == placeholder {
			args = append(args, value)
		} else {
			args = append(args, arg)
		}
	}
	return ConfigExec{
		Path: c.Path,
		Args: args,
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
		Replaygain ConfigExec
	}
}
