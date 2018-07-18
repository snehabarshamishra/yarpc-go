package chooserbenchmark

type ResponseWriter chan struct{}

type RequestWriter chan ResponseWriter
