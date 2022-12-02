// Package syscmd TODO
package syscmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Mkdir TODO
func Mkdir(path string) error {
	return execCmd(fmt.Sprintf("mkdir -p %s", path))
}

// Rmdir TODO
func Rmdir(path string) error {
	return execCmd(fmt.Sprintf("rm -rf %s", path))
}

// Cp TODO
func Cp(src, dst string) error {
	return execCmd(fmt.Sprintf("cp -r %s %s", src, dst))
}

// Mv TODO
func Mv(src, dst string) error {
	return execCmd(fmt.Sprintf("mv %s %s", src, dst))
}

// Rm TODO
func Rm(file string) error {
	return execCmd(fmt.Sprintf("rm -f %s", file))
}

// Ln TODO
func Ln(src, dst string) error {
	return execCmd(fmt.Sprintf("ln -s %s %s", dst, src))
}

// Ls TODO
func Ls(path string) ([]string, error) {
	res, err := queryCmd(fmt.Sprintf("ls %s", path))
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(res[0:len(res)-1]), "\n"), nil
}

// GitStatus TODO
func GitStatus() ([]string, error) {
	res, err := Query("git status")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(res[0:len(res)-1]), "\n"), nil
}

// GitAdd TODO
func GitAdd() error {
	return Exec("git add --all")
}

// GitCommit TODO
func GitCommit(msg string) error {
	return Exec("git commit -m'%s'", msg)
}

// GitPush TODO
func GitPush() error {
	return Exec("git push")
}

// MysqlDump TODO
func MysqlDump(user string, password string, database string, path string) error {
	return Exec("mysqldump -u'%s' -p'%s' %s > %s/%s.sql", user, password, database, path, database)
}

func cmdoutput(command string) {
	host, _ := getHostname()
	fmt.Printf("%s@%s:~$ sudo %s\n", strings.Replace(os.Args[0], "./", "", -1), host, command)
}

// Exec TODO
func Exec(format string, args ...interface{}) error {
	return execCmd(fmt.Sprintf(format, args...))
}

// Query TODO
func Query(format string, args ...interface{}) ([]byte, error) {
	return queryCmd(fmt.Sprintf(format, args...))
}

func execCmd(command string) error {
	cmdoutput(command)
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s %s", err.Error(), stderr.String())
	}

	if len(stderr.String()) != 0 {
		fmt.Println(stderr.String())
	}

	return nil
}

func queryCmd(command string) ([]byte, error) {
	// cmdoutput(command)
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stderr = &stderr
	res, err := cmd.Output()
	if err != nil {
		return res, fmt.Errorf("%s %s", err.Error(), stderr.String())
	}

	if len(stderr.String()) != 0 {
		fmt.Println(stderr.String())
	}

	return res, nil
}

func getHostname() (string, error) {
	host := exec.Command("bash", "-c", "hostname")
	hostname, err := host.Output()
	return string(hostname[0 : len(hostname)-1]), err
}
