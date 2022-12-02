package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/qianlnk/codidoc/syscmd"
	"github.com/qianlnk/log"
)

// PushLoop TODO
func (s *Server) PushLoop() {
	for {
		s.Push()
		time.Sleep(s.cfg.GitPushFrequency)
	}
}

// Push TODO
func (s *Server) Push() error {

	os.Chdir("/root/blog")

	status, err := syscmd.GitStatus()
	if err != nil {
		log.Errorf("git status err: %v", err)
		return err
	}

	log.Info(status)

	if len(status) <= 5 {
		return nil
	}

	err = syscmd.GitAdd()
	if err != nil {
		log.Errorf("git add err: %v", err)
		return err
	}

	msg := fmt.Sprintf("--auto=%s\n%s", time.Now().Format("2006-01-02 15:04"), strings.Join(status[6:len(status)-1], "\n"))

	err = syscmd.GitCommit(msg)
	if err != nil {
		log.Errorf("git commit err: %v", err)
		return err
	}

	err = syscmd.GitPush()
	if err != nil {
		log.Errorf("git push err: %v", err)
		return err
	}

	return nil

}
