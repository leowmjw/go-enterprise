# SCENARIO

A Temporal workflow example that show how can make below sub-scenarios by wrapping simple commands using go script; so that a batch run of 10 can be done and it leverages the built in retry from Temporal.  Make it default exponential backoff of max 2 times before mark as fatal.
If it is a fatal error or terminal error; there should return immeditely and no retries.

  Sub-scenario 1a:

    - Use a go script wrapper to open up a non-existent file at about 80% of the time.  It will show a fatal error

  Sub-scenario 1b:

    - Use a go script wrapper will call a mock endpoint (implemented via function method instaed of interface); that will either fail or take longer than 2s timeout 20% of the time

Sub-scenrio using uv to use httpie to call the http function

  Sub-scenario 2a:

    - Call a localhost port 123456 which is not listened so will fail at about 80% of the time; otherwise it will call port 8080 which works


  Sub-scenario 2b:

    - Call a dummy endpoint where there is incorrect auth; so 401 about 20% of time, but otherwise will return a result


## HTTP Batch Scenario

A special scenario that demonstrates batch HTTP requests using the uv command-line utility:

  HTTP Batch Scenario:

    - Executes 10 HTTP GET requests to various websites using the uv command
    - Records the response from each request in the workflow result
    - Demonstrates the ability to execute external commands and process their outputs
    - Can be triggered by setting the ExecutorCmd to "uv" and APIFunction to "http_batch"
    - No script path is needed for this special scenario
