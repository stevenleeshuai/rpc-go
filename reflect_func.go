package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Strategy struct {
	StrategyName string
	StrategyParma string
}

func CheckStrategyOne(ctx context.Context, strategyParam string) error {
	log.DebugContextf(ctx, "CheckStrategyOne check Success")
	return errors.New("CheckStrategyOne check fail")
}

func CheckStrategyTwo(ctx context.Context, strategyParam string) error {
	log.DebugContextf(ctx, "CheckStrategyTwo check Success")
	return nil
}

var strategyFuncs = map[string]func(ctx context.Context, strategyParam string) error {
	"CheckStrategyOne" : CheckStrategyOne,
	"CheckStrategyTwo" : CheckStrategyTwo,
}

func VerifyStrategy(ctx context.Context, strategyList []Strategy ) error {
	for _, strategy := range strategyList {
		if _, ok := strategyFuncs[strategy.StrategyName]; !ok {
			errMsg := fmt.Sprintf("[策略日志][%s]不存在该策略名", strategy.StrategyName)
			log.ErrorContextf(ctx, errMsg)
			return errors.New(errMsg)
		}

		err := InvokeStrategy(ctx, strategy.StrategyName, strategy.StrategyParma)
		if err != nil {
			return err
		}
	}
	return nil
}

func InvokeStrategy(ctx context.Context, strategyName, strategyParams string) error {
	funcValue := reflect.ValueOf(strategyFuncs[strategyName])
	log.DebugContextf(ctx, "InvokeStrategy: %+v", funcValue.String())
	if !funcValue.IsValid() {
		return errors.New("InvokeStrategy err")
	}

	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(strategyParams),
	}
	result := funcValue.Call(args)
	if len(result) > 0 && !result[0].IsNil() {
		log.ErrorContextf(ctx, "InvokeStrategy err:%s", result[0].Interface().(error))
		return result[0].Interface().(error)
	}
	log.InfoContextf(ctx, "[策略日志][%v]执行成功-退出策略, 入参: %v", strategyName, strategyParams)
	return nil
}

func main() {
	strategyList := []Strategy{
		{StrategyName: "CheckStrategyOne", StrategyParma: "hello"},
		{StrategyName: "CheckStrategyTwo", StrategyParma: "world"},
	}

	VerifyStrategy(context.Background(), strategyList)
}
