package main

import "context"

func DataHandlerCount(){
	globalRedisClient.Incr(context.Background(),"count:copa:"+globalConfig.NodeInfo.Name)
}

