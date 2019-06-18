
package handlers


import (
    "encoding/json"
    "net/http"
    "fmt"
    "os"

    s "github.com/zilard/mongotix/structs"
    u "github.com/zilard/mongotix/handlers/utils"

    "github.com/juju/fslock"
    "github.com/gorilla/mux"
)


var NodeMetricsFilePath string
var ProcessMetricsHistoryFilePath string
var NodeMetricsReadLock *fslock.Lock
var NodeMetricsWriteLock *fslock.Lock
var ProcessMetricsHistoryReadLock *fslock.Lock
var ProcessMetricsHistoryWriteLock *fslock.Lock




// CreateProcessMetrics - handler for path /v1/metrics/nodes/{nodename}/process/{processname}/
// Populates the NodeMetricsMap struct and it's internals with the received data, process metrics like
// timeselice, cpu usage and mem usage regarding the specific Node and Process
func CreateProcessMetrics(w http.ResponseWriter, r *http.Request) {

    var nodeMetricsMap s.NodeMetricsMap
    var processMetricsArray s.ProcessMetricsArray

    nodeMetricsMap = make(s.NodeMetricsMap)
    processMetricsArray = s.ProcessMetricsArray{}

    params := mux.Vars(r)
    nodeName := params["nodename"]
    processName := params["processname"]

    var processMeasurement s.ProcessMeasurement
    result := json.NewDecoder(r.Body).Decode(&processMeasurement)

    if result != nil {
        fmt.Fprintf(os.Stderr, "result=%v\n", result)
        return
    }

    //fmt.Printf("GOT => processMeasurement %v for NODE %v and PROCESS %v\n", processMeasurement, nodeName, processName)


    processMetricsData := s.ProcessMetricsByName{}
    processMetricsData.ProcessName = processName
    processMetricsData.MetricsData = processMeasurement



    // processMetricsArray
    u.GetLock(ProcessMetricsHistoryWriteLock, ProcessMetricsHistoryFilePath)
    u.ReadFromFile(ProcessMetricsHistoryFilePath, &processMetricsArray)

    processMetricsArray = append(processMetricsArray, processMetricsData)

    // processMetricsArray
    u.GetLock(ProcessMetricsHistoryReadLock, ProcessMetricsHistoryFilePath)

    u.WriteToFile(ProcessMetricsHistoryFilePath, &processMetricsArray)

    u.ReleaseLock(ProcessMetricsHistoryReadLock, ProcessMetricsHistoryFilePath)
    u.ReleaseLock(ProcessMetricsHistoryWriteLock, ProcessMetricsHistoryFilePath)




    // nodeMetricsMap
    u.GetLock(NodeMetricsWriteLock, NodeMetricsFilePath)
    u.ReadFromFile(NodeMetricsFilePath, &nodeMetricsMap)



    if _, ok := nodeMetricsMap[nodeName]; ok {

        nodeData := nodeMetricsMap[nodeName]

        processMetricsMap := nodeData.ProcessMeasurementMap

        if _, ok := processMetricsMap[processName]; ok {

            processMeasurementArray := processMetricsMap[processName]
            processMeasurementArray = append(processMeasurementArray, processMeasurement)

            processMetricsMap[processName] = processMeasurementArray

        } else {

           processMeasurementArray := []s.ProcessMeasurement{}
           processMeasurementArray = append(processMeasurementArray, processMeasurement)

           processMetricsMap[processName] = processMeasurementArray

        }

        nodeData.ProcessMeasurementMap = processMetricsMap

        nodeMetricsMap[nodeName] = nodeData

    } else {

        nodeData := s.NodeData{}
        processMetricsMap := make(s.ProcessMetricsMap)

        processMeasurementArray := []s.ProcessMeasurement{}
        processMeasurementArray = append(processMeasurementArray, processMeasurement)

        processMetricsMap[processName] = processMeasurementArray

        nodeData.NodeMeasurementArray = []s.NodeMeasurement{}
        nodeData.ProcessMeasurementMap = processMetricsMap

        nodeMetricsMap[nodeName] = nodeData

    }



    // nodeMetricsMap
    u.GetLock(NodeMetricsReadLock, NodeMetricsFilePath)

    u.WriteToFile(NodeMetricsFilePath, &nodeMetricsMap)

    u.ReleaseLock(NodeMetricsReadLock, NodeMetricsFilePath)
    u.ReleaseLock(NodeMetricsWriteLock, NodeMetricsFilePath)



    //fmt.Printf("SET nodeMetricsMap %v\n", nodeMetricsMap)

    json.NewEncoder(w).Encode(processMeasurement)


}




// CreateNodeMetrics - handler for path /v1/metrics/node/{nodename}/
// Populates the NodeMetricsMap struct and it's internals with the received data, node metrics like
// timeselice, cpu percentage and mem percentage regarding the specific Node
func CreateNodeMetrics(w http.ResponseWriter, r *http.Request) {

    var nodeMetricsMap s.NodeMetricsMap
    nodeMetricsMap = make(s.NodeMetricsMap)

    params := mux.Vars(r)
    nodeName := params["nodename"]

    var nodeMeasurement s.NodeMeasurement
    result := json.NewDecoder(r.Body).Decode(&nodeMeasurement)

    if result != nil {
        fmt.Fprintf(os.Stderr, "result=%v\n", result)
        return
    }


    //fmt.Printf("GOT => nodeMeasurement %v for NODE %v\n", nodeMeasurement, nodeName)


    // nodeMetricsMap
    u.GetLock(NodeMetricsWriteLock, NodeMetricsFilePath)
    u.ReadFromFile(NodeMetricsFilePath, &nodeMetricsMap)


    if _, ok := nodeMetricsMap[nodeName]; ok {

        nodeData := nodeMetricsMap[nodeName]
        nodeData.NodeMeasurementArray = append(nodeData.NodeMeasurementArray, nodeMeasurement)
        nodeMetricsMap[nodeName] = nodeData

    } else {

        nodeData := s.NodeData{}
        nodeData.NodeMeasurementArray = append(nodeData.NodeMeasurementArray, nodeMeasurement)
        processMetricsMap := make(s.ProcessMetricsMap)
        nodeData.ProcessMeasurementMap = processMetricsMap
        nodeMetricsMap[nodeName] = nodeData

    }



    // nodeMetricsMap
    u.GetLock(NodeMetricsReadLock, NodeMetricsFilePath)

    u.WriteToFile(NodeMetricsFilePath, &nodeMetricsMap)

    u.ReleaseLock(NodeMetricsReadLock, NodeMetricsFilePath)
    u.ReleaseLock(NodeMetricsWriteLock, NodeMetricsFilePath)



    //fmt.Printf("SET nodeMetricsMap %v\n", nodeMetricsMap)

    json.NewEncoder(w).Encode(nodeMeasurement)

}



