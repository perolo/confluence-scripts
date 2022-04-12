package schedulerutil

import (
	"encoding/json"
	"fmt"
	"github.com/perolo/jira-scripts/jirautils"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

type ScheduleType struct {
	Report    string          `properties:"report"`
	LastTime  time.Time       `properties:"lasttime"`
	StartTime time.Time       `properties:"starttime"`
	ExecTimes []time.Duration `properties:"exectimes"`
}

type SchedFunc func(string)

var AllSchedules map[string]ScheduleType
var filename = "C:\\Users\\perolo\\Downloads\\schedules.json"

func DummyFunc(string) {
	//fmt.Printf("DummyFunc\n")
}

func SetSchedulerFile(name string) {
	filename = name
}

func CheckSchedule(duration time.Duration, reset bool, dofunc SchedFunc, propfile string) bool {
	if filename == "" {
		panic(nil)
	}
	funcname := runtime.FuncForPC(reflect.ValueOf(dofunc).Pointer()).Name()

	return CheckScheduleDetail(funcname, duration, reset, dofunc, propfile)
}

func CheckScheduleDetail(report string, duration time.Duration, reset bool, dofunc SchedFunc, propfile string) bool {
	fmt.Println("--> ", report)
	if AllSchedules == nil {
		AllSchedules = readSched()
	}
	if _, ok := AllSchedules[report]; !ok {
		var newsched ScheduleType
		newsched.Report = report
		newsched.LastTime = time.Now().Add(-(duration + time.Hour))
		newsched.StartTime = time.Now()
		AllSchedules[report] = newsched
	}
	if (time.Since(AllSchedules[report].LastTime) > duration) || reset {
		uppdsched := AllSchedules[report]
		uppdsched.StartTime = time.Now()
		dofunc(propfile) // add bool return?
		uppdsched.LastTime = time.Now()
		uppdsched.ExecTimes = append(uppdsched.ExecTimes, uppdsched.LastTime.Sub(uppdsched.StartTime))
		AllSchedules[report] = uppdsched
		saveSched()
		return true
	}
	return false
}

func saveSched() {
	//	theFile := "schedules.json"
	body, err := json.Marshal(AllSchedules)
	if err != nil {
		log.Fatal(err)
	}

	//	f, err := ioutil.TempFile(os.TempDir(), theFile)
	/*
		name := filepath.Join("C:\\Users\\perolo\\Downloads", theFile)
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0600)

		Check(err)
		_, err = f.Write(buf)
		Check(err)
		err = f.Close()
		Check(err)
	*/

	err = ioutil.WriteFile(filename, body, 0600)
	jirautils.Check(err)
}

func readSched() map[string]ScheduleType {
	var tmp map[string]ScheduleType

	jsonFile, err := os.Open(filename)

	//	theFile := "schedules.json"
	//	name := filepath.Join("C:\\Users\\perolo\\Downloads", theFile)
	//	f, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err == nil {
		byteValue, err := ioutil.ReadAll(jsonFile)
		jirautils.Check(err)

		err = json.Unmarshal(byteValue, &tmp)
		if err != nil {
			log.Fatal(err)
		}
		//os.Close

	} else {
		tmp = make(map[string]ScheduleType, 10)
	}

	return tmp
}
