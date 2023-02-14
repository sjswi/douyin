package initialize

//
//func InitJaeger(service string) (server.Suite, io.Closer) {
//	cfg, _ := jaegercfg.FromEnv()
//	cfg.ServiceName = service
//	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
//	if err != nil {
//		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
//	}
//	opentracing.InitGlobalTracer(tracer)
//	return internal_opentracing.NewDefaultServerSuite(), closer
//}
