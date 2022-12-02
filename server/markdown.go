package server

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/qianlnk/codidoc/syscmd"
	"github.com/qianlnk/log"
)

// Download TODO
func (s *Server) Download() {
	for {
		docs, err := s.getDocs()
		if err != nil {
			log.Errorf("get docs err: %v", err)
			time.Sleep(s.cfg.DownloadFrequency)
			continue
		}

		for _, doc := range docs {
			md, err := toMarkdown(doc)
			if err != nil {
				if !strings.Contains(err.Error(), "empty doc") {
					log.Errorf("toMarkdown err: %v", err)
				}
				continue
			}

			err = md.localizaiton(s.cfg.MarkdownPath)
			if err != nil {
				log.Errorf("localizaiton %s err: %v", md.Name, err)
				continue
			}

			md.localizionImage(s.cfg.SourceImagePath, s.cfg.MarkdownPath)
		}

		// 暂时不做数据库自动备份github
		// s.backupMysql()

		time.Sleep(s.cfg.DownloadFrequency)
	}
}

// getDocs 文档本地化
func (s *Server) getDocs() ([]string, error) {
	sqltxt := fmt.Sprintf(`SELECT content FROM Notes WHERE ownerId = '%s'`, s.cfg.Owner)

	ssession := s.mysqlCli.NewSession()
	defer ssession.Close()

	rows, err := ssession.QueryString(sqltxt)
	if err != nil {
		return nil, err
	}

	var docs []string
	for _, row := range rows {
		docs = append(docs, row["content"])
	}

	return docs, nil
}

// Markdown TODO
type Markdown struct {
	Name    string
	Content string
	Images  []string
}

// toMarkdown 转markdown
func toMarkdown(doc string) (*Markdown, error) {
	r := bytes.NewReader([]byte(doc))
	scanner := bufio.NewScanner(r)

	var buf []string
	var imgs []string
	for scanner.Scan() {
		line := scanner.Text()
		buf = append(buf, line)

		content := strings.TrimSpace(line)
		if ok := strings.Contains(content, "![]("); ok {
			img := content[13 : len(content)-1]
			imgs = append(imgs, img)
		}
	}

	if len(buf) == 0 || len(buf[0]) == 0 {
		return nil, fmt.Errorf("empty doc")
	}

	return &Markdown{
		Name:    strings.TrimSpace(buf[0]),
		Content: strings.Join(buf[1:], "\n"),
		Images:  imgs,
	}, nil
}

// localizaiton 本地化
func (m *Markdown) localizaiton(path string) error {
	if m.Name[0] != '!' {
		return nil
	}

	filename := path + "/" + m.Name[1:]
	dir := getdir(filename)
	if err := mkdir(dir); err != nil {
		return err
	}

	file, err := os.Create(path + "/" + m.Name[1:])
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(m.Content)
	if err != nil {
		return err
	}

	return nil
}

func (m *Markdown) localizionImage(src string, path string) error {
	if m.Name[0] != '!' {
		return nil
	}

	if len(m.Images) == 0 {
		return nil
	}

	filename := path + "/" + m.Name[1:]
	dir := getdir(filename) + "/uploads"
	if err := mkdir(dir); err != nil {
		return err
	}

	srcs := make([]string, 0, len(m.Images))

	for _, img := range m.Images {
		srcs = append(srcs, src+"/"+img)
	}

	strSrcs := strings.Join(srcs, " ")
	err := syscmd.Cp(strSrcs, dir+"/")
	if err != nil {
		log.Errorf("cp %s %s err: %v", strSrcs, dir+"/", err)
		return err
	}

	return nil
}

func getdir(filename string) string {
	dirs := strings.Split(filename, "/")
	return strings.Join(dirs[:len(dirs)-1], "/")
}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		if err := os.Chmod(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) backupMysql() error {
	up := strings.Split(s.cfg.MysqlDSN, "@")[0]

	ups := strings.Split(up, ":")
	if len(ups) != 2 {
		return fmt.Errorf("git user password err")
	}

	db := strings.Split(s.cfg.MysqlDSN, "/")[1]
	dbs := strings.Split(db, "?")

	return syscmd.MysqlDump(ups[0], ups[1], dbs[0], s.cfg.MarkdownPath)
}
