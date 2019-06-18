package dao


import (
    "context"
    "log"
    "time"
    "fmt"

    s "github.com/zilard/mongotix/structs"
    c "github.com/zilard/mongotix/config"

    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/mongodb/mongo-go-driver/mongo"
)


var client *mongo.Client
var dao DAO


const (
    NM = "node_metrics"
    PH = "process_history"
)



type NodeMetrics struct {
    ID                    primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
    NodeMetricsMap        s.NodeMetricsMap       `json:"node_metrics_map,omitempty" bson:"node_metrics_map,omitempty"`
}


type ProcessHistory struct {
    ID                    primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
    ProcessMetricsHistory s.ProcessMetricsArray  `json:"process_metrics_history,omitempty" bson:"process_metrics_history,omitempty"`
}




type DAO struct {
    Server    string
    Port      string
    Database  string
}





func (d *DAO) Connect() {

    ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
    client, _ = mongo.Connect(ctx, fmt.Sprintf("mongodb://%s:%s", d.Server, d.Port))

}




func (d *DAO) CeateNodeMetricsEndpoint(nodeMetricsMap s.NodeMetricsMap) *InsertOneResult {

    nm := NodeMetrics{}
    nm.NodeMetricsMap = nodeMetricsMap

    collection := client.Database(d.Database).Collection(NM)
    ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    result, _ := collection.InsertOne(ctx, nm)

    return result

}



func (d *DAO) CreateProcessHistoryEndpoint(processHistory s.ProcessMetricsHistory) *InsertOneResult {

    ph := ProcessHistory{}
    ph.ProcessHistory= processHistory

    collection := client.Database(d.Database).Collection(PH)
    ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
    result, _ := collection.InsertOne(ctx, ph)

    return result


}













