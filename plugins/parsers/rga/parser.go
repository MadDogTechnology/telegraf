package rga

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"io"
	"strconv"
	"strings"
	"time"
)

type RGAParser struct {
	agentIds map[string]struct{}	// whitelisted agents
	defaultTags map[string]string	// might need it sometime...
}

type HistrnxMsg struct {
	Headers struct {
		Protocol string `json:"protocol"`
		CustomerKey string  `json:"customerKey"`
		AgentId string `json:"agentId"`
		File string `json:"file"`
		MessageType string `json:"messageType"`
		Source string `json:"source"`
		TransactionId string `json:"transactionId"`
		TransmissionAttempts string `json:"transmissionAttempts"`
		CustomerId string `json:"customerId"`
		SampleId string `json:"sampleId"`
		SampleTs string `json:"sampleTs"`
		SampleSize string `json:"sampleSize"`
		FragmentCnt string `json:"fragmentCnt"`
		FragmentSeq string `json:"fragmentSeq"`
		LastTransmission string `json:"lastTransmission"`
		LineRange string `json:"lineRange"`
	} `json:"headers"`
	Body string `json:"body"`
}

func NewParser(defaultTags map[string]string) *RGAParser {

	wl := make(map[string]struct{})
	wl["495ac95c-adc6-4089-9c5a-71f60100e3e9"] = struct{}{} // McLaren, for testing
	return &RGAParser{agentIds: wl, defaultTags: defaultTags}
}

func (p *RGAParser) Parse(buf []byte) ([]telegraf.Metric, error) {
	metrics := make([]telegraf.Metric, 0)

	// Deserialize the message into header and body
	var m HistrnxMsg
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	// If the agent ID doesn't match one on the whitelist, ignore the message
	_, ok := p.agentIds[m.Headers.AgentId]
	if !ok {
		return metrics, nil
	}

	// The message body (which contains all the tags) is gzipped to save transmission time
	// and space
	var r io.Reader = strings.NewReader(m.Body)
	r = base64.NewDecoder(base64.StdEncoding, r)
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := gz.Close(); err != nil {
			_ = fmt.Errorf("close of zip reader, %v", err)
		}
	}()

	// First four bytes is the number of records in the
	var numRecords int32
	err = binary.Read(gz, binary.BigEndian, &numRecords)
	if err != nil {
		_ = fmt.Errorf("couldn't read number of records while parsing message: %v", err)
		return nil, err
	}

	// Loop through the records...collecting the metrics that are properly formatted
	var recordLen int32
	for i := int32(0); i < numRecords; i++ {
		err = binary.Read(gz, binary.BigEndian, &recordLen)
		if err != nil {
			fmt.Printf("E! Couldn't read record length while parsing message: %v", err)
			return metrics, err
		}
		rbuf := make([]byte, recordLen)
		numRead, err := io.ReadFull(gz, rbuf)
		if err != nil || int32(numRead) != recordLen{
			fmt.Printf("E! Couldn't read record length while parsing message: %v", err)
			return metrics, err
		}
		r, err := p.ParseLine(string(rbuf))
		if err != nil {
			fmt.Printf("I! Badly formated history record...ignored")
			continue
		}
		metrics = append(metrics, r)
	}
	return metrics, nil
}

// History topic records are bar ("|") separated. To parse, the buffer is split into fields
// converted into their proper type.
func (p *RGAParser) ParseLine(rbuf string) (telegraf.Metric, error) {

	t := strings.Split(rbuf, "|")
	if len(t) != 6 {
		 return nil, errors.New("history record truncated, ignored")
	}

	// history topics record time stamps are defined as  unix epoch in milliseconds
	ts, err := strconv.ParseInt(t[4], 0, 64)
	if err != nil {
		return nil, err
	}
	tm := time.Unix(ts/1000, ts % 1000)

	// The third field contains the type of the value, defined as single characters
	var v interface{}
	switch t[2] {
	case "b":
		v, err = strconv.ParseBool(t[5])
	case "n":
		v, err = strconv.ParseFloat(t[5], 64)
	case "e":
		v, err = strconv.ParseInt(t[5], 0, 32)
	case "t":
		v, err = strconv.ParseInt(t[5], 0, 64)
	case "s":
		v = t[5]
	default:
		v = t[5]
	}
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)
	tags["customer"] = t[0]
	tags["agent"] = t[1]

	fields := make(map[string]interface{})
	fields["value"] = v

	m, err := metric.New(t[3], tags, fields, tm)
	if err != nil {
		fmt.Printf("error allocation metric: %v", err)
		return nil, err
	}
	return m, nil
}

// SetDefaultTags set the DefaultTags
func (p *RGAParser) SetDefaultTags(tags map[string]string) {
	p.defaultTags = tags
}
