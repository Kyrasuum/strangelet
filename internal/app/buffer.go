package app

import (
	"io/ioutil"
	"os"
	"regexp"
	// "strconv"

	isatty "github.com/mattn/go-isatty"
	"strangelet/internal/buffer"
	// "strangelet/internal/cursor"
)

var ()

// LoadInput determines which files should be loaded into buffers
// based on the input stored in flag.Args()
func (app application) LoadInput(args []string) []*buffer.Buffer {
	// There are a number of ways strangelet should start given its input

	// 1. If it is given a files in flag.Args(), it should open those

	// 2. If there is no input file and the input is not a terminal, that means
	// something is being piped in and the stdin should be opened in an
	// empty buffer

	// 3. If there is no input file and the input is a terminal, an empty buffer
	// should be opened

	var filename string
	var input []byte
	var err error
	buffers := make([]*buffer.Buffer, 0, len(args))

	btype := buffer.BTDefault
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		btype = buffer.BTStdout
	}

	files := make([]string, 0, len(args))
	// flagStartPos := cursor.Loc{-1, -1}
	flagr := regexp.MustCompile(`^\+(\d+)(?::(\d+))?$`)
	for _, a := range args {
		match := flagr.FindStringSubmatch(a)
		if len(match) == 3 && match[2] != "" {
			// line, err := strconv.Atoi(match[1])
			if err != nil {
				app.TermMessage(err)
				continue
			}
			// col, err := strconv.Atoi(match[2])
			if err != nil {
				app.TermMessage(err)
				continue
			}
			// flagStartPos = cursor.Loc{col - 1, line - 1}
		} else if len(match) == 3 && match[2] == "" {
			// line, err := strconv.Atoi(match[1])
			if err != nil {
				app.TermMessage(err)
				continue
			}
			// flagStartPos = cursor.Loc{0, line - 1}
		} else {
			files = append(files, a)
		}
	}

	if len(files) > 0 {
		// Option 1
		// We go through each file and load it
		for i := 0; i < len(files); i++ {
			buf, err := buffer.NewBufferFromFile(files[i], btype) //, flagStartPos
			if err != nil {
				app.TermMessage(err)
				continue
			}
			// If the file didn't exist, input will be empty, and we'll open an empty buffer
			buffers = append(buffers, buf)
		}
	} else if !isatty.IsTerminal(os.Stdin.Fd()) {
		// Option 2
		// The input is not a terminal, so something is being piped in
		// and we should read from stdin
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			app.TermMessage("Error reading from stdin: ", err)
			input = []byte{}
		}
		buffers = append(buffers, buffer.NewBufferFromString(string(input), filename, btype)) //, flagStartPos
	} else {
		// Option 3, just open an empty buffer
		buffers = append(buffers, buffer.NewBufferFromString(string(input), filename, btype)) //, flagStartPos
	}

	return buffers
}
