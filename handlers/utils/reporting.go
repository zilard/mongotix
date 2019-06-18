
package utils

import (
    s "github.com/zilard/mongotix/structs"
)



// CreateNodeAverageReport - creates cpu/mem average report for all the nodes
// and returns back the result to the GetAllNodeAverageMetrics API handler
// The analytics data is extracted from nodeMetricsMap where all Process and Node metrics are stored
func CreateNodeAverageReport(nodeMetricsMap s.NodeMetricsMap, timeSlice float64) s.NodeAverageReport {

    //fmt.Printf("TIMESLICE %v\n", timeSlice)

    //fmt.Printf("NODE METRICS %v\n", nodeMetricsMap)

    nodeAverageAnalytics := make(s.NodeAverageAnalytics)

    for nodeName, nodeData := range nodeMetricsMap {

        if len(nodeData.NodeMeasurementArray) == 0 {
            continue
        }

        //fmt.Printf("NODE NAME %v\n", nodeName)


        nodeAverageReport := s.NodeAverageReport{
                                 TimeSlice: 0,
                                 CpuUsed: 0,
                                 MemUsed: 0,
                             }

        timeS := timeSlice

        for i := range nodeData.NodeMeasurementArray {
            nodeMeasurement := nodeData.NodeMeasurementArray[len(nodeData.NodeMeasurementArray) - 1 - i]
            //fmt.Printf("__NODE MEASUREMENT %v\n", nodeMeasurement)

            if nodeMeasurement.TimeSlice >= timeS {
                nodeAverageReport.TimeSlice += timeS
                nodeAverageReport.CpuUsed += nodeMeasurement.Cpu * timeS
                nodeAverageReport.MemUsed += nodeMeasurement.Mem * timeS
                break
            } else {
                timeS -= nodeMeasurement.TimeSlice
                nodeAverageReport.TimeSlice += nodeMeasurement.TimeSlice
                nodeAverageReport.CpuUsed += nodeMeasurement.Cpu * nodeMeasurement.TimeSlice
                nodeAverageReport.MemUsed += nodeMeasurement.Mem * nodeMeasurement.TimeSlice
            }
        }

        nodeAverageReport.CpuUsed = nodeAverageReport.CpuUsed/nodeAverageReport.TimeSlice
        nodeAverageReport.MemUsed = nodeAverageReport.MemUsed/nodeAverageReport.TimeSlice


        nodeAverageAnalytics[nodeName] = nodeAverageReport

    }

    //fmt.Printf("PER NODE AVERAGE ANALYTICS %v\n", nodeAverageAnalytics)



    totalNodeAverageReport := s.NodeAverageReport{
                                  TimeSlice: 0,
                                  CpuUsed: 0,
                                  MemUsed: 0,
                              }


    for _, nodeAverageReport := range nodeAverageAnalytics {

        if totalNodeAverageReport.TimeSlice == 0 {
            totalNodeAverageReport.TimeSlice = nodeAverageReport.TimeSlice
        } else {
            if totalNodeAverageReport.TimeSlice > nodeAverageReport.TimeSlice {
                totalNodeAverageReport.TimeSlice = nodeAverageReport.TimeSlice
            }
        }

        totalNodeAverageReport.CpuUsed += nodeAverageReport.CpuUsed
        totalNodeAverageReport.MemUsed += nodeAverageReport.MemUsed
    }

    totalNodeAverageReport.CpuUsed = totalNodeAverageReport.CpuUsed/float64(len(nodeAverageAnalytics))
    totalNodeAverageReport.MemUsed = totalNodeAverageReport.MemUsed/float64(len(nodeAverageAnalytics))


    return totalNodeAverageReport

}




// CreateProcessAverageReport - creates cpu/mem average report for a specific process running on multiple nodes
// and returns back the result to the GetMostRecentProcesses API handler
// The analytics data is extracted from nodeMetricsMap where all Process and Node metrics are stored
func CreateProcessAverageReport(nodeMetricsMap s.NodeMetricsMap, processName string, timeSlice float64) s.ProcessAverageReport {

    //fmt.Printf("TIMESLICE %v\n", timeSlice)

    //fmt.Printf("NODE METRICS %v\n", nodeMetricsMap)


    allProcessMetricsArrays := [][]s.ProcessMeasurement{}
    minAvailableTimeSlice := float64(0)

    for _, nodeData := range nodeMetricsMap {

        if _, ok := nodeData.ProcessMeasurementMap[processName]; ok {


            processMeasurementArray := nodeData.ProcessMeasurementMap[processName]

            allProcessMetricsArrays = append(allProcessMetricsArrays, processMeasurementArray)


            timeS := timeSlice

            availableTimeSlice := float64(0)

            for i := range processMeasurementArray {

                processMeasurement := processMeasurementArray[len(processMeasurementArray) - 1 - i]

                if processMeasurement.TimeSlice >= timeS {
                    availableTimeSlice += timeS
                    break
                } else {
                    timeS -= processMeasurement.TimeSlice
                    availableTimeSlice += processMeasurement.TimeSlice
                }
            }


            if minAvailableTimeSlice == 0 {
                minAvailableTimeSlice = availableTimeSlice
            } else {
                if minAvailableTimeSlice > availableTimeSlice {
                    minAvailableTimeSlice = availableTimeSlice
                }
            }


        }
    }


    processAverageReport := s.ProcessAverageReport{
                                TimeSlice: minAvailableTimeSlice,
                                CpuUsed: 0,
                                MemUsed: 0,
                                NumInstances: float64(len(allProcessMetricsArrays)),
                            }


    for _, processMeasurementArray := range allProcessMetricsArrays {

        timeS := minAvailableTimeSlice

        for i := range processMeasurementArray {

            processMeasurement := processMeasurementArray[len(processMeasurementArray) - 1 - i]

            if processMeasurement.TimeSlice >= timeS {
                processAverageReport.CpuUsed += processMeasurement.CpuUsed * timeS
                processAverageReport.MemUsed += processMeasurement.MemUsed * timeS
                break
            } else {
                timeS -= processMeasurement.TimeSlice
                processAverageReport.CpuUsed += processMeasurement.CpuUsed * processMeasurement.TimeSlice
                processAverageReport.MemUsed += processMeasurement.MemUsed * processMeasurement.TimeSlice
            }
        }

    }


    processAverageReport.CpuUsed = processAverageReport.CpuUsed / minAvailableTimeSlice / processAverageReport.NumInstances
    processAverageReport.MemUsed = processAverageReport.MemUsed / minAvailableTimeSlice / processAverageReport.NumInstances


    //fmt.Printf("PROCESS AVERAGE REPORT: %v\n", processAverageReport)


    return processAverageReport

}




// CreateProcessHistoryReport - creates history report for the most recent process metrics submitted to the API Server
// within the timeSlice timeframe, and returns back the result to the GetProcessAverage
// The analytics data is extracted from processMetricsArray where the Process Metrics History is stored
func CreateProcessHistoryReport(processMetricsArray []s.ProcessMetricsByName, timeSlice float64) s.ProcessHistoryReport {

    timeS := timeSlice

    processHistoryReport := s.ProcessHistoryReport{
                                   TimeSlice: 0,
                                   Processes: []s.ProcessInfo{},
                            }

    processes := []s.ProcessInfo{}


    for i := range processMetricsArray {

        processMetrics := processMetricsArray[len(processMetricsArray) - 1 - i]

        processes = append(processes, s.ProcessInfo{
                                          Name: processMetrics.ProcessName,
                                          Url: "/v1/analytics/processes/" + processMetrics.ProcessName + "/",
                                      })

        if processMetrics.MetricsData.TimeSlice >= timeS {
            processHistoryReport.TimeSlice += timeS
            break
        } else {
            timeS -= processMetrics.MetricsData.TimeSlice
            processHistoryReport.TimeSlice += processMetrics.MetricsData.TimeSlice
        }

    }

    processHistoryReport.Processes = processes

    //fmt.Printf("PROCESS HISTORY REPORT: %v\n", processHistoryReport)

    return processHistoryReport


}



