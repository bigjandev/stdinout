package main

import (
	"bytes"
	"context"
	// "errors"
	"fmt"
	"os"
	"encoding/csv"
	"strings"
	"github.com/Jeffail/benthos/v3/public/service"

	_ "github.com/Jeffail/benthos/v3/public/components/all"
)

type procer struct {
	text string
}

func(p *procer) Process(ctx context.Context, msg *service.Message) (service.MessageBatch, error){
	b, err := msg.AsBytes()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err = buf.Write(b); err != nil {
		return nil, err
	}

	// we can do processing here like addition or something
	filenames := strings.Split(p.text, ",")
	// if _, err = buf.WriteString(strings.Repeat(p.text, p.count)); err != nil {
	// 	return nil, err
	// }

	lines := make([]string, 0)
	fmt.Println("test text is " + p.text)
	for _, filename := range filenames {
		f, err := os.Open(filename)


		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer f.Close()

		csvReader := csv.NewReader(f)

		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}

		for _, ln := range records {
			lines = append(lines, strings.Join(ln," "))
		}
	}
	

	// parse the string as comma delimited and treat as file name, do some stuff

	for _, line := range lines {
		buf.WriteString("line")
		if _, err = buf.WriteString(line); err != nil {
			return nil, err
		}
	}
	
	msg.SetBytes(buf.Bytes())
	return service.MessageBatch{msg}, nil
}

func(p *procer) Close(ctx context.Context) error{
	return nil
}

func main() {
	spec := service.NewConfigSpec().
		Field(service.NewStringField("text").Default("test"))

	service.RegisterProcessor(
		"proc", spec, 
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error){			
			text, err := conf.FieldString("text")
			if err != nil {
				return nil, err
			}
			return &procer{text}, nil
	})

	// run benthos
	service.RunCLI(context.Background())
}