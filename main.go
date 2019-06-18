
package main

import (
    "log"
    "net/http"
    "fmt"
    "strconv"
    "os"

    h "github.com/zilard/mongotix/handlers"
    u "github.com/zilard/mongotix/handlers/utils"
    d "github.com/zilard/mongotix/dao"
    c "github.com/zilard/mongotix/config"

    "github.com/gorilla/mux"
    "github.com/spf13/cobra"
    "github.com/juju/fslock"

)

const PORT int = 8080
var Port int
var StoragePath string
var config c.Config



const (
    nodeMetricsFile string = "node_metrics.txt"
    nodeMetricsReadLockFile string = "node_metrics_read_lock.txt"
    nodeMetricsWriteLockFile string = "node_metrics_write_lock.txt"
    processMetricsHistoryFile string = "process_metrics_history.txt"
    processMetricsHistoryReadLockFile string = "process_metrics_history_read_lock.txt"
    processMetricsHistoryWriteLockFile string = "process_metrics_history_write_lock.txt"
)


func initialize() {

    u.HostName, _ = os.Hostname()

    setStorageDir()

    h.NodeMetricsReadLock = fslock.New(StoragePath + "/" + nodeMetricsReadLockFile)
    h.NodeMetricsWriteLock = fslock.New(StoragePath + "/" + nodeMetricsWriteLockFile)
    h.ProcessMetricsHistoryReadLock = fslock.New(StoragePath + "/" + processMetricsHistoryReadLockFile)
    h.ProcessMetricsHistoryWriteLock = fslock.New(StoragePath + "/" + processMetricsHistoryWriteLockFile)

    h.NodeMetricsFilePath = StoragePath + "/" + nodeMetricsFile
    h.ProcessMetricsHistoryFilePath = StoragePath + "/" + processMetricsHistoryFile

    checkOrCreateFile(h.NodeMetricsFilePath)
    checkOrCreateFile(h.ProcessMetricsHistoryFilePath)

}


func setStorageDir() {
    dataStorage := os.Getenv("METRIX_DATA_STORAGE")
    if dataStorage != "" {
        if StoragePath == "" {
            StoragePath = dataStorage
        }
    } else {
        fmt.Printf("Environment Variable METRIX_DATA_STORAGE is not set\n")
        if StoragePath == "" {
            StoragePath = "."
        }
    }
    fmt.Printf("Final Metrix Data Storage is: %s\n", StoragePath)
}


func checkOrCreateFile(filePath string) {
    file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
    file.Close()
    if err != nil {
        fmt.Printf("ERROR opening or creating file %s => %v\n", filePath, err)
        os.Exit(1)
    }
}


var RootCmd = &cobra.Command{
    Use: "mongotix",
    Run: func(cmd *cobra.Command, args []string) {

        if Port < 1024 || Port > 65535 {
            fmt.Printf("Port number out of range!\nPlease use a port number between 1024 and 65535\n")
            return
        }

        Run()

    },
}


// Adding flag for optional port number that can be specified
func init() {

    config = c.Config{}
    d.dao = d.DAO{}

    config.Read()

    d.dao.Server = config.Server
    d.dao.Port = config.Port
    d.dao.Database = config.Database

    d.dao.Connect()

    RootCmd.PersistentFlags().IntVarP(&Port, "port", "p", PORT, "port")
    RootCmd.PersistentFlags().StringVarP(&StoragePath, "datadir", "d", "", "datadir")

}



func main() {

    if err := RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

}


// Run calling HTTP Handlers functions for each specific API path
func Run() {

    initialize()

    router := mux.NewRouter()

    router.HandleFunc("/v1/metrics/node/{nodename}", h.CreateNodeMetrics).Methods("POST")
    router.HandleFunc("/v1/metrics/node/{nodename}/", h.CreateNodeMetrics).Methods("POST")
    router.HandleFunc("/v1/metrics/nodes/{nodename}/process/{processname}", h.CreateProcessMetrics).Methods("POST")
    router.HandleFunc("/v1/metrics/nodes/{nodename}/process/{processname}/", h.CreateProcessMetrics).Methods("POST")

    router.HandleFunc("/v1/analytics/nodes/average", h.GetAllNodeAverageMetrics).Methods("GET")

    router.HandleFunc("/v1/analytics/processes/{processname}", h.GetProcessAverageMetricsAllNodes).Methods("GET")
    router.HandleFunc("/v1/analytics/processes/{processname}/", h.GetProcessAverageMetricsAllNodes).Methods("GET")
    router.HandleFunc("/v1/analytics/processes", h.GetMostRecentProcesses).Methods("GET")
    router.HandleFunc("/v1/analytics/processes/", h.GetMostRecentProcesses).Methods("GET")

    fmt.Printf("SERVER LISTENING ON :%d\n", Port)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(Port), router))
}



