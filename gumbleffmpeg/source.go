package gumbleffmpeg

import (
	"io"
	"os/exec"
)

// Source is a Stream source.
type Source interface {
	// must include the -i <filename>
	Arguments() []string
	Start(*exec.Cmd) error
	Done()
}

// sourceFile

type sourceFile string

// SourceFile is standard file source.
func SourceFile(filename string) Source {
	return sourceFile(filename)
}

func (s sourceFile) Arguments() []string {
	return []string{"-i", string(s)}
}

func (sourceFile) Start(*exec.Cmd) error {
	return nil
}

func (sourceFile) Done() {
}

// sourceReader

type sourceReader struct {
	r io.ReadCloser
}

// SourceReader is a ReadCloser source.
func SourceReader(r io.ReadCloser) Source {
	return &sourceReader{r}
}

func (*sourceReader) Arguments() []string {
	return []string{"-i", "-"}
}

func (s *sourceReader) Start(cmd *exec.Cmd) error {
	cmd.Stdin = s.r
	return nil
}

func (s *sourceReader) Done() {
	s.r.Close()
}

// sourceExec

type sourceExec struct {
	name string
	arg  []string

	cmd *exec.Cmd
}

// SourceExec uses the output of the given command and arguments as source
// data.
func SourceExec(name string, arg ...string) Source {
	return &sourceExec{
		name: name,
		arg:  arg,
	}
}

func (*sourceExec) Arguments() []string {
	return []string{"-i", "-"}
}

func (s *sourceExec) Start(cmd *exec.Cmd) error {
	s.cmd = exec.Command(s.name, s.arg...)
	r, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stdin = r
	if err := s.cmd.Start(); err != nil {
		cmd.Stdin = nil
		return err
	}
	return nil
}

func (s *sourceExec) Done() {
	if s.cmd != nil {
		if p := s.cmd.Process; p != nil {
			p.Kill()
		}
		s.cmd.Wait()
	}
}
