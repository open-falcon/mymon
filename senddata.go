/*
* Open-Falcon
*
* Copyright (c) 2014-2018 Xiaomi, Inc. All Rights Reserved.
*
* This product is licensed to you under the Apache License, Version 2.0 (the "License").
* You may not use this product except in compliance with the License.
*
* This product may include a number of subcomponents with separate copyright notices
* and license terms. Your use of these subcomponents is subject to the terms and
* conditions of the subcomponent's license, as noted in the LICENSE file.
 */

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/open-falcon/mymon/common"
)

// SendData Post the json of all result to falcon-agent
func SendData(conf *common.Config, data []*MetaData) ([]byte, error) {
	data = filterIgnoreData(conf, data)
	js, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	Log.Info("Send to %s, size: %d", conf.Base.FalconClient, len(data))
	for _, m := range data {
		Log.Info("%s", m)
	}

	res, err := http.Post(conf.Base.FalconClient, "Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()
	return ioutil.ReadAll(res.Body)
}

func parseLine(line string) map[string]string {
	var parseRes map[string]string
	// remove space and \n
	line = strings.Replace(strings.TrimSpace(line), "\n", "", -1)

	// match metric, tag and value
	reMetricTagValue, _ := regexp.Compile("^([0-9A-Za-z_,]+)" + TagSplitChar + "?([0-9A-Za-z_,=]*)" + ValueSplitChar + "?([0-9A-Za-z_]*)$")
	matchMetricTagValue := reMetricTagValue.FindSubmatch([]byte(line))
	if len(matchMetricTagValue) > 0 {
		parseRes = map[string]string{
			"metric": string(matchMetricTagValue[1]),
			"tag":    string(matchMetricTagValue[2]),
			"value":  string(matchMetricTagValue[3]),
		}
	} else {
		Log.Info("Error format of ignorefile: %s", line)
	}
	return parseRes
}

func filterIgnoreData(conf *common.Config, data []*MetaData) []*MetaData {
	ignoreFile := conf.Base.IgnoreFile
	if ignoreFile == "" {
		return data
	}
	f, err := os.OpenFile(ignoreFile, os.O_RDONLY, 0644)
	// ignorefile does not exists
	if err != nil {
		return data
	}
	inputReader := bufio.NewReader(f)
	for {
		// get a line
		line, err := inputReader.ReadString('\n')
		// over if EOF
		if err != nil {
			break
		}
		// jump head
		if strings.Contains(line, "FalconIgnore") {
			continue
		}
		// jump annotation
		if strings.HasPrefix(line, "#") {
			continue
		}

		metricTagValue := parseLine(line)
		metric := metricTagValue["metric"]
		// wrong format of ignorefile
		if metric == "" {
			continue
		}
		tag := metricTagValue["tag"]
		value := metricTagValue["value"]

		for i, eachData := range data {
			if metric != eachData.Metric && metric != "*" {
				continue
			}
			if tag != "" && !tagSame(tag, eachData.Tags) {
				continue
			}
			if value != "" {
				data[i].SetValue(value)
				continue
			}
			eachData.SetName("_" + eachData.Metric)
		}
	}
	return data
}

func tagSame(tag1, tag2 string) bool {
	x, y := strings.Split(tag1, ","), strings.Split(tag2, ",")
	sort.Strings(x)
	sort.Strings(y)
	return reflect.DeepEqual(x, y)
}

// Snapshot make a record of note, some metric should be noted before sending
func Snapshot(conf *common.Config, note string, fileNameDay string, fileNameOldDay string) error {
	if conf.Base.SnapshotDay < 0 {
		// Just remind but do not stop
		Log.Info("snapshot_day setted error!")
	}
	f, err := os.OpenFile(fileNameDay, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(note)
	if err != nil {
		return err
	}
	e := os.Remove(fileNameOldDay)
	if e != nil {
		// Just remind but do not stop
		Log.Info("Error remove %s, %s", fileNameOldDay, e.Error())
	}
	return err
}
