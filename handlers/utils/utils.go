package utils

import (
    "encoding/json"
    "math/rand"
    "io/ioutil"
    "time"

    s "github.com/zilard/mongotix/structs"

    "github.com/juju/fslock"
)



const (
    RandFloor = 1
    RandCeiling = 100
)


var HostName string



func WaitUntilIsUnLocked(lock *fslock.Lock, fileName string) {
    i := 0
    for {
        if locked := IsLocked(lock); locked {
            //fmt.Printf("\n%s - File %s is locked for reading, re-checking soon ... %v\n",
            //           HostName, fileName, locked)
            randValue := RandFloor + rand.Intn(RandCeiling-RandFloor+1)
            sleepTime := time.Duration(randValue) * time.Millisecond
            //fmt.Printf("\n%s - %d Sleeping %v %d\n", HostName, i, sleepTime, int(sleepTime))
            time.Sleep(sleepTime)
            i += 1
        } else {
            //fmt.Printf("\n%s - File %s not locked\n", HostName, fileName)
            break
        }
    }
}


func IsLocked(lock *fslock.Lock) bool {
    lockErr := lock.TryLock()
    if lockErr == fslock.ErrLocked {
        return true
    }
    lock.Unlock()
    return false
}


func GetLock(lock *fslock.Lock, fileName string) {
    i := 0
    for {
        lockErr := lock.TryLock()
        if lockErr != nil {
             //fmt.Printf("\n%s - Lock %s not available now, retrying ... %v\n",
             //           HostName, fileName, lockErr.Error())
             randValue := RandFloor + rand.Intn(RandCeiling-RandFloor+1)
             sleepTime := time.Duration(randValue) * time.Millisecond
             //fmt.Printf("\n%s - %d Sleeping %v %d\n", HostName, i, sleepTime, int(sleepTime))
             time.Sleep(sleepTime)
             i += 1
         } else {
             //fmt.Printf("\n%s - Got the lock %s for writing\n", HostName, fileName)
             break
         }
    }
}


func ReleaseLock(lock *fslock.Lock, fileName string ) {
    //fmt.Printf("%s - releasing lock %s\n", HostName, fileName)
    lock.Unlock()
}


func ReadFromFile(filePath string, v interface{}) {
    var jsonData []byte
    jsonData, _ = ioutil.ReadFile(filePath)
    if processMetricsArrayPtr, ok := v.(*s.ProcessMetricsArray); ok {
         json.Unmarshal([]byte(jsonData), processMetricsArrayPtr)
    } else if nodeMetricsMapPtr, ok := v.(*s.NodeMetricsMap); ok {
         json.Unmarshal([]byte(jsonData), nodeMetricsMapPtr)
    }
}


func WriteToFile(filePath string, v interface{}) {
    var jsonData []byte
    if processMetricsArrayPtr, ok := v.(*s.ProcessMetricsArray); ok {
        jsonData, _ = json.MarshalIndent(*processMetricsArrayPtr, "", " ")
    } else if nodeMetricsMapPtr, ok := v.(*s.NodeMetricsMap); ok {
        jsonData, _ = json.MarshalIndent(*nodeMetricsMapPtr, "", " ")
    }
    if len(jsonData) != 0 {
        ioutil.WriteFile(filePath, jsonData, 0666)
    }
}


