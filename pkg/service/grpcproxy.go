package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kmpp/pkg/dto"
	"github.com/kmpp/pkg/util/grpcproxy/core"
	"github.com/kmpp/pkg/util/grpcproxy/handler"
	"strings"
)


type GrpcproxyService interface {
	GetActiveConns() []string
	CloseActiveConns(host string) error
	GetLists(host, service string, useTLS, restart bool) ([]string, error)
	DescribeFunction(host, funcName string, useTLS bool) (interface{}, error)
	InvokeFunction(host, funcName string, useTLS bool, metadataAddr string, vars map[string]interface{}) (interface{}, error)
}

type grpcproxyService struct {
	g *core.GrpCox
}

func NewGrpcproxyService() GrpcproxyService{
	return &grpcproxyService {
		g: core.InitGrpCox(),
	}
}

func (gs *grpcproxyService) GetActiveConns() []string{
	return gs.g.GetActiveConns(context.TODO())
}

func (gs *grpcproxyService) CloseActiveConns(host string) error{
	err := gs.g.CloseActiveConns(strings.Trim(host, " "))
	if err != nil {
		return err
	}
	return errors.New("success")
}

func (gs *grpcproxyService) GetLists(host, service string, useTLS, restart bool) ([]string, error){
	res, err := gs.g.GetResource(context.Background(), host, !useTLS, restart)
	if err != nil {
		return nil, err
	}
	result, err := res.List(service)
	if err != nil {
		return nil, err
	}
	gs.g.Extend(host)
	return result, nil
}

func (gs *grpcproxyService) DescribeFunction(host, funcName string, useTLS bool) (interface{}, error) {
	res, err := gs.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		return nil, err
	}

	result, _, err := res.Describe(funcName)
	if err != nil {
		return nil, err
	}
	match := handler.ReGetFuncArg.FindStringSubmatch(result)
	if len(match) < 2 {
		return nil, err
	}

	// describe func
	result, template, err := res.Describe(match[1])
	if err != nil {
		return nil, err
	}

	type desc struct {
		Schema   string `json:"schema"`
		Template string `json:"template"`
	}

	gs.g.Extend(host)

	var desc1 dto.Desc
	var mapSchema    map[string]interface{}
	var mapTemplate  map[string]interface{}

	err = json.Unmarshal([]byte(result), &mapSchema)
	if err != nil {
		desc1.Schema = result
	} else {
		desc1.MapSchema = mapSchema
	}

	err = json.Unmarshal([]byte(template), &mapTemplate)
	if err != nil {
		desc1.Template = result
	} else {
		desc1.MapTemplate = mapTemplate
	}
	return desc1, nil
}

func (gs *grpcproxyService) InvokeFunction(host, funcName string, useTLS bool, metadataHeader string, vars map[string]interface{}) (interface{}, error) {
	res, err := gs.g.GetResource(context.Background(), host, !useTLS, false)
	if err != nil {
		return nil, err
	}

	metadataArr := strings.Split(metadataHeader, ",")
	var metadata []string
	var metadataStr string
	for i, m := range metadataArr {
		i += 1
		if isEven := i % 2 == 0; isEven {
			metadataStr = metadataStr+m
			metadata = append(metadata, metadataStr)
			metadataStr = ""
			continue
		}
		metadataStr = fmt.Sprintf("%s:", m)
	}

	mjson,_ :=json.Marshal(vars)
	mString :=string(mjson)
	in := strings.NewReader(mString)

	// get param
	result, timer, err := res.Invoke(context.Background(), metadata, funcName, in)
	if err != nil {
		return nil, err
	}

	gs.g.Extend(host)

	invRes1 := dto.InvRes{
		Time:   timer.String(),
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal([]byte(result), &mapResult)
	if err != nil {
		invRes1.Result = result
	} else {
		invRes1.MapResult = mapResult
	}

	return invRes1, nil
}