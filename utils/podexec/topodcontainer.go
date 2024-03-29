package podexec

import (
	"bytes"
	"fmt"
	"strings"
)

// ToPodContainer 从io.Writer拷贝到pod
// tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
func (p *PodExec) ToPodContainer(destPath string) error {
	if p.NoPreserve {
		p.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		p.Command = []string{"tar", "-xmf", "-"}
	}
	//destPath = path.Dir(destPath)
	if len(destPath) > 0 {
		p.Command = append(p.Command, "-C", destPath)
	}
	p.Tty = false
	var stderr bytes.Buffer
	p.Stderr = &stderr
	err := p.Exec()
	if err != nil {
		return fmt.Errorf("%s(%s)", err.Error(), stderr.String())
	}
	if len(stderr.Bytes()) != 0 {
		for _, line := range strings.Split(stderr.String(), "\n") {
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			if !strings.Contains(strings.ToLower(line), "removing") {
				return fmt.Errorf(line)
			}
		}
	}
	return nil
}
