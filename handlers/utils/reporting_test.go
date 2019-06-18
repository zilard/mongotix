
package utils

import (
    "bufio"
    "bytes"
    "flag"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "testing"

    s "github.com/zilard/mongotix/structs"
)


var testNodeMetricsMap = make(s.NodeMetricsMap)
var testProcessMetricsArray = []s.ProcessMetricsByName{}


var update = flag.Bool("update", false, "update .golden files")

const TIMESLICE = 360


// TestCreateNodeAverageReport - Golden test for CreateNodeAverageReport function
func TestCreateNodeAverageReport(t *testing.T) {

    timeSlice := float64(TIMESLICE)

    testtable := []struct {
        tname string
    }{
        {
            tname: "ok",
        },
    }

    CreateDummyNodeMetrics(testNodeMetricsMap)

    for _, tc := range testtable {

        t.Run(tc.tname, func(t *testing.T) {

            var buffer bytes.Buffer
            writer := bufio.NewWriter(&buffer)

            err := json.NewEncoder(writer).Encode(CreateNodeAverageReport(testNodeMetricsMap, timeSlice))
            if err != nil {
                t.Fatalf("failed writing json: %s", err)
            }
            writer.Flush()

            goldenPath := filepath.Join("testdata", filepath.FromSlash(t.Name()) + ".golden")


            if *update {

                t.Log("update golden file")
                if err := ioutil.WriteFile(goldenPath, buffer.Bytes(), 0644); err != nil {
                    t.Fatalf("failed to update golden file %s: %s", goldenPath, err)
                }

             }


             goldenData, err := ioutil.ReadFile(goldenPath)

             if err != nil {
                 t.Fatalf("failed reading .golden file %s: %s", goldenPath, err)
             }

             t.Log(string(buffer.Bytes()))

             if !bytes.Equal(buffer.Bytes(), goldenData) {
                 t.Errorf("bytes do not match .golden file %s", goldenPath)
             }

         })
    }

}



// TestCreateProcessAverageReport - Golden test for CreateProcessAverageReport function
func TestCreateProcessAverageReport(t *testing.T) {

    timeSlice := float64(TIMESLICE)
    processName := "proc5"

    testtable := []struct {
        tname string
    }{
        {
            tname: "ok",
        },
    }

    CreateDummyProcessMetrics(testNodeMetricsMap)

    for _, tc := range testtable {

        t.Run(tc.tname, func(t *testing.T) {

            var buffer bytes.Buffer
            writer := bufio.NewWriter(&buffer)

            err := json.NewEncoder(writer).Encode(CreateProcessAverageReport(testNodeMetricsMap, processName, timeSlice))
            if err != nil {
                t.Fatalf("failed writing json: %s", err)
            }
            writer.Flush()

            goldenPath := filepath.Join("testdata", filepath.FromSlash(t.Name()) + ".golden")


            if *update {

                t.Log("update golden file")
                if err := ioutil.WriteFile(goldenPath, buffer.Bytes(), 0644); err != nil {
                    t.Fatalf("failed to update golden file %s: %s", goldenPath, err)
                }

             }


             goldenData, err := ioutil.ReadFile(goldenPath)

             if err != nil {
                 t.Fatalf("failed reading .golden file %s: %s", goldenPath, err)
             }

             t.Log(string(buffer.Bytes()))

             if !bytes.Equal(buffer.Bytes(), goldenData) {
                 t.Errorf("bytes do not match .golden file %s", goldenPath)
             }

         })
    }

}


// TestCreateProcessHistoryReport - Golden test for CreateProcessHistoryReport function
func TestCreateProcessHistoryReport(t *testing.T) {

    timeSlice := float64(TIMESLICE)

    testtable := []struct {
        tname string
    }{
        {
            tname: "ok",
        },
    }

    CreateDummyProcessMetricsHistory(&testProcessMetricsArray)

    for _, tc := range testtable {

        t.Run(tc.tname, func(t *testing.T) {

            var buffer bytes.Buffer
            writer := bufio.NewWriter(&buffer)

            err := json.NewEncoder(writer).Encode(CreateProcessHistoryReport(testProcessMetricsArray, timeSlice))
            if err != nil {
                t.Fatalf("failed writing json: %s", err)
            }
            writer.Flush()

            goldenPath := filepath.Join("testdata", filepath.FromSlash(t.Name()) + ".golden")


            if *update {

                t.Log("update golden file")
                if err := ioutil.WriteFile(goldenPath, buffer.Bytes(), 0644); err != nil {
                    t.Fatalf("failed to update golden file %s: %s", goldenPath, err)
                }

             }


             goldenData, err := ioutil.ReadFile(goldenPath)

             if err != nil {
                 t.Fatalf("failed reading .golden file %s: %s", goldenPath, err)
             }

             t.Log(string(buffer.Bytes()))

             if !bytes.Equal(buffer.Bytes(), goldenData) {
                 t.Errorf("bytes do not match .golden file %s", goldenPath)
             }

         })
    }

}





