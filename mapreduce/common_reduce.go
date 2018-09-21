package mapreduce

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//

	// Read input and sort by key
	reduceInput := make(map[string][]string)
	for m := 0; m < nMap; m++ {
		inFileName := reduceName(jobName, m, reduceTask)
		inFile, err := os.Open(inFileName)
		if err != nil {
			log.Fatal("open intermediate file:", err)
		}

		dec := json.NewDecoder(inFile)
		var kv KeyValue
		for {
			err := dec.Decode(&kv)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("decode:", err)
			}
			if reduceInput[kv.Key] == nil {
				reduceInput[kv.Key] = make([]string, 0)
			}
			reduceInput[kv.Key] = append(reduceInput[kv.Key], kv.Value)

		}
		inFile.Close()
	}
	// Open output file
	out, err := os.Create(outFile)
	defer out.Close()
	if err != nil {
		log.Fatal("open merge file:", err)
	}
	enc := json.NewEncoder(out)

	// Reduce per key
	for key, values := range reduceInput {
		reduceOutput := reduceF(key, values)
		kv := KeyValue{key, reduceOutput}
		err := enc.Encode(&kv)
		if err != nil {
			log.Fatal("reduce encode:", err)
		}
	}
	log.Println("reduce done")
}
