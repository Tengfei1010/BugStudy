package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// parse a particular dir
func ParseDir(projName string, pathToDir string, pathToMainDir string) PackageCounter {

	var fileSet *token.FileSet = token.NewFileSet()
	var counter PackageCounter = PackageCounter{
		Counter: Counter{
			Go_count:     0,
			Send_count:   0,
			Rcv_count:    0,
			Chan_count:   0,
			IsPackage:    true,
			Project_name: projName},
		File_counters: []*Counter{}}

	f, err := parser.ParseDir(fileSet, pathToDir, nil, parser.AllErrors)

	if projName == "test" {
		ast.Print(fileSet, f)
	}
	if err != nil {
		fmt.Printf("An error was found in package %s : %v", filepath.Base(pathToDir), err)
	}

	if len(f) == 0 {
		return counter
	}

	for packName, node := range f {
		var packageCounterChan chan Counter = make(chan Counter)
		counter.Counter.Package_name = strings.TrimPrefix(strings.TrimPrefix(pathToDir, pathToMainDir)+"/"+packName, "/")
		counter.Counter.Package_path = pathToDir
		// Analyse each file
		for name, file := range node.Files {
			filename := strings.TrimPrefix(strings.TrimPrefix(pathToDir, pathToMainDir)+"/"+filepath.Base(name), "/")
			go AnalyseAst(fileSet, packName, filename, file, packageCounterChan, name) // launch a goroutine for each file
		}

		// Receive the results of the analysis of each file
		for range node.Files {

			var newCounter Counter = <-packageCounterChan

			newCounter.IsPackage = false
			newCounter.Project_name = projName
			if len(newCounter.Features) > 0 {
				newCounter.Has_feature = true
			}
			counter.Counter.Go_count += newCounter.Go_count
			counter.Counter.Send_count += newCounter.Send_count
			counter.Counter.Rcv_count += newCounter.Rcv_count
			counter.Counter.Chan_count += newCounter.Chan_count
			counter.Counter.Go_in_for_count += newCounter.Go_in_for_count
			counter.Counter.Range_over_chan_count += newCounter.Range_over_chan_count
			counter.Counter.Go_in_constant_for_count += newCounter.Go_in_constant_for_count
			counter.Counter.Array_of_channels_count += newCounter.Array_of_channels_count
			counter.Counter.Sync_Chan_count += newCounter.Sync_Chan_count
			counter.Counter.Known_chan_depth_count += newCounter.Known_chan_depth_count
			counter.Counter.Unknown_chan_depth_count += newCounter.Unknown_chan_depth_count
			counter.Counter.Make_chan_in_for_count += newCounter.Make_chan_in_for_count
			counter.Counter.Make_chan_in_constant_for_count += newCounter.Make_chan_in_constant_for_count
			counter.Counter.Constant_chan_array_count += newCounter.Constant_chan_array_count
			counter.Counter.Chan_slice_count += newCounter.Chan_slice_count
			counter.Counter.Chan_map_count += newCounter.Chan_map_count
			counter.Counter.Close_chan_count += newCounter.Close_chan_count
			counter.Counter.Select_count += newCounter.Select_count
			counter.Counter.Default_select_count += newCounter.Default_select_count
			counter.Counter.Assign_chan_in_for_count += newCounter.Assign_chan_in_for_count
			counter.Counter.Chan_of_chans_count += newCounter.Chan_of_chans_count
			counter.Counter.Send_chan_count += newCounter.Send_chan_count
			counter.Counter.Receive_chan_count += newCounter.Receive_chan_count
			counter.Counter.Param_chan_count += newCounter.Param_chan_count

			counter.File_counters = append(counter.File_counters, &newCounter)

		}

	}
	return counter
}
