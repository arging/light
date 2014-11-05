// Copyright 2014 li. All rights reserved.

package view

import (
	"bufio"
	"github.com/roverli/light/conf"
	"github.com/roverli/light/log"
	"github.com/roverli/utils/errors"
	"github.com/roverli/utils/slice"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DIR       = "/view/" // Template directory relative to root.
	ScreenDIR = "/view/screen/"
	SUFFIX    = ".tpl" // Template file name suffix.
)

var (
	ErrViewNotFound = errors.New("light/view: view not found error.")
)

var (
	tplHeaders       map[string]string
	tpls             map[string]*template.Template
	tmpTpls          map[string]*template.Template
	screenFiles      []string
	hasDefaultScreen bool
)

func Render(tpl string, data Context, wr io.Writer) errors.Error {
	t := tpls[tpl]
	if t == nil {
		return ErrViewNotFound
	}

	err := t.Execute(wr, data)
	if err != nil {
		return errors.Wrapf(err, "light/view: parse tpl %s error.", tpl)
	}
	return nil
}

type Context map[string]interface{}

func init() {
	initScreen()
	initView()

	// TODO change to file watcher. This is just for develop test.
	// NOTE: Delete later
	if conf.IsDebug {
		go func() {
			t := time.NewTicker(time.Second * time.Duration(15))
			for _ = range t.C {
				log.Debug("light/view: Reload tpls.")
				initScreen()
				initView()
			}
		}()
	}
}

func initScreen() {
	screenNames := conf.List(ScreenDIR, func(fname string) bool {
		return strings.HasSuffix(fname, SUFFIX)
	})
	screenFiles = slice.MapString(screenNames,
		func(name string) string {
			if strings.HasPrefix(name, "default") {
				hasDefaultScreen = true
			}
			return conf.ROOT + ScreenDIR + name
		})
}

func initView() {
	tmpTpls := make(map[string]*template.Template)

	fInfos, err := listFile("")
	if err != nil {
		log.Warn("light/view: No view templates found.")
		return
	}

	tplHeaders = make(map[string]string)
	readTplsHeader(fInfos, "")

	for name, header := range tplHeaders {
		switch {
		case strings.HasPrefix(header, "$$screen="):
			parseByScreen(name, header[len("$$screen="):])

		case hasDefaultScreen:
			parseByScreen(name, "default")

		default:
			tpl, err := template.ParseFiles(conf.ROOT + DIR + name + SUFFIX)
			switch err {
			case nil:
				log.Infof("light/view: Load view succeed, view: %s .", name)
				tmpTpls[name] = tpl
			default:
				log.Errorf("light/view: Load view fail, view: %s. Error:%v", name, err)
			}
		}

	}

	tpls = tmpTpls

	// Reset vars
	tplHeaders = nil
	tmpTpls = nil
	screenFiles = nil
	hasDefaultScreen = false
}

func parseByScreen(name string, screen string) {
	parseFiles := make([]string, len(screenFiles)+2, len(screenFiles)+2)
	i := 1
	found := false
	for _, screenFile := range screenFiles {
		if strings.HasSuffix(screenFile, "/"+screen+SUFFIX) {
			parseFiles[0] = screenFile
			found = true
		} else {
			parseFiles[i] = screenFile
		}
		i++
	}

	if !found {
		log.Errorf("light/view: Cann't find screen %s for template %s.tpl", screen, name)
		return
	}

	parseFiles[i] = conf.ROOT + ScreenDIR + name + SUFFIX
	tpl, err := template.ParseFiles(parseFiles...)
	if err != nil {
		log.Errorf("light/view: Parse template %s error.", name)
		return
	}

	tmpTpls[name] = tpl
}

func readTplsHeader(fInfos []os.FileInfo, parent string) {

	for _, finfo := range fInfos {
		switch {
		case finfo.IsDir():
			childFInfos, err := listFile(parent + finfo.Name())
			if err != nil {
				log.Error("light/view: %v.", err)
			} else {
				readTplsHeader(childFInfos, parent+finfo.Name()+"/")
			}
		case strings.HasSuffix(finfo.Name(), SUFFIX):
			readHeaderData(parent + finfo.Name())
		}
	}
}

func listFile(dir string) ([]os.FileInfo, error) {
	dirFile, err := os.Open(filepath.Clean(conf.ROOT + DIR + dir))
	if err != nil {
		return nil, err
	}

	defer dirFile.Close()
	fileInfos, err := dirFile.Readdir(-1)
	if err != nil {
		return nil, err
	}

	filteredInfos := make([]os.FileInfo, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		// Screen files not load.
		if dir == "" && fileInfo.Name() == "screen" {
			continue
		}
		filteredInfos = append(filteredInfos, fileInfo)
	}
	return filteredInfos, nil
}

func readHeaderData(fname string) {
	f, err := os.Open(filepath.Clean(conf.ROOT + DIR + fname))
	if err != nil {
		log.Infof("light/view: Opean file %s error. Error: %v.", fname, err)
		return
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	// Read the first line
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Infof("light/view: %v.", err)
		return
	}

	// strip the ".tpl" suffix
	tplHeaders[fname[:(len(fname)-len(SUFFIX))]] = line
}
