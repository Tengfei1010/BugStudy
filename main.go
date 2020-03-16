package main

import (
	"github.com/Tengfei1010/BugStudy/analyzer"
	"github.com/google/logger"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init()  {
	//logPath := "./bug_study_log.txt"
	//lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	//if err != nil {
	//	panic(err)
	//}
	defer logger.Init("BugStudy", true, false, ioutil.Discard).Close()
	logger.SetFlags(log.LstdFlags)
	logger.Info("Logger works well.")
}


func main() {
	os.Mkdir("results", 0755)
	//os.Mkdir(csv_result_dir, 0755)
	//os.Mkdir(html_results_dir, 0755)
	//
	//if os.Args[1] == "test" {
	//	var new_counter PackageCounter = ParseDir("test", "tests", "")
	//	var test_counter Counter = HtmlOutputCounters([]*PackageCounter{&new_counter}, "test", "test", nil, "")
	//	OutputCounters("tests", []*PackageCounter{&new_counter}, "", test_counter)
	//	return
	//}

	//data, e := ioutil.ReadFile(os.Args[1])
	var repoPathStr = "C:/Go/GOPATH/src/github.com/Tengfei1010/BugStudy/repo.txt"
	data, e := ioutil.ReadFile(repoPathStr)
	if e != nil {
		logger.Errorf("prevent panic by handling " +
			"failure accessing a path %q: %v", repoPathStr, e)
		return
	}
	proLists := strings.Split(string(data), "\n")
	//var project_counters []analyzer.Counter
//
//	var indexData *analyzer.IndexFileData = &analyzer.IndexFileData{Indexes: []*analyzer.IndexData{}}

	for _, projectName := range proLists {
		if projectName == "" {
			continue
		}
		//projectPath := filepath.Base(string(projectName))

		var pathToDir string
		var commitHash string
		pathToDir, commitHash = analyzer.CloneRepo(string(projectName))
		if commitHash == "" {
			logger.Infof("Failed to clone the repo: %s", projectName)
			continue
		}
		// find go mod in pathToDir
		goModFile := filepath.Join(pathToDir, "go.mod")
		if _, err := os.Stat(goModFile); os.IsNotExist(err) {
			// TODO: if has not mod, run go mod init; go mod tidy ?
			logger.Infof("Project %v has not go mod file", projectName)
			continue
			// run go mod init
			//outStr, errStr := analyzer.RunModuleCmd(pathToDir, "go", "mod", "init")
			//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
			//
			//// run go mod vendor
			//outStr, errStr = analyzer.RunModuleCmd(pathToDir, "go", "mod", "tidy")
			//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
		}
		// run go mod vendor
		outStr, errStr := analyzer.RunModuleCmd(pathToDir, "go", "mod", "vendor")
		if outStr != "" {
			logger.Infof("out: %s", outStr)
		}
		if errStr != "" {
			logger.Errorf("err: %s", errStr)
		}

		sAPathToDir := filepath.Join(pathToDir, "...")
		outStr, errStr = analyzer.RunCmd("C:\\Go\\GOPATH\\staticcheck_windows_amd64\\staticcheck\\staticcheck.exe", sAPathToDir)
		logger.Infof("out: %s", outStr)
		if errStr != "" {
			logger.Errorf("err: %s", errStr)
		}
		//var packages []*analyzer.PackageCounter

		//err := filepath.Walk(pathToDir, func(path string, info os.FileInfo, err error) error {
		//
		//	if err != nil {
		//		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", pathToDir, err)
		//		return err
		//	}
		//	if info.IsDir() {
		//		if info.Name() == "vendor" || info.Name() == "tests" || info.Name() == "test" {
		//			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		//			return filepath.SkipDir
		//		}
		//		var packageCounter analyzer.PackageCounter = analyzer.ParseDir(projectPath, path, pathToDir)
		//		packages = append(packages, &packageCounter)
		//		return nil
		//	}
		//	return nil
		//})
		//
		//if err != nil {
		//	fmt.Printf("error walking the path %q: %v\n", pathToDir, err)
		//}
	//	var projectCounter analyzer.Counter = analyzer.HtmlOutputCounters(packages, commitHash, projectName, indexData, pathToDir) // html
	//
	//	analyzer.OutputCounters(projectName, packages, pathToDir, projectCounter) // csvs
	//	defer os.RemoveAll(pathToDir)                                             // clean up
	//	// project_counters = append(project_counters, projectCounter)
	}
	//createIndexFile(indexData) // index html

}
//func createIndexFile(index_data *analyzer.IndexFileData) {
//	f, err := os.Create("index.html")
//	if err != nil {
//		panic(err)
//	}
//	tmpl := template.Must(template.ParseFiles("../analyser/index_layout.html"))
//	tmpl.Execute(f, index_data) // write the index page
//}
