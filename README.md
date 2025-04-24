# How to Run
To run 


`go run server.go --kubeconfig=<kubeconfig_path> --port=<port> --max-concurrency=<concurrency>`

# Web APIs

`curl --location 'http://localhost:<port>/jobs/running'`

`curl --location 'http://localhost:<port>/jobs' \
--form 'jobFile=@"<file-path>/job.yaml"' \
--form 'priority="1"'`


`curl --location 'http://localhost:<port>/jobs/pending'`
