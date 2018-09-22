// +build integration

package e2e

import (
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

type cbProcess struct {
	Cmd    *exec.Cmd
	DBFile string
}

func (p *cbProcess) Destroy() {
	p.Cmd.Process.Kill()
	os.Remove(p.DBFile)
}

func spawnProcess() (*cbProcess, error) {

	f, err := ioutil.TempFile(baseTmpDir, "contactbook")
	if err != nil {
		return nil, err
	}
	f.Close()

	if err := os.Remove(f.Name()); err != nil {
		return nil, err
	}

	cb := &cbProcess{
		DBFile: f.Name(),
	}

	cb.Cmd = exec.Command(path.Join(binDir, "contactbook"),
		"-user", testUser,
		"-password", testPassword,
		"-db-file", cb.DBFile)

	if err := cb.Cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		cb.Cmd.Wait()
	}()

	time.Sleep(1 * time.Second)

	return cb, nil
}

func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
}

var runeSet = []rune("987654321")

func randNum(n int) string {
	r := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		b[i] = runeSet[r.Intn(len(runeSet))]
	}

	return string(b)
}
